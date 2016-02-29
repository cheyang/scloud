// softlayer.go
package softlayer

import (
	"fmt"
	"os"
	"strings"

	"github.com/cheyang/scloud/pkg/drivers"
	"github.com/cheyang/scloud/pkg/host"
	"github.com/cheyang/scloud/pkg/host/errs"
	"github.com/cheyang/scloud/pkg/state"
	slclient "github.com/cheyang/softlayer-go/client"
	datatypes "github.com/maximilien/softlayer-go/data_types"
)

var (
	apiUser   string
	apiKey    string
	apiClient *client.SoftLayerClient
)

type Driver struct {
	VirtualGuestTemplate *datatypes.SoftLayer_Virtual_Guest_Template
	*drivers.BaseDriver
	Id int
}

func init() {
	apiUser := os.Getenv("SL_API_USER")
	apiKey := os.Getenv("SL_API_KEY")

	if apiUser == "" || apiKey == "" {

		fmt.Println("Please don't forget to set SL_API_USER and SL_API_KEY before running command")
		os.Exit(1)
	}

	apiClient = &slclient.NewSoftLayerClient(apiUser, apiKey)
}

func NewDriver(hostname, storePath string) drivers.Driver {

	return &Driver{
		VirtualGuestTemplate: &datatypes.SoftLayer_Virtual_Guest_Template{},
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			storePath:   storePath,
		},
	}
}

func (d *Driver) DriverName() string {
	return "softlayer"
}

func (d *Driver) SetCreateConfigs(config interface{}) {
	if createConfig, ok := config.(*datatypes.SoftLayer_Virtual_Guest_Template); ok {
		d.VirtualGuestTemplate = createConfig
		d.VirtualGuestTemplate.Hostname = d.MachineName
	}
}

func (d *Driver) Create() error {
	virtualGuestService, err := apiClient.GetSoftLayer_Virtual_Guest_Service()

	if err != nil {
		return err
	}

	virtualGuest, err := virtualGuestService.CreateObject(virtualGuestTemplate)

	if err != nil {
		return err
	}

	d.Id = virtualGuest.Id

	return nil
}

// Check the VM has no active transcation and status is "RUNNING"
func (d *Driver) GetState() (state.State, error) {
	virtualGuestService, err := apiClient.GetSoftLayer_Virtual_Guest_Service()

	if err != nil {
		return state.None, err
	}

	activeTransactions, err := virtualGuestService.GetActiveTransactions(d.Id)

	if err != nil {
		return state.None, err
	}

	if len(activeTransactions) > 0 {
		fmt.Printf("active transactions for %s are %s", d.MachineName, activeTransactions)
		return state.Starting, err
	}

	vgPowerState, err := virtualGuestService.GetPowerState(d.Id)

	var vmState state.State
	switch strings.ToLower(vgPowerState.KeyName) {
	case "running":
		vmState = state.Running
	case "halted":
		vmState = state.Stopped
	default:
		vmState = state.None
	}
	return vmState, nil

}

func (d *Driver) PreCreateCheck() error {

	return validateCreateTemplate(d.VirtualGuestTemplate)
}

func validateCreateTemplate(createVirtualTemplate *datatypes.SoftLayer_Virtual_Guest_Template) error {

	if !host.validHostName(createVirtualTemplate.Hostname) {
		return errs.ErrInvalidHostname
	}

	if createVirtualTemplate.Datacenter.Name == "" {
		return fmt.Errorf("Missing required setting -- data center")
	}

	if createVirtualTemplate.Domain == "" {
		return fmt.Errorf("Missing required setting -- domain name")
	}

	if createVirtualTemplate.OperatingSystemReferenceCode == "" && createVirtualTemplate.BlockDeviceTemplateGroup.GlobalIdentifier == "" {

		return fmt.Errorf("Missing required setting -- OperationSystemReference Doe or Template Id")

	}
}
