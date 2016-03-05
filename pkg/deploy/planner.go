// planner.go
package deploy

import (
	"errors"
	"fmt"
	"os"

	"github.com/cheyang/scloud/pkg/msg"
)

var FailedToAddHostToPlanError error = errors.New("Failed to add host to the deployment plan!")

type Planner struct {
	ProvisionerObserver msg.Receiver
	msg.Sender

	Deployment // Operate on map, no need pointer

	*DeploymentSpec
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

// Add the host entry to plan
func (p *Planner) AddHostToPlan(h *host.Host) error {

	added := false

	for _, role := range p.DeploymentSpec.Roles {
		if role.Match(h) {

			p.Deployment.Add(role.Name, h)

			if role.groupName != "" {
				gMembers := p.DeploymentSpec.FindReuseGroupByName(role.groupName)

				for _, gMember := range gMembers {
					p.Deployment.Add(gMember.Name, h)
				}

			}

			added = true

			break
		}
	}

	if !added {
		return FailedToAddHostToPlanError
	}

	return nil

}
