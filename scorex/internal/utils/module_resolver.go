package utils

import (
    "scorex/internal/model"
    "scorex/internal/service/module"
)

// Deprecated: use scorex/internal/service/module.ResolveModuleWithFallback.
func ResolveModuleWithFallback(name string, knownGood map[string]model.ModuleInfo) (model.ModuleInfo, error) {
    return module.ResolveModuleWithFallback(name, knownGood)
}

// Deprecated: use scorex/internal/service/module.ResolveModules.
func ResolveModules(modules []string, knownGood map[string]model.ModuleInfo) (map[string]model.ModuleInfo, error) {
    return module.ResolveModules(modules, knownGood)
}