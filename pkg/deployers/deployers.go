// deployers.go
package deployers

type Deployer interface {
	// if this deployer plugin installed
	IsSupported() bool

	UnSupportedMsg() string

	// do the deployment work
	Deploy(d deploy.Deployment) error
}
