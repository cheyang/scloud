// pkg.go
package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cheyang/scloud/pkg/drivers"
	"github.com/cheyang/scloud/pkg/host"
	//	"github.com/cheyang/scloud/pkg/host/errs"
	"github.com/cheyang/scloud/pkg/persist"
	"github.com/cheyang/scloud/pkg/state"
	"github.com/cheyang/scloud/pkg/utils"
)

func GetDefaultStore(clusterName string) *persist.FileStore {
	homeDir := utils.GetHomeDir()
	clusterDir := filepath.Join(homeDir, ".scloud", clusterName)
	return &persist.FileStore{
		Path: clusterDir,
	}
}

func Create(store persist.Store, host *host.Host) error {

	fmt.Fprintf(os.Stderr, "Running pre-create check for %s ... \n", host.Name)

	if host.Driver.DriverName() != "None" {
		return fmt.Errorf("Not an implmented cloud driver")
	}

	if err := host.Driver.PreCreateCheck(); err != nil {

		return fmt.Errorf("Error with precheck for machien %s : %s", host.Name, err)
	}

	fmt.Fprintf(os.Stderr, "Creating machine... for %s ...", host.Name)

	if err := host.Driver.Create(); err != nil {
		return fmt.Errorf("Error in driver during machine %s creation: %s", host.Name, err)
	}

	if err := store.NewHost(host); err != nil {
		return fmt.Errorf("Error with saving meta data for %s", host.Name)
	}

	if err := waitForReady(host); err != nil {
		return fmt.Errorf("Error with waiting for %s: %s", host.Name, err)
	}

	if err := store.Update(host); err != nil {
		return fmt.Errorf("Error with Saving for %s: %s", host.Name, err)
	}

	return nil

}

func waitForReady(host *host.Host) error {
	fmt.Fprintf(os.Stderr, "Waiting for machine %s to be running, this may take a few minutes...\n", host.Name)

	if err := utils.WaitFor(drivers.MachineInState(host.Driver, state.Running)); err != nil {
		return fmt.Errorf("Error waiting for machine %s to be running: %s", host.Name, err)
	}

	fmt.Fprintf(os.Stderr, "Machine %s is running, waiting for SSH to be available ...\n", host.Name)

	if err := drivers.WaitForSSH(host.Driver); err != nil {
		return fmt.Errorf("Error waiting %s for SSH: %s", host.Name, err)
	}

	return nil
}
