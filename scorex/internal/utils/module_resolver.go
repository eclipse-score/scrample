package utils

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "scorex/internal/model"
)

// ResolveModuleWithFallback tries to resolve a module from knownGood.
// If not present, it falls back to GitHub "latest release" for the given repo name.
func ResolveModuleWithFallback(
    name string,
    knownGood map[string]model.ModuleInfo,
) (model.ModuleInfo, error) {
    if mi, ok := knownGood[name]; ok {
        return mi, nil
    }

    // Fallback: neuester Commit auf main von GitHub
    latestCommit, err := fetchLatestGithubMainCommit("eclipse-score", name)
    if err != nil {
        return model.ModuleInfo{}, fmt.Errorf(
            "module %q not in known_good and GitHub lookup failed: %w",
            name, err,
        )
    }

    mi := model.ModuleInfo{
        Version: "0.1.0",
        Hash:    latestCommit,
        Repo:    fmt.Sprintf("https://github.com/eclipse-score/%s.git", name),
        Branch:  "main",
    }

    return mi, nil
}

func fetchLatestGithubMainCommit(owner, repo string) (string, error) {
    url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/main", owner, repo)

    client := http.Client{
        Timeout: 5 * time.Second,
    }

    resp, err := client.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("unexpected status %d from GitHub", resp.StatusCode)
    }

    var payload struct {
        SHA string `json:"sha"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
        return "", err
    }

    if payload.SHA == "" {
        return "", fmt.Errorf("no sha in GitHub response")
    }
    return payload.SHA, nil
}