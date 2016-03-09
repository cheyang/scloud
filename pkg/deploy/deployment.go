// deployment.go
package deploy

import (
	"fmt"
	"os"

	"github.com/cheyang/scloud/pkg/host"
	"github.com/cheyang/scloud/pkg/utils"
)

type Deployment struct {
	Nodes map[string][]*host.Host // role name: host
}

func NewDeployment(capacity int) Deployment {
	return Deployment{Nodes: make(map[string][]*host.Host, capacity)}
}

//
func (d Deployment) Equals(t Deployment) bool {
	equal := true

	if len(d.Nodes) != len(t.Nodes) {
		fmt.Fprintf(os.StdOut, "len(d.Nodes) %d != len(t.Nodes) %d\n", len(d.Nodes), len(t.Nodes))
		return false
	}

	for k, v := range d.Nodes {

		if tk, ok := t.Nodes[k]; ok {
			if len(v) != len(tk) {
				fmt.Fprintf(os.StdOut, "len(v) %v %d != len(tk) %v %d\n", v, len(v), tk, len(tk))
				return false
			}

			for _, value := range v {

				has := utils.Contains([]interface{}{tk}, value)

				if !has {
					fmt.Fprintf(os.StdOut, " utils.Contains([]interface{}{tk}, value) %v %v\n", tk, value)
					return false
				}
			}

		} else {
			fmt.Fprintf(os.StdOut, "  t.Nodes[k] doesn't exist:%v, %v \n", t.Nodes, k)

			return false
		}

	}

	return equal
}

func (d Deployment) Add(name string, h *host.Host) {

	if v, ok := d.Nodes[name]; ok {
		d.Nodes[name] = append(v, h)
	} else {
		d.Nodes[name] = []*host.Host{h}
	}
}

func (d Deployment) FindHostsByName(name string) []*host.Host {

	v, _ := d.Nodes[name]

	return v

}

//Determine if the deployment is empty
func (d Deployment) IsEmpty() bool {
	return d.Nodes == nil
}

// caculate the number of hosts by name
func (d *Deployment) GetHostNumberByName(name string) int {

	hosts := d.FindHostsByName(name)

	if hosts == nil {
		return 0
	}

	return len(hosts)

}

// Get the size of deployment
func (d *Deployment) Size() int {

	uniqueHostMap := make(map[string]bool)

	for _, nodes := range d.Nodes {

		for _, node := range nodes {
			uniqueHostMap[node.Driver.GetMachineName()] = true
		}

	}

	return len(uniqueHostMap)
}
