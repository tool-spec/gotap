package cmd

import (
	"fmt"
	"strings"

	toolspec "github.com/hydrocode-de/tool-spec-go"

	"github.com/hydrocode-de/gotap/internal/config"
	"github.com/hydrocode-de/gotap/internal/generate/bindings"
	"github.com/hydrocode-de/gotap/internal/io"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate parameter bindings for a target language",
	Long:  `Generates typed parameter binding code that calls gotap parse and returns validated parameters.`,
	Run:   generateRun,
}

func generateRun(cmd *cobra.Command, args []string) {
	target, _ := cmd.Flags().GetString("target")
	output, _ := cmd.Flags().GetString("output")

	if target == "" {
		cobra.CheckErr(fmt.Errorf("--target is required (python, r, matlab, node-js, node-ts)"))
	}
	if output == "" {
		cobra.CheckErr(fmt.Errorf("--output is required"))
	}

	gen, ok := bindings.GetGenerator(target)
	if !ok {
		cobra.CheckErr(fmt.Sprintf("unsupported target %q; supported: %s",
			target, strings.Join(bindings.SupportedTargets(), ", ")))
	}

	spec, err := resolveSpecForGenerate(args)
	cobra.CheckErr(err)

	err = gen.Generate(*spec, output)
	cobra.CheckErr(err)
}

func init() {
	generateCmd.Flags().String("target", "", "Target language (python, r, matlab, node-js, node-ts)")
	generateCmd.Flags().String("output", "", "Output file path")
	generateCmd.MarkFlagRequired("target")
	generateCmd.MarkFlagRequired("output")
	rootCmd.AddCommand(generateCmd)
}

func resolveSpecForGenerate(args []string) (*toolspec.ToolSpec, error) {
	v := config.GetViper()
	specFile := v.GetString("spec_file")

	spec, err := io.ReadSpecFile(specFile)
	if err != nil {
		return nil, err
	}

	toolname, err := config.ResolveToolname(args, toolspec.InputFile{})
	if err != nil {
		if len(spec.Tools) == 1 {
			for name := range spec.Tools {
				toolname = name
				break
			}
		} else {
			return nil, fmt.Errorf("tool name required when spec has multiple tools: pass as argument or set RUN_TOOL")
		}
	}

	toolSpec, err := spec.GetTool(toolname)
	if err != nil {
		return nil, err
	}

	return &toolSpec, nil
}
