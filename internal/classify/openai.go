/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package classify

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/rs/zerolog/log"
)

type OpenAIClassifier struct {
	Client openai.Client
	Model  string
}

func (c *OpenAIClassifier) Classify(
	ctx context.Context,
	chunks []string,
) ([]Item, error) {

	if len(chunks) == 0 {
		return nil, nil
	}

	prompt := BuildPrompt(chunks)

	resp, err := c.Client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model: c.Model,
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemPrompt),
				openai.UserMessage(prompt),
			},
			Temperature: param.NewOpt(0.2),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("openai completion: %w", err)
	}

	raw := resp.Choices[0].Message.Content
	log.Debug().Str("raw", raw).Msg("LLM output")

	var items []Item
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil, fmt.Errorf(
			"invalid LLM JSON output: %w\nraw:\n%s",
			err,
			raw,
		)
	}

	return items, nil
}
