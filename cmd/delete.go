package cmd

import (
	"log"
	"os"

	"github.com/JosephSalisbury/vm/provider"
	"github.com/JosephSalisbury/vm/provider/providerset"
	"github.com/spf13/cobra"
)

var (
	deleteCmd = &cobra.Command{
		Use:   "delete [vm-id]",
		Short: "Delete a VM",
		Long: `Delete a VM.
If more than one VM is available, the VM ID must be passed.`,
		Aliases:      []string{"rm"},
		SilenceUsage: true,
		RunE:         deleteRun,
	}
)

// TODO: Add --add.

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func deleteRun(cmd *cobra.Command, args []string) error {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	c := provider.Config{
		Logger: logger,
	}

	p, err := providerset.New(provider.Name(providerName), c)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		statuses, err := p.List()
		if err != nil {
			return err
		}

		if len(statuses) == 1 {
			if err := p.Delete(statuses[0].ID); err != nil {
				return err
			}
		}

		if len(statuses) != 1 {
			return MultipleVMError
		}
	}
	if len(args) == 1 {
		if err := p.Delete(args[0]); err != nil {
			return err
		}

	}

	return nil
}
