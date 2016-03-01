package pkg_test

import (
	sl_cloud "github.com/cheyang/scloud/drivers/softlayer"
	lib "github.com/cheyang/scloud/pkg"
	"github.com/cheyang/scloud/pkg/persist"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pkg", func() {

	Context("#Create", func() {

		It("create a new VM on Softlayer", func() {
			store := lib.GetDefaultStore("mytest")

			hostname := "mytesthost"

			sl_driver := sl_cloud.NewDriver(hostname, store.Path)

			sl_driver.PreCreateCheck()
		})
	})
})