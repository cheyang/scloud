// List.go
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cheyang/scloud/pkg/state"

	slclient "github.com/maximilien/softlayer-go/client"
	datatypes "github.com/maximilien/softlayer-go/data_types"
	softlayer "github.com/maximilien/softlayer-go/softlayer"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var (
	TIMEOUT          time.Duration = 30 * time.Minute
	POLLING_INTERVAL time.Duration = 2 * time.Minute
)

func main() {
	name := "SIDRK8SNODE"
	hosts, err := FindGuestByHostname(name)

	if err != nil {
		fmt.Println("err:", err)
	}

	fmt.Printf("hosts = %v", hosts)

	service, err := CreateVirtualGuestService()

	if err != nil {
		fmt.Println("err:", err)
	}

	reload_OS_Config = datatypes.Image_Template_Config{
		ImageTemplateId: "c3b41ce1-21f0-41d5-8e4d-d10be596d4f3",
	}

	for _, h := range hosts {
		fmt.Println(h.PrimaryBackendIpAddress)
		service.ReloadOperatingSystem(h.Id, reload_OS_Config)
	}
}

func GetClient() (softlayer.Client, error) {

	apiUser := os.Getenv("SL_API_USER")
	apiKey := os.Getenv("SL_API_KEY")

	if apiUser == "" || apiKey == "" {

		fmt.Println("Please don't forget to set SL_API_USER and SL_API_KEY before running command")
		return nil, fmt.Errorf("apiUser and key are not setting.")
	}

	return slclient.NewSoftLayerClient(apiUser, apiKey), nil
}

func CreateVirtualGuestService() (softlayer.SoftLayer_Virtual_Guest_Service, error) {

	client, err := GetClient()

	if err != nil {
		return nil, err
	}

	virtualGuestService, err := client.GetSoftLayer_Virtual_Guest_Service()
	if err != nil {
		return nil, err
	}

	return virtualGuestService, nil
}

func CreateAccountService() (softlayer.SoftLayer_Account_Service, error) {
	client, err := GetClient()

	if err != nil {
		return nil, err
	}

	virtualAccountService, err := client.GetSoftLayer_Account_Service()
	if err != nil {
		return nil, err
	}

	return virtualAccountService, nil
}

func FindGuestByHostname(name string) ([]datatypes.SoftLayer_Virtual_Guest, error) {
	accountService, err := CreateAccountService()
	if err != nil {
		return []datatypes.SoftLayer_Virtual_Guest{}, err
	}

	virtualGuests, err := accountService.GetVirtualGuests()
	if err != nil {
		return []datatypes.SoftLayer_Virtual_Guest{}, err
	}

	targetVirtualGuests := []datatypes.SoftLayer_Virtual_Guest{}
	for _, vGuest := range virtualGuests {
		if strings.Contains(vGuest.Hostname, name) {
			targetVirtualGuests = append(targetVirtualGuests, vGuest)
		}
	}

	return targetVirtualGuests, nil
}

func GetState(id int) (state.State, error) {

	apiClient, err := GetClient()

	if err != nil {
		return nil, err
	}

	virtualGuestService, err := apiClient.GetSoftLayer_Virtual_Guest_Service()

	if err != nil {
		return state.None, err
	}

	activeTransactions, err := virtualGuestService.GetActiveTransactions(id)

	if err != nil {
		return state.None, err
	}

	if len(activeTransactions) > 0 {
		fmt.Fprintf(os.Stderr, "active transactions for %d are %s", id, activeTransactions)
		return state.Starting, err
	}

	vgPowerState, err := virtualGuestService.GetPowerState(id)

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

func WaitForVirtualGuestToRunning(virtualGuestId int) {

	Eventually(func() int {
		state, err := GetState(virtualGuestId)
		Expect(err).ToNot(HaveOccurred())
		fmt.Printf("----> virtual guest: %d, has state %s\n", virtualGuestId, state)
		return state
	}, TIMEOUT, POLLING_INTERVAL).Should(BeEquivalentTo(state.Running), "failed waiting for virtual guest to run")
}
