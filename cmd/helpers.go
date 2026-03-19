package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
)

// runUpdateDB runs update-desktop-database on the given directory.
func runUpdateDB(dir string) error {
	return exec.Command("update-desktop-database", dir).Run()
}

// desktopDirForScope returns the applications directory for the given scope.
func desktopDirForScope(scope string) string {
	if scope == "global" {
		return "/usr/share/applications"
	}
	return filepath.Join(os.Getenv("HOME"), ".local", "share", "applications")
}
