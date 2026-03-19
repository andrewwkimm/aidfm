package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aidfm",
	Short: "AppImage Desktop File Manager",
	Long:  "Manages AppImages as first-class packages. Creates and maintains .desktop files and tracks them via a local registry.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// global flags go here later e.g. --verbose
}
