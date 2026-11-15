package settings

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/projectious-work/aibox/internal/project"
	"github.com/tailscale/hujson"
)

// ProjectLayer is the project layer of one selected project: the three
// members of customizations.aibox.workspace (CFG-020).
type ProjectLayer struct {
	result layerResult
}

// projectMembers maps the read members to the only keys they may set.
var projectMembers = map[string]string{"theme": WorkspaceTheme, "mode": WorkspaceMode, "layout": WorkspaceLayout}

// ReadProject reads the project layer from <root>/.devcontainer/devcontainer.json.
// It returns nil when the declaration cannot be read or parsed: inspection
// then fails and there is no view (effectiveConfig.settingsSources note).
// Only theme, mode and layout of customizations.aibox.workspace are read;
// nothing else in the project can change a setting (CFG-022).
func ReadProject(root string) *ProjectLayer {
	canonical, err := project.CanonicalRoot(root)
	if err != nil {
		return nil
	}
	path, err := project.ConfinedRegularFile(canonical, filepath.Join(canonical, projectDeclarationPath))
	if err != nil {
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	contents, err := io.ReadAll(io.LimitReader(file, maxProjectFileBytes+1))
	if err != nil || len(contents) > maxProjectFileBytes {
		return nil
	}
	standard, err := hujson.Standardize(contents)
	if err != nil {
		return nil
	}
	document, err := parseJSON(standard)
	if err != nil || document.kind != jsonObject {
		return nil
	}
	layer := &ProjectLayer{result: layerResult{layer: Project, source: &SourceState{Layer: Project, Path: path, State: "loaded"}}}
	workspace := workspaceObject(document)
	if workspace == nil {
		return layer
	}
	for _, member := range workspace.members {
		name, read := projectMembers[member.name]
		if !read {
			continue
		}
		key, _ := Lookup(name)
		var value any
		ok := member.value.kind == jsonString
		if ok {
			value, ok = convertText(key, member.value.text)
		}
		if !ok {
			layer.result.rejection = reject(InvalidValue, Project, path, key.Name, fmt.Sprintf("%s has a value that is not allowed", key.Name))
			return layer
		}
		assignment := Assignment{Key: key.Name, Layer: Project, Source: path, Value: value}
		replaced := false
		for index := range layer.result.assignments {
			if layer.result.assignments[index].Key == key.Name {
				layer.result.assignments[index], replaced = assignment, true
			}
		}
		if !replaced {
			layer.result.assignments = append(layer.result.assignments, assignment)
		}
	}
	return layer
}

// workspaceObject returns customizations.aibox.workspace when every step is
// an object; any other shape leaves the project layer empty, because the rest
// of customizations.aibox is not validated by this feature.
func workspaceObject(document *jsonNode) *jsonNode {
	node := document
	for _, name := range []string{"customizations", "aibox", "workspace"} {
		node = node.member(name)
		if node == nil || node.kind != jsonObject {
			return nil
		}
	}
	return node
}
