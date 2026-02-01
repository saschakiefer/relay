/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package cmd

import (
	"fmt"

	"github.com/saschakiefer/relay/internal/ocr"
	"github.com/spf13/viper"
)

func resolveOCREngine(engine string) (ocr.Engine, error) {
	switch engine {

	case "google":
		creds := viper.GetString("ocr.google.credentials")
		if creds == "" {
			return nil, fmt.Errorf(
				"google ocr selected but ocr.google.credentials not configured",
			)
		}

		return &ocr.GoogleVisionEngine{
			CredentialsPath: creds,
		}, nil

	default:
		return nil, fmt.Errorf("unknown ocr engine: %s", engine)
	}
}
