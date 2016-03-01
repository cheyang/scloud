package pkg_test

import (
	"github.com/cheyang/scloud/drivers/softlayer"
	sl_cloud "github.com/cheyang/scloud/drivers/softlayer"
	lib "github.com/cheyang/scloud/pkg"
	"github.com/cheyang/scloud/pkg/drivers"
	"github.com/cheyang/scloud/pkg/persist"
	datatypes "github.com/maximilien/softlayer-go/data_types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pkg", func() {

	var (
		sl_driver drivers.Driver
		store     persist.Store
	)

	BeforeEach(func() {
		store := lib.GetDefaultStore("mytest")

		hostname := "mytesthost"

		sl_driver, err := sl_cloud.NewDriver(hostname, store.Path)

		Expect(err).To(BeNil())

		//			Expect(err.Error()).To(ContainSubstring("Failed to init sl client!"))

		Expect(sl_driver).ToNot(BeNil())

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
				},
				NetworkComponents: []datatypes.NetworkComponents{datatypes.NetworkComponents{
					MaxSpeed: 1000,
				}},
				SshKeys:                  []datatypes.SshKey{datatypes.SshKey{Id: 3922}},
				HourlyBillingFlag:        true,
				LocalDiskFlag:            true,
				BlockDeviceTemplateGroup: &datatypes.BlockDeviceTemplateGroup{GlobalIdentifier: "00b8c96d-287a-4dba-b253-dab68ffdf56a"},
				PrimaryNetworkComponent:  &datatypes.PrimaryNetworkComponent{NetworkVlan: datatypes.NetworkVlan{Id: 282238}},
			}

			sl_driver.SetCreateConfigs(virtualGuestTemplate)

			real_driver, ok := sl_driver.(*softlayer.Driver)

			Expect(ok).To(BeTrue())

			Expect(real_driver.MachineName).To(Equal("mytesthost"))

		})
	})
})
