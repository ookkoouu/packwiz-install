package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "packwiz-install",
	Short:   "A tool to install and update modpack made by packwiz.",
	Version: version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate("{{.Version}}\n")
	rootCmd.Flags().BoolP("version", "v", false, "Print version and quit")
}
