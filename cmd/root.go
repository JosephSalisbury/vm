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

	googleCredentialsFilePath string
	googlePreemptible         bool
	googleProject             string
	googleZone                string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&providerName,
		"provider",
		"p",
		string(providerset.DefaultProvider),
		fmt.Sprintf("which VM provider should be used, options are: %s", providerset.Names()),
	)

	rootCmd.PersistentFlags().StringVar(&googleCredentialsFilePath, "google-credentials-file-path", "/Users/joseph/secrets/vm-credentials.json", "path to Google Cloud credentials JSON file")
	rootCmd.PersistentFlags().BoolVar(&googlePreemptible, "google-preemptible", true, "Whether to use premptible instances on Google Cloud")
	rootCmd.PersistentFlags().StringVar(&googleProject, "google-project", "dev-vm-216913", "ID of the Google Cloud project")
	rootCmd.PersistentFlags().StringVar(&googleZone, "google-zone", "europe-west1-b", "Google Cloud zone to use")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
