package azurewrapper

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
)

// ClientInfo Information loaded from the authorization file to identify the client
type ClientInfo struct {
	SubscriptionID        string
	ResourceGroupName     string
	ResourceGroupLocation string
}

var (
	ctx        = context.Background()
	clientData ClientInfo
	authorizer autorest.Authorizer
)
