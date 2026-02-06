package config

import (
	"fmt"
	"os"
)

func RequireObsidianVault() (string, error) {
	if AppConfig == nil {
		return "", fmt.Errorf("config not initialized")
	}

	vault := AppConfig.Obsidian.Vault
	if vault == "" {
		return "", fmt.Errorf(
			"obsidian vault not configured. Set obsidian.vault in config",
		)
	}

	st, err := os.Stat(vault)
	if err != nil {
		return "", fmt.Errorf("obsidian vault does not exist: %s", vault)
	}
	if !st.IsDir() {
		return "", fmt.Errorf("obsidian vault is not a directory: %s", vault)
	}

	return vault, nil
}

func RequireObsidianIndex() (string, error) {
	index := AppConfig.Obsidian.Index
	if index == "" {
		return "", fmt.Errorf(
			"obsidian index path not configured. Set obsidian.index in config",
		)
	}
	return index, nil
}
