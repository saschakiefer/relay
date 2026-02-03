/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package classify

import (
	"fmt"
	"strings"
)

// BuildPrompt constructs the user prompt containing the concrete OCR chunks.
// This prompt is dynamic and varies per capture.
func BuildPrompt(chunks []string) string {
	var b strings.Builder

	b.WriteString(`
These are OCR chunks extracted from handwritten personal notes.

Chunks:
`)

	for i, c := range chunks {
		fmt.Fprintf(&b, "[%d] %s\n", i+1, c)
	}

	b.WriteString(`
Your task:
- Interpret the meaning of the notes.
- Identify todos, notes, ideas, and projects.
- Reconstruct intent where possible.

Return a JSON array.

Each item must have:
- type (todo | note | idea | project)
- text
- confidence (0.0 - 1.0)
`)

	return strings.TrimSpace(b.String())
}
