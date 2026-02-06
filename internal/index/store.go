package index

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Save(path string, idx *VaultIndex) error {
	b, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}

	// Create parent directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	// Write file (creates if doesn't exist, overwrites if it does)
	return os.WriteFile(path, b, 0o644)
}

func Load(path string) (*VaultIndex, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var idx VaultIndex
	if err := json.Unmarshal(b, &idx); err != nil {
		return nil, err
	}
	return &idx, nil
}
