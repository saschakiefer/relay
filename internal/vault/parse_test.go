/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package vault

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFilesystemVaultReader_ListNotes(t *testing.T) {
	root := filepath.Join("testdata", "vault")
	r := NewFilesystemVaultReader(root)

	ctx := context.Background()
	notes, err := r.ListNotes(ctx)
	if err != nil {
		t.Fatalf("ListNotes error: %v", err)
	}

	if len(notes) != 3 {
		t.Fatalf("expected 3 notes, got %d", len(notes))
	}

	// Build map by path for stable assertions
	m := map[string]NoteMeta{}
	for _, n := range notes {
		m[n.Path] = n
	}

	assertHas := func(path string) NoteMeta {
		n, ok := m[path]
		if !ok {
			t.Fatalf("expected note path %q in results", path)
		}
		return n
	}

	n1 := assertHas("Projects/relay.md")
	if n1.Title != "relay" {
		t.Fatalf("expected title 'relay', got %q", n1.Title)
	}
	if len(n1.Tags) != 2 || n1.Tags[0] != "project" || n1.Tags[1] != "relay" {
		t.Fatalf("unexpected tags for relay.md: %#v", n1.Tags)
	}
	if len(n1.Links) != 2 {
		t.Fatalf("expected 2 links for relay.md, got %d: %#v", len(n1.Links), n1.Links)
	}

	n2 := assertHas("Notes/Agentic AI.md")
	if n2.Title != "Agentic AI" {
		t.Fatalf("expected title 'Agentic AI', got %q", n2.Title)
	}
	if len(n2.Tags) != 2 || n2.Tags[0] != "ai" || n2.Tags[1] != "notes" {
		t.Fatalf("unexpected tags for Agentic AI.md: %#v", n2.Tags)
	}
	if len(n2.Links) != 1 || n2.Links[0] != "relay" {
		t.Fatalf("unexpected links for Agentic AI.md: %#v", n2.Links)
	}

	n3 := assertHas("Daily/2026-01-31.md")
	if n3.Title != "2026-01-31" {
		t.Fatalf("expected title '2026-01-31', got %q", n3.Title)
	}

	// ModifiedAt should be set (non-zero)
	if n3.ModifiedAt.IsZero() {
		t.Fatalf("expected ModifiedAt to be set")
	}
}

func TestFilesystemVaultReader_ReadNote(t *testing.T) {
	root := filepath.Join("testdata", "vault")
	r := NewFilesystemVaultReader(root)

	ctx := context.Background()
	n, err := r.ReadNote(ctx, "Projects/relay.md")
	if err != nil {
		t.Fatalf("ReadNote error: %v", err)
	}

	if n.Meta.Path != "Projects/relay.md" {
		t.Fatalf("expected path 'Projects/relay.md', got %q", n.Meta.Path)
	}
	if n.Body == "" {
		t.Fatalf("expected non-empty body")
	}
	if n.Raw == "" {
		t.Fatalf("expected non-empty raw")
	}
}

func TestParseFrontmatter_NoFrontmatter(t *testing.T) {
	fm, body := ParseFrontmatter("# Title\nhello")
	if len(fm) != 0 {
		t.Fatalf("expected empty frontmatter")
	}
	if body != "# Title\nhello" {
		t.Fatalf("unexpected body: %q", body)
	}
}

// sanity for windows path normalization in reader (we store with forward slashes)
func TestPathNormalization(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-only")
	}
	root := filepath.Join("testdata", "vault")
	r := NewFilesystemVaultReader(root)

	ctx := context.Background()
	notes, err := r.ListNotes(ctx)
	if err != nil {
		t.Fatalf("ListNotes error: %v", err)
	}
	for _, n := range notes {
		if filepath.Separator == '\\' && containsBackslash(n.Path) {
			t.Fatalf("expected vault-relative path to use forward slashes, got %q", n.Path)
		}
	}
}

func containsBackslash(s string) bool {
	for _, ch := range s {
		if ch == '\\' {
			return true
		}
	}
	return false
}
