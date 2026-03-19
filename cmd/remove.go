package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/andrewwkimm/aidfm/internal/registry"
)

var (
	removeYes   bool
	removePurge bool
)

var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a managed entry and its desktop file",
	Args:  cobra.ExactArgs(1),
	RunE:  runRemove,
}

func init() {
	removeCmd.Flags().BoolVar(&removeYes, "yes", false, "skip confirmation prompt")
	removeCmd.Flags().BoolVar(&removePurge, "purge", false, "also delete the AppImage binary and icon")
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	name := args[0]

	reg, err := registry.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	entry := reg.Find(name)
	if entry == nil {
		return fmt.Errorf("no entry found for %q", name)
	}

	if !removeYes {
		confirmed := false
		prompt := fmt.Sprintf("Remove %s and its desktop file?", name)
		if removePurge {
			prompt = fmt.Sprintf("Remove %s, its desktop file, binary, and icon?", name)
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title(prompt).
					Value(&confirmed),
			),
		)

		if err := form.Run(); err != nil {
			return fmt.Errorf("form error: %w", err)
		}

		if !confirmed {
			fmt.Println("aborted")
			return nil
		}
	}

	// Delete desktop file
	if err := os.Remove(entry.Desktop); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove desktop file: %w", err)
	}

	// Purge binary and icon if requested
	if removePurge {
		if err := os.Remove(entry.Binary); err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "warning: failed to remove binary: %v\n", err)
		}
		if entry.Icon != "" {
			if err := os.Remove(entry.Icon); err != nil && !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "warning: failed to remove icon: %v\n", err)
			}
		}
	}

	desktopDir := desktopDirForScope(entry.Scope)
	if err := runUpdateDB(desktopDir); err != nil {
		fmt.Fprintf(os.Stderr, "warning: update-desktop-database failed: %v\n", err)
	}

	reg.Remove(name)
	if err := reg.Save(); err != nil {
		return fmt.Errorf("failed to save registry: %w", err)
	}

	fmt.Printf("removed %s\n", name)
	return nil
}
