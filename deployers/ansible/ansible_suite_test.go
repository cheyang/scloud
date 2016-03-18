package ansible_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAnsible(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ansible Suite")
}
