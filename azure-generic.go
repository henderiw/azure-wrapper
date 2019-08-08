package azurewrapper

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/Azure/go-autorest/autorest"
)

// ClientInfo Information loaded from the authorization file to identify the client
type ClientInfo struct {
	SubscriptionID        string
	ResourceGroupName     string
	ResourceGroupLocation string
}

var (
	ctx = context.Background()
	// ClientData initialization info
	ClientData ClientInfo
	// Authorizer initialization info
	Authorizer autorest.Authorizer
)

// ReadJSON generic fucntion to read azure credentials from file
func ReadJSON(path string) (*map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	contents := make(map[string]interface{})
	json.Unmarshal(data, &contents)
	return &contents, nil
}
