package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/yourusername/aidfm/internal/registry"
)

var editCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Open the desktop file for an entry in $EDITOR",
	Args:  cobra.ExactArgs(1),
	RunE:  runEdit,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func runEdit(cmd *cobra.Command, args []string) error {
	name := args[0]

	reg, err := registry.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	entry := reg.Find(name)
	if entry == nil {
		return fmt.Errorf("no entry found for %q", name)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("$EDITOR is not set")
	}

	c := exec.Command(editor, entry.Desktop)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return fmt.Errorf("editor exited with error: %w", err)
	}

	desktopDir := desktopDirForScope(entry.Scope)
	if err := runUpdateDB(desktopDir); err != nil {
		fmt.Fprintf(os.Stderr, "warning: update-desktop-database failed: %v\n", err)
	}

	return nil
}
