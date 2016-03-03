// Store.go
package persist

import (
	"errors"

	"github.com/cheyang/scloud/pkg/host"
	"github.com/cheyang/scloud/pkg/drivers"
)

type Store interface {

	// Exist returns whether a machine exists or not
	Exists(name string) (bool, error)

	// NewHost will initialize a new host machine
	NewHost(driver drivers.Driver) (host *host.Host, error)

	// Update persists with existing host
	Update(host *host.Host) error

	// List returns a list of all hosts in the store
	List() ([]*host.Host, error)

	// Get loads a host by name
	Load(name string) (*host.Host, error)

	// Remove removes a machine from the store
	Remove(name string) error

	// the Direcotory of the Store
	MyDir() string
}

var (
	HostEntryNotExistError = errors.New("Host Entry Not Exists!")

	HostEntryAlreadyExistError = errors.New("Host Entry Already Exists!")
)
