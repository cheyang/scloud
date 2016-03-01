// utils.go
package drivers

import (
	"fmt"

	"os"

	"github.com/cheyang/scloud/pkg/utils"
)

func WaitForSSH(d Driver) error {
	if err := utils.WaitFor(sshAvailableFunc(d)); err != nil {
		fmt.Errorf("Too many retries waiting for SSH to be available.  Last error: %s", err)
	}

	return nil
}

func sshAvailableFunc(d Driver) func() bool {
	return func() bool {
		fmt.Fprintf(os.Stderr, "Getting to waitForSSH function for %s...\n", d.GetMachineName())

		if _, err := RunSSHCommandFromDriver(d, "exit 0"); err != nil {

			fmt.Fprintf(os.Stderr, "Error getting ssh command 'exit 0' : %s", err)
			return false
		}

		return true
	}
}

// To be implemented in future
func RunSSHCommandFromDriver(d Driver, command string) (string, error) {

	return command, nil
}
