package config

import (
    "encoding/json"
    "os"
    "path/filepath"
)

type ProjectConfig struct {
    ProjectName  string   `json:"project_name"`
    Template     string   `json:"template"`
    BazelVersion string   `json:"bazel_version"`
    KnownGoodURL string   `json:"known_good_url"`
    Modules      []string `json:"modules"`
}

const DefaultConfigFileName = "scorex.json"

func WriteProjectConfig(dir string, cfg *ProjectConfig) error {
    if err := os.MkdirAll(dir, 0o755); err != nil {
        return err
    }
    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }
    path := filepath.Join(dir, DefaultConfigFileName)
    return os.WriteFile(path, data, 0o644)
}

func ReadProjectConfig(dir string) (*ProjectConfig, error) {
    path := filepath.Join(dir, DefaultConfigFileName)
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var cfg ProjectConfig
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}