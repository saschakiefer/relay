/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package hierarchy

import (
	"github.com/saschakiefer/relay/internal/ocr"
)

type Node struct {
	Text     string
	Level    int
	Children []*Node
}

// BuildHierarchy builds a tree from ordered lines and their levels.
func BuildHierarchy(lines []ocr.Line, levelMap map[int]int) []*Node {
	var roots []*Node
	stack := []*Node{}

	for _, line := range lines {
		level := levelMap[line.Indent]
		node := &Node{
			Text:  line.Text,
			Level: level,
		}

		// Pop stack until we find a parent
		for len(stack) > 0 && stack[len(stack)-1].Level >= level {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			// root node
			roots = append(roots, node)
		} else {
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, node)
		}

		stack = append(stack, node)
	}

	return roots
}
