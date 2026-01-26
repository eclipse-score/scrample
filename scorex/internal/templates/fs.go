package templates

import "embed"

// FS embeds all template files under internal/templates.
//go:embed application/** module/**
var FS embed.FS
