// Package mcpserver exposes the bounded V1-03 read-only operation over stdio.
package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/projectious-work/aibox/internal/contract"
	"github.com/projectious-work/aibox/internal/usecase"
)

// InspectWorkspaceInput is the normalized MCP request for the local preview.
// It deliberately has no actor, policy, executable, or runtime selector.
type InspectWorkspaceInput struct {
	RequestID   string `json:"requestId"`
	ProjectRoot string `json:"projectRoot"`
}

var inspectWorkspaceSchema = json.RawMessage(`{
  "$schema":"https://json-schema.org/draft/2020-12/schema",
  "type":"object",
  "additionalProperties":false,
  "properties":{
    "requestId":{"type":"string","pattern":"^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$"},
    "projectRoot":{"type":"string","pattern":"^/","maxLength":4096}
  },
  "required":["requestId","projectRoot"]
}`)

// New creates the fixed V1-03 local registry. It has no operator mode or
// runtime client and never registers a mutation tool.
func New(version string) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "aibox", Version: version}, nil)
	closedWorld := false
	mcp.AddTool(server, &mcp.Tool{
		Name:        "inspect_workspace",
		Description: "Inspect a canonical local workspace without contacting a host runtime.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: &closedWorld},
		InputSchema: inspectWorkspaceSchema,
	}, func(_ context.Context, _ *mcp.CallToolRequest, input InspectWorkspaceInput) (*mcp.CallToolResult, contract.Result, error) {
		result := usecase.InspectWorkspace(input.RequestID, "mcp", input.ProjectRoot, true)
		if err := result.ValidateEnvelope(); err != nil {
			return nil, contract.Result{}, fmt.Errorf("invalid inspect result: %w", err)
		}
		return &mcp.CallToolResult{IsError: result.Outcome != contract.Succeeded}, result, nil
	})
	return server
}

// Serve runs the local server over stdin/stdout. The caller reserves stdout
// for the MCP protocol and sends process diagnostics to stderr.
func Serve(ctx context.Context, version string) error {
	return New(version).Run(ctx, &mcp.StdioTransport{})
}
