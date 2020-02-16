package main

import (
	"fmt"
	"github.com/labulaka521/crocodile/core/cmd"
	"github.com/spf13/cobra"
)

var (
	version   string
	commit    string
	builddate string
)


// @title Crocidle API
// @version 1.0
// @description Crocodile Swaager JSON API
// @termsOfService https://github.com/labulaka521/crocodile

// @contact.name labulaka521
// @contact.url http://www.swagger.io/support
// @contact.email labulakalia@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.basic BasicAuth
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	rootCmd := &cobra.Command{Use: "crocodile"}
	rootCmd.AddCommand(cmd.Client(version))
	rootCmd.AddCommand(cmd.Server())
	rootCmd.AddCommand(cmd.Version(version, commit, builddate))
	rootCmd.AddCommand(cmd.GeneratePemKey())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("rootCmd.Execute failed", err.Error())
	}
}
