// planner.go
package deploy

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/cheyang/scloud/pkg/host"
	"github.com/cheyang/scloud/pkg/msg"
)

const (
	defaultStepSize = 1
	cStepSize       = "STEP_SIZE"
)

//var FailedToAddHostToPlanError error = errors.New("Failed to add host to the deployment plan!")

type Planner struct {
	ProvisionerObserver msg.Receiver
	msg.Sender
	Deployment // Operate on map, no need pointer
	*DeploymentSpec
	wait         *sync.WaitGroup
	stepSize     int
	FinishReport chan interface{}
}

func NewPlanner(spec *DeploymentSpec, wait *sync.WaitGroup) *Planner {

	size, err := strconv.Atoi(os.Getenv(cStepSize))

	if err != nil {
		size = defaultStepSize
	}

	return &Planner{
		//		TargetSize:     spec.GetTargetSize(),
		DeploymentSpec:      spec,
		Sender:              msg.NewQueue(spec.GetTargetSize()),
		ProvisionerObserver: nil,
		wait:                wait,
		stepSize:            size,
	}
}

// Register to cloud Provisioner
func (p *Planner) RegisterOberserver(r msg.Receiver) {
	p.ProvisionerObserver = r
}

// Loop processing host->role assignment
func (p *Planner) Run() {

	readyToPublish := false

	lastPublishSize := 0

	if p.ProvisionerObserver == nil {
		fmt.Fprintf(os.Stderr, "p.ProvisionerObserver is not set, exit!")
		return
	}

	for i := 0; i < p.DeploymentSpec.GetTargetSize(); i++ {
		fmt.Fprintf(os.Stderr, "Begin receiving %d times provision notification", i+1)
		host, err := p.WaitForHostReady()
		if err != nil {
			fmt.Fprintf(os.Stdout, "Finish receiving %d times provision: %s\n", i+1, err)
			continue
		}

		fmt.Fprintf(os.Stderr, "Finish receiving %d times provision notification entry %v", i+1, host)

		//		fmt.Fprintf(os.Stderr, "Add receiving %d times provision notification entry %v",  entry)
		p.AddHostToPlan(host)

		fmt.Fprintf(os.Stderr, "Begin Notifying the deploymnet manager with new deployment design %v\n", p.Deployment)

		// Check if there are enough VMs which are provisioned by cloud api
		if !readyToPublish {
			readyToPublish = p.CheckReadyToPublish()
		}

		if readyToPublish {

			if p.Deployment.Size()-lastPublishSize < p.stepSize {
				gap := p.Deployment.Size() - lastPublishSize
				fmt.Fprintf(os.Stderr, "Current deploymnet size is %d, the last deployment size %d, the gap between them is %d, while the expected gap is %d\n ", p.Deployment.Size(), lastPublishSize, gap, cStepSize)
				fmt.Fprintln(os.Stderr, "jump out the publish process. ")
				continue
			}

			fmt.Fprintln(os.Stderr, "Start the publish process. ")

			err := p.PublishDeployment()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to publish the deployment due to %s\n", err)
			}

			fmt.Fprintln(os.Stderr, "Finish the publish process. ")

		} else {
			fmt.Fprintf(os.Stderr, "The current deployment size is %d, which is less than minimum deployment size %d\n", p.Deployment.Size(), p.DeploymentSpec.GetLeastDeployableSize())
		}

		fmt.Fprintf(os.Stderr, "End Notifying the deploymnet manager with new deployment design %v\n", p.Deployment)
	}

	p.wait.Done()

	if !readyToPublish {
		panic("Failed to create deployment plan due to not enough machines are created.")
	}

}

/** Check if it's ready to publish the deployment plan to deployment manager,
* Basically, it depends on, the deployment reaches the mininal deployment size
 */

func (p *Planner) CheckReadyToPublish() bool {
	return p.Deployment.Size() >= p.DeploymentSpec.GetLeastDeployableSize()
}

func (p *Planner) PublishDeployment() error {
	fmt.Fprintf(os.Stdout, "Do publish for %v\n", p.Deployment.Size())

	return nil
}

// Add the host entry to plan
func (p *Planner) AddHostToPlan(h *host.Host) error {

	added := false

	fmt.Fprintf(os.Stderr, "Begin to Add host %s to plan ", h.Driver.GetMachineName())

	for _, role := range p.DeploymentSpec.Roles {
		if role.Match(h) {

			p.Deployment.Add(role.Name, h)

			// If groupName is empty, so no share with other role
			if role.groupName != "" {
				memberNames := p.DeploymentSpec.FindReuseGroupByName(role.groupName)

				for _, gMember := range memberNames {
					p.Deployment.Add(gMember, h)
				}

			}

			added = true

			break
		}
	}

	if !added {
		return fmt.Errorf("Failed to add host %s to the deployment plan!", h.Driver.GetMachineName())
	}

	return nil

}

func (p *Planner) WaitForHostReady() (*host.Host, error) {

	fmt.Fprintf(os.Stderr, "Begin receiving  provision notification\n")
	entry := p.ProvisionerObserver.Recieve()
	fmt.Fprintf(os.Stderr, "Finish receiving provision notification entry %v\n", entry)

	// Check if it's an error
	if err, ok := entry.(error); ok {
		fmt.Fprintf(os.Stderr, "Cloud provision error: %s\n", err)
		fmt.Fprintf(os.Stdout, "Cloud provision failed for: %s\n", err)
		return nil, err
	}

	if h, ok := entry.(*host.Host); ok {
		return h, nil
	}

}
