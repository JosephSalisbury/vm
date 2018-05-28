package cmd

import (
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/JosephSalisbury/vm/provider"
	"github.com/JosephSalisbury/vm/provider/providerset"
	"github.com/spf13/cobra"
)

var (
	connectCmd = &cobra.Command{
		Use:   "connect [vm-id]",
		Short: "Connect to a VM",
		Long: `Connect to a running VM.
If more than one VM is available, the VM ID must be passed.

If mosh is available locally it will be used to connect, otherwise ssh will be used.
If mosh fails to connect, ssh will be used as a fallback.
If --force-mosh or --force-ssh are specified, those methods will be used specifically.`,
		Aliases:      []string{"ssh"},
		SilenceUsage: true,
		RunE:         connectRun,
	}

	forceMosh bool
	forceSSH  bool
)

func init() {
	connectCmd.Flags().BoolVar(&forceMosh, "force-mosh", false, "if true, only mosh will be used for the connection")
	connectCmd.Flags().BoolVar(&forceSSH, "force-ssh", false, "if true, only ssh will be used for the connection")

	rootCmd.AddCommand(connectCmd)
}

func connectRun(cmd *cobra.Command, args []string) error {
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

	var ip string
	var port int

	if len(args) == 0 {
		if len(statuses) == 1 {
			ip = statuses[0].IP
			port = statuses[0].Port
		}

		if len(statuses) > 1 {
			return MultipleVMError
		}
	}

	if len(args) == 1 {
		for _, status := range statuses {
			if status.ID == args[0] {
				ip = status.IP
				port = status.Port
			}
		}
	}

	logger.Printf("connecting to %s:%v", ip, port)

	// TODO: Take username from args
	// TODO: Take key from args
	// TODO: ignore user keys etc.
	sshCommand := exec.Command("ssh", ip, "-p", strconv.Itoa(port), "-l", "joe", "-i", "~/.ssh/dev-vm")
	sshCommand.Stdin = os.Stdin
	sshCommand.Stdout = os.Stdout
	sshCommand.Stderr = os.Stderr

	if err := sshCommand.Run(); err != nil {
		return err
	}

	return nil
}
