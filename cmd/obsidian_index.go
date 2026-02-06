package cmd

import (
	"github.com/saschakiefer/relay/internal/index"
	"github.com/saschakiefer/relay/internal/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var obsidianIndexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index the Obsidian vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		vaultPath := viper.GetString("obsidian.vault")
		indexPath := viper.GetString("obsidian.index")

		reader := vault.NewFilesystemVaultReader(vaultPath)
		builder := index.NewBuilder(reader)

		idx, err := builder.Build(ctx)
		if err != nil {
			return err
		}

		return index.Save(indexPath, idx)
	},
}

func init() {
	rootCmd.AddCommand(obsidianIndexCmd)
}
