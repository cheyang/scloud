package utils_test

import (
	. "github.com/cheyang/scloud/pkg/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hostname Helper", func() {
	var (
		nums []string
	)

	BeforeEach(func() {
		nums = []string{"001", "002", "010", "00f", "xxy"}
	})

	Context("#GetCurrentMaxExt", func() {
		It("Get the Max Value", func() {
			value := GetCurrentMaxExt(nums)
			Expect(value).To(Equal("010"))
		})
	})

})
