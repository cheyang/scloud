// deployment_spec.go
package deploy

import (
	"strings"

	"github.com/cheyang/scloud/pkg/host"
)

type DeploymentSpec struct {
	Roles []*DeploymentRole
	*ReuseGroup
	RolesMap   map[string]*DeploymentRole
	targetSize int
	leastSize  int // least size for deployable
}

// Initialize the role maps for the future usage
func (d *DeploymentSpec) InitRoleMaps() {

	// skip if the role maps already init
	if d.RolesMap != nil {
		return
	}

	d.RolesMap = make(map[string]*DeploymentRole)

	for _, role := range d.Roles {
		d.RolesMap[role.Name] = role
	}

}

// Find slice of group by name
func (d *DeploymentSpec) FindReuseGroupByName(name string) (members []string) {

	if d.ReuseGroup == nil {
		return members
	}

	// Need init groupMap
	if len(d.ReuseGroup.GroupMap) == 0 {
		if len(d.ReuseGroup.Group) == 0 {
			return members
		}

		d.InitRoleMaps()
		d.ReuseGroup.InitGroupMaps()

	}

	return d.ReuseGroup.GroupMap[name]
}

// Define the Reuse group for the machines which can be shared
type ReuseGroup struct {
	Group    []*GroupMember
	GroupMap map[string][]string
}

type GroupMember struct {
	GroupName string
	Members   []string
}

// Init the group maps
func (r ReuseGroup) InitGroupMaps() {

	if r.GroupMap != nil {
		return
	}

	r.GroupMap = make(map[string][]string)

	for _, g := range r.Group {

		r.GroupMap[g.GroupName] = g.Members

	}

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

	roleMap := make(map[string](*DeploymentRole))

	if spec.targetSize <= 0 {
		totalCount := 0
		sharedCount := 0

		for _, role := range spec.Roles {
			totalCount += role.MaxNum

			roleMap[role.Name] = role
		}

		for _, group := range spec.ReuseGroup.Group {

			max := 0

			for _, member := range group.Members {
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

// Get the least deployable size of the deployment
func (d *DeploymentSpec) GetLeastDeployableSize() int {
	roleMap := make(map[string](*DeploymentRole))

	if d.leastSize <= 0 {
		totalCount := 0
		sharedCount := 0

		for _, role := range d.Roles {

			num := role.MaxNum

			if role.MinNum > 0 {
				num := role.MinNum
			}

			totalCount += num

			roleMap[role.Name] = role
		}

		groups := d.ReuseGroup.Group

		for _, group := range groups {

			// Found the max of min num
			max := 0

			for _, member := range group.Members {

				num := roleMap[member].MaxNum

				if roleMap[member].MinNum > 0 {
					num := roleMap[member].MinNum
				}

				sharedCount += num

				if num > max {
					max = num
				}
			}

			sharedCount = sharedCount - max

		}

		d.targetSize = totalCount - sharedCount
	}

	return d.leastSize
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
