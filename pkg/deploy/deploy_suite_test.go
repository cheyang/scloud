package deploy

import (
	scloudLog "github.com/cheyang/scloud/pkg/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDeploy(t *testing.T) {

	err := scloudLog.InitLog()

	if err != nil {
		t.Errorf("Failed to init log")
	}

	defer scloudLog.CloseLog()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Deploy Suite")
}
