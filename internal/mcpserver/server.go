// Package mcpserver exposes the bounded read-only operation over stdio.
package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/projectious-work/aibox/internal/contract"
	"github.com/projectious-work/aibox/internal/logging"
	"github.com/projectious-work/aibox/internal/settings"
	"github.com/projectious-work/aibox/internal/usecase"
)

// InspectWorkspaceInput is the normalized MCP request for the local preview.
// It deliberately has no actor, policy, executable, runtime selector or
// settings key: no tool input changes a setting (CFG-055, D14).
type InspectWorkspaceInput struct {
	RequestID       string `json:"requestId"`
	ProjectRoot     string `json:"projectRoot"`
	EffectiveConfig bool   `json:"effectiveConfig,omitempty"`
}

var inspectWorkspaceSchema = json.RawMessage(`{
  "$schema":"https://json-schema.org/draft/2020-12/schema",
  "type":"object",
  "additionalProperties":false,
  "properties":{
    "requestId":{"type":"string","pattern":"^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$"},
    "projectRoot":{"type":"string","pattern":"^/","maxLength":4096},
    "effectiveConfig":{"type":"boolean","default":false}
  },
  "required":["requestId","projectRoot"]
}`)

// Dependencies are fixed when mcp serve starts.
type Dependencies struct {
	// Version is the server version.
	Version string
	// Settings are the layers loaded at start; each call adds only the
	// project layer of its project.
	Settings *settings.Loaded
	// Logger is the process logger; each call gets the scope of its request ID.
	Logger *logging.Logger
	// Workspace is the inspect_workspace use case.
	Workspace usecase.Workspace
}

// New creates the fixed local registry. It has no operator mode or runtime
// client and never registers a mutation tool.
func New(deps Dependencies) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "aibox", Version: deps.Version}, nil)
	closedWorld := false
	mcp.AddTool(server, &mcp.Tool{
		Name:        "inspect_workspace",
		Description: "Inspect a canonical local workspace without contacting a host runtime.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: &closedWorld},
		InputSchema: inspectWorkspaceSchema,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input InspectWorkspaceInput) (*mcp.CallToolResult, contract.Result, error) {
		result := deps.inspect(ctx, input)
		if err := result.ValidateEnvelope(); err != nil {
			return nil, contract.Result{}, fmt.Errorf("invalid inspect result: %w", err)
		}
		return &mcp.CallToolResult{IsError: result.Outcome != contract.Succeeded}, result, nil
	})
	return server
}

// inspect resolves the project layer of the call's project on top of the
// layers fixed at start, then runs the use case. A rejected project layer
// fails the call before the operation starts, as in the CLI.
func (deps Dependencies) inspect(ctx context.Context, input InspectWorkspaceInput) contract.Result {
	scope := deps.Logger.Scope(input.RequestID)
	resolved, rejection := deps.Settings.Resolve(settings.ReadProject(input.ProjectRoot))
	if rejection != nil {
		scope.EmitStderr(logging.RejectedEvent(rejection.Fields()))
		return contract.InvalidRequestResult(input.RequestID, "mcp", "the settings were rejected: "+string(rejection.Reason))
	}
	request := usecase.InspectRequest{RequestID: input.RequestID, Initiator: "mcp", ProjectRoot: input.ProjectRoot,
		CanonicalOnly: true, EffectiveConfig: input.EffectiveConfig}
	return deps.Workspace.InspectWorkspace(ctx, resolved, request, scope)
}

// Serve runs the local server over stdin/stdout. The caller reserves stdout
// for the MCP protocol and sends events to stderr as JSON Lines (LOG-011).
func Serve(ctx context.Context, deps Dependencies) error {
	return New(deps).Run(ctx, &mcp.StdioTransport{})
}
