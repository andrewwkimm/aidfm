package detect

import (
	"os"
	"path/filepath"
	"strings"
)

// iconPrecedence defines the order in which icon formats are preferred.
var iconPrecedence = []string{".svg", ".png", ".jpg", ".jpeg", ".webp", ".ico", ".gif"}

// Result holds the detected AppImage binary and icon paths for a given directory.
type Result struct {
	Binary string
	Icon   string
}

// FromDirectory scans the given directory and attempts to detect an AppImage binary
// and an icon file. Binary detection checks for .AppImage files first, then falls
// back to any executable file excluding hidden files and known non-binary extensions.
// Icon detection follows the order defined in iconPrecedence.
// Returns a Result with whichever paths were found — Binary or Icon may be empty if not detected.
func FromDirectory(dir string) (Result, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return Result{}, err
	}

	var result Result

	for _, e := range entries {
		if e.IsDir() || isHidden(e.Name()) {
			continue
		}
		if strings.HasSuffix(e.Name(), ".AppImage") {
			result.Binary = filepath.Join(dir, e.Name())
			break
		}
	}

	if result.Binary == "" {
		result.Binary = findExecutable(dir, entries)
	}

	result.Icon = findIcon(dir, entries)

	return result, nil
}

// FromFile takes a direct path to an AppImage file and detects an icon in the same directory.
func FromFile(path string) (Result, error) {
	dir := filepath.Dir(path)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return Result{}, err
	}

	result := Result{
		Binary: path,
		Icon:   findIcon(dir, entries),
	}

	return result, nil
}

func findExecutable(dir string, entries []os.DirEntry) string {
	skipExts := []string{".so", ".py"}

	for _, e := range entries {
		if e.IsDir() || isHidden(e.Name()) {
			continue
		}
		ext := filepath.Ext(e.Name())
		if isSkippedExt(ext, skipExts) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.Mode()&0111 != 0 {
			return filepath.Join(dir, e.Name())
		}
	}
	return ""
}

func findIcon(dir string, entries []os.DirEntry) string {
	byExt := make(map[string]string)
	for _, e := range entries {
		if e.IsDir() || isHidden(e.Name()) {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if _, ok := byExt[ext]; !ok {
			byExt[ext] = filepath.Join(dir, e.Name())
		}
	}

	for _, ext := range iconPrecedence {
		if path, ok := byExt[ext]; ok {
			return path
		}
	}
	return ""
}

func isHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

func isSkippedExt(ext string, skip []string) bool {
	for _, s := range skip {
		if ext == s {
			continue
		}
		if strings.HasPrefix(ext, s+".") {
			return true
		}
	}
	return false
}
