package knowngood

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"

    "scorex/internal/model"
)

// Load loads a KnownGood specification from a local file or HTTP(S) URL.
func Load(urlOrPath string) (*model.KnownGood, error) {
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

    var kg model.KnownGood
    if err := json.Unmarshal(data, &kg); err != nil {
        return nil, err
    }

    return &kg, nil
}
