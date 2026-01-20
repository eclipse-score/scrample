package skeleton

import "scorex/internal/model"

// Properties holds all data required to render a project skeleton.
type Properties struct {
    ProjectName     string
    SelectedModules map[string]model.ModuleInfo
    BazelVersion    string
    TargetDir       string
    IsApplication   bool
    UseFeo          bool
}
