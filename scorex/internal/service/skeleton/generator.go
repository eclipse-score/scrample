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
package skeleton

import (
    "io/fs"
    "os"
    "path/filepath"
    "strings"
    "text/template"

    templatesfs "scorex/internal/templates"
)

type moduleTemplateData struct {
    ProjectName     string
    SelectedModules map[string]any
    BazelVersion    string
}

func renderTemplate(tmplPath, dstPath string, data any) error {
    if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
        return err
    }

    t, err := template.ParseFS(templatesfs.FS, tmplPath)
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

func undotifyPath(rel string) string {
    // Work in slash-form, rewrite each segment, then convert back.
    parts := strings.Split(filepath.ToSlash(rel), "/")
    for i, p := range parts {
        if strings.HasPrefix(p, "point.") {
            parts[i] = "." + strings.TrimPrefix(p, "point.")
        }
    }
    return filepath.FromSlash(strings.Join(parts, "/"))
}

// Generate creates a project skeleton based on the provided properties.
func Generate(props Properties) error {
    targetDir := props.TargetDir

    if err := os.MkdirAll(targetDir, 0o755); err != nil {
        return err
    }

    data := moduleTemplateData{
        ProjectName:     props.ProjectName,
        SelectedModules: toAnyMap(props.SelectedModules),
        BazelVersion:    props.BazelVersion,
    }

    templatePath := "module"

    if props.IsApplication {
        if props.UseFeo {
            templatePath = filepath.Join("application", "feo_app")
        } else {
            templatePath = filepath.Join("application", "daal_app")
        }
    }

    err := fs.WalkDir(templatesfs.FS, templatePath, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            return nil
        }
        if !strings.HasSuffix(path, ".tmpl") {
            return nil
        }

        rel, err := filepath.Rel(templatePath, path)
        if err != nil {
            return err
        }

        // Optional: only include .devcontainer when requested.
        slashRel := filepath.ToSlash(rel)
        if !props.IncludeDevcontainer && strings.HasPrefix(slashRel, "point.devcontainer/") {
            return nil
        }

        outRel := strings.TrimSuffix(rel, ".tmpl")
		outRel = undotifyPath(outRel)

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

func toAnyMap[T any](in map[string]T) map[string]any {
    out := make(map[string]any, len(in))
    for k, v := range in {
        out[k] = v
    }
    return out
}
