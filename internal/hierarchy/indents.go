/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package hierarchy

import "github.com/saschakiefer/relay/internal/ocr"

func ExtractIndents(lines []ocr.Line) []int {
	seen := make(map[int]struct{})
	var indents []int
	for _, l := range lines {
		if _, ok := seen[l.Indent]; ok {
			continue
		}
		seen[l.Indent] = struct{}{}
		indents = append(indents, l.Indent)
	}
	return indents
}
