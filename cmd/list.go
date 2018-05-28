package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/JosephSalisbury/vm/provider"
	"github.com/JosephSalisbury/vm/provider/providerset"
)

var (
	listCmd = &cobra.Command{
		Use:          "list",
		Short:        "List all VMs",
		Aliases:      []string{"ls"},
		SilenceUsage: true,
		RunE:         listRun,
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
}

func listRun(cmd *cobra.Command, args []string) error {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	c := provider.Config{
		Logger: logger,
	}

	p, err := providerset.New(provider.Name(providerName), c)
	if err != nil {
		return err
	}

	statuses, err := p.List()
	if err != nil {
		return err
	}

	if len(statuses) == 0 {
		return nil
	}

	fmt.Printf("ID\tIP\t\tPORT\n")
	for _, status := range statuses {
		fmt.Printf("%s\t%s\t%v\n", status.ID, status.IP, status.Port)
	}

	return nil
}
