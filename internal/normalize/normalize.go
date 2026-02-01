/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package normalize

import "strings"

// Lines normalizes raw OCR output into clean, non-empty lines.
func Lines(input string) []string {
	rawLines := strings.Split(input, "\n")

	var result []string
	for _, line := range rawLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		result = append(result, line)
	}

	return result
}
