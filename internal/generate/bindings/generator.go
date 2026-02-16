package bindings

import (
	toolspec "github.com/hydrocode-de/tool-spec-go"
)

type Generator interface {
	Generate(spec toolspec.ToolSpec, outputPath string) error
}

var generators = map[string]Generator{
	"python":    &PythonGenerator{},
	"r":         &RGenerator{},
	"javascript": &JavaScriptGenerator{},
	"matlab":    &MatlabGenerator{},
}

func GetGenerator(target string) (Generator, bool) {
	g, ok := generators[target]
	return g, ok
}

func SupportedTargets() []string {
	targets := make([]string, 0, len(generators))
	for t := range generators {
		targets = append(targets, t)
	}
	return targets
}
