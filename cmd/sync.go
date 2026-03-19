package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/aidfm/internal/registry"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Reconcile registry against disk state",
	RunE:  runSync,
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

func runSync(cmd *cobra.Command, args []string) error {
	reg, err := registry.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	if err := syncRegistry(reg); err != nil {
		return err
	}

	if err := reg.Save(); err != nil {
		return fmt.Errorf("failed to save registry: %w", err)
	}

	fmt.Println("sync complete")
	return nil
}

// syncRegistry updates the status of every entry in the registry based on
// current disk state. It does not save — callers are responsible for saving.
func syncRegistry(reg *registry.Registry) error {
	for i := range reg.Entries {
		reg.Entries[i].Status = resolveStatus(&reg.Entries[i])
	}
	return nil
}

// resolveStatus checks disk state for a single entry and returns the appropriate status.
func resolveStatus(e *registry.Entry) registry.Status {
	binaryExists := fileExists(e.Binary)
	desktopExists := fileExists(e.Desktop)

	if !binaryExists && desktopExists {
		return registry.StatusOrphaned
	}

	if !binaryExists || !desktopExists {
		return registry.StatusBroken
	}

	if !isExecutable(e.Binary) {
		return registry.StatusBroken
	}

	return registry.StatusHealthy
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0111 != 0
}
