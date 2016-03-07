// deployment.go
package deploy

import (
	"strings"

	"github.com/cheyang/scloud/pkg/host"
)

type Deployment struct {
	Nodes map[string][]*host.Host // role name: host
}

func (d *Deployment) Add(name string, h *host.Host) {

	if v, ok := d.Nodes[role.Name]; ok {
		d.Nodes[role.Name] = append(v, h)
	} else {
		d.Nodes[role.Name] = []*host.Host{h}
	}
}

// Get the size of deployment
func (d *Deployment) Size() int {

	uniqueHostMap := make(map[string]bool)

	for _, v := range d.Nodes {
		uniqueHostMap[v.GetMachineName()] = true
	}

	return len(uniqueHostMap)
}

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

	for _, role = range d.Roles {
		d.RolesMap[role.Name] = role
	}

}

// Find slice of group by name
func (d *DeploymentSpec) FindReuseGroupByName(name string) (members []*DeploymentRole) {

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

		r.Group[g.GroupName] = g.Members
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

// Get the least deployable size of th
func (d *DeploymentSpec) GetLeastDeployableSize() {
	var roleMap map[string](*DeploymentRole) = make(map[string]([]*DeploymentRole))

	if spec.leastSize <= 0 {
		totalCount := 0
		sharedCount := 0

		for _, role := range Roles {

			num := role.MaxNum

			if role.MinNum > 0 {
				num := role.MinNum
			}

			totalCount += num

			roleMap[role.Name] = role
		}

		for _, group := range spec.ReuseGroup {

			// Found the max of min num
			max := 0

			for _, member := range group.members {

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

		spec.targetSize = totalCount - sharedCount
	}

	return spec.leastSize
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
