/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package classify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildPrompt_BasicStructure(t *testing.T) {
	chunks := []string{
		"Projekt Alpha",
		"Peter anrufen wegen Budget",
	}

	prompt := BuildPrompt(chunks)

	require.Contains(t, prompt, "Chunks")
	require.Contains(t, prompt, "[1] Projekt Alpha")
	require.Contains(t, prompt, "[2] Peter anrufen wegen Budget")
	require.Contains(t, prompt, "confidence (0.0 - 1.0)")
}
