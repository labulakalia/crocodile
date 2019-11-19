package main

import (
	"fmt"
	"github.com/labulaka521/crocodile/core/cmd"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	rootCmd := &cobra.Command{Use: "crocodile"}
	rootCmd.AddCommand(cmd.Server())
	rootCmd.AddCommand(cmd.Client())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Execute failed", err)
		os.Exit(1)
	}
}
