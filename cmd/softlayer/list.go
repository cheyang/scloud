// List.go
package main

import (
	"fmt"
	"os"

	datatypes "github.com/maximilien/softlayer-go/data_types"
	softlayer "github.com/maximilien/softlayer-go/softlayer"
)

func main() {
	name := "SIDRK8SNODE"
	hosts, err := FindGuestByHostname(name)

	fmt.Printf("hosts = %v", hosts)
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
