package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	slclient "github.com/maximilien/softlayer-go/client"
	//	datatypes "github.com/maximilien/softlayer-go/data_types"
)

const (
	ApiEndpoint = "https://api.softlayer.com/rest/v3"
	ApiUser     = "SL_API_USER"
	ApiKey      = "SL_API_KEY"
)

func main() {

	apiUser := os.Getenv(ApiUser)
	apiKey := os.Getenv(ApiKey)

	id := 12345

	client := slclient.NewSoftLayerClient(apiUser, apiKey)

	virtualGuestService, err := client.GetSoftLayer_Virtual_Guest_Service()

	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("No error found")
	}

	virtualGuest, err := virtualGuestService.GetObject(id)

	if err != nil {
		fmt.Println("errors:", err)
		return
	} else {
		fmt.Println("No error found")
	}

	spew.Printf("virtualGuest =%#+v\n", virtualGuest)

	vgPowerState, err := virtualGuestService.GetPowerState(id)

	spew.Printf("vgPowerState =%#+v\n", vgPowerState)
}

func FindHostname(hostname string) (bool, error) {

	apiUser := os.Getenv(ApiUser)
	apiKey := os.Getenv(ApiKey)

	client := slclient.NewSoftLayerClient(apiUser, apiKey)

	accountService, err := client.GetSoftLayer_Account_Service()
	if err != nil {
		return false, err
	}

	virtualGuests, err := accountService.GetVirtualGuests()

	if err != nil {
		return false, err
	}

	for _, guest := range virtualGuests {
		if strings.Contains(guest.Hostname, hostname) {
			fmt.Printf("Found guest %v ")
			return true, nil
		}
	}

}
