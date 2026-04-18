package template

import (
	"bytes"
	"fmt"
	"strings"
	gotmpl "text/template"
)

// Renderer resolves template expressions in values using secret data.
// Supported syntax: {{ secret "path" "key" }}
type Renderer struct {
	secrets map[string]map[string]string
}

// NewRenderer creates a Renderer pre-loaded with resolved secret data.
func NewRenderer(secrets map[string]map[string]string) *Renderer {
	return &Renderer{secrets: secrets}
}

// Render evaluates a template string and returns the resolved value.
func (r *Renderer) Render(tmpl string) (string, error) {
	if !strings.Contains(tmpl, "{{}") && !strings.Contains(tmpl, "{{") {
		return tmpl, nil
	}

	funcMap := gotmpl.FuncMap{
		"secret": func(path, key string) (string, error) {
			data, ok := r.secrets[path]
			if !ok {
				return "", fmt.Errorf("secret path %q not found", path)
			}
			val, ok := data[key]
			if !ok {
				return "", fmt.Errorf("key %q not found in secret %q", key, path)
			}
			return val, nil
		},
	}

	t, err := gotmpl.New("").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, nil); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderMap applies Render to every value in the provided map.
func (r *Renderer) RenderMap(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		resolved, err := r.Render(v)
		if err != nil {
			return nil, fmt.Errorf("env var %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}
