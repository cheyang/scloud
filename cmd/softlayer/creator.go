// creator.go
package main

import (
	sl_cloud "github.com/cheyang/scloud/drivers/softlayer"

	datatypes "github.com/maximilien/softlayer-go/data_types"
)

func main() {
	fmt.Println("Hello World!")

	err := scloudLog.InitLog()

	if err != nil {
		t.Errorf("Failed to init log")
	}

	defer scloudLog.CloseLog()

	var (
		sl_driver drivers.Driver
		store     persist.Store
	)
	names := string{"kubemastergo-001", "kubenodego-001", "kubenodego-002", "kubenodego-003", "kubenodego-004"}

	for _, name := range names {
		store = lib.GetDefaultStore(name)

		hostname := name

		sl_driver, err = sl_cloud.NewDriver(hostname, store.MyDir())
	}

}
