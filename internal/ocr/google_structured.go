/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package ocr

import (
	"context"
	"fmt"
	"os"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "cloud.google.com/go/vision/v2/apiv1/visionpb"
	"google.golang.org/api/option"
)

type Line struct {
	Text   string
	Indent int
	Page   int
}

// ExtractStructured extracts OCR text together with layout information
// (indentation derived from bounding boxes).
func (g *GoogleVisionEngine) ExtractStructured(
	ctx context.Context,
	inputPath string,
) ([]Line, error) {

	if g.CredentialsPath == "" {
		return nil, fmt.Errorf("google vision credentials not configured")
	}

	client, err := vision.NewImageAnnotatorClient(
		ctx,
		option.WithCredentialsFile(g.CredentialsPath),
	)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	f, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := vision.NewImageFromReader(f)
	if err != nil {
		return nil, err
	}

	doc, err := client.DetectDocumentText(ctx, img, nil)
	if err != nil {
		return nil, err
	}

	var result []Line

	for pIdx, page := range doc.Pages {
		for _, block := range page.Blocks {
			for _, para := range block.Paragraphs {

				text := extractParagraphText(para)
				if text == "" {
					continue
				}

				indent := minX(para.BoundingBox.Vertices)

				result = append(result, Line{
					Text:   text,
					Indent: indent,
					Page:   pIdx,
				})
			}
		}
	}

	return result, nil
}

func extractParagraphText(p *visionpb.Paragraph) string {
	var b strings.Builder

	for _, w := range p.Words {
		for _, s := range w.Symbols {
			b.WriteString(s.Text)
		}
		b.WriteString(" ")
	}

	return strings.TrimSpace(b.String())
}

func minX(vertices []*visionpb.Vertex) int {
	min := int(^uint(0) >> 1) // max int

	for _, v := range vertices {
		if int(v.X) < min {
			min = int(v.X)
		}
	}

	return min
}
