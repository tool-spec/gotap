package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hydrocode-de/gotap/internal/config"
	"github.com/hydrocode-de/gotap/internal/generate/bindings"
	"github.com/hydrocode-de/gotap/internal/input"
	"github.com/hydrocode-de/gotap/internal/io"
	"github.com/hydrocode-de/gotap/internal/logging"
	"github.com/hydrocode-de/gotap/internal/validation"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute this tool",
	Long: `Validates input data and starts the tool entrypoint.

The run command optionally creates the input.json, if the parameters are 
provided as command line arguments. Next, the inputs are validated and 
finally the tool is executed.`,
	DisableFlagParsing: true,
	Run:                execute,
}

func execute(cmd *cobra.Command, args []string) {
	var failOnWarnings bool
	var generateBindings bool
	filteredArgs := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == "--fail-on-warnings" {
			failOnWarnings = true
			continue
		}
		if arg == "--generate-bindings" {
			generateBindings = true
			continue
		}
		filteredArgs = append(filteredArgs, arg)
	}

	dry, err := PrepareInputs(cmd, filteredArgs)
	cobra.CheckErr(err)

	result, err := validation.LoadAndValidateSpec(filteredArgs)
	cobra.CheckErr(err)

	if result.WarningCount() > 0 && failOnWarnings || result.ErrorCount() > 0 {
		fmt.Println("FAIL")

		for _, warning := range result.Warnings {
			fmt.Println(io.WriteValidationError(warning, true))
		}
		for _, err := range result.Errors {
			fmt.Println(io.WriteValidationError(err, true))
		}
		os.Exit(result.ErrorCount())
	}

	command, err := input.ResolveCommand(result.ToolSpec)
	cobra.CheckErr(err)

	if generateBindings {
		target, outputFile, ok := input.InferBindingTarget(command.Executable)
		if ok {
			specFile := config.GetViper().GetString("spec_file")
			outputPath := filepath.Join(filepath.Dir(specFile), outputFile)
			gen, ok := bindings.GetGenerator(target)
			if ok {
				if err := gen.Generate(result.ToolSpec, outputPath); err != nil {
					cobra.CheckErr(fmt.Errorf("failed to generate bindings: %w", err))
				}
			} else {
				fmt.Fprintln(os.Stderr, "WARNING: --generate-bindings: no generator for target", target)
			}
		} else {
			fmt.Fprintf(os.Stderr, "WARNING: --generate-bindings: unsupported executable %q (python, Rscript, node, matlab, octave)\n", command.Executable)
		}
	}

	if dry {
		fmt.Println(command.Command)
		return
	}

	// execute the command finally. This can later be replaced by
	// by logging, tracing, etc.
	outputFolder := config.GetViper().GetString("output_folder")
	if err := os.MkdirAll(outputFolder, 0755); err != nil {
		cobra.CheckErr(err)
	}

	runContext, err := logging.NewRunContext(outputFolder)
	cobra.CheckErr(err)

	cmdResult, err := input.ExecuteCommand(command, runContext.Env())
	cobra.CheckErr(err)
	cmdResult.RunID = runContext.RunID
	cmdResult.LogFile = runContext.LogFile
	cmdResult.LogFormat = runContext.LogFormat
	cmdResult.LogLevel = runContext.LogLevel

	if cmdResult.Stderr != nil {
		os.WriteFile(filepath.Join(outputFolder, "STDERR"), cmdResult.Stderr, 0644)
	}
	if cmdResult.Stdout != nil {
		os.WriteFile(filepath.Join(outputFolder, "STDOUT"), cmdResult.Stdout, 0644)
	}
	jsonResult, err := json.MarshalIndent(cmdResult, "", "  ")
	if err == nil {
		os.WriteFile(filepath.Join(outputFolder, "_metadata.json"), jsonResult, 0644)
	}
}

func init() {
	runCmd.Flags().Bool("dry", false, "Dry run the tool, returning the new inputs.json, instead of executing the tool.")
	runCmd.Flags().Bool("update-inputs", false, "Update the inputs.json if arguments are provided and the file already exists.")

	runCmd.Flags().Bool("fail-on-warnings", false, "Fail the tool if there are warnings.")
	runCmd.Flags().Bool("generate-bindings", false, "Generate parameter bindings before running (inferred from tool command).")
	rootCmd.AddCommand(runCmd)
}
