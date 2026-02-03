/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/saschakiefer/relay/internal/chunk"
	"github.com/saschakiefer/relay/internal/normalize"
	"github.com/spf13/cobra"
)

var (
	captureEngine string
	captureInput  string
)

var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture a handwritten note using Google Vision OCR and process it",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// 1. OCR Engine
		engine, err := resolveOCREngine(captureEngine)
		if err != nil {
			return err
		}

		// 2. Classifier (LLM)
		classifier, err := resolveClassifier()
		if err != nil {
			return err
		}

		// 3. OCR
		rawText, err := engine.Extract(ctx, captureInput)
		if err != nil {
			return err
		}

		// 4. Normalize
		lines := normalize.Lines(rawText)

		for _, line := range lines {
			log.Debug().Msg(line)
		}

		// 5. Chunk
		chunks := chunk.FromLines(lines)

		for _, c := range chunks {
			log.Debug().
				Int("chunk_id", c.ID).
				Msg(c.Text)
		}

		// 6. Extract chunk texts
		texts := make([]string, len(chunks))
		for i, c := range chunks {
			texts[i] = c.Text
		}

		// 7. Classify (LLM)
		items, err := classifier.Classify(ctx, texts)
		if err != nil {
			return err
		}

		// 8. Output (for now)
		for _, item := range items {
			log.Debug().
				Str("type", string(item.Type)).
				Float64("confidence", item.Confidence).
				Msg(item.Text)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(captureCmd)

	captureCmd.Flags().StringVarP(
		&captureEngine,
		"engine",
		"e",
		"google",
		"OCR engine to use (google, <more to come>)",
	)

	captureCmd.Flags().StringVarP(
		&captureInput,
		"input",
		"i",
		"",
		"Path to input image",
	)

	_ = captureCmd.MarkFlagRequired("input")
}
