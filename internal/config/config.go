/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	OCR struct {
		Google struct {
			Credentials string `mapstructure:"credentials"`
		} `mapstructure:"google"`
	} `mapstructure:"ocr"`

	Classify struct {
		Provider string `mapstructure:"provider"`
	} `mapstructure:"classify"`

	OpenAI struct {
		APIKey string `mapstructure:"api_key"`
		Model  string `mapstructure:"model"`
	} `mapstructure:"openai"`

	Obsidian struct {
		Vault string `mapstructure:"vault"`
		Index string `mapstructure:"index"`
	} `mapstructure:"obsidian"`
}

var AppConfig *Config

func Init() error {
	viper.SetConfigName("relay-ocr-config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.config/")
	viper.AddConfigPath("$HOME/.config/relay-ocr")

	// Defaults (safe ones only)
	viper.SetDefault("classify.provider", "openai")
	viper.SetDefault("openai.model", "gpt-4.1-mini")

	if err := viper.ReadInConfig(); err != nil {
		log.Warn().Msg("No config file found, using defaults")
	}

	AppConfig = &Config{}
	return viper.Unmarshal(&AppConfig)
}
