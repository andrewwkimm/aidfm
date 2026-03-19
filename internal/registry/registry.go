package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Status represents the health state of a managed AppImage entry.
type Status string

const (
	// StatusHealthy means the binary, desktop file, and icon all exist and are valid.
	StatusHealthy Status = "healthy"
	// StatusBroken means the binary exists but something is missing or invalid (e.g. no desktop file, not executable).
	StatusBroken Status = "broken"
	// StatusOrphaned means the desktop file exists but the binary it points to is missing.
	StatusOrphaned Status = "orphaned"
)

// Entry represents a single managed AppImage and its associated files.
type Entry struct {
	Name      string `json:"name"`
	Binary    string `json:"binary"`
	Desktop   string `json:"desktop"`
	Icon      string `json:"icon"`
	ParentDir string `json:"parent_dir"`
	Scope     string `json:"scope"` // "user" or "global"
	Status    Status `json:"status"`
}

// Registry holds all entries managed by aidfm.
type Registry struct {
	Entries []Entry `json:"entries"`
}

func registryPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", "aidfm", "registry.json"), nil
}

// Load reads the registry from disk. If no registry file exists yet, it returns an empty registry.
func Load() (*Registry, error) {
	path, err := registryPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Registry{}, nil
	}
	if err != nil {
		return nil, err
	}

	var r Registry
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Save writes the registry to disk, creating the directory if it does not exist.
func (r *Registry) Save() error {
	path, err := registryPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Add appends a new entry to the registry. Call Save to persist the change.
func (r *Registry) Add(e Entry) {
	r.Entries = append(r.Entries, e)
}

// Find returns a pointer to the entry with the given name, or nil if not found.
func (r *Registry) Find(name string) *Entry {
	for i := range r.Entries {
		if r.Entries[i].Name == name {
			return &r.Entries[i]
		}
	}
	return nil
}

// Remove deletes the entry with the given name from the registry and returns true if it was found.
// Call Save to persist the change.
func (r *Registry) Remove(name string) bool {
	for i, e := range r.Entries {
		if e.Name == name {
			r.Entries = append(r.Entries[:i], r.Entries[i+1:]...)
			return true
		}
	}
	return false
}
