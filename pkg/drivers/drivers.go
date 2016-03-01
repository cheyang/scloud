// driver.go
package drivers

import (
	"errors"
	"fmt"
	"os"

	"github.com/cheyang/scloud/pkg/state"
)

type Driver interface {

	// Create a host using driver's config
	Create() error

	// DriverName returns the name of the driver as it is registered
	DriverName() string

	// Set the config for creating VM
	SetCreateConfigs(config interface{})

	//Precheck before the create request
	PreCreateCheck() error

	// GetIP returns an IP or hostname that this host is available at
	// e.g. 1.2.3.4 or abc.com
	GetIP() (string, error)

	// GetMachineName returns the name of the machine
	GetMachineName() string

	// GetSSHHostname returns hostname for use with ssh
	GetSSHHostname() (string, error)

	// GetSSHKeyPath returns key path for use with ssh
	GetSSHKeyPath() string

	// GetSSHPort returns port for use with ssh
	GetSSHPort() (int, error)

	// GetSSHUsername returns username for use with ssh
	GetSSHUsername() string

	// Remove a host
	//	Remove() error

	// GetState returns the state that the host is in (running, stopped, etc)
	GetState() (state.State, error)
}

var (
	ErrHostIsNotReachable = errors.New("Host is not reachable by ssh")
)

func MachineInState(d Driver, desireState state.State) func() bool {

	return func() bool {
		currentState, err := d.GetState()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in getting machine %s state: %s\n", d.GetMachineName(), err)
		}

		if currentState == desireState {
			return true
		}

		return false
	}

}
