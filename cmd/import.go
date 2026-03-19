package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/andrewwkimm/aidfm/internal/desktop"
	"github.com/andrewwkimm/aidfm/internal/registry"
)

var importCmd = &cobra.Command{
	Use:   "import <name|path>",
	Short: "Take ownership of an existing desktop file",
	Args:  cobra.ExactArgs(1),
	RunE:  runImport,
}

func init() {
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	input := args[0]

	// Resolve path — either a bare name or an absolute path
	desktopPath := resolveDesktopPath(input)

	if _, err := os.Stat(desktopPath); err != nil {
		return fmt.Errorf("desktop file not found: %s", desktopPath)
	}

	df, err := desktop.Read(desktopPath)
	if err != nil {
		return fmt.Errorf("failed to read desktop file: %w", err)
	}

	// Parse Exec= line to extract binary
	execLine := desktop.ParseExec(df.Get("Exec"))
	if execLine.Binary == "" {
		return fmt.Errorf("could not detect binary from Exec= line: %q", df.Get("Exec"))
	}

	icon := df.Get("Icon")
	name := df.Get("Name")
	if name == "" {
		name = strings.TrimSuffix(filepath.Base(desktopPath), ".desktop")
	}

	// Check if already managed
	reg, err := registry.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	if existing := reg.Find(name); existing != nil {
		return fmt.Errorf("%q is already managed by aidfm", name)
	}

	// Show detected values and confirm
	confirmed := false
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Detected values").
				Description(fmt.Sprintf(
					"Name:    %s\nBinary:  %s\nIcon:    %s\nDesktop: %s",
					name, execLine.Binary, icon, desktopPath,
				)),
			huh.NewConfirm().
				Title("Import this entry?").
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

	scope := "user"
	if strings.HasPrefix(desktopPath, "/usr/share") {
		scope = "global"
	}

	reg.Add(registry.Entry{
		Name:    name,
		Binary:  execLine.Binary,
		Desktop: desktopPath,
		Icon:    icon,
		Scope:   scope,
		Status:  registry.StatusHealthy,
	})

	if err := reg.Save(); err != nil {
		return fmt.Errorf("failed to save registry: %w", err)
	}

	fmt.Printf("imported %s\n", name)
	return nil
}

// resolveDesktopPath resolves a bare name like "krokiet" to its full desktop
// file path under ~/.local/share/applications, or returns the input unchanged
// if it looks like an absolute path.
func resolveDesktopPath(input string) string {
	if filepath.IsAbs(input) {
		return input
	}
	name := strings.TrimSuffix(input, ".desktop")
	return filepath.Join(os.Getenv("HOME"), ".local", "share", "applications", name+".desktop")
}
