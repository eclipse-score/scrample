package model

type ModuleInfo struct {
    Version string `json:"version"`
    Hash    string `json:"hash"`
    Repo    string `json:"repo"`
    Branch  string `json:"branch,omitempty"`
}