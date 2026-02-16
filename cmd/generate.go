/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/hydrocode-de/gotap/internal/config"
	"github.com/hydrocode-de/gotap/internal/io"
	"github.com/hydrocode-de/gotap/internal/metadata"
	"github.com/hydrocode-de/gotap/internal/metadata/converters"
	"github.com/hydrocode-de/gotap/internal/validation"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate metadata for this tool",
	Long:  ``,
	Run:   generate,
}

func generate(cmd *cobra.Command, args []string) {
	v := config.GetViper()
	citationFile := v.GetString("citation_file")
	format, _ := cmd.Flags().GetString("format")
	if format == "" {
		format = "schema.org"
	}

	spec, err := validation.LoadSpec(args)
	cobra.CheckErr(err)

	citation, err := io.ReadCitationFile(citationFile)
	if err == nil {
		spec.Citation = citation
	}

	var converter metadata.Converter
	switch format {
	case "schema.org":
		converter = &converters.SchemaOrgConverter{}
	case "nfdi4earth":
		converter = &converters.NFDI4EarthConverter{}
	}

	converter.Ingest(spec)
	data, err := converter.Serialize("")
	cobra.CheckErr(err)

	fmt.Println(string(data))
}

func init() {

	generateCmd.Flags().String("format", "", "Format to generate the metadata in; defaults to schema.org (supported: schema.org, nfdi4earth)")
	rootCmd.AddCommand(generateCmd)
}
