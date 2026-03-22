package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aidfm",
	Short: "AppImage Desktop File Manager",
	Long:  "An AppImage manager that automates desktop integration and registry tracking.",
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
