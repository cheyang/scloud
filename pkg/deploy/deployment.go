// deployment.go
package deploy

import (
	"github.com/cheyang/scloud/pkg/host"
)

type Deployment struct {
	Nodes map[string][]*host.Host // role name: host
}

func (d *Deployment) Add(name string, h *host.Host) {

	if v, ok := d.Nodes[name]; ok {
		d.Nodes[name] = append(v, h)
	} else {
		d.Nodes[name] = []*host.Host{h}
	}
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
