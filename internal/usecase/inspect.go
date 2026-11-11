// Package usecase holds operations shared by the CLI and MCP adapters.
package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"path/filepath"

	"github.com/projectious-work/aibox/internal/contract"
	"github.com/projectious-work/aibox/internal/inspection"
)

// NewRequestID creates a correlation identifier without granting authority.
func NewRequestID() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}

// InspectWorkspace executes the same read-only operation for CLI and MCP.
// The initiator is an attribution label, never an authorization decision.
func InspectWorkspace(requestID, initiator, projectRoot string, canonicalOnly bool) contract.Result {
	result := contract.Result{
		SchemaVersion:    contract.ResultSchemaVersion,
		Operation:        contract.InspectWorkspace,
		RequestID:        requestID,
		Actors:           contract.Actors{Initiator: initiator, Executor: "aibox"},
		ChangedResources: []string{}, Warnings: []string{}, Evidence: []string{},
	}
	var err error
	if canonicalOnly && (!filepath.IsAbs(projectRoot) || filepath.Clean(projectRoot) != projectRoot) {
		err = errors.New("MCP project root must be a canonical absolute path")
	}
	var absProject string
	if err == nil {
		absProject, err = filepath.Abs(projectRoot)
	}
	if err == nil {
		root, digest, data, inspectErr := inspection.Workspace(absProject)
		err = inspectErr
		if err == nil && canonicalOnly && root != projectRoot {
			err = errors.New("MCP project root must resolve to itself")
		}
		if err == nil {
			result.Target = &contract.Target{Scope: "local", WorkspaceRoot: root}
			result.InputDigest = digest
			result.Outcome = contract.Succeeded
			result.Data = data
		}
	}
	if err != nil {
		result.Outcome = contract.Failed
		result.Error = &contract.Error{
			Code: "invalid_input", Category: "invalid_input",
			Message:   "cannot inspect the native project declaration",
			Retryable: false, NextAction: "edit_configuration",
		}
	}
	return result
}
