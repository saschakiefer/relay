/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package vault

import (
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Frontmatter map[string]any

var (
	reHeading1 = regexp.MustCompile(`(?m)^\s*#\s+(.+?)\s*$`)
	// [[Note]] or [[Note|Alias]] or [[path/to/note|Alias]]
	reWikilink = regexp.MustCompile(`\[\[([^\]\|#]+)(?:\|[^\]]+)?\]\]`)
	// YAML frontmatter: starts with --- at beginning
)

// ParseFrontmatter splits raw markdown into YAML frontmatter (if present) and body.
func ParseFrontmatter(raw string) (Frontmatter, string) {
	s := strings.TrimLeft(raw, "\ufeff") // strip BOM if present
	if !strings.HasPrefix(s, "---") {
		return Frontmatter{}, strings.TrimSpace(raw)
	}

	// Find end delimiter "\n---" on its own line
	// We keep it simple and robust.
	lines := strings.Split(s, "\n")
	if len(lines) < 3 {
		return Frontmatter{}, strings.TrimSpace(raw)
	}
	if strings.TrimSpace(lines[0]) != "---" {
		return Frontmatter{}, strings.TrimSpace(raw)
	}

	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			end = i
			break
		}
	}
	if end == -1 {
		return Frontmatter{}, strings.TrimSpace(raw)
	}

	yamlPart := strings.Join(lines[1:end], "\n")
	body := strings.Join(lines[end+1:], "\n")

	fm := Frontmatter{}
	if strings.TrimSpace(yamlPart) != "" {
		_ = yaml.Unmarshal([]byte(yamlPart), &fm) // tolerate YAML errors; treat as empty
	}

	return fm, strings.TrimSpace(body)
}

// ExtractTitle: frontmatter title -> first H1 -> filename without extension.
func ExtractTitle(fm Frontmatter, body string, vaultRelPath string) string {
	// 1) frontmatter: title
	if v, ok := fm["title"]; ok {
		if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}

	// 2) first H1 heading
	if m := reHeading1.FindStringSubmatch(body); len(m) == 2 {
		h := strings.TrimSpace(m[1])
		if h != "" {
			return h
		}
	}

	// 3) filename
	base := vaultRelPath
	if idx := strings.LastIndex(base, "/"); idx >= 0 {
		base = base[idx+1:]
	}
	if strings.HasSuffix(strings.ToLower(base), ".md") {
		base = base[:len(base)-3]
	}
	return base
}

func ExtractTags(fm Frontmatter) []string {
	var tags []string

	// frontmatter tags can be string, []any, []string
	if v, ok := fm["tags"]; ok {
		switch t := v.(type) {
		case string:
			for _, part := range splitTagsString(t) {
				tags = append(tags, normalizeTag(part))
			}
		case []any:
			for _, it := range t {
				if s, ok := it.(string); ok {
					tags = append(tags, normalizeTag(s))
				}
			}
		case []string:
			for _, s := range t {
				tags = append(tags, normalizeTag(s))
			}
		}
	}

	tags = uniqueNonEmpty(tags)
	sort.Strings(tags)
	return tags
}

func splitTagsString(s string) []string {
	// allow "a, b" or "a b"
	s = strings.ReplaceAll(s, ",", " ")
	fields := strings.Fields(s)
	return fields
}

func normalizeTag(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	// Obsidian tags usually without '#', but keep consistent: store without '#'
	s = strings.TrimPrefix(s, "#")
	return s
}

func ExtractLinks(body string) []string {
	matches := reWikilink.FindAllStringSubmatch(body, -1)
	if len(matches) == 0 {
		return nil
	}

	var links []string
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		target := strings.TrimSpace(m[1])
		if target == "" {
			continue
		}
		links = append(links, target)
	}

	links = uniqueNonEmpty(links)
	sort.Strings(links)
	return links
}

func uniqueNonEmpty(in []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
