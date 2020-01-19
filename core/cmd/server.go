package cmd

import (
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/router"
	"github.com/labulaka521/crocodile/core/schedule"
	"github.com/labulaka521/crocodile/core/utils/define"
	mylog "github.com/labulaka521/crocodile/core/utils/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

// Server crocodile server
func Server() *cobra.Command {
	var (
		cfg string
	)
	cmdServer := &cobra.Command{
		Use:   "server",
		Short: "Start Run crocodile server",
		Run: func(cmd *cobra.Command, args []string) {
			if len(cfg) == 0 {
				cmd.Help()
				os.Exit(0)
			}
			config.Init(cfg)
			mylog.Init()
			model.InitDb()
			model.InitRabc()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			lis, err := router.GetListen(define.Server)
			if err != nil {
				log.Fatal("listen failed", zap.String("error", err.Error()))
			}

			err = schedule.Init()
			if err != nil {
				log.Fatal("init schedule failed", zap.String("error", err.Error()))
			}
			err = router.Run(define.Server, lis)
			if err != nil {
				log.Error("router.Run error", zap.Error(err))
			}
			return nil
		},
	}
	cmdServer.Flags().StringVarP(&cfg, "conf", "c", "", "server config [toml]")
	return cmdServer
}
