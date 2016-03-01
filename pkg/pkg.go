// pkg.go
package pkg

import (
	"fmt"
	"path/filepath"

	"github.com/cheyang/scloud/pkg/drivers"
	"github.com/cheyang/scloud/pkg/host"
	"github.com/cheyang/scloud/pkg/persist"
	"github.com/cheyang/scloud/pkg/state"
	"github.com/cheyang/scloud/pkg/utils"
)

func GetDefaultStore(clusterName string) *persist.FileStore {
	homeDir := utils.GetHomeDir()
	clusterDir := filepath.Join(homeDir, ".scloud", clusterName)
	return &persist.FileStore{
		Path: filepath.Join(homeDir),
	}
}

func Create(store persist.Store, host *host.Host) error {

	fmt.Println("Running pre-create check for ", host.Name, "...")

	if host.Driver.DriverName() != "None" {
		return fmt.Errorf("Not an implmented cloud driver")
	}

	if err := host.Driver.PreCreateCheck(); err != nil {

		return fmt.Errorf("Error with precheck for machien %s : %s", host.Name, err)
	}

	fmt.Println("Creating machine... for", host.Name, "...")

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
	fmt.Printf("Waiting for machine %s to be running, this may take a few minutes...\n", host.Name)

	if err := utils.WaitFor(utils.WaitFor(drivers.MachineInState(host.Driver, state.Running))); err != nil {
		return fmt.Errorf("Error waiting for machine %s to be running: %s", host.Name, err)
	}

	fmt.Println("Machine %s is running, waiting for SSH to be available ...\n", host.Name)

	if err := utils.WaitForSSH(host.Driver); err != nil {
		return fmt.Errorf("Error waiting %s for SSH: %s", host.Name, err)
	}

	return nil
}
