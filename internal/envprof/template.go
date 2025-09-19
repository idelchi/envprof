package envprof

import (
	"bytes"
	"text/template"

	sprig "github.com/go-task/slim-sprig/v3"

	"github.com/idelchi/godyl/pkg/env"
)

// Template renders a Go text/template using the provided env map.
func Template(data []byte, env env.Env) ([]byte, error) {
	tmpl, err := template.New("env").
		Funcs(sprig.FuncMap()).
		Option("missingkey=zero").
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
