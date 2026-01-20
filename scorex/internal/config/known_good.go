package config

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"

    "scorex/internal/model"
)

// KnownGood represents the structure of known_good.json.
type KnownGood struct {
    Timestamp      string                     `json:"timestamp"`
    Modules        map[string]model.ModuleInfo `json:"modules"`
    ManifestSHA256 string                     `json:"manifest_sha256"`
    Suite          string                     `json:"suite"`
    DurationS      int                        `json:"duration_s"`
}

// LoadKnownGood loads a KnownGood specification from a local file or HTTP(S) URL.
func LoadKnownGood(urlOrPath string) (*KnownGood, error) {
    var data []byte
    var err error

    if strings.HasPrefix(urlOrPath, "http://") || strings.HasPrefix(urlOrPath, "https://") {
        resp, err := http.Get(urlOrPath)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            return nil, fmt.Errorf("HTTP %s", resp.Status)
        }

        data, err = io.ReadAll(resp.Body)
        if err != nil {
            return nil, err
        }
    } else {
        data, err = os.ReadFile(urlOrPath)
        if err != nil {
            return nil, err
        }
    }

    var kg KnownGood
    if err := json.Unmarshal(data, &kg); err != nil {
        return nil, err
    }

    return &kg, nil
}
