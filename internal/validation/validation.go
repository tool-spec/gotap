package validation

import (
	"fmt"

	"github.com/hydrocode-de/gotap/internal/config"
	"github.com/hydrocode-de/gotap/internal/io"
	toolspec "github.com/hydrocode-de/tool-spec-go"
	"github.com/hydrocode-de/tool-spec-go/validate"
)

type ValidationResult struct {
	ToolSpec  toolspec.ToolSpec
	ToolInput toolspec.ToolInput
	Warnings  []*validate.ValidationError
	Errors    []*validate.ValidationError
}

func (r *ValidationResult) ErrorCount() int {
	return len(r.Errors)
}

func (r *ValidationResult) WarningCount() int {
	return len(r.Warnings)
}

func LoadAndValidateSpec(args []string) (ValidationResult, error) {
	v := config.GetViper()
	specFile := v.GetString("spec_file")
	inputFile := v.GetString("input_file")
	citationFile := v.GetString("citation_file")
	licenseFile := v.GetString("license_file")

	warnings := make([]*validate.ValidationError, 0)
	errors := make([]*validate.ValidationError, 0)

	spec, err := io.ReadSpecFile(specFile)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("critical. failed to read tool.yml file: %w", err)
	}

	input, err := io.ReadInputFile(inputFile)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("critical. failed to read inputs.json file: %w", err)
	}

	toolname, err := config.ResolveToolname(args, input)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("critical. %w", err)
	}

	toolSpec, err := spec.GetTool(toolname)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("critical. a tool named %w is not specified in %s", err, specFile)
	}

	toolInput, err := input.GetToolInput(toolname)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("critical. a tool named %w is not specified in %s", err, inputFile)
	}

	citation, err := io.ReadCitationFile(citationFile)
	if err != nil {
		warnings = append(warnings, &validate.ValidationError{
			Field:    "Files",
			Name:     "CITATION.cff",
			Message:  "No citation file found. We recommend adding one.",
			Type:     validate.ErrorType("warning"),
			Expected: "/src/CITATION.cff",
			Actual:   "None",
		})
	} else {
		toolSpec.Citation = citation
	}

	hasErrors, errs := validate.ValidateInputs(toolSpec, toolInput)
	if hasErrors {
		for _, err := range errs {
			errors = append(errors, err)
		}
	}

	_, err = io.ReadLicenseFile(licenseFile)
	if err != nil {
		warnings = append(warnings, &validate.ValidationError{
			Field:    "Files",
			Name:     "LICENSE",
			Message:  "No license file found. We recommend adding one.",
			Type:     validate.ErrorType("warning"),
			Expected: "/src/LICENSE",
			Actual:   "None",
		})
	}

	return ValidationResult{
		ToolSpec:  toolSpec,
		ToolInput: toolInput,
		Warnings:  warnings,
		Errors:    errors,
	}, nil
}

func LoadSpec(args []string) (toolspec.ToolSpec, error) {
	v := config.GetViper()
	specFile := v.GetString("spec_file")

	spec, err := io.ReadSpecFile(specFile)
	if err != nil {
		return toolspec.ToolSpec{}, fmt.Errorf("critical. failed to read tool.yml file: %w", err)
	}

	toolname, err := config.ResolveToolname(args, toolspec.InputFile{})
	if err != nil {
		return toolspec.ToolSpec{}, fmt.Errorf("critical. %w", err)
	}

	toolSpec, err := spec.GetTool(toolname)
	if err != nil {
		return toolspec.ToolSpec{}, fmt.Errorf("critical. a tool named %w is not specified in %s", err, specFile)
	}

	return toolSpec, nil
}
