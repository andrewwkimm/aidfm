package cmd

import "os/exec"

// runUpdateDB runs update-desktop-database on the given directory.
func runUpdateDB(dir string) error {
	return exec.Command("update-desktop-database", dir).Run()
}
