/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package chunk

import "strings"

// Chunk represents a semantically related text fragment.
// IMPORTANT: Chunks are hypotheses, not guarantees.
type Chunk struct {
	ID   int
	Text string
}

// FromLines groups normalized OCR lines into rough semantic chunks.
//
// Design goals:
// - tolerate OCR noise
// - NOT rely on perfect line breaks
// - produce chunks suitable for LLM reinterpretation
func FromLines(lines []string) []Chunk {
	var chunks []Chunk
	var current []string
	id := 1

	flush := func() {
		if len(current) == 0 {
			return
		}

		chunks = append(chunks, Chunk{
			ID:   id,
			Text: strings.Join(current, " "),
		})
		id++
		current = nil
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Empty lines strongly suggest a new thought,
		// but OCR may be noisy — we treat them as soft boundaries.
		if line == "" {
			flush()
			continue
		}

		current = append(current, line)
	}

	flush()
	return chunks
}
