package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/JosephSalisbury/vm/provider/providerset"
)

var (
	rootCmd = &cobra.Command{
		Use:   "vm [command]",
		Short: "vm is a tool for managing CoreOS Container Linux Virtual Machines",
	}

	providerName string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&providerName,
		"provider",
		"p",
		string(providerset.DefaultProvider),
		fmt.Sprintf("which VM provider should be used, options are: %s", providerset.Names()),
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
