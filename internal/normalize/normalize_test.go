/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package normalize

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLines(t *testing.T) {
	input := `
	First line

	  Second line  

	  
	Third line
	`

	lines := Lines(input)

	require.Equal(t, []string{
		"First line",
		"Second line",
		"Third line",
	}, lines)
}
