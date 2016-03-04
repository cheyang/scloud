// deployment.go
package deployment

import (
	"github.com/cheyang/scloud/pkg/host"
)

type Deployment struct {
	Nodes map[string]*host.Host // role name: host
}

type DeploymentSpec struct {
	Roles      []*DeploymentRole
	ReuseGroup map[string]([]string)
	targetSize int
}

type ReuseGroup struct {
	*GroupMember
}

type GroupMemnber struct {
	GroupName string
	Members   []string
}

// Get the target size of the deployment spec
func (spec *DeploymentSpec) GetTargetSize() int {

	var groups map[string]([]*DeploymentRole) = make(map[string]([]*DeploymentRole))

	if spec.targetSize <= 0 {
		totalCount := 0
		sharedCount := 0

		for _, role := range Roles {
			totalCount += role.MaxNum

			groups[role.Name] = role

		}

		for _, value := range ReuseGroup {

		}

		spec.targetSize = totalCount - sharedCount
	}

	return spec.targetSize
}

type DeploymentRole struct {
	MaxNum int // The Number to deploy

	MinNum int // the minium number which can do the deploy

	Name string

	HostnamePrefix string

	IpAddresses []string // ip addresses which can be defined before provisioning

	ShareWith *[]DeploymentRole
}

// Check if this host should be used as
func (role Role) ShouldBeAssign(h *host.Host) bool {

}
