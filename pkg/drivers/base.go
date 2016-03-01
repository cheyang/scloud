package drivers

import (
	"errors"
	"fmt"
	"net"
	"path/filepath"
)

const (
	DefaultSSHPort     = 22
	DefaultSSHUserName = "root"
)

type BaseDriver struct {
	IPAddress   string
	MachineName string
	SSHUser     string
	SSHPort     int
	SSHKeyPath  string
	StorePath   string
}

func (d *BaseDriver) DriverName() string {
	return "None"
}

// GetMachineName return the machine name, by default, it's hostname
func (d *BaseDriver) GetMachineName() string {
	return d.MachineName
}

func (d *BaseDriver) GetIP() (string, error) {

	if d.IPAddress == "" {
		return "", errors.New("IP Address is not set.")
	}

	ip := net.ParseIP(d.IPAddress)

	if ip == nil {
		return "", fmt.Errorf("IP Address is invalid: %s", d.IPAddress)
	}

	return d.IPAddress, nil
}

// GetSSHKeyPath returns the SSH Key path
func (d *BaseDriver) GetSSHKeyPath() string {

	if d.SSHKeyPath == "" {
		d.SSHKeyPath = d.ResolveStorePath("id_rsa")
	}

	return d.SSHKeyPath

}

// GetSSHUsername returns the ssh user name, root if not specified
func (d *BaseDriver) GetSSHUsername() string {

	if d.SSHUser == "" {
		d.SSHUser = DefaultSSHUserName
	}

	return d.SSHUser
}

// GetSSHPort returns the ssh port, 22 if not specified
func (d *BaseDriver) GetSSHPort() (int, error) {
	if d.SSHPort == 0 {
		d.SSHPort = DefaultSSHPort
	}

	return d.SSHPort, nil
}

func (d *BaseDriver) PreCreateCheck() error {
	return nil
}

// ResolveStorePath returns the store path where the machine is
func (d *BaseDriver) ResolveStorePath(file string) string {
	return filepath.Join(d.StorePath, "machines", d.MachineName, file)
}
