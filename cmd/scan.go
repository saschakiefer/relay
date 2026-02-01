/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	scanEngine string
	scanInput  string
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a handwritten note using Google Vision OCR",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		engine, err := resolveOCREngine(scanEngine)
		if err != nil {
			return err
		}

		text, err := engine.Extract(ctx, scanInput)
		if err != nil {
			return err
		}

		// Output without log to be directly usable
		fmt.Printf("%s\n", text)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringVarP(
		&scanEngine,
		"engine",
		"e",
		"google",
		"OCR engine to use (google, <more to come>)",
	)

	scanCmd.Flags().StringVarP(
		&scanInput,
		"input",
		"i",
		"",
		"Path to input image",
	)

	_ = scanCmd.MarkFlagRequired("input")
}
