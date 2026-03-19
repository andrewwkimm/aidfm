package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/yourusername/aidfm/internal/registry"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all managed AppImage entries",
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

var (
	styleHealthy  = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // green
	styleBroken   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // red
	styleOrphaned = lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // yellow
	styleHeader   = lipgloss.NewStyle().Bold(true).Underline(true)
)

func runList(cmd *cobra.Command, args []string) error {
	reg, err := registry.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	if len(reg.Entries) == 0 {
		fmt.Println("no managed entries found. run aidfm add <path> to get started.")
		return nil
	}

	// Run sync before displaying
	if err := syncRegistry(reg); err != nil {
		fmt.Fprintf(os.Stderr, "warning: sync failed: %v\n", err)
	}

	// Header
	fmt.Printf("%s\n",
		styleHeader.Render(fmt.Sprintf("%-20s %-8s %-10s %-10s %-10s",
			"NAME", "SCOPE", "BINARY", "DESKTOP", "ICON")),
	)

	for _, e := range reg.Entries {
		style := statusStyle(e.Status)
		binaryStatus := checkMark(e.Binary)
		desktopStatus := checkMark(e.Desktop)
		iconStatus := checkMark(e.Icon)

		fmt.Println(style.Render(fmt.Sprintf("%-20s %-8s %-10s %-10s %-10s",
			e.Name, e.Scope, binaryStatus, desktopStatus, iconStatus)))
	}

	return nil
}

func statusStyle(s registry.Status) lipgloss.Style {
	switch s {
	case registry.StatusHealthy:
		return styleHealthy
	case registry.StatusBroken:
		return styleBroken
	case registry.StatusOrphaned:
		return styleOrphaned
	default:
		return lipgloss.NewStyle()
	}
}

func checkMark(path string) string {
	if path == "" {
		return "✗ missing"
	}
	if _, err := os.Stat(path); err != nil {
		return "✗ missing"
	}
	return "✓"
}
