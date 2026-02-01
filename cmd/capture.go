/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
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

		engine, err := resolveOCREngine(captureEngine)
		if err != nil {
			return err
		}

		rawText, err := engine.Extract(ctx, captureInput)
		if err != nil {
			return err
		}

		lines := normalize.Lines(rawText)

		for _, line := range lines {
			log.Debug().Msg(line)
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
