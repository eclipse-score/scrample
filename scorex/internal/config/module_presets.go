package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

type ModulePreset struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	ProjectType string   `json:"projectType,omitempty"`
	AppType     string   `json:"appType,omitempty"`
	Modules     []string `json:"modules"`
}

type modulePresetFile struct {
	Presets []ModulePreset `json:"presets"`
}

//go:embed module_presets.json
var modulePresetsJSON []byte

func LoadModulePresets() ([]ModulePreset, error) {
	var f modulePresetFile
	if err := json.Unmarshal(modulePresetsJSON, &f); err != nil {
		return nil, fmt.Errorf("parsing embedded module presets: %w", err)
	}

	seen := make(map[string]struct{}, len(f.Presets))
	for i := range f.Presets {
		p := &f.Presets[i]
		p.ID = strings.TrimSpace(p.ID)
		p.Label = strings.TrimSpace(p.Label)
		p.ProjectType = strings.TrimSpace(p.ProjectType)
		p.AppType = strings.TrimSpace(p.AppType)

		if p.ID == "" {
			return nil, fmt.Errorf("module preset missing id")
		}
		if _, ok := seen[p.ID]; ok {
			return nil, fmt.Errorf("duplicate module preset id %q", p.ID)
		}
		seen[p.ID] = struct{}{}

		if p.Label == "" {
			p.Label = p.ID
		}

		// Normalize module names to the score_ prefix (consistent with module resolver).
		for j := range p.Modules {
			p.Modules[j] = normalizeModuleName(p.Modules[j])
		}
		p.Modules = dedupeStrings(p.Modules)
	}

	return f.Presets, nil
}

func ApplicableModulePresets(all []ModulePreset, projectType, appType string) []ModulePreset {
	var out []ModulePreset
	for _, p := range all {
		if p.ProjectType != "" && p.ProjectType != projectType {
			continue
		}
		if p.AppType != "" && p.AppType != appType {
			continue
		}
		out = append(out, p)
	}
	return out
}

func FindModulePreset(all []ModulePreset, id string) (ModulePreset, bool) {
	id = strings.TrimSpace(id)
	for _, p := range all {
		if p.ID == id {
			return p, true
		}
	}
	return ModulePreset{}, false
}

func normalizeModuleName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	if strings.HasPrefix(name, "score_") {
		return name
	}
	return "score_" + name
}

func dedupeStrings(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	// Preserve stable output order: first occurrence wins.
	return out
}

func ValidateModulePresetUsage(modules []string, presetID string) error {
	if presetID == "" {
		return nil
	}
	if len(modules) > 0 {
		return fmt.Errorf("--module and --module-preset are mutually exclusive")
	}
	return nil
}

func KnownPresetIDs(all []ModulePreset) []string {
	ids := make([]string, 0, len(all))
	for _, p := range all {
		ids = append(ids, p.ID)
	}
	slices.Sort(ids)
	return ids
}
