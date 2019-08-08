package azurewrapper

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/go-autorest/autorest/to"
)

//GetGroup wrapper using the azure go SDK
func GetGroup() (group resources.Group, err error) {
	groupsClient := resources.NewGroupsClient(ClientData.SubscriptionID)
	groupsClient.Authorizer = Authorizer

	return groupsClient.CreateOrUpdate(
		ctx,
		ClientData.ResourceGroupName,
		resources.Group{
			Location: to.StringPtr(ClientData.ResourceGroupLocation)})
}
