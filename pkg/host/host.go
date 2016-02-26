// host.go
package host

import (
	"errors"
	"regexp"

	"github.com/cheyang/pkg/drivers"
)

var (
	validHostNameChars                = `^[a-zA-Z0-9][a-zA-Z0-9\-\.]*$`
	validHostNamePattern              = regexp.MustCompile(validHostNameChars)
	errMachineMustBeRunningForUpgrade = errors.New("Error: machine must be running to upgrade.")
)

type Host struct {
	Driver      drivers.Driver
	DriverName  string
	HostOptions *Options
	Name        string
	RawDriver   []byte
}

type Options struct {
	Driver string
	Memory int
	Disk   int
}

func validHostName(name string) bool {
	return validHostNamePattern.MatchString(name)
}

func (h *Host) RunSSHCommand(command string) (string, error) {
	return drivers.RunSSHCommandFromDriver(h.Driver, command)
}
