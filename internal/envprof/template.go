package envprof

import (
	"bytes"
	"text/template"

	"github.com/idelchi/godyl/pkg/env"
)

// Template renders a Go text/template using the provided env map.
// Template usage: {{ .FOO }} or {{ .FOO | default "fallback" }}.
func Template(data []byte, env env.Env) ([]byte, error) {
	functions := template.FuncMap{
		"default": func(val, def string) string {
			if val == "" {
				return def
			}

			return val
		},
	}

	tmpl, err := template.New("env").
		Option("missingkey=zero").
		Funcs(functions).
		Parse(string(data))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, env); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
