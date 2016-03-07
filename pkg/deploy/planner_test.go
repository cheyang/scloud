// deployment_test.go
package deploy

import (
	"fmt"
	//	. "github.com/cheyang/scloud/pkg/deploy"
	"os"
	"sync"

	"strconv"

	"github.com/cheyang/scloud/pkg/drivers"
	"github.com/cheyang/scloud/pkg/host"
	"github.com/cheyang/scloud/pkg/msg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Planner Test", func() {

	var (
		spec  *DeploymentSpec
		roles []*DeploymentRole
		//		h       []*host.Host
		deployment Deployment
		planner    *Planner
		waitgroup  *sync.WaitGroup
		hosts      []*host.Host
		num        int
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
				MaxNum: 5,
			},
			&DeploymentRole{
				Name:   "etcd",
				MaxNum: 3,
				//				MinNum:    1,
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

		spec = &DeploymentSpec{
			Roles:      roles,
			ReuseGroup: reuseGroup,
		}

		deployment = NewDeployment(spec.CountOfRoles())

		num = 2

		hosts = make([]*host.Host, num)

		for i := 0; i < num; i++ {
			hosts[i] = &host.Host{Name: strconv.Itoa(i),
				Driver: &drivers.BaseDriver{IPAddress: strconv.Itoa(i),
					MachineName: fmt.Sprintf("kubemaster-", strconv.Itoa(i))},
			}
			fmt.Fprintf(os.Stdout, "exec method GetMachineName for %s\n", hosts[i].Driver.GetMachineName())
		}
	})

	Context("#Planner", func() {
		It("Monitoring Provision work and add to the deployment plan", func() {
			//			Expect(roles[3].Name).To(Equal("registry"))
			//			Expect(spec.GetTargetSize()).To(BeEquivalentTo(8))
			//			Expect(spec.GetLeastDeployableSize()).To(BeEquivalentTo(8))

			waitgroup.Add(1)

			queue := msg.NewQueue(num)

			defer queue.Close()

			planner = NewPlanner(spec, waitgroup)

			planner.RegisterOberserver(queue)

			go planner.Run()

			for i := 0; i < num; i++ {
				go func(n int) {
					fmt.Fprintf(os.Stdout, "call %d\n", n)
					queue.Send(hosts[n])
				}(i)
			}

			waitgroup.Wait()

			Expect(roles[3].Name).To(Equal("registry"))
		})
	})
})
