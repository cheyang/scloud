// utils.go
package drivers

import (
	"fmt"

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
		fmt.Printf("Getting to waitForSSH function for %s...", d.GetMachineName())

		if _, err := RunSSHCommandFromDriver(d, "exit 0"); err != nil {

			fmt.Printf("Error getting ssh command 'exit 0' : %s", err)
			return false
		}

		return true
	}
}

// To be implemented in future
func RunSSHCommandFromDriver(d Driver, command string) (string, error) {

	return command, nil
}
