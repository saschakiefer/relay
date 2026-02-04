/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package vault

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type FilesystemVaultReader struct {
	Root       string
	IgnoreDirs []string
}

func NewFilesystemVaultReader(root string) *FilesystemVaultReader {
	return &FilesystemVaultReader{
		Root: root,
		IgnoreDirs: []string{
			".obsidian",
			".git",
			"node_modules",
		},
	}
}

func (r *FilesystemVaultReader) ListNotes(ctx context.Context) ([]NoteMeta, error) {
	if strings.TrimSpace(r.Root) == "" {
		return nil, errors.New("vault root is empty")
	}

	rootAbs, err := filepath.Abs(r.Root)
	if err != nil {
		return nil, err
	}

	var out []NoteMeta

	err = filepath.WalkDir(rootAbs, func(absPath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// Respect cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if d.IsDir() {
			if r.shouldIgnoreDir(absPath, rootAbs) {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}

		meta, err := r.readMeta(absPath, rootAbs)
		if err != nil {
			return err
		}

		out = append(out, meta)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *FilesystemVaultReader) ReadNote(ctx context.Context, vaultRelPath string) (NoteContent, error) {
	if strings.TrimSpace(r.Root) == "" {
		return NoteContent{}, errors.New("vault root is empty")
	}
	if strings.TrimSpace(vaultRelPath) == "" {
		return NoteContent{}, errors.New("vaultRelPath is empty")
	}

	// Normalize user path to OS path
	relOS := filepath.FromSlash(vaultRelPath)

	absPath := filepath.Join(r.Root, relOS)
	absPath, err := filepath.Abs(absPath)
	if err != nil {
		return NoteContent{}, err
	}

	// Respect cancellation
	select {
	case <-ctx.Done():
		return NoteContent{}, ctx.Err()
	default:
	}

	b, err := os.ReadFile(absPath)
	if err != nil {
		return NoteContent{}, err
	}

	raw := string(b)
	fm, body := ParseFrontmatter(raw)

	st, err := os.Stat(absPath)
	if err != nil {
		return NoteContent{}, err
	}

	rootAbs, err := filepath.Abs(r.Root)
	if err != nil {
		return NoteContent{}, err
	}

	meta := BuildNoteMeta(absPath, rootAbs, fm, body, st.ModTime())

	return NoteContent{
		Meta: meta,
		Body: body,
		Raw:  raw,
	}, nil
}

func (r *FilesystemVaultReader) readMeta(absPath string, rootAbs string) (NoteMeta, error) {
	b, err := os.ReadFile(absPath)
	if err != nil {
		return NoteMeta{}, err
	}

	raw := string(b)
	fm, body := ParseFrontmatter(raw)

	st, err := os.Stat(absPath)
	if err != nil {
		return NoteMeta{}, err
	}

	return BuildNoteMeta(absPath, rootAbs, fm, body, st.ModTime()), nil
}

func (r *FilesystemVaultReader) shouldIgnoreDir(absPath string, rootAbs string) bool {
	rel, err := filepath.Rel(rootAbs, absPath)
	if err != nil {
		return false
	}

	// rel == "." is root; don't ignore root.
	if rel == "." {
		return false
	}

	// Compare first path segment against ignore list.
	parts := splitPath(rel)
	if len(parts) == 0 {
		return false
	}
	first := parts[0]

	for _, ign := range r.IgnoreDirs {
		if ign == first {
			return true
		}
	}
	return false
}

func splitPath(p string) []string {
	p = filepath.Clean(p)
	if p == "." || p == string(filepath.Separator) {
		return nil
	}
	parts := strings.Split(p, string(filepath.Separator))
	// On Windows, sometimes clean can return with drive letter; but here p is Rel() so fine.
	var out []string
	for _, s := range parts {
		if s != "" && s != "." {
			out = append(out, s)
		}
	}
	return out
}

// unix-ish path for storage in index / prompts
func toVaultRelPath(absPath string, rootAbs string) (string, error) {
	rel, err := filepath.Rel(rootAbs, absPath)
	if err != nil {
		return "", err
	}
	rel = filepath.Clean(rel)

	// Ensure forward slashes even on Windows
	if runtime.GOOS == "windows" {
		rel = strings.ReplaceAll(rel, `\`, `/`)
	} else {
		rel = filepath.ToSlash(rel)
	}
	return rel, nil
}

func BuildNoteMeta(absPath string, rootAbs string, fm Frontmatter, body string, modifiedAt time.Time) NoteMeta {
	rel, err := toVaultRelPath(absPath, rootAbs)
	if err != nil {
		// fallback: basename (still stable-ish)
		rel = filepath.ToSlash(filepath.Base(absPath))
	}

	title := ExtractTitle(fm, body, rel)
	tags := ExtractTags(fm)
	links := ExtractLinks(body)

	return NoteMeta{
		Path:       rel,
		Title:      title,
		Tags:       tags,
		Links:      links,
		ModifiedAt: modifiedAt,
	}
}
