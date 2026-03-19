package desktop

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// File represents a parsed .desktop file.
type File struct {
	Path   string
	Fields map[string]string
}

// New returns a File with the default fields required for an AppImage desktop entry.
func New(path string) *File {
	return &File{
		Path: path,
		Fields: map[string]string{
			"Type":          "Application",
			"Terminal":      "false",
			"StartupNotify": "true",
		},
	}
}

// Read parses a .desktop file from disk and returns a File.
func Read(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	file := &File{
		Path:   path,
		Fields: make(map[string]string),
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "[") || line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		file.Fields[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return file, nil
}

// Write writes the .desktop file to disk.
func (f *File) Write() error {
	out, err := os.Create(f.Path)
	if err != nil {
		return err
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	fmt.Fprintln(w, "[Desktop Entry]")

	// Write known fields in a consistent order
	order := []string{"Type", "Name", "Exec", "Icon", "Terminal", "StartupNotify"}
	written := make(map[string]bool)

	for _, key := range order {
		if val, ok := f.Fields[key]; ok {
			fmt.Fprintf(w, "%s=%s\n", key, val)
			written[key] = true
		}
	}

	// Write any remaining fields not in the ordered list
	for key, val := range f.Fields {
		if !written[key] {
			fmt.Fprintf(w, "%s=%s\n", key, val)
		}
	}

	return w.Flush()
}

// Get returns the value of a field.
func (f *File) Get(key string) string {
	return f.Fields[key]
}

// Set sets the value of a field.
func (f *File) Set(key, value string) {
	f.Fields[key] = value
}
