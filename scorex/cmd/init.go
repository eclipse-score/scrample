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
        IsApplication: true,
        UseFeo: false,
    }

    if err := generateSkeleton(props); err != nil {
        return err
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
