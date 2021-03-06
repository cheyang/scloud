// deployment_test.go
package deploy

import (
	"fmt"
	//	. "github.com/cheyang/scloud/pkg/deploy"
	"os"
	//	"sync"

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
		//		waitgroup  *sync.WaitGroup
		hosts []*host.Host
		num   int
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
				MinNum: 1,
				MaxNum: 5,
			},
			&DeploymentRole{
				Name:   "etcd",
				MinNum: 1,
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
			GroupMap: make(map[string][]string),
		}

		spec = &DeploymentSpec{
			Roles:      roles,
			ReuseGroup: reuseGroup,
		}

		deployment = NewDeployment(spec.CountOfRoles())

		num = spec.GetTargetSize()

		spec.InitGroupMaps()

		fmt.Println("spec.groupMap:", spec.ReuseGroup)

		fmt.Println("The target size num: ", num)

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

			//			var waitgroup sync.WaitGroup

			//			waitgroup.Add(1)

			planReport := make(chan interface{}, 1)

			queue := msg.NewQueue(num)

			defer queue.Close()

			planner = NewPlanner(spec, planReport)

			fmt.Println("spec.groupMap:", spec.ReuseGroup)

			planner.RegisterOberserver(queue)

			go planner.Run()

			for i := 0; i < num; i++ {
				go func(n int) {
					fmt.Fprintf(os.Stdout, "call %d\n", n)
					queue.Send(hosts[n])
				}(i)
			}

			result := <-planReport

			fmt.Println("result is", result)

			Expect(roles[3].Name).To(Equal("registry"))
		})
	})
})
