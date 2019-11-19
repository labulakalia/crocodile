package cmd

import (
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/router"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/labulaka521/crocodile/core/utils/log"
	"github.com/spf13/cobra"
	"os"
)

func Client() *cobra.Command {
	var (
		cfg string
	)
	cmdClient := &cobra.Command{
		Use:   "client",
		Short: "Start Run crocodile client",
		Run: func(cmd *cobra.Command, args []string) {
			if len(cfg) == 0 {
				cmd.Help()
				os.Exit(0)
			}
			config.Init(cfg)
			log.Init()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return router.Run(define.Client)
		},
	}
	cmdClient.Flags().StringVarP(&cfg, "conf", "c", "", "server config [toml]")
	return cmdClient
}
