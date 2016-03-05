// deployment.go
package deployment

import (
	"strings"

	"github.com/cheyang/scloud/pkg/host"
)

type Deployment struct {
	Nodes map[string]*host.Host // role name: host
}

type DeploymentSpec struct {
	Roles []*DeploymentRole
	*ReuseGroup
	targetSize int
}

// Define the Reuse group for the machines which can be shared
type ReuseGroup struct {
	Group []*GroupMember
}

type GroupMemnber struct {
	GroupName string
	Members   []string
}

type DeploymentRole struct {
	MaxNum int // The Number to deploy

	MinNum int // the minium number which can do the deploy

	Name string

	HostnamePrefix string

	IpAddresses []string // ip addresses which can be defined before provisioning

	//	ShareWith *[]DeploymentRole

	groupName string // if it's not set, means not shared with other role
}

// Get the target size of the deployment spec
func (spec *DeploymentSpec) GetTargetSize() int {

	var roleMap map[string](*DeploymentRole) = make(map[string]([]*DeploymentRole))

	if spec.targetSize <= 0 {
		totalCount := 0
		sharedCount := 0

		for _, role := range Roles {
			totalCount += role.MaxNum

			roleMap[role.Name] = role
		}

		for _, group := range spec.ReuseGroup {

			max := 0

			for _, member := range group.members {
				sharedCount += roleMap[member].MaxNum

				if roleMap[member].MaxNum > max {
					max = roleMap[member].MaxNum
				}
			}

			sharedCount = sharedCount - max

		}

		spec.targetSize = totalCount - sharedCount
	}

	return spec.targetSize
}

// Check if this host should be used as
func (r DeploymentRole) Match(h *host.Host) bool {
	match := true

	if r.HostnamePrefix != "" {
		if !strings.HasPrefix(h.GetMachineName(), r.HostnamePrefix) {
			match = false
		}
	}

	return match
}
