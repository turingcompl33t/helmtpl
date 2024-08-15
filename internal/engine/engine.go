package engine

import (
	"bufio"
	"bytes"
	"fmt"
	"helmtpl/internal/logger"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Engine is responsible for performing templating logic,
// transforming input bytes that encode a YAML document
// to an equivalent YAML document with templating applied.
type Engine struct {
	// The key that contains variables for templating
	varKey string

	// The logger instance for the engine
	logger *logger.Logger
}

// Create a new engine instance.
func New(varKey string, logger *logger.Logger) *Engine {
	return &Engine{varKey, logger}
}

// Execute templating logic.
func (e Engine) Run(input []byte) ([]byte, error) {
	e.logger.Debug("executing engine.Run()")

	var data map[string]interface{}
	if err := yaml.Unmarshal(input, &data); err != nil {
		return nil, err
	}

	// Grab the variables key from data
	vars, ok := data[e.varKey].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("could not find variables key %s", e.varKey)
	}
	vars = map[string]interface{}{
		e.varKey: vars,
	}
	delete(data, e.varKey)

	// Recursively template input data, ignoring the variables key
	r, err := e.templateMap(data, vars)
	if err != nil {
		return nil, err
	}

	out, err := yaml.Marshal(r)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// Switch based on the type of the value.
func (e Engine) templateAny(v interface{}, vars map[string]interface{}) (interface{}, error) {
	e.logger.Debug("templating any")

	switch value := v.(type) {
	case map[string]interface{}:
		return e.templateMap(value, vars)
	case string:
		return e.templateString(value, vars)
	default:
		return v, nil
	}
}

// Recursively template a map.
func (e Engine) templateMap(m map[string]interface{}, vars map[string]interface{}) (map[string]interface{}, error) {
	e.logger.Debug("templating map")

	for k, v := range m {
		var err error
		m[k], err = e.templateAny(v, vars)
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

// Execute templating logic on a string.
func (e Engine) templateString(s string, vars map[string]interface{}) (string, error) {
	e.logger.Debug("templating string")

	t, err := template.New("engine").Parse(s)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	w := bufio.NewWriter(&buf)
	if err := t.Execute(w, vars); err != nil {
		return "", err
	}
	w.Flush()

	return buf.String(), nil
}
