// planner.go
package deploy

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cheyang/scloud/pkg/host"
	"github.com/cheyang/scloud/pkg/msg"
)

const (
	defaultStepSize = 3
	cStepSize       = "STEP_SIZE"
)

//var FailedToAddHostToPlanError error = errors.New("Failed to add host to the deployment plan!")

type Planner struct {
	ProvisionerObserver msg.Receiver
	msg.Sender
	Deployment // Operate on map, no need pointer
	*DeploymentSpec
	//	wait         *sync.WaitGroup
	stepSize     int
	FinishReport chan interface{}
}

func NewPlanner(spec *DeploymentSpec, chs chan interface{}) *Planner {

	size, err := strconv.Atoi(os.Getenv(cStepSize))

	if err != nil {
		size = defaultStepSize
	}

	deployment := NewDeployment(spec.CountOfRoles())

	//	chs := make(chan interface{}, capacity)

	return &Planner{
		//		TargetSize:     spec.GetTargetSize(),
		DeploymentSpec:      spec,
		Sender:              msg.NewQueue(spec.GetTargetSize()),
		ProvisionerObserver: nil,
		Deployment:          deployment,
		//		wait:                wait,
		FinishReport: chs,
		stepSize:     size,
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
		fmt.Fprintf(os.Stderr, "p.ProvisionerObserver is not set, exit!\n")
		return
	}

	for i := 0; i < p.DeploymentSpec.GetTargetSize(); i++ {
		fmt.Fprintf(os.Stderr, "Begin receiving %d times provision notification\n", i+1)
		host, err := p.WaitForHostReady()
		if err != nil {
			fmt.Fprintf(os.Stdout, "Finish receiving %d times provision: %s\n", i+1, err)
			continue
		}

		fmt.Fprintf(os.Stderr, "Finish receiving %d times provision notification entry %v\n", i+1, host)

		//		fmt.Fprintf(os.Stderr, "Add receiving %d times provision notification entry %v",  entry)
		err = p.AddHostToPlan(host)

		if err != nil {
			fmt.Println("Add host to plan error:", err)
		}

		fmt.Fprintf(os.Stderr, "Begin Notifying the deploymnet manager with new deployment design %v\n", p.Deployment)

		// Check if there are enough VMs which are provisioned by cloud api
		if !readyToPublish {
			readyToPublish = p.CheckReadyToPublish()
		}

		if readyToPublish {

			if p.Deployment.Size()-lastPublishSize < p.stepSize {
				gap := p.Deployment.Size() - lastPublishSize
				fmt.Fprintf(os.Stderr, "Current deploymnet size is %d, the last published deployment size %d, the gap between them is %d, while the expected gap is %d\n ", p.Deployment.Size(), lastPublishSize, gap, p.stepSize)
				fmt.Fprintln(os.Stderr, "jump out the publish process. ")
				if i+1 < p.DeploymentSpec.GetTargetSize() {
					continue
				}
			}

			fmt.Fprintln(os.Stderr, "Start the publish process. ")

			err := p.PublishDeployment()

			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to publish the deployment due to %s\n", err)
			}

			lastPublishSize = p.Deployment.Size()

			fmt.Fprintln(os.Stderr, "Finish the publish process. ")

		} else {
			fmt.Fprintf(os.Stderr, "The current deployment size is %d, which is less than minimum deployment size %d\n", p.Deployment.Size(), p.DeploymentSpec.GetLeastDeployableSize())
		}

		fmt.Fprintf(os.Stderr, "End Notifying the deploymnet manager with new deployment design %v\n", p.Deployment)
	}

	if !readyToPublish {
		//		panic("Failed to create deployment plan due to not enough machines are created.")

		p.FinishReport <- fmt.Errorf("Failed to create deployment plan due to not enough machines are created.")
	} else if p.Deployment.Size() < p.DeploymentSpec.GetTargetSize() {
		p.FinishReport <- fmt.Errorf("The least size of the deployment plan can be created, but not all of them are created.")
	} else {
		p.FinishReport <- "Deployment plan done successfully."
	}

}

/** Check if it's ready to publish the deployment plan to deployment manager,
* Basically, it depends on, the deployment reaches the mininal deployment size
 */

func (p *Planner) CheckReadyToPublish() bool {

	ready := true

	// If least deployment size can't reach, return false
	if p.Deployment.Size() < p.DeploymentSpec.GetLeastDeployableSize() {
		fmt.Fprintf(os.Stderr, "Deployment.Size() is %d, while DeploymentSpec.GetLeastDeployableSize is %d\n", p.Deployment.Size(), p.DeploymentSpec.GetLeastDeployableSize())
		ready = false
		return ready
	}

	//	for k, v := range p.Deployment {
	//		role := p.DeploymentSpec.FindRoleByName(k)

	//		if role
	//	}

	for k, _ := range p.DeploymentSpec.GetRoleMaps() {

		if p.Deployment.GetHostNumberByName(k) < p.DeploymentSpec.GetDeployableSizeByName(k) {
			fmt.Fprintf(os.Stderr, "%s 's Deployment.GetHostNumberByName(k) is %d, while DeploymentSpec.GetDeployableSizeByName(k) is %d\n", k, p.Deployment.GetHostNumberByName(k), p.DeploymentSpec.GetDeployableSizeByName(k))

			ready = false
			return ready
		}

	}

	return ready
}

func (p *Planner) PublishDeployment() error {
	fmt.Fprintf(os.Stderr, "Do publish for %v\n", p.Deployment.Size())

	return nil
}

// Add the host entry to plan
func (p *Planner) AddHostToPlan(h *host.Host) error {

	added := false

	fmt.Fprintf(os.Stderr, "Begin to Add host %v to plan \n", h.Driver.GetMachineName())

	for _, role := range p.DeploymentSpec.Roles {
		if role.Match(h) && role.MaxNum > p.Deployment.GetHostNumberByName(role.Name) {

			p.Deployment.Add(role.Name, h)

			fmt.Fprintf(os.Stderr, "Added host %v to role %v\n", h.Driver.GetMachineName(), role.Name)

			// If groupName is empty, so no share with other role
			if role.groupName != "" {
				fmt.Fprintf(os.Stderr, "Found group name %v of role %v\n", role.groupName, role.Name)

				memberNames := p.DeploymentSpec.FindReuseGroupByName(role.groupName)

				fmt.Fprintf(os.Stderr, "Found members %s of group %s\n", role.groupName, memberNames)

				for _, gMember := range memberNames {

					if gMember != role.Name && p.Deployment.GetHostNumberByName(gMember) < p.DeploymentSpec.FindRoleByName(gMember).MaxNum {
						fmt.Fprintf(os.Stderr, "Added host %v to role %v\n", h.Driver.GetMachineName(), gMember)
						p.Deployment.Add(gMember, h)
					}
				}

			}

			added = true

			break
		}
	}

	if !added {
		return fmt.Errorf("Failed to add host %s to the deployment plan %s!", h.Driver.GetMachineName(), p.Deployment)
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
	} else {
		return nil, fmt.Errorf("return unknown object %s", entry)
	}

}
