/*
Copyright © 2025 Mirko Mälicke

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/hydrocode-de/gotap/internal/config"
	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags.
var Version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gotap",
	Short: "Shim to tap tool-spec",
	Long: `Shim to tap tool-spec compliant metadata.

This tool is used inside docker containers, which were
set up with tool-specification compliant metadata. 
It can be used to verify and convert the metadata.
More info can be found at:

https://vforwater.github.io/tool-specs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		showVersion, _ := cmd.Flags().GetBool("version")
		if showVersion {
			fmt.Fprintln(cmd.OutOrStdout(), cmd.Version)
			return nil
		}
		return cmd.Help()
	},
}

func Execute() {
	// first init the config
	config.Init()
	bindFlags()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = Version
	rootCmd.Flags().BoolP("version", "v", false, "Print gotap version")
	rootCmd.PersistentFlags().String("spec-file", "", "Path to the tool.yml metadata file")
	rootCmd.PersistentFlags().String("input-file", "", "Path to the inputs.json file")
	rootCmd.PersistentFlags().String("citation-file", "", "Path to the CITATION.cff file")
	rootCmd.PersistentFlags().String("license-file", "", "Path to the LICENSE file")
	rootCmd.PersistentFlags().String("output-folder", "", "Output folder for the tool execution metadata")
}

func bindFlags() {
	v := config.GetViper()

	v.BindPFlag("spec_file", rootCmd.PersistentFlags().Lookup("spec-file"))
	v.BindPFlag("input_file", rootCmd.PersistentFlags().Lookup("input-file"))
	v.BindPFlag("citation_file", rootCmd.PersistentFlags().Lookup("citation-file"))
	v.BindPFlag("license_file", rootCmd.PersistentFlags().Lookup("license-file"))
	v.BindPFlag("output_folder", rootCmd.PersistentFlags().Lookup("output-folder"))
}
