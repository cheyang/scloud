// deployer.go
package ansible

import (
	"fmt"
	"os"
	osexec "os/exec"

	"github.com/cheyang/scloud/pkg/deploy"
)

const (
	BinName = "ansible"
)

type Deployer struct {
	CmdPath     string
	PlaybookDir string
}

// If support current Deployer
func (d *Deployer) IsSupported() error {
	_, err := osexec.LookPath(d.CmdPath)

	if err != nil {
		return fmt.Errorf("CmdPath %s is not found.", d.CmdPath)
	}

	dir, err := os.Stat(d.PlaybookDir)

	if err != nil {
		return fmt.Errorf("Playbook Path %s can't be found", d.PlaybookDir)
	}

	if !dir.IsDir() {
		return fmt.Errorf("Playbook Path %s is not an directory", d.PlaybookDir)
	}

	return nil
}

// do the deployment work
func (d *Deployer) Deploy(d deploy.Deployment, workingDir string) error {

}
