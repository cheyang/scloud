// filestore.go
package persist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
