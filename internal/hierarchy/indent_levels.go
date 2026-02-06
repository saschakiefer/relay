/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package hierarchy

import "sort"

// NormalizeIndents maps raw indent values to logical levels (0,1,2,...).
// threshold is the minimum delta between indent clusters (e.g. 30-50 px).
func NormalizeIndents(indents []int, threshold int) map[int]int {
	sort.Ints(indents)

	levelMap := make(map[int]int)
	level := 0

	var lastCluster int
	hasCluster := false

	for _, v := range indents {
		if !hasCluster {
			levelMap[v] = level
			lastCluster = v
			hasCluster = true
			continue
		}
		if v-lastCluster > threshold {
			level++
			lastCluster = v
		}
		levelMap[v] = level
	}

	return levelMap
}
