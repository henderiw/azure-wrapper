package azurewrapper

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/go-autorest/autorest/to"
)

//GetGroup wrapper using the azure go SDK
func GetGroup() (group resources.Group, err error) {
	groupsClient := resources.NewGroupsClient(clientData.SubscriptionID)
	groupsClient.Authorizer = authorizer

	return groupsClient.CreateOrUpdate(
		ctx,
		clientData.ResourceGroupName,
		resources.Group{
			Location: to.StringPtr(clientData.ResourceGroupLocation)})
}
