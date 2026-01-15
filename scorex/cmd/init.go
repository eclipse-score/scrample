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
)

type SelectedModule struct {
    Name   string
    Info   ModuleInfo
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
    Modules       map[string]ModuleInfo  `json:"modules"`
    ManifestSHA256 string                `json:"manifest_sha256"`
    Suite         string                 `json:"suite"`
    DurationS     int                    `json:"duration_s"`
}

type ModuleInfo struct {
    Version string `json:"version"`
    Hash    string `json:"hash"`
    Repo    string `json:"repo"`
    Branch  string `json:"branch,omitempty"`
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
	initCmd.Flags().StringVar(&initName, "name", "score_app", "name of the generated project")
    initCmd.Flags().StringVar(&initTargetDir, "dir", ".", "targetdirectory of the generated project")
    initCmd.Flags().StringVar(
        &initKGURL,
        "known-good-url",
        "https://raw.githubusercontent.com/eclipse-score/reference_integration/main/known_good.json",
        "URL or path to known_good.json",
    )
	initCmd.Flags().StringVar(&initBazelVersion, "bazel-version", "8.3.0", "bazel version to be used in project")
}

func runInit() error {
    if len(initModules) == 0 {
        return fmt.Errorf("at least one --module must be set")
    }

    kg, err := loadKnownGood(initKGURL)
    if err != nil {
        return fmt.Errorf("error loading known_good.json: %w", err)
    }

    selected := map[string]ModuleInfo{}
    for _, m := range initModules {
        name := strings.TrimSpace(m)
        info, ok := kg.Modules[name]
        if !ok {
            return fmt.Errorf("module %q is not defined in known_good.json", name)
        }
        selected[name] = info
    }

	initTargetDir = initTargetDir + "/" + initName

    if err := generateSkeleton(initTargetDir, selected, initBazelVersion); err != nil {
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
