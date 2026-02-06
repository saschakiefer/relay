/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package chunk

import (
	"strings"

	"github.com/saschakiefer/relay/internal/hierarchy"
)

func writeNode(b *strings.Builder, n *hierarchy.Node, level int) {
	indent := strings.Repeat("  ", level)
	b.WriteString(indent)
	b.WriteString("- ")
	b.WriteString(n.Text)
	b.WriteString("\n")

	for _, child := range n.Children {
		writeNode(b, child, level+1)
	}
}
