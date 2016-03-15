// deployer.go
package deploy

type Deployer interface {
	// if this deployer plugin installed
	IsSupported() error

	//	UnSupportedMsg() string

	// do the deployment work
	Deploy(d Deployment, workingDir string) error
}
