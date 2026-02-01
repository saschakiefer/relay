/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package ocr

import (
	"context"
	"fmt"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

type GoogleVisionEngine struct {
	CredentialsPath string
}

func (g *GoogleVisionEngine) Extract(
	ctx context.Context,
	inputPath string,
) (string, error) {
	log.Info().Str("image", inputPath).Msg("Extracting text")

	if g.CredentialsPath == "" {
		return "", fmt.Errorf("google vision credentials not configured")
	}

	if _, err := os.Stat(g.CredentialsPath); err != nil {
		return "", fmt.Errorf("credentials file not found: %w", err)
	}

	client, err := vision.NewImageAnnotatorClient(
		ctx,
		option.WithCredentialsFile(g.CredentialsPath),
	)
	if err != nil {
		return "", fmt.Errorf("vision client: %w", err)
	}
	defer client.Close()

	file, err := os.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("open image: %w", err)
	}
	defer file.Close()

	img, err := vision.NewImageFromReader(file)
	if err != nil {
		return "", fmt.Errorf("image from reader: %w", err)
	}

	doc, err := client.DetectDocumentText(ctx, img, nil)
	if err != nil {
		return "", fmt.Errorf("detect document text: %w", err)
	}

	log.Debug().Str("text", doc.Text).Msg("")
	log.Info().Msg("Text successfully extracted")
	return doc.GetText(), nil
}
