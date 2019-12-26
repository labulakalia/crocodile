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

func main() {
	rootCmd := &cobra.Command{Use: "crocodile"}
	rootCmd.AddCommand(cmd.CmdClient(version))
	rootCmd.AddCommand(cmd.CmdServer())
	rootCmd.AddCommand(cmd.CmdVersion(version, commit, builddate))
	rootCmd.AddCommand(cmd.GeneratePemKey())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("rootCmd.Execute failed", err.Error())
	}
}
