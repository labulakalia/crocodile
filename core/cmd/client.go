package cmd

import (
	"time"
	"context"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/router"
	"github.com/labulaka521/crocodile/core/schedule"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/labulaka521/crocodile/core/utils/resp"
	mylog "github.com/labulaka521/crocodile/core/utils/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net"
	"os"
	"strconv"
)

// Client crocodile client
func Client(version string) *cobra.Command {
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
			var maxretry = 10
			for i:=0;i<maxretry;i++{
				err = schedule.RegistryClient(version, intport)
				if err != nil {
					if err == context.DeadlineExceeded {
						err = resp.GetMsgErr(resp.ErrCtxDeadlineExceeded)
					}
					log.Error("registryClient failed", zap.Int("trytime", i+1),zap.Error(err))
					time.Sleep(time.Second)
					if i == maxretry - 1 {
						log.Fatal("registry client failed,already try 10 time")
					}
				} else {
					log.Info("registry success from server")
					break
				}
			}
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
