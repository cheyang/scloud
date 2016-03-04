// planner.go
package deploy

import (
	"fmt"
	"os"

	"github.com/cheyang/scloud/pkg/msg"
)

type Planner struct {
	ProvisionerObserver msg.Receiver
	msg.Sender

	Deployment // Operate on map, no need pointer

	*DeploymentSpec

	TargetSize int
}

func NewPlanner(spec *DeploymentSpec) *Planner {

	return &Planner{
		TargetSize:     spec.GetTargetSize(),
		DeploymentSpec: spec,
		Sender:         msg.NewQueue(targetSize),
		Receiver:       nil,
	}
}

// Register to cloud Provisioner
func (p *Planner) RegisterOberserver(r msg.Receiver) {
	p.ProvisionerObserver = r
}

func (p *Planner) OnPlanning() {
	if p.ProvisionerObserver == nil {
		fmt.Fprintf(os.Stderr, "p.ProvisionerObserver is not set, exit!")
		return
	}

	for i := 0; i < p.TargetSize; i++ {
		fmt.Fprintf(os.Stderr, "Begin to receive %d times provision notification", i+1)
		entry := p.ProvisionerObserver.Recieve()
		fmt.Fprintf(os.Stderr, "Finish receiving %d times provision notification entry %v", i+1, entry)
		p.AddHostToPlan(entry)
	}

}

func (p *Planner) AddHostToPlan(h *host.Host) {
		
		
		for _, role range p.DeploymentSpec.Roles{
			
		}
}
