package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/andrewwkimm/aidfm/internal/desktop"
	"github.com/andrewwkimm/aidfm/internal/registry"
)

var showCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Pretty-print all fields of a managed entry",
	Args:  cobra.ExactArgs(1),
	RunE:  runShow,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func runShow(cmd *cobra.Command, args []string) error {
	name := args[0]

	reg, err := registry.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	entry := reg.Find(name)
	if entry == nil {
		return fmt.Errorf("no entry found for %q", name)
	}

	df, err := desktop.Read(entry.Desktop)
	if err != nil {
		return fmt.Errorf("failed to read desktop file: %w", err)
	}

	exec := desktop.ParseExec(df.Get("Exec"))

	fmt.Printf("%-12s %s\n", "Name:", entry.Name)
	fmt.Printf("%-12s %s\n", "Binary:", entry.Binary)
	fmt.Printf("%-12s %s\n", "Icon:", entry.Icon)
	fmt.Printf("%-12s %s\n", "Desktop:", entry.Desktop)
	fmt.Printf("%-12s %s\n", "Scope:", entry.Scope)
	fmt.Printf("%-12s %s\n", "Status:", entry.Status)
	fmt.Printf("%-12s %s\n", "Exec:", df.Get("Exec"))

	if len(exec.Env) > 0 {
		fmt.Printf("%-12s", "Env:")
		for k, v := range exec.Env {
			fmt.Printf(" %s=%s", k, v)
		}
		fmt.Println()
	}

	return nil
}
