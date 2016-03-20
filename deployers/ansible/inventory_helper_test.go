package ansible_test

import (
	"fmt"

	. "github.com/cheyang/scloud/deployers/ansible"
	. "github.com/cheyang/scloud/pkg/engine"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Inventory Helper", func() {

	var (
		sections        map[string][]string
		invetoryManager *InventoryFile
	)

	BeforeEach(func() {

		invetoryManager = NewInventory()

		invetoryManager.AddSection("kube-masters", []string{"10.53.14.181"})

		invetoryManager.AddSection("kube-nodes", []string{"10.53.14.195", "10.53.14.166", "10.52.36.27"})

		invetoryManager.AddSection("etcd", []string{"10.53.14.181"})

		invetoryManager.AddSection("registry", []string{"10.62.71.77"})

		invetoryManager.AddSection("pop_svt:children", []string{"kube-masters", "kube-nodes", "etcd"})

	})

	Context("#SaveTo", func() {

		It("Create a file to save ansible inventories", func() {
			manager := NewInventory()

			err := manager.SaveTo("/tmp/result.txt")

			Expect(err).To(BeNil())

			cmd := NewCommand("diff", "target.txt", "/tmp/result.txt")

			err = cmd.Run()

			if err != nil {
				fmt.Println("error:", err)
			}

			Expect(err).To(BeNil())

			fmt.Println(cmd.GetPeriod())

		})

	})
})
