package cmd

import (
    "embed"
    "os"
    "path/filepath"
    "text/template"
)

//go:embed templates/**
var templatesFS embed.FS

type moduleTemplateData struct {
    ProjectName     string
    SelectedModules map[string]ModuleInfo
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

func generateSkeleton(targetDir string, selected map[string]ModuleInfo, bazelVersion string) error {
    if err := os.MkdirAll(targetDir, 0o755); err != nil {
        return err
    }

    data := moduleTemplateData{
        ProjectName:     filepath.Base(targetDir),
        SelectedModules: selected,
		BazelVersion:	 bazelVersion,
    }

    if err := renderTemplate(
        "templates/application/MODULE.bazel.tmpl",
        filepath.Join(targetDir, "MODULE.bazel"),
        data,
    ); err != nil {
        return err
    }

    if err := renderTemplate(
        "templates/application/BUILD.tmpl",
        filepath.Join(targetDir, "BUILD"),
        data,
    ); err != nil {
        return err
    }

    if err := renderTemplate(
        "templates/application/src/BUILD.tmpl",
        filepath.Join(targetDir, "src", "BUILD"),
        data,
    ); err != nil {
        return err
    }

    if err := renderTemplate(
        "templates/application/src/main.cpp.tmpl",
        filepath.Join(targetDir, "src", "main.cpp"),
        data,
    ); err != nil {
        return err
    }

	if err := renderTemplate(
        "templates/application/bazelrc.tmpl",
        filepath.Join(targetDir, ".bazelrc"),
        data,
    ); err != nil {
        return err
    }

	if err := renderTemplate(
        "templates/application/bazelversion.tmpl",
        filepath.Join(targetDir, ".bazelversion"),
        data,
    ); err != nil {
        return err
    }

    return nil
}
