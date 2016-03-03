package msg_test

import (
	scloudLog "github.com/cheyang/scloud/pkg/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMsg(t *testing.T) {

	err := scloudLog.InitLog()

	if err != nil {
		t.Errorf("Failed to init log")
	}

	defer scloudLog.CloseLog()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Msg Suite")
}
