package module

import (
    "scorex/internal/model"
    servicemodule "scorex/internal/service/module"
)

// Deprecated: use scorex/internal/service/module.ResolveModuleWithFallback.
func ResolveModuleWithFallback(name string, knownGood map[string]model.ModuleInfo) (model.ModuleInfo, error) {
    return servicemodule.ResolveModuleWithFallback(name, knownGood)
}

// Deprecated: use scorex/internal/service/module.ResolveModules.
func ResolveModules(modules []string, knownGood map[string]model.ModuleInfo) (map[string]model.ModuleInfo, error) {
    return servicemodule.ResolveModules(modules, knownGood)
}
