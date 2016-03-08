// deployment.go
package deploy

import (
	"github.com/cheyang/scloud/pkg/host"
)

type Deployment struct {
	Nodes map[string][]*host.Host // role name: host
}

func NewDeployment(capacity int) Deployment {
	return Deployment{Nodes: make(map[string][]*host.Host, capacity)}
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
