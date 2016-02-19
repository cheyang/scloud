// filestore.go
package persist

import (
	"io/ioutil"
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

}
