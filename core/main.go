package main

import (
	"fmt"
	"github.com/labulaka521/crocodile/core/cmd"
	version2 "github.com/labulaka521/crocodile/core/version"
	"github.com/spf13/cobra"
)

var (
	version string
	commit  string
)

func main() {
	rootCmd := &cobra.Command{Use: "crocodile"}
	rootCmd.AddCommand(cmd.CmdClient())
	rootCmd.AddCommand(cmd.CmdServer())
	rootCmd.AddCommand(cmdVersion())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("rootCmd.Execute failed", err.Error())
	}
}

func cmdVersion() *cobra.Command {
	cmdClient := &cobra.Command{
		Use:   "version",
		Short: "crocodile version",
		Run: func(cmd *cobra.Command, args []string) {
			version2.Version = version
			version2.Commit = commit
			fmt.Printf("Version:    %s\n", version2.Version)
			fmt.Printf("Git Commit: %s\n", version2.Commit)
		},
	}
	return cmdClient

}
