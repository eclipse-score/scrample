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
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "bufio"
    "sort"
    "strconv"

    "github.com/spf13/cobra"
    "scorex/internal/config"
    "scorex/internal/utils"
    "scorex/internal/model"
)

type SelectedModule struct {
    Name   string
    Info   model.ModuleInfo
}

var (
    initModules   		[]string
    initTargetDir 		string
	initName			string
    initKGURL     		string
	initBazelVersion	string
    initProjectType     string
    initAppType         string
)

// KnownGood - known_good.json
type KnownGood struct {
    Timestamp     string                 `json:"timestamp"`
    Modules       map[string]model.ModuleInfo  `json:"modules"`
    ManifestSHA256 string                `json:"manifest_sha256"`
    Suite         string                 `json:"suite"`
    DurationS     int                    `json:"duration_s"`
}

// initCmd represents the init command
var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Generates an S-CORE skeleton application",
    Long:  `Generates a new S-CORE project with selected modules.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        if len(initModules) == 0 {
            return runInitInteractive()
        }
        return runInit()
    },
}

func init() {
    rootCmd.AddCommand(initCmd)

    initCmd.Flags().StringSliceVar(&initModules, "module", nil, "S-CORE Module(s), e.g.: score_communication, score_baselibs")
	initCmd.Flags().StringVar(&initName, "name", config.DefaultProjectName, "name of the generated project")
    initCmd.Flags().StringVar(&initTargetDir, "dir", config.DefaultTargetDir, "targetdirectory of the generated project")
    initCmd.Flags().StringVar(
        &initKGURL,
        "known-good-url",
        config.DefaultKnownGoodURL,
        "URL or path to known_good.json",
    )
	initCmd.Flags().StringVar(&initBazelVersion, "bazel-version", config.DefaultBazelVersion, "bazel version to be used in project")
}

func runInit() error {
    if len(initModules) == 0 {
        return fmt.Errorf("at least one --module must be set")
    }

    

    kg, err := loadKnownGood(initKGURL)
    if err != nil {
        return fmt.Errorf("error loading known_good.json: %w", err)
    }

    selected := map[string]model.ModuleInfo{}
    knownGoodModules := kg.Modules

    for _, name := range initModules {
        moduleName := name
        if !strings.HasPrefix(moduleName, "score_") {
            moduleName = "score_" + moduleName
        }

        mi, err := utils.ResolveModuleWithFallback(moduleName, knownGoodModules)
        if err != nil {
            return fmt.Errorf("resolving module %q failed: %w", moduleName, err)
        }
        selected[moduleName] = mi
    }

	initTargetDir = initTargetDir + "/" + initName

    props := utils.SkeletonProperties{
        ProjectName: initName, 
        SelectedModules: selected,
        BazelVersion: initBazelVersion,
        TargetDir: initTargetDir,
        IsApplication: initProjectType == "Application",
        UseFeo: initAppType == "feo",
    }

    if err := generateSkeleton(props); err != nil {
        return err
    }

    cfg := &config.ProjectConfig{
        ProjectName:  initName,
        Template:     "daal_app",
        BazelVersion: initBazelVersion,
        KnownGoodURL: initKGURL,
        Modules:      initModules,
    }
    if err := config.WriteProjectConfig(initTargetDir, cfg); err != nil {
        return fmt.Errorf("writing scorex config: %w", err)
    }

    fmt.Println("Generating skeleton in", initTargetDir, "with modules:", selected)
    return nil
}

func loadKnownGood(urlOrPath string) (*KnownGood, error) {
    var data []byte
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
        var err error
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

func runInitInteractive() error {
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
        initProjectType = "Application"
    case moduleChar:
        initProjectType = "Module"
    default:
        return fmt.Errorf("invalid project type %q (use %s or %s)", v, appChar, moduleChar)
    }

    // application type (nur bei Application)
    if initProjectType == "Application" {
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
            initAppType = "daal"
        case feoChar:
            initAppType = "feo"
        default:
            return fmt.Errorf("invalid application type %q (use %s or %s)", v, feoChar, daalChar)
        }
    }

    // project name
    fmt.Printf("Project name [%s]: ", initName)
    if v, err := readLine(reader); err != nil {
        return err
    } else if v != "" {
        initName = v
    }

    // target directory
    fmt.Printf("Target directory [%s]: ", initTargetDir)
    if v, err := readLine(reader); err != nil {
        return err
    } else if v != "" {
        initTargetDir = v
    }

    // load known-good
    kg, err := loadKnownGood(initKGURL)
    if err != nil {
        return fmt.Errorf("error loading known_good.json: %w", err)
    }

    // choose modules
    modules, err := promptModules(reader, kg.Modules)
    if err != nil {
        return err
    }
    if len(modules) == 0 {
        return fmt.Errorf("no modules selected")
    }
    initModules = modules

    return runInit()
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
    fmt.Print("Select modules (comma-separated indices, e.g. 1,3,5): ")

    sel, err := readLine(r)
    if err != nil {
        return nil, err
    }
    if sel == "" {
        return nil, nil
    }

    parts := strings.Split(sel, ",")
    var result []string
    for _, p := range parts {
        p = strings.TrimSpace(p)
        if p == "" {
            continue
        }
        idx, err := strconv.Atoi(p)
        if err != nil || idx < 1 || idx > len(names) {
            return nil, fmt.Errorf("invalid module index: %q", p)
        }
        result = append(result, names[idx-1])
    }
    return result, nil
}
