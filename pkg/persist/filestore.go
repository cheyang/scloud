// filestore.go
package persist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FileStore struct {
	Path string
}

func (this FileStore) getMachinesDir() string {
	return filepath.Join(this.Path, "machines")
}

func (this FileStore) saveToFile(data []byte, file string) error {
	return ioutil.WriteFile(file, data, 0600)
}

func (this FileStore) save(host *host.Host) error {

	if rpcClientDriver, ok := host.Driver.(RpcClientDriver); ok {
		data, err := rpcClientDriver.GetConfigRaw()

		if err != nil {
			return fmt.Errorf("Error getting raw config for driver: %s", err)
		}

		host.RawDriver = data
	}

	data, err = json.MarshalIndent(host, "", "    ")

	if err != nil {
		return err
	}

	hostpath := filepath.Join(this.getMachinesDir(), host.Name)

	if err := os.MkdirAll(hostpath); err != nil {
		return err
	}

	return this.saveToFile(data, filepath.Join(hostpath, config.json))

}

// Determine if the name exists
func (this FileStore) Exists(name string) (bool, error) {
	_, err := os.Stat(filepath.Join(this.getMachinesDir(), name))

	if err == nil {
		return true, nil
	}

	return false, err
}

// Remove the host with this name
func (this FileStore) Remove(name string) error {
	hostpath := filepath.Join(this.getMachinesDir(), name)

	return os.RemoveAll(hostpath)
}

func (this FileStore) List() ([]*host.Host, error) {
	fileInfos, err := ioutil.ReadDir(this.getMachinesDir())

	if err != nil {
		return nil, err
	}

	hosts := make([]*host.Host, len(fileInfos))

	for i, file := range fileInfos {

		if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			host, err := Load(file.Name())

			if err != nil {
				fmt.Printf("error loading host %q: %s", file.Name(), err)
				continue
			}

			hosts[i] = host
		}
	}

	return hosts, nil

}

func (this FileStore) Load(name string) (*host.Host, error) {
	exist, err := this.Exists(name)

	if !exist {
		return nil, err
	}

	host := &host.Host{name: name}

	host, err = this.loadConfig(host)

	return host, err
}

func (this FileStore) loadConfig(h *host.Host) (*host.Host, error) {
	data, err := ioutil.ReadFile(filepath.Join(s.getMachinesDir(), h.Name, "config.json"))
	if err != nil {
		return err
	}

	// Remember the machine name so we don't have to pass it through each
	// struct in the migration.
	name := h.Name

	if err := json.Unmarshal(data, h); err != nil {
		return &host.Host{}, err
	}

	h.Name = name

	return h, err
}

func (this FileStore) NewHost(host *host.Host) error {
	hostDir := filepath.Join(this.getMachinesDir(), host.Name)

	_, err := os.Stat(hostDir)

	// if the directory has already existed
	if !os.IsNotExist(err) {
		return HostEntryAlreadyExistError
	}

	err = this.save(host)

	return err

}
