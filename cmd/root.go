/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package cmd

import (
	"github.com/saschakiefer/relay/internal/config"
	"github.com/saschakiefer/relay/internal/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Version = "0.0.1"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: Version,
	Use:     "relay",
	Short:   "Relay bridges handwritten notes to your digital systems",
	Long: `relay captures handwritten notes, understands them using an agentic AI,
and routes the results into tools like Obsidian and Reminders.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		debug := viper.GetBool("debug")
		logging.Init(debug)

		return nil
	},
}

func Execute() {
	// minimal logger for bootstrap
	logging.Init(false)

	cobra.CheckErr(config.Init())
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug logging")

	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}
