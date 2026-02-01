/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var captureCmd = &cobra.Command{
	Use:   "capture [file]",
	Short: "Capture a handwritten note and process it",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug().
			Str("agent", "capture").
			Str("file", args[0]).
			Msg("capturing handwritten note")

		// TODO: OCR → Agent → MCP
	},
}

func init() {
	rootCmd.AddCommand(captureCmd)
}
