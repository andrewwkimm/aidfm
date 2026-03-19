package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/aidfm/internal/desktop"
	"github.com/yourusername/aidfm/internal/registry"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environment variables on the Exec= line of a desktop entry",
}

var envSetCmd = &cobra.Command{
	Use:   "set <name> <key> <value>",
	Short: "Set an environment variable on the Exec= line",
	Args:  cobra.ExactArgs(3),
	RunE:  runEnvSet,
}

var envUnsetCmd = &cobra.Command{
	Use:   "unset <name> <key>",
	Short: "Remove an environment variable from the Exec= line",
	Args:  cobra.ExactArgs(2),
	RunE:  runEnvUnset,
}

var envListCmd = &cobra.Command{
	Use:   "list <name>",
	Short: "List all environment variables on the Exec= line",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvList,
}

func init() {
	envCmd.AddCommand(envSetCmd)
	envCmd.AddCommand(envUnsetCmd)
	envCmd.AddCommand(envListCmd)
	rootCmd.AddCommand(envCmd)
}

func runEnvSet(cmd *cobra.Command, args []string) error {
	name, key, value := args[0], args[1], args[2]

	entry, df, err := loadEntryAndDesktop(name)
	if err != nil {
		return err
	}

	exec := desktop.ParseExec(df.Get("Exec"))

	if !fileExists(exec.Binary) {
		return fmt.Errorf("binary %q does not exist on disk, refusing to rewrite Exec= line", exec.Binary)
	}

	exec.SetEnv(key, value)
	df.Set("Exec", exec.String())

	if err := df.Write(); err != nil {
		return fmt.Errorf("failed to write desktop file: %w", err)
	}

	if err := runUpdateDB(desktopDirForScope(entry.Scope)); err != nil {
		fmt.Fprintf(os.Stderr, "warning: update-desktop-database failed: %v\n", err)
	}

	fmt.Printf("set %s=%s on %s\n", key, value, name)
	return nil
}

func runEnvUnset(cmd *cobra.Command, args []string) error {
	name, key := args[0], args[1]

	entry, df, err := loadEntryAndDesktop(name)
	if err != nil {
		return err
	}

	exec := desktop.ParseExec(df.Get("Exec"))
	exec.UnsetEnv(key)
	df.Set("Exec", exec.String())

	if err := df.Write(); err != nil {
		return fmt.Errorf("failed to write desktop file: %w", err)
	}

	if err := runUpdateDB(desktopDirForScope(entry.Scope)); err != nil {
		fmt.Fprintf(os.Stderr, "warning: update-desktop-database failed: %v\n", err)
	}

	fmt.Printf("unset %s on %s\n", key, name)
	return nil
}

func runEnvList(cmd *cobra.Command, args []string) error {
	_, df, err := loadEntryAndDesktop(args[0])
	if err != nil {
		return err
	}

	exec := desktop.ParseExec(df.Get("Exec"))

	if len(exec.Env) == 0 {
		fmt.Println("no environment variables set")
		return nil
	}

	for k, v := range exec.Env {
		fmt.Printf("%s=%s\n", k, v)
	}

	return nil
}

// loadEntryAndDesktop is a helper that loads the registry entry and desktop file
// for a given name, used by all env subcommands.
func loadEntryAndDesktop(name string) (*registry.Entry, *desktop.File, error) {
	reg, err := registry.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load registry: %w", err)
	}

	entry := reg.Find(name)
	if entry == nil {
		return nil, nil, fmt.Errorf("no entry found for %q", name)
	}

	df, err := desktop.Read(entry.Desktop)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read desktop file: %w", err)
	}

	return entry, df, nil
}
