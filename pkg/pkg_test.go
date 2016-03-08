package pkg_test

import (
	"fmt"

	"github.com/cheyang/scloud/drivers/softlayer"
	sl_cloud "github.com/cheyang/scloud/drivers/softlayer"
	lib "github.com/cheyang/scloud/pkg"
	"github.com/cheyang/scloud/pkg/drivers"
	//	"github.com/cheyang/scloud/pkg/log"
	"github.com/cheyang/scloud/pkg/persist"
	"github.com/cheyang/scloud/pkg/state"
	"github.com/cheyang/scloud/pkg/utils"
	datatypes "github.com/maximilien/softlayer-go/data_types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pkg", func() {

	var (
		sl_driver drivers.Driver
		store     persist.Store
		err       error
		name      string
	)

	BeforeEach(func() {

		name = "apmwdc-001"

		store = lib.GetDefaultStore(name)

		hostname := name

		sl_driver, err = sl_cloud.NewDriver(hostname, store.MyDir())

		Expect(err).To(BeNil())

		//			Expect(err.Error()).To(ContainSubstring("Failed to init sl client!"))

		Expect(sl_driver).ToNot(BeNil())

		fmt.Println(sl_driver)

		//			err = sl_driver.PreCreateCheck()

		//			Expect(err).To(HaveOccurred())
	})

	Context("#Create", func() {

		It("create a new VM on Softlayer", func() {
			virtualGuestTemplate := &datatypes.SoftLayer_Virtual_Guest_Template{
				Domain:    "softlayergo.com",
				StartCpus: 2,
				MaxMemory: 2048,
				Datacenter: datatypes.Datacenter{
					Name: "dal05",
					//					Name: "wdc04",
				},
				NetworkComponents: []datatypes.NetworkComponents{datatypes.NetworkComponents{
					MaxSpeed: 1000,
				}},
				SshKeys:                  []datatypes.SshKey{datatypes.SshKey{Id: 3922}},
				HourlyBillingFlag:        true,
				LocalDiskFlag:            true,
				BlockDeviceTemplateGroup: &datatypes.BlockDeviceTemplateGroup{GlobalIdentifier: "00b8c96d-287a-4dba-b253-dab68ffdf56a"},
				//				PrimaryBackendNetworkComponent: &datatypes.PrimaryBackendNetworkComponent{NetworkVlan: datatypes.NetworkVlan{Id: 282238}},
				PrimaryBackendNetworkComponent: &datatypes.PrimaryBackendNetworkComponent{NetworkVlan: datatypes.NetworkVlan{Id: 28223}},
				//				PrimaryBackendNetworkComponent: &datatypes.PrimaryBackendNetworkComponent{NetworkVlan: datatypes.NetworkVlan{Id: 1191337}},
				PrivateNetworkOnlyFlag: true,
			}

			fmt.Println(sl_driver)

			sl_driver.SetCreateConfigs(virtualGuestTemplate)

			real_driver, ok := sl_driver.(*softlayer.Driver)

			Expect(ok).To(BeTrue())

			Expect(real_driver.MachineName).To(Equal(name))

			fmt.Println(real_driver.VirtualGuestTemplate)

			err = sl_driver.PreCreateCheck()

			fmt.Println("PreCheck...", err)

			err = sl_driver.Create()

			//			real_driver.Id = 16407243

			fmt.Println("Create...", err)

			Expect(err).To(BeNil())

			if err != nil {
				fmt.Println("Create Error", err)
			}

			if err == nil {
				if err = utils.WaitFor(drivers.MachineInState(sl_driver, state.Running)); err != nil {
					fmt.Printf("Error waiting for machine %s to be running: %s\n", real_driver.MachineName, err)
				}
			}

			sshHostname, err := sl_driver.GetSSHHostname()

			fmt.Println(sshHostname)

			fmt.Println(err)

		})
	})
})
