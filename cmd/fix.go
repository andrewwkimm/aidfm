package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/andrewwkimm/aidfm/internal/desktop"
	"github.com/andrewwkimm/aidfm/internal/registry"
)

var fixCmd = &cobra.Command{
	Use:   "fix [name]",
	Short: "Re-apply setup for a managed entry or all broken entries",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runFix,
}

func init() {
	rootCmd.AddCommand(fixCmd)
}

func runFix(cmd *cobra.Command, args []string) error {
	reg, err := registry.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	if err := syncRegistry(reg); err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	if len(args) == 1 {
		entry := reg.Find(args[0])
		if entry == nil {
			return fmt.Errorf("no entry found for %q", args[0])
		}
		return fixEntry(entry, reg)
	}

	// No args — fix everything that is not healthy
	fixed := 0
	for i := range reg.Entries {
		if reg.Entries[i].Status != registry.StatusHealthy {
			if err := fixEntry(&reg.Entries[i], reg); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to fix %s: %v\n", reg.Entries[i].Name, err)
				continue
			}
			fixed++
		}
	}

	if fixed == 0 {
		fmt.Println("nothing to fix")
	} else {
		fmt.Printf("fixed %d entries\n", fixed)
	}

	return nil
}

func fixEntry(entry *registry.Entry, reg *registry.Registry) error {
	// chmod +x
	if err := os.Chmod(entry.Binary, 0755); err != nil {
		return fmt.Errorf("failed to set executable bit on %s: %w", entry.Binary, err)
	}

	// Rewrite desktop file preserving existing exec line
	df, err := desktop.Read(entry.Desktop)
	if err != nil {
		// Desktop file missing — recreate it
		df = desktop.New(entry.Desktop)
		df.Set("Name", entry.Name)
		df.Set("Exec", entry.Binary)
		if entry.Icon != "" {
			df.Set("Icon", entry.Icon)
		}
	}

	if err := df.Write(); err != nil {
		return fmt.Errorf("failed to write desktop file: %w", err)
	}

	desktopDir := desktopDirForScope(entry.Scope)
	if err := runUpdateDB(desktopDir); err != nil {
		fmt.Fprintf(os.Stderr, "warning: update-desktop-database failed: %v\n", err)
	}

	entry.Status = registry.StatusHealthy
	if err := reg.Save(); err != nil {
		return fmt.Errorf("failed to save registry: %w", err)
	}

	fmt.Printf("fixed %s\n", entry.Name)
	return nil
}
