/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package cmd

import (
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/saschakiefer/relay/internal/classify"
	"github.com/spf13/viper"
)

// resolveClassifier initializes and returns a Classifier implementation
// based on configuration.
func resolveClassifier() (classify.Classifier, error) {
	provider := viper.GetString("classify.provider")
	if provider == "" {
		provider = "openai"
	}

	switch provider {

	case "openai":
		apiKey := viper.GetString("openai.api_key")
		if apiKey == "" {
			return nil, fmt.Errorf(
				"openai classifier selected but openai.api_key not configured",
			)
		}

		model := viper.GetString("openai.model")
		if model == "" {
			model = "gpt-4.1-mini"
		}

		client := openai.NewClient(
			option.WithAPIKey(apiKey),
		)

		return &classify.OpenAIClassifier{
			Client: client,
			Model:  model,
		}, nil

	default:
		return nil, fmt.Errorf(
			"unknown classifier provider: %s",
			provider,
		)
	}
}
