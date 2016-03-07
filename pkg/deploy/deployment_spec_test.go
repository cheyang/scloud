package deploy

import (
	//	. "github.com/cheyang/scloud/pkg/deploy"
	//	"github.com/cheyang/scloud/pkg/host"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deployment Spec", func() {

	var (
		spec  DeploymentSpec
		roles []*DeploymentRole
		//		h       []*host.Host

	)

	BeforeEach(func() {

		roles = []*DeploymentRole{
			&DeploymentRole{
				Name:           "kube-master",
				MaxNum:         1,
				MinNum:         1,
				groupName:      "k8s1",
				HostnamePrefix: "kubemaster",
			},
			&DeploymentRole{
				Name:   "kube-nodes",
				MaxNum: 2,
			},
			&DeploymentRole{
				Name:      "etcd",
				MaxNum:    3,
				groupName: "k8s1",
			},
			&DeploymentRole{
				Name:        "registry",
				IpAddresses: []string{"10.62.71.77"},
			},
		}

		reuseGroup := &ReuseGroup{
			Group: []*GroupMember{
				&GroupMember{
					GroupName: "k8s1",
					Members:   []string{"kube-master", "etcd"},
				},
			},
		}

		spec = DeploymentSpec{
			Roles:      roles,
			ReuseGroup: reuseGroup,
		}

	})

	Context("#Generate deployment spec", func() {
		It("create a new VM on Softlayer", func() {
			Expect(roles[3].Name).To(Equal("registry"))
			Expect(spec.GetTargetSize()).To(BeEquivalentTo(5))
			Expect(spec.GetLeastDeployableSize()).To(BeEquivalentTo(5))
		})
	})
})
