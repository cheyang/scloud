package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	slclient "github.com/maximilien/softlayer-go/client"
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
	}

	virtualGuest, err := virtualGuestService.GetObject(12345)

	if err != nil {
		fmt.Println(err)
		return
	}

	spew.Printf("virtualGuest =%#+v\n", virtualGuest)

}
