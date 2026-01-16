package utils

import (
	"scorex/internal/model"
)

type SkeletonProperties struct {
    ProjectName     string
    SelectedModules map[string]model.ModuleInfo
    BazelVersion    string
    TargetDir       string
    IsApplication   bool
    UseFeo          bool
}