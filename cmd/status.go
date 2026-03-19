package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourusername/aidfm/internal/registry"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print a summary of managed, broken, and orphaned entries",
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	reg, err := registry.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	if err := syncRegistry(reg); err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	if err := reg.Save(); err != nil {
		return fmt.Errorf("failed to save registry: %w", err)
	}

	var healthy, broken, orphaned int
	for _, e := range reg.Entries {
		switch e.Status {
		case registry.StatusHealthy:
			healthy++
		case registry.StatusBroken:
			broken++
		case registry.StatusOrphaned:
			orphaned++
		}
	}

	fmt.Printf("%d managed, %d broken, %d orphaned\n",
		len(reg.Entries), broken, orphaned)

	_ = healthy // included in total managed count
	return nil
}
