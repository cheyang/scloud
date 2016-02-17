package main

import (
	"os"

	slclient "github.com/maximilien/softlayer-go/client"
	//	datatypes "github.com/maximilien/softlayer-go/data_types"
)

const (
	ApiEndpoint = "https://api.softlayer.com/rest/v3"
	ApiUser     = "SL_USERNAME"
	ApiKey      = "SL_API_KEY"
)

func main() {

	username := os.Getenv(ApiUser)
	apiKey := os.Getenv(ApiKey)

	client := sclient.NewSoftLayerClient(username, apiKey)

}
