package cmd

import (
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/router"
	"github.com/labulaka521/crocodile/core/schedule"
	"github.com/labulaka521/crocodile/core/utils/define"
	mylog "github.com/labulaka521/crocodile/core/utils/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net"
	"os"
	"strconv"
)

func CmdClient(version string) *cobra.Command {
	var (
		cfg string
	)
	cmdClient := &cobra.Command{
		Use:   "client",
		Short: "Start Run crocodile client",
		Run: func(cmd *cobra.Command, args []string) {
			if len(cfg) == 0 {
				_ = cmd.Help()
				os.Exit(0)
			}
			config.Init(cfg)
			mylog.Init()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			lis, err := router.GetListen(define.Client)
			if err != nil {
				return err
			}
			_, port, _ := net.SplitHostPort(lis.Addr().String())
			intport, _ := strconv.Atoi(port)
			err = schedule.RegistryClient(version, intport)
			if err != nil {
				log.Fatal("RegistryClient failed", zap.String("error", err.Error()))
			}
			err = router.Run(define.Client, lis)
			if err != nil {
				log.Fatal("router.Run failed", zap.Error(err))
			}
			return nil
		},
	}
	cmdClient.Flags().StringVarP(&cfg, "conf", "c", "", "server config [toml]")
	return cmdClient
}
