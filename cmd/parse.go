/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hydrocode-de/gotap/internal/io"
	"github.com/hydrocode-de/gotap/internal/validation"
	"github.com/spf13/cobra"
)

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse and validate inputs.json, output validated ToolInput as JSON to stdout",
	Long:  ``,
	Run:   parse,
}

func parse(cmd *cobra.Command, args []string) {
	_, err := ParseInputs(cmd, args)
	cobra.CheckErr(err)
	os.Exit(0)
}

func ParseInputs(cmd *cobra.Command, args []string) (bool, error) {
	result, err := validation.LoadAndValidateSpec(args)
	if err != nil {
		return false, err
	}

	if result.ErrorCount() > 0 {
		for _, err := range result.Errors {
			fmt.Fprintln(os.Stderr, io.WriteValidationError(err, true))
		}
		return false, fmt.Errorf("validation failed with %d errors", result.ErrorCount())
	}

	json.NewEncoder(os.Stdout).Encode(result.ToolInput)
	return true, nil
}

func init() {
	rootCmd.AddCommand(parseCmd)

}
