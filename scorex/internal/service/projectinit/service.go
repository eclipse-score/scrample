package projectinit

import (
    "fmt"
    "path/filepath"

    "scorex/internal/config"
    "scorex/internal/model"
    "scorex/internal/service/knowngood"
    "scorex/internal/service/module"
    "scorex/internal/service/skeleton"
)

// Options represents all inputs required to initialize a new project.
type Options struct {
    Modules      []string
    TargetDir    string
    Name         string
    KnownGoodURL string
    BazelVersion string
    ProjectType  string // "Application" or "Module"
    AppType      string // "feo" or "daal"
    Template     string // "simple" or "activities" (for feo apps)
	IncludeDevcontainer bool
}

// Result contains information about the generated project.
type Result struct {
    TargetDir       string
    SelectedModules map[string]model.ModuleInfo
}

// Run performs the full project initialization flow based on the provided options.
func Run(opts Options) (*Result, error) {
    if len(opts.Modules) == 0 {
        return nil, fmt.Errorf("at least one module must be set")
    }

    kg, err := knowngood.Load(opts.KnownGoodURL)
    if err != nil {
        return nil, fmt.Errorf("error loading known_good.json: %w", err)
    }

    selected, err := module.ResolveModules(opts.Modules, kg.Modules)
    if err != nil {
        return nil, err
    }

    targetDir := filepath.Join(opts.TargetDir, opts.Name)

    props := skeleton.Properties{
        ProjectName:     opts.Name,
        SelectedModules: selected,
        BazelVersion:    opts.BazelVersion,
        TargetDir:       targetDir,
        IsApplication:   opts.ProjectType == "Application",
        UseFeo:          opts.AppType == "feo",
        Template:        opts.Template,
		IncludeDevcontainer: opts.IncludeDevcontainer,
    }

    if err := skeleton.Generate(props); err != nil {
        return nil, err
    }

    cfg := &config.ProjectConfig{
        ProjectName:  opts.Name,
        Template:     templateFor(opts.ProjectType, opts.AppType, opts.Template),
        BazelVersion: opts.BazelVersion,
        KnownGoodURL: opts.KnownGoodURL,
        Modules:      opts.Modules,
    }

    if err := config.WriteProjectConfig(targetDir, cfg); err != nil {
        return nil, fmt.Errorf("writing scorex config: %w", err)
    }

    return &Result{
        TargetDir:       targetDir,
        SelectedModules: selected,
    }, nil
}

func templateFor(projectType, appType, template string) string {
    if projectType != "Application" {
        return "module"
    }
    switch appType {
    case "feo":
        if template == "activities" {
            return "feo_app_activities"
        }
        return "feo_app_simple"
    case "daal", "":
        return "daal_app"
    default:
        return "daal_app"
    }
}
