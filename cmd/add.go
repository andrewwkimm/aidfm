package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/andrewwkimm/aidfm/internal/desktop"
	"github.com/andrewwkimm/aidfm/internal/detect"
	"github.com/andrewwkimm/aidfm/internal/registry"
)

var globalFlag bool

var addCmd = &cobra.Command{
	Use:   "add <path>",
	Short: "Register an AppImage and create a .desktop file",
	Args:  cobra.ExactArgs(1),
	RunE:  runAdd,
}

func init() {
	addCmd.Flags().BoolVar(&globalFlag, "global", false, "install desktop file to /usr/share/applications instead of ~/.local/share/applications")
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	info, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("cannot access path: %w", err)
	}

	var result detect.Result
	if info.IsDir() {
		result, err = detect.FromDirectory(inputPath)
	} else {
		result, err = detect.FromFile(inputPath)
	}
	if err != nil {
		return fmt.Errorf("detection failed: %w", err)
	}

	if result.Binary == "" {
		return fmt.Errorf("no AppImage or executable found in %s", inputPath)
	}

	parentDir := inputPath
	if !info.IsDir() {
		parentDir = filepath.Dir(inputPath)
	}

	name := strings.ToLower(filepath.Base(parentDir))
	confirmed := false

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Value(&name),
			huh.NewNote().
				Title("Detected").
				Description(fmt.Sprintf(
					"Binary: %s\nIcon:   %s\nDir:    %s",
					result.Binary,
					result.Icon,
					parentDir,
				)),
			huh.NewConfirm().
				Title("Add this entry?").
				Value(&confirmed),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("form error: %w", err)
	}

	if !confirmed {
		fmt.Fprintln(os.Stderr, "aborted")
		return nil
	}

	if err := os.Chmod(result.Binary, 0755); err != nil {
		return fmt.Errorf("failed to set executable bit: %w", err)
	}

	desktopDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "applications")
	if globalFlag {
		desktopDir = "/usr/share/applications"
	}
	desktopPath := filepath.Join(desktopDir, strings.ToLower(name)+".desktop")

	execLine := desktop.ExecLine{
		Binary: result.Binary,
		Env:    make(map[string]string),
	}

	df := desktop.New(desktopPath)
	df.Set("Name", name)
	df.Set("Exec", execLine.String())
	if result.Icon != "" {
		df.Set("Icon", result.Icon)
	} else {
		fmt.Fprintln(os.Stderr, "warning: no icon found, desktop entry will have no icon")
	}

	if err := df.Write(); err != nil {
		return fmt.Errorf("failed to write desktop file: %w", err)
	}

	reg, err := registry.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	reg.Add(registry.Entry{
		Name:      name,
		Binary:    result.Binary,
		Desktop:   desktopPath,
		Icon:      result.Icon,
		ParentDir: parentDir,
		Scope:     scope(),
		Status:    registry.StatusHealthy,
	})

	if err := reg.Save(); err != nil {
		return fmt.Errorf("failed to save registry: %w", err)
	}

	if err := runUpdateDB(desktopDir); err != nil {
		fmt.Fprintf(os.Stderr, "warning: update-desktop-database failed: %v\n", err)
	}

	fmt.Printf("added %s\n", name)
	return nil
}

func scope() string {
	if globalFlag {
		return "global"
	}
	return "user"
}
