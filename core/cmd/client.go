package cmd

import (
	"net"
	"os"
	"strconv"

	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/router"
	"github.com/labulaka521/crocodile/core/schedule"
	"github.com/labulaka521/crocodile/core/utils/define"
	mylog "github.com/labulaka521/crocodile/core/utils/log"
	"github.com/labulaka521/crocodile/core/version"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Client crocodile client
func Client() *cobra.Command {
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
			schedule.InitWorker()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			lis, err := router.GetListen(define.Client)
			if err != nil {
				return err
			}
			_, port, _ := net.SplitHostPort(lis.Addr().String())
			intport, _ := strconv.Atoi(port)
			go schedule.RegistryClient(version.Version, intport)
			err = router.Run(define.Client, lis)
			if err != nil {
				log.Error("router.Run error", zap.Error(err))
			}
			return nil
		},
	}
	cmdClient.Flags().StringVarP(&cfg, "conf", "c", "", "server config [toml]")
	return cmdClient
}
