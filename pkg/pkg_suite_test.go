package pkg_test

import (
	scloudLog "github.com/cheyang/scloud/pkg/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPkg(t *testing.T) {

	err := scloudLog.InitLog()

	if err != nil {
		t.Errorf("Failed to init log")
	}

	defer scloudLog.CloseLog()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Pkg Suite")
}
