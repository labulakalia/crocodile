package cmd

import (
	"fmt"

	"github.com/labulaka521/crocodile/core/cert"
	"github.com/labulaka521/crocodile/core/version"
	"github.com/spf13/cobra"
)

// Version return current build message
func Version() *cobra.Command {
	cmdClient := &cobra.Command{
		Use:   "version",
		Short: "crocodile version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version   : %s\n", version.Version)
			fmt.Printf("Commit    : %s\n", version.Commit)
			fmt.Printf("BuildDate : %s\n", version.BuildDate)
		},
	}
	return cmdClient
}

// GeneratePemKey generate cert pem and key
func GeneratePemKey() *cobra.Command {
	var outdir string
	cmdClient := &cobra.Command{
		Use:   "cert",
		Short: "generate cert key",
		Run: func(cmd *cobra.Command, args []string) {
			if outdir == "" {
				outdir = "."
			}
			err := cert.GenerateCert(outdir)
			if err != nil {
				fmt.Println("GenerateCert failed: ", err)
			}
		},
	}
	// ioutil.TempDir(dir string, prefix string)
	cmdClient.Flags().StringVarP(&outdir, "dir", "d", "", "output pem dir[default current dir]")
	return cmdClient

}
