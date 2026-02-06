/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/saschakiefer/relay/internal/chunk"
	"github.com/saschakiefer/relay/internal/hierarchy"
	"github.com/saschakiefer/relay/internal/normalize"
	"github.com/saschakiefer/relay/internal/ocr"
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

		var texts []string

		// If the OCR engine supports structured extraction, prefer it.
		type structured interface {
			ExtractStructured(context.Context, string) ([]ocr.Line, error)
		}

		if s, ok := engine.(structured); ok {
			log.Debug().Msg("using structured OCR")

			// 3a. Structured OCR
			lines, err := s.ExtractStructured(ctx, captureInput)
			if err != nil {
				return err
			}

			for _, l := range lines {
				log.Debug().
					Int("indent", l.Indent).
					Msg(l.Text)
			}

			// 4a. Build hierarchy
			indents := hierarchy.ExtractIndents(lines)
			levelMap := hierarchy.NormalizeIndents(indents, 40)
			tree := hierarchy.BuildHierarchy(lines, levelMap)

			// 5a. Chunk from hierarchy
			chunks := chunk.FromHierarchy(tree)

			for _, c := range chunks {
				log.Debug().
					Int("chunk_id", c.ID).
					Msg(c.Text)
			}

			// 6a. Extract texts
			texts = make([]string, len(chunks))
			for i, c := range chunks {
				texts[i] = c.Text
			}

		} else {
			log.Debug().Msg("using flat OCR")

			// 3b. Flat OCR
			rawText, err := engine.Extract(ctx, captureInput)
			if err != nil {
				return err
			}

			// 4b. Normalize
			lines := normalize.Lines(rawText)

			for _, line := range lines {
				log.Debug().Msg(line)
			}

			// 5b. Chunk from lines
			chunks := chunk.FromLines(lines)

			for _, c := range chunks {
				log.Debug().
					Int("chunk_id", c.ID).
					Msg(c.Text)
			}

			// 6b. Extract texts
			texts = make([]string, len(chunks))
			for i, c := range chunks {
				texts[i] = c.Text
			}
		}

		// 7. Classify (LLM)
		items, err := classifier.Classify(ctx, texts)
		if err != nil {
			return err
		}

		// 8. Output (for now)
		for _, item := range items {
			log.Debug().
				Int("level", item.Level).
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
