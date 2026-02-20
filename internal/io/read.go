package io

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexander-lindner/go-cff"
	toolspec "github.com/hydrocode-de/tool-spec-go"
)

func ReadSpecFile(path string) (toolspec.SpecFile, error) {
	specBuffer, err := os.ReadFile(path)
	if err != nil {
		return toolspec.SpecFile{}, fmt.Errorf("failed to read tool spec file: %w", err)
	}

	fileSpec, err := toolspec.LoadToolSpec(specBuffer)
	if err != nil {
		return toolspec.SpecFile{}, fmt.Errorf("failed to load tool spec file: %w", err)
	}

	return fileSpec, nil
}

func ReadInputFile(path string) (toolspec.InputFile, error) {
	inputBuffer, err := os.ReadFile(path)
	if err != nil {
		// Backward-compatible fallback for tools that still mount /in/inputs.json.
		if os.IsNotExist(err) && filepath.Base(path) == "input.json" {
			fallbackPath := filepath.Join(filepath.Dir(path), "inputs.json")
			fallbackBuffer, fallbackErr := os.ReadFile(fallbackPath)
			if fallbackErr == nil {
				input, parseErr := toolspec.LoadInputs(fallbackBuffer)
				if parseErr != nil {
					return toolspec.InputFile{}, fmt.Errorf("failed to load input file: %w", parseErr)
				}
				return input, nil
			}
		}
		return toolspec.InputFile{}, fmt.Errorf("failed to read input file: %w", err)
	}

	input, err := toolspec.LoadInputs(inputBuffer)
	if err != nil {
		return toolspec.InputFile{}, fmt.Errorf("failed to load input file: %w", err)
	}

	return input, nil
}

func ReadCitationFile(path string) (cff.Cff, error) {
	citationBuffer, err := os.ReadFile(path)
	if err != nil {
		return cff.Cff{}, fmt.Errorf("failed to read citation file: %w", err)
	}

	citation, err := cff.Parse(string(citationBuffer))
	if err != nil {
		return cff.Cff{}, fmt.Errorf("failed to parse citation file: %w", err)
	}

	return citation, nil
}

func ReadLicenseFile(path string) (string, error) {
	licenseBuffer, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read license file: %w", err)
	}

	return string(licenseBuffer), nil
}
