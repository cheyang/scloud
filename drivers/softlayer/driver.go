// softlayer.go
package softlayer

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/cheyang/scloud/pkg/drivers"
	"github.com/cheyang/scloud/pkg/host"
	"github.com/cheyang/scloud/pkg/host/errs"
	"github.com/cheyang/scloud/pkg/state"
	slclient "github.com/maximilien/softlayer-go/client"
	datatypes "github.com/maximilien/softlayer-go/data_types"
)

var (
	apiUser   string
	apiKey    string
	apiClient *slclient.SoftLayerClient
	initDone  bool = false
	once      sync.Once
)

type Driver struct {
	VirtualGuestTemplate *datatypes.SoftLayer_Virtual_Guest_Template
	*drivers.BaseDriver
	Id int
}

func setup() {
	apiUser := os.Getenv("SL_API_USER")
	apiKey := os.Getenv("SL_API_KEY")

	if apiUser == "" || apiKey == "" {

		fmt.Println("Please don't forget to set SL_API_USER and SL_API_KEY before running command")
		return
	}

	apiClient = slclient.NewSoftLayerClient(apiUser, apiKey)

	accountService, err := apiClient.GetSoftLayer_Account_Service()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Init sl cloud failed due to %s \n", err)
		return
	}

	accountStatus, err := accountService.GetAccountStatus()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Init sl cloud failed due to %s \n", err)
		return
	} else if strings.ToLower(accountStatus.Name) != "active" {
		fmt.Fprintf(os.Stderr, "Account status is %s, not as Active expected \n", accountStatus.Name)
		return
	}

	initDone = true
}

func NewDriver(hostName, storePath string) (drivers.Driver, error) {

	once.Do(setup)

	if !initDone {
		return &Driver{}, fmt.Errorf("Failed to init sl client!")
	}

	return &Driver{
		VirtualGuestTemplate: &datatypes.SoftLayer_Virtual_Guest_Template{},
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
	}, nil
}

func (d *Driver) DriverName() string {
	return "softlayer"
}

func (d *Driver) SetCreateConfigs(config interface{}) error {
	if createConfig, ok := config.(*datatypes.SoftLayer_Virtual_Guest_Template); ok {
		d.VirtualGuestTemplate = createConfig
		d.VirtualGuestTemplate.Hostname = d.MachineName
	} else {
		return fmt.Errorf("Not compatible interface %s", config)
	}

	return nil
}

func (d *Driver) Create() error {
	virtualGuestService, err := apiClient.GetSoftLayer_Virtual_Guest_Service()

	if err != nil {
		return err
	}

	virtualGuestTemplate := d.VirtualGuestTemplate

	virtualGuest, err := virtualGuestService.CreateObject(*virtualGuestTemplate)

	if err != nil {
		return err
	}

	if virtualGuest.Id <= 0 {
		return fmt.Errorf("Failed to retrieve the instance id of %s", d.GetMachineName())
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
		fmt.Fprintf(os.Stderr, "active transactions for %s are %s", d.MachineName, activeTransactions)
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

func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

func (d *Driver) GetIP() (string, error) {
	if d.IPAddress != "" {
		return d.IPAddress, nil
	}

	virtualGuestService, err := apiClient.GetSoftLayer_Virtual_Guest_Service()

	if err != nil {
		return "", err
	}

	virtualGuest, err := virtualGuestService.GetObject(d.Id)

	if err != nil {
		return "", err
	}

	d.IPAddress = virtualGuest.PrimaryBackendIpAddress

	return d.IPAddress, nil
}

func validateCreateTemplate(createVirtualTemplate *datatypes.SoftLayer_Virtual_Guest_Template) error {

	if !host.ValidHostName(createVirtualTemplate.Hostname) {
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

	return nil
}
