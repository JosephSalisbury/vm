package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/JosephSalisbury/vm/provider"
	"github.com/JosephSalisbury/vm/provider/providerset"
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

	forceMosh          bool
	forceSSH           bool
	username           string
	keyPath            string
	connectTimeout     time.Duration
	connectRetries     int
	connectRetryDelay  time.Duration
	connectSuccessTime time.Duration
)

// TODO: Add support for mosh.

func init() {
	connectCmd.Flags().BoolVar(&forceMosh, "force-mosh", false, "if true, only mosh will be used for the connection")
	connectCmd.Flags().BoolVar(&forceSSH, "force-ssh", false, "if true, only ssh will be used for the connection")
	connectCmd.Flags().StringVar(&username, "username", "joe", "username to use for ssh / mosh")
	connectCmd.Flags().StringVar(&keyPath, "key", "~/.ssh/dev-vm", "path to key to use for ssh")
	connectCmd.Flags().DurationVar(&connectTimeout, "connect-timeout", 3*time.Second, "connect timeout for ssh")
	connectCmd.Flags().IntVar(&connectRetries, "connect-retries", 5, "number of retries for ssh")
	connectCmd.Flags().DurationVar(&connectRetryDelay, "connect-retry-delay", 10*time.Second, "delay between retrying ssh connections")

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

	switch len(args) {
	case 0:
		if len(statuses) == 1 {
			ip = statuses[0].IP
			port = statuses[0].Port
		} else {
			return MultipleVMError
		}
	case 1:
		for _, status := range statuses {
			if status.ID == args[0] {
				ip = status.IP
				port = status.Port
			}
		}
	}

	for i := 0; i < connectRetries; i++ {
		logger.Printf("connecting to %s:%v", ip, port)

		sshCommand := exec.Command(
			"ssh",
			ip,
			"-p", strconv.Itoa(port),
			"-l", username,
			"-i", keyPath,
			"-o", "ForwardAgent=yes",
			"-o", "StrictHostKeyChecking=no",
			"-o", "UserKnownHostsFile=/dev/null",
			"-o", fmt.Sprintf("ConnectTimeout=%v", connectTimeout.Seconds()),
		)
		sshCommand.Stdin = os.Stdin
		sshCommand.Stdout = os.Stdout
		sshCommand.Stderr = os.Stderr

		err := sshCommand.Run()

		if err != nil {
			logger.Printf("could not connect, retrying in %v", connectRetryDelay)
			time.Sleep(connectRetryDelay)
		} else {
			return nil
		}
	}

	return nil
}
