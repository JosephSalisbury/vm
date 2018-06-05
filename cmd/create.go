package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/JosephSalisbury/vm/ignition"
	"github.com/JosephSalisbury/vm/provider"
	"github.com/JosephSalisbury/vm/provider/providerset"
)

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a VM",
		Long: `Create a new CoreOS Container Linux Virtual Machine.

With cloud providers, a VM type with resources closest to the requested resources will be chosen.
The correct CoreOS image for the specified channel will be chosen automatically, as well.

If a Container Linux Config is specified (with --ignition) instead of an Ignition Config,
ct will be called (if available) to compile the Container Linux Config into an Ignition Config.
If a URL is specified instead of a local path, the file will be downloaded and used.`,
		SilenceUsage: true,
		RunE:         createRun,
	}

	channel      string
	cpu          int
	ram          int
	ignitionPath string
)

// TODO: Add support for passing a Container Linux Config.
// TODO: Add support for passing a URL.
// TODO: Add support for specifying secrets.

func init() {
	createCmd.Flags().StringVar(&channel, "channel", "stable", "channel for CoreOS Container Linux")
	createCmd.Flags().IntVar(&cpu, "cpu", 2, "number of CPU cores for the VM")
	createCmd.Flags().IntVar(&ram, "ram", 4, "amount of RAM (in GB) for the VM")
	createCmd.Flags().StringVar(&ignitionPath, "ignition", "./config.ign", "path to Ignition Config")

	rootCmd.AddCommand(createCmd)
}

func createRun(cmd *cobra.Command, args []string) error {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	i, err := ignition.New(ignition.Config{
		Logger: logger,
		Path:   ignitionPath,
	})
	if err != nil {
		return err
	}

	c := provider.Config{
		Logger: logger,
	}

	p, err := providerset.New(provider.Name(providerName), c)
	if err != nil {
		return err
	}

	if err := p.Create(channel, i, cpu, ram); err != nil {
		return err
	}

	return nil
}
