package main

import (
	"fmt"
	"os"

	slclient "github.com/cheyang/softlayer-go/client"
	"github.com/davecgh/go-spew/spew"
	//	datatypes "github.com/maximilien/softlayer-go/data_types"
)

const (
	ApiEndpoint = "https://api.softlayer.com/rest/v3"
	ApiUser     = "SL_USERNAME"
	ApiKey      = "SL_API_KEY"
)

func main() {

	apiUser := os.Getenv(ApiUser)
	apiKey := os.Getenv(ApiKey)

	client := slclient.NewSoftLayerClient(apiUser, apiKey)

	virtualGuestService, err := client.GetSoftLayer_Virtual_Guest_Service()

	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("No error found")
	}

	virtualGuest, err := virtualGuestService.GetObject(12345)

	if err != nil {
		fmt.Println("errors:", err)
		return
	} else {
		fmt.Println("No error found")
	}

	spew.Printf("virtualGuest =%#+v\n", virtualGuest)

}
