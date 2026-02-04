/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package vault

import (
	"context"
	"time"
)

// NoteMeta is lightweight metadata for indexing and selection.
type NoteMeta struct {
	// Path is vault-relative path using forward slashes, e.g. "Projects/relay.md".
	Path string

	Title      string
	Tags       []string
	Links      []string
	ModifiedAt time.Time
}

// NoteContent is returned when you need full body text.
type NoteContent struct {
	Meta NoteMeta

	// Body is markdown body without frontmatter.
	Body string

	// Raw is the full file content (frontmatter + body). Useful for debugging.
	Raw string
}

// VaultReader provides read-only access to an Obsidian vault.
type VaultReader interface {
	ListNotes(ctx context.Context) ([]NoteMeta, error)
	ReadNote(ctx context.Context, vaultRelPath string) (NoteContent, error)
}
