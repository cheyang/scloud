// deployment.go
package deploy

type Deployment struct {
	Members map[string][]string // role name: sshhostname
}
