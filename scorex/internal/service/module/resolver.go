package module

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
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

    repoName := strings.TrimPrefix(name, "score_")

    // Fallback: neuester Commit auf main von GitHub
    latestCommit, err := fetchLatestGithubMainCommit("eclipse-score", repoName)
    if err != nil {
        return model.ModuleInfo{}, fmt.Errorf(
            "module %q not in known_good and GitHub lookup failed: %w",
            name, err,
        )
    }

    mi := model.ModuleInfo{
        Version: "0.1.0", //TODO get correct version number
        Hash:    latestCommit,
        Repo:    fmt.Sprintf("https://github.com/eclipse-score/%s.git", repoName),
        Branch:  "main",
    }

    return mi, nil
}

// ResolveModules resolves a list of module names against the known-good set,
// automatically prefixing names with "score_" when missing and falling back
// to GitHub if a module is not present in known_good.
func ResolveModules(
    modules []string,
    knownGood map[string]model.ModuleInfo,
) (map[string]model.ModuleInfo, error) {
    selected := make(map[string]model.ModuleInfo, len(modules))

    for _, name := range modules {
        moduleName := name
        if !strings.HasPrefix(moduleName, "score_") {
            moduleName = "score_" + moduleName
        }

        mi, err := ResolveModuleWithFallback(moduleName, knownGood)
        if err != nil {
            return nil, fmt.Errorf("resolving module %q failed: %w", moduleName, err)
        }
        selected[moduleName] = mi
    }

    return selected, nil
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
