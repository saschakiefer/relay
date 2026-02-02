/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package chunk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromLines_SimpleParagraphs(t *testing.T) {
	lines := []string{
		"Projekt Alpha",
		"Status unklar",
		"",
		"Peter anrufen",
		"wegen Budget",
	}

	chunks := FromLines(lines)

	require.Len(t, chunks, 2)

	require.Equal(t, "Projekt Alpha Status unklar", chunks[0].Text)
	require.Equal(t, "Peter anrufen wegen Budget", chunks[1].Text)
}

func TestFromLines_NoEmptyLines(t *testing.T) {
	lines := []string{
		"Heute Meeting mit Peter",
		"Budget kritisch",
		"Entscheidung nächste Woche",
	}

	chunks := FromLines(lines)

	require.Len(t, chunks, 1)
	require.Equal(
		t,
		"Heute Meeting mit Peter Budget kritisch Entscheidung nächste Woche",
		chunks[0].Text,
	)
}

func TestFromLines_OnlyEmptyLines(t *testing.T) {
	lines := []string{
		"",
		"",
		"",
	}

	chunks := FromLines(lines)
	require.Empty(t, chunks)
}

func TestFromLines_MultipleBlocks(t *testing.T) {
	lines := []string{
		"Idee neues Pricing",
		"",
		"TODO",
		"Angebot überarbeiten",
		"",
		"Random Gedanke",
	}

	chunks := FromLines(lines)

	require.Len(t, chunks, 3)

	require.Equal(t, "Idee neues Pricing", chunks[0].Text)
	require.Equal(t, "TODO Angebot überarbeiten", chunks[1].Text)
	require.Equal(t, "Random Gedanke", chunks[2].Text)
}
