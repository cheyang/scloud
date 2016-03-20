// deployer.go
package ansible

import (
	"fmt"
	"os"
	osexec "os/exec"
	"strings"

	"github.com/cheyang/scloud/pkg/deploy"
)

const (
	BinName = "ansible"
)

type Deployer struct {
	cmdPath      string
	playbookDir  string
	groupVarFile string
	environment  string // dev, test, staging, production
	keyFile      string // ssh key filename
}

// If support current Deployer
func (d *Deployer) IsSupported() error {
	_, err := osexec.LookPath(d.cmdPath)

	if err != nil {
		return fmt.Errorf("CmdPath %s is not found.", d.cmdPath)
	}

	dir, err := os.Stat(d.playbookDir)

	if err != nil {
		return fmt.Errorf("Playbook Path %s can't be found", d.playbookDir)
	}

	if !dir.IsDir() {
		return fmt.Errorf("Playbook Path %s is not an directory", d.playbookDir)
	}

	file, err := os.Stat(d.groupVarFile)

	if err != nil {
		return fmt.Errorf("groupVarFile %s can't be found", d.groupVarFile)
	}

	if file.IsDir() {
		return fmt.Errorf("groupVarFile %s is not a file", d.groupVarFile)
	}

	return nil
}

// do the deployment work
func (d *Deployer) Deploy(deployment deploy.Deployment, workingDir string) error {
	return nil
}

//Generate configuration file, and return string content
func (d *Deployer) createInventoryfile(deployment deploy.Deployment, filename string) error {

	keys := make([]string)

	sections := make(map[string]([]string))

	for key, hosts := range deployment {
		keys = append(keys, key)

		ips := make([]string)

		for _, h := range hosts {
			ip, err := h.Driver.GetIP()
			if err != nil {
				return err
			}

			ips = append(ips, ip)
		}

		sections[key] = ips
	}

	childrenKey := fmt.Sprintf("[%s:children]", d.environment)

	sections[childrenKey] = keys

	invetory_file := NewInventory(sections)

	inventory_file.SaveTo(filename)

}
