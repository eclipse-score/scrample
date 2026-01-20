package model

// KnownGood represents the structure of known_good.json.
type KnownGood struct {
	Timestamp      string                `json:"timestamp"`
	Modules        map[string]ModuleInfo `json:"modules"`
	ManifestSHA256 string                `json:"manifest_sha256"`
	Suite          string                `json:"suite"`
	DurationS      int                   `json:"duration_s"`
}
