package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/spf13/cobra"

	"github.com/JosephSalisbury/vm/ignition"
	"github.com/JosephSalisbury/vm/ignition/ignitionset"
	"github.com/JosephSalisbury/vm/provider"
	"github.com/JosephSalisbury/vm/provider/providerset"
	"github.com/JosephSalisbury/vm/secrets"
)

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a VM",
		Long: `Create a new CoreOS Container Linux Virtual Machine.

With cloud providers, a VM type with resources closest to the requested resources will be chosen.
The correct CoreOS image for the specified channel will be chosen automatically, as well.

If a Container Linux Config is specified (with --ignition-path) instead of an Ignition Config,
ct will be called (if available) to compile the Container Linux Config into an Ignition Config.
If a URL is specified instead of a local path, the file will be downloaded and used.`,
		SilenceUsage: true,
		RunE:         createRun,
	}

	channel              string
	cpu                  int
	ensure               bool
	ram                  int
	ignitionName         string
	ignitionPath         string
	ignitionURL          string
	secretsDirectoryPath string
)

// TODO: Add support for passing a Container Linux Config.
// TODO: Add support for passing a URL.
// TODO: Add support for specifying secrets.

func init() {
	createCmd.Flags().StringVar(&channel, "channel", "stable", "channel for CoreOS Container Linux")
	createCmd.Flags().IntVar(&cpu, "cpu", 4, "number of CPU cores for the VM")
	createCmd.Flags().BoolVar(&ensure, "ensure", false, "ensure that only one VM is running (useful with 'vm create --ensure && vm connect')")
	createCmd.Flags().IntVar(&ram, "ram", 8, "amount of RAM (in GB) for the VM")
	createCmd.Flags().StringVar(
		&ignitionName,
		"ignition",
		string(ignitionset.DefaultIgnition),
		fmt.Sprintf("which Igniton should be used, options are: %s", ignitionset.Names()),
	)
	createCmd.Flags().StringVar(&ignitionPath, "ignition-path", "", "path to Ignition Config")
	createCmd.Flags().StringVar(&ignitionURL, "ignition-url", "https://raw.githubusercontent.com/JosephSalisbury/ignition/master/config.yaml", "url to Container Linux Config")
	createCmd.Flags().StringVar(&secretsDirectoryPath, "secrets", "/Users/joseph/secrets/", "path to directory containing secrets")

	rootCmd.AddCommand(createCmd)
}

func createRun(cmd *cobra.Command, args []string) error {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	s, err := secrets.New(secrets.Config{
		Directory: secretsDirectoryPath,
		Logger:    logger,
	})
	if err != nil {
		return err
	}

	url, err := url.Parse(ignitionURL)
	if err != nil {
		return err
	}

	i, err := ignitionset.New(ignition.Name(ignitionName), ignition.Config{
		Logger:  logger,
		Path:    ignitionPath,
		URL:     *url,
		Secrets: s,
	})
	if err != nil {
		return err
	}

	p, err := providerset.New(provider.Name(providerName), provider.Config{
		Logger: logger,
	})
	if err != nil {
		return err
	}

	if ensure {
		statuses, err := p.List()
		if err != nil {
			return err
		}

		switch len(statuses) {
		case 0:
			break
		case 1:
			return nil
		default:
			return MultipleVMError
		}
	}

	if err := p.Create(channel, i, cpu, ram); err != nil {
		return err
	}

	return nil
}
