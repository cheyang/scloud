package deploy

import (
	//	. "github.com/cheyang/scloud/pkg/deploy"
	"os"
	"strconv"

	"github.com/cheyang/scloud/pkg/drivers"
	"github.com/cheyang/scloud/pkg/host"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func GenerateDeployment(maps map[string]int, h []*host.Host) Deployment {
	d := Deployment{Nodes: make(map[string][]*host.Host)}

	point := 0
	for k, v := range maps {
		d.Nodes[k] = h[point:v]
		point = v
	}

	return d
}

func ReverseDeployment(d Deployment) Deployment {

	target := Deployment{Nodes: make(map[string][]*host.Host)}

	for k, v := range d.Nodes {
		target.Nodes[k] = v
	}

	for k, v := range target.Nodes {

		length := len(v)

		slice := make([]*host.Host, 0, length)

		for i := length - 1; i > 0; i-- {
			slice = append(slice, v[i])
		}

		target.Nodes[k] = slice
	}

	return target
}

var _ = Describe("Deployment", func() {

	var (
		deployment1 Deployment
		deployment2 Deployment

		hosts1 []*host.Host
		//		hosts2 []*host.Host

		num int
	)

	BeforeEach(func() {

	})

	Context("#Compare the data ", func() {
		It("order is different", func() {
			num = 10

			hosts1 = make([]*host.Host, num)

			for i := 0; i < num; i++ {
				hosts1[i] = &host.Host{Name: strconv.Itoa(i),
					Driver: &drivers.BaseDriver{IPAddress: strconv.Itoa(i),
						MachineName: fmt.Sprintf("kubemaster-", strconv.Itoa(i))},
				}
				fmt.Fprintf(os.Stdout, "exec method GetMachineName for %s\n", hosts1[i].Driver.GetMachineName())
			}

			map1 := map[string]int{"kube-master": 1, "kube-node": 8, "etcd": 1}

			deployment1 = GenerateDeployment(map1, hosts1)

			deployment2 = ReverseDeployment(deployment1)

			fmt.Printf("deployment1 %v", deployment1)
			fmt.Printf("deployment2 %v", deployment2)

			Expect(deployment1.Equals(deployment2)).To(BeTrue())
		})
	})
})
