package deploy

import (
	//	. "github.com/cheyang/scloud/pkg/deploy"
	"github.com/cheyang/scloud/pkg/host"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deploy", func() {

	var (
		spec    DeploymentSpec
		roles   []*DeploymentRole
		h       []*host.Host
		specNum int
	)

	BeforeEach(func() {

		roles = []*DeploymentRole{
			DeploymentRole{
				Name:           "kube-master",
				MaxNum:         1,
				MinNum:         1,
				groupName:      "k8s1",
				hostnamePrefix: "kubemaster",
			},
			DeploymentRole{
				Name:   "kube-nodes",
				MaxNum: 2,
			},
			DeploymentRole{
				Name:      "etcd",
				MaxNum:    3,
				groupName: "k8s1",
			},
			DeploymentRole{
				Name:        "registry",
				IpAddresses: []string{"10.62.71.77"},
			},
		}

	})

	Context("#Generate deployment spec", func() {
		Expect(roles[3].Name).To(Equal(registry))
	})
})
