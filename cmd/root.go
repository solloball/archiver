package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "It's a simple archiver",
}

func handleErr(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		handleErr(err)
	}

}
