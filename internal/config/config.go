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
		Language string
	}
}

var AppConfig *Config

func Init() error {
	viper.SetConfigName("relay-ocr-config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.config/")
	viper.AddConfigPath("$HOME/.config/relay-ocr")

	if err := viper.ReadInConfig(); err != nil {
		log.Warn().Msg("No config file found, using defaults")
	}

	AppConfig = &Config{}
	return viper.Unmarshal(&AppConfig)
}
