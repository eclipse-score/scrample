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
    "embed"
    "io/fs"
    "os"
    "path/filepath"
    "strings"
    "text/template"
    "scorex/internal/utils"
    "scorex/internal/model"
)

//go:embed templates/**
var templatesFS embed.FS

type moduleTemplateData struct {
    ProjectName     string
    SelectedModules map[string]model.ModuleInfo
	BazelVersion    string
}

func renderTemplate(tmplPath, dstPath string, data any) error {
    if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
        return err
    }

    t, err := template.ParseFS(templatesFS, tmplPath)
    if err != nil {
        return err
    }

    f, err := os.Create(dstPath)
    if err != nil {
        return err
    }
    defer f.Close()

    return t.Execute(f, data)
}

func generateSkeleton(props utils.SkeletonProperties) error {
    targetDir := props.TargetDir
    
    if err := os.MkdirAll(targetDir, 0o755); err != nil {
        return err
    }

    data := moduleTemplateData{
        ProjectName:     props.ProjectName,
        SelectedModules: props.SelectedModules,
		BazelVersion:	 props.BazelVersion,
    }

    templatePath := "templates"

    if props.IsApplication {
        templatePath = filepath.Join(templatePath, "application")

        if props.UseFeo {
            templatePath = filepath.Join(templatePath, "feo_app")
        } else {
            templatePath = filepath.Join(templatePath, "daal_app")
        }
    }


    err := fs.WalkDir(templatesFS, templatePath, func(path string, d fs.DirEntry, err error ) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            return nil
        }
        if !strings.HasSuffix(path, ".tmpl") {
            return nil
        }

        // path relative to template's directory bestimmen
        rel, err :=  filepath.Rel(templatePath, path)
        if err != nil {
            return err
        }

        outRel := strings.TrimSuffix(rel, ".tmpl")
        
        // special case: "point.<name>" -> ".<name>"
        base := filepath.Base(outRel)
        if strings.HasPrefix(base, "point.") {
            dir := filepath.Dir(outRel)
            base = "." + strings.TrimPrefix(base, "point.")
            outRel = filepath.Join(dir, base)
        }

        dstPath := filepath.Join(targetDir, outRel)
        return renderTemplate(path, dstPath, data)
    })
    if err != nil {
        return err
    }

    return nil
}
