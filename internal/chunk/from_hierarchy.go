/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package chunk

import (
	"strings"

	"github.com/saschakiefer/relay/internal/hierarchy"
)

// FromHierarchy converts a hierarchy tree into LLM-friendly chunks.
// Each root node becomes one chunk containing its full subtree.
func FromHierarchy(roots []*hierarchy.Node) []Chunk {
	var chunks []Chunk
	id := 1

	for _, root := range roots {
		var b strings.Builder
		writeNode(&b, root, 0)

		chunks = append(chunks, Chunk{
			ID:   id,
			Text: strings.TrimSpace(b.String()),
		})
		id++
	}

	return chunks
}
