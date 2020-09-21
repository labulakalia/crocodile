package main

import (
	"fmt"

	"github.com/labulaka521/crocodile/core/cmd"
	"github.com/labulaka521/crocodile/core/version"
	"github.com/spf13/cobra"
)

var (
	v string
	c string
	d string
)

// @title Crocidle API
// @version 1.0
// @description Crocodile Swaager JSON API
// @termsOfService https://github.com/labulaka521/crocodile

// @contact.name labulaka521
// @contact.url http://www.swagger.io/support
// @contact.email labulakalia@gmail.com

// @license.name MIT 2.0
// @license.url https://github.com/labulaka521/crocodile/blob/master/LICENSE

// @securityDefinitions.basic BasicAuth
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main111() {
	version.Commit = c
	version.Version = v
	version.BuildDate = d
	rootCmd := &cobra.Command{Use: "crocodile"}
	rootCmd.AddCommand(cmd.Client())
	rootCmd.AddCommand(cmd.Server())
	rootCmd.AddCommand(cmd.Version())
	rootCmd.AddCommand(cmd.GeneratePemKey())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("rootCmd.Execute failed", err.Error())
	}
}

type a []string

func main() {
	b := a{}
	fmt.Println(b)
}
