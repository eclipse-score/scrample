/********************************************************************************
* Copyright (c) 2025 Contributors to the Eclipse Foundation
*
* See the NOTICE file(s) distributed with this work for additional
* information regarding copyright ownership.
*
* This program and the accompanying materials are made available under the
* terms of the Apache License Version 2.0 which is available at
* https://www.apache.org/licenses/LICENSE-2.0
*
* SPDX-License-Identifier: Apache-2.0
********************************************************************************/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"scorex/internal/config"
	"scorex/internal/model"
	"scorex/internal/service/knowngood"
	"scorex/internal/service/projectinit"
)

type initOptions struct {
	Modules      []string
	TargetDir    string
	Name         string
	KnownGoodURL string
	BazelVersion string
	ProjectType  string // Application|Module
	AppType      string // daal|feo
	IncludeDevcontainer bool
	ModulePreset string
}

var initOpts = initOptions{}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generates an S-CORE skeleton application",
	Long:  `Generates a new S-CORE project with selected modules.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.ValidateModulePresetUsage(initOpts.Modules, initOpts.ModulePreset); err != nil {
			return err
		}

		if err := validateInitOptions(initOpts); err != nil {
			return err
		}

		if len(initOpts.Modules) == 0 {
			if initOpts.ModulePreset != "" {
				if err := applyPresetNonInteractive(&initOpts); err != nil {
					return err
				}
				return runInit(initOpts)
			}
			return runInitInteractive(&initOpts)
		}
		return runInit(initOpts)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initOpts.ProjectType = "Application"
	initOpts.AppType = "daal"

	initCmd.Flags().StringSliceVar(&initOpts.Modules, "module", nil, "S-CORE Module(s), e.g.: score_communication, score_baselibs")
	initCmd.Flags().StringVar(&initOpts.Name, "name", config.DefaultProjectName, "name of the generated project")
	initCmd.Flags().StringVar(&initOpts.TargetDir, "dir", config.DefaultTargetDir, "targetdirectory of the generated project")
	initCmd.Flags().StringVar(
		&initOpts.KnownGoodURL,
		"known-good-url",
		config.DefaultKnownGoodURL,
		"URL or path to known_good.json",
	)
	initCmd.Flags().StringVar(&initOpts.BazelVersion, "bazel-version", config.DefaultBazelVersion, "bazel version to be used in project")
	initCmd.Flags().StringVar(&initOpts.ProjectType, "project-type", initOpts.ProjectType, "project type: Application or Module")
	initCmd.Flags().StringVar(&initOpts.AppType, "app-type", initOpts.AppType, "application type (for Application projects): daal or feo")
	initCmd.Flags().BoolVar(&initOpts.IncludeDevcontainer, "devcontainer", false, "include a .devcontainer folder")
	initCmd.Flags().StringVar(&initOpts.ModulePreset, "module-preset", "", "use a predefined module preset (e.g. feo-standard, daal-standard)")
}

func runInit(opts initOptions) error {
	piOpts := projectinit.Options{
		Modules:             opts.Modules,
		TargetDir:           opts.TargetDir,
		Name:                opts.Name,
		KnownGoodURL:         opts.KnownGoodURL,
		BazelVersion:        opts.BazelVersion,
		ProjectType:         opts.ProjectType,
		AppType:             opts.AppType,
		IncludeDevcontainer: opts.IncludeDevcontainer,
	}

	result, err := projectinit.Run(piOpts)
	if err != nil {
		return err
	}

	fmt.Println("Generating skeleton in", result.TargetDir, "with modules:", result.SelectedModules)
	return nil
}

func runInitInteractive(opts *initOptions) error {
	reader := bufio.NewReader(os.Stdin)

	appChar := "a"
	moduleChar := "m"

	// project type
	fmt.Printf("Project type (%s = application, %s = module): ", appChar, moduleChar)
	v, err := readLine(reader)
	if err != nil {
		return err
	}
	v = strings.ToLower(v)

	switch v {
	case "", appChar: // Default: application
		opts.ProjectType = "Application"
	case moduleChar:
		opts.ProjectType = "Module"
	default:
		return fmt.Errorf("invalid project type %q (use %s or %s)", v, appChar, moduleChar)
	}

	// application type (nur bei Application)
	if opts.ProjectType == "Application" {
		feoChar := "f"
		daalChar := "d"

		fmt.Printf("Application type (%s = FEO, %s = DAAL): ", feoChar, daalChar)
		v, err := readLine(reader)
		if err != nil {
			return err
		}
		v = strings.ToLower(v)

		switch v {
		case "", daalChar: // Default: DAAL
			opts.AppType = "daal"
		case feoChar:
			opts.AppType = "feo"
		default:
			return fmt.Errorf("invalid application type %q (use %s or %s)", v, feoChar, daalChar)
		}
	}

	// project name
	fmt.Printf("Project name [%s]: ", opts.Name)
	if v, err := readLine(reader); err != nil {
		return err
	} else if v != "" {
		opts.Name = v
	}

	// target directory
	fmt.Printf("Target directory [%s]: ", opts.TargetDir)
	if v, err := readLine(reader); err != nil {
		return err
	} else if v != "" {
		opts.TargetDir = v
	}

    // devcontainer
    devcontainer, err := confirm(reader, "Use .devcontainer?")
    if err != nil {
        return err
    } else {
        opts.IncludeDevcontainer = devcontainer
    }

	// load known-good
	kg, err := knowngood.Load(opts.KnownGoodURL)
	if err != nil {
		return fmt.Errorf("error loading known_good.json: %w", err)
	}

	// presets (optional)
	if err := applyPresetInteractive(reader, opts, kg.Modules); err != nil {
		return err
	}
	if len(opts.Modules) > 0 {
		if err := validateInitOptions(*opts); err != nil {
			return err
		}
		return runInit(*opts)
	}

	// choose modules
	modules, err := promptModules(reader, kg.Modules)
	if err != nil {
		return err
	}
	if len(modules) == 0 {
		return fmt.Errorf("no modules selected")
	}
	opts.Modules = modules

	if err := validateInitOptions(*opts); err != nil {
		return err
	}
	return runInit(*opts)
}

func applyPresetNonInteractive(opts *initOptions) error {
	all, err := config.LoadModulePresets()
	if err != nil {
		return err
	}
	p, ok := config.FindModulePreset(all, opts.ModulePreset)
	if !ok {
		return fmt.Errorf("unknown --module-preset %q (known: %s)", opts.ModulePreset, strings.Join(config.KnownPresetIDs(all), ", "))
	}
	if p.ProjectType != "" && p.ProjectType != opts.ProjectType {
		return fmt.Errorf("module preset %q is not applicable to project type %q", p.ID, opts.ProjectType)
	}
	if p.AppType != "" && p.AppType != opts.AppType {
		return fmt.Errorf("module preset %q is not applicable to app type %q", p.ID, opts.AppType)
	}
	if len(p.Modules) == 0 {
		return fmt.Errorf("module preset %q has no modules", p.ID)
	}
	opts.Modules = append([]string(nil), p.Modules...)
	return nil
}

func applyPresetInteractive(reader *bufio.Reader, opts *initOptions, known map[string]model.ModuleInfo) error {
	all, err := config.LoadModulePresets()
	if err != nil {
		return err
	}
	applicable := config.ApplicableModulePresets(all, opts.ProjectType, opts.AppType)
	if len(applicable) == 0 {
		return nil
	}

	fmt.Println("\nModule presets:")
	fmt.Println("  [0] Custom (select manually)")
	for i, p := range applicable {
		fmt.Printf("  [%d] %s (%s)\n", i+1, p.Label, p.ID)
	}
	fmt.Print("Select preset [0]: ")

	v, err := readLine(reader)
	if err != nil {
		return err
	}
	v = strings.TrimSpace(v)
	if v == "" {
		v = "0"
	}

	idx, err := strconv.Atoi(v)
	if err != nil || idx < 0 || idx > len(applicable) {
		return fmt.Errorf("invalid preset selection: %q", v)
	}
	if idx == 0 {
		return nil
	}

	preset := applicable[idx-1]
	if len(preset.Modules) == 0 {
		return fmt.Errorf("selected preset %q has no modules", preset.ID)
	}

	opts.Modules = append([]string(nil), preset.Modules...)

	addMore, err := confirm(reader, "Add more modules on top of the preset?")
	if err != nil {
		return err
	}
	if !addMore {
		return nil
	}

	extra, err := promptModules(reader, known)
	if err != nil {
		return err
	}
	opts.Modules = mergeUnique(opts.Modules, extra)
	return nil
}

func mergeUnique(a, b []string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	out := make([]string, 0, len(a)+len(b))
	for _, s := range a {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	for _, s := range b {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func validateInitOptions(opts initOptions) error {
	validProjectTypes := []string{"Application", "Module"}
	if opts.ProjectType == "" {
		return fmt.Errorf("--project-type must be set (Application or Module)")
	}
	if !slices.Contains(validProjectTypes, opts.ProjectType) {
		return fmt.Errorf("invalid --project-type %q (use Application or Module)", opts.ProjectType)
	}

	if opts.ProjectType == "Application" {
		validAppTypes := []string{"daal", "feo"}
		if opts.AppType == "" {
			return fmt.Errorf("--app-type must be set for Application projects (daal or feo)")
		}
		if !slices.Contains(validAppTypes, opts.AppType) {
			return fmt.Errorf("invalid --app-type %q (use daal or feo)", opts.AppType)
		}
	}

	if opts.Name == "" {
		return fmt.Errorf("--name must be set")
	}
	if opts.TargetDir == "" {
		return fmt.Errorf("--dir must be set")
	}
	if opts.KnownGoodURL == "" {
		return fmt.Errorf("--known-good-url must be set")
	}
	if opts.BazelVersion == "" {
		return fmt.Errorf("--bazel-version must be set")
	}

	return nil
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func promptModules(r *bufio.Reader, known map[string]model.ModuleInfo) ([]string, error) {
	if len(known) == 0 {
		return nil, fmt.Errorf("no modules in known_good.json")
	}

	// sorted list of module names
	names := make([]string, 0, len(known))
	for n := range known {
		names = append(names, n)
	}
	sort.Strings(names)

	fmt.Println("\nAvailable S-CORE modules:")
	for i, n := range names {
		fmt.Printf("  [%d] %s\n", i+1, n)
	}
	fmt.Print("Select modules (comma-separated indices or names, e.g. 1,3 or score_foo): ")

	sel, err := readLine(r)
	if err != nil {
		return nil, err
	}
	if sel == "" {
		return nil, nil
	}

	parts := strings.Split(sel, ",")
	var result []string
	seen := make(map[string]struct{})
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		// Either an index into the known-good list, or a module name.
		if idx, err := strconv.Atoi(p); err == nil {
			if idx < 1 || idx > len(names) {
				return nil, fmt.Errorf("invalid module index: %q", p)
			}
			name := names[idx-1]
			if _, ok := seen[name]; !ok {
				result = append(result, name)
				seen[name] = struct{}{}
			}
			continue
		}

		// Treat as module name.
		name := p
		if _, ok := known[name]; !ok {
			ok, err := confirm(r, fmt.Sprintf("Module %q is not in known_good.json. Add anyway?", name))
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
		}
		if _, ok := seen[name]; !ok {
			result = append(result, name)
			seen[name] = struct{}{}
		}
	}
	return result, nil
}

func confirm(r *bufio.Reader, prompt string) (bool, error) {
	for {
		fmt.Printf("%s (y/N): ", prompt)
		v, err := readLine(r)
		if err != nil {
			return false, err
		}
		v = strings.TrimSpace(strings.ToLower(v))
		switch v {
		case "", "n", "no":
			return false, nil
		case "y", "yes":
			return true, nil
		default:
			// keep asking
		}
	}
}
