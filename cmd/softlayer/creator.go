// creator.go
package main

import (
	"github.com/cheyang/scloud/drivers/softlayer"
	sl_cloud "github.com/cheyang/scloud/drivers/softlayer"
	lib "github.com/cheyang/scloud/pkg"
	"github.com/cheyang/scloud/pkg/drivers"
	//	"github.com/cheyang/scloud/pkg/log"
	"github.com/cheyang/scloud/pkg/persist"
	"github.com/cheyang/scloud/pkg/state"
	"github.com/cheyang/scloud/pkg/utils"
	datatypes "github.com/cheyang/softlayer-go/data_types"
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
