package azurewrapper

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/to"
)

//AzureVWanCfg Configuration structure of the download file from Azure VHUB
type AzureVWanCfg struct {
	ConfigurationVersion struct {
		LastUpdatedTime time.Time `json:"LastUpdatedTime"`
		Version         string    `json:"Version"`
	} `json:"configurationVersion"`
	VpnSiteConfiguration struct {
		Name      string `json:"Name"`
		IPAddress string `json:"IPAddress"`
	} `json:"vpnSiteConfiguration"`
	VpnSiteConnections []struct {
		HubConfiguration struct {
			AddressSpace     string   `json:"AddressSpace"`
			ConnectedSubnets []string `json:"ConnectedSubnets"`
		} `json:"hubConfiguration"`
		GatewayConfiguration struct {
			IPAddresses struct {
				Instance0 string `json:"Instance0"`
				Instance1 string `json:"Instance1"`
			} `json:"IpAddresses"`
		} `json:"gatewayConfiguration"`
		ConnectionConfiguration struct {
			IsBgpEnabled    bool   `json:"IsBgpEnabled"`
			PSK             string `json:"PSK"`
			IPsecParameters struct {
				SADataSizeInKilobytes int `json:"SADataSizeInKilobytes"`
				SALifeTimeInSeconds   int `json:"SALifeTimeInSeconds"`
			} `json:"IPsecParameters"`
		} `json:"connectionConfiguration"`
	} `json:"vpnSiteConnections"`
}

// GetVirtualWansClient wrapper using the azure go SDK
func GetVirtualWansClient() network.VirtualWansClient {
	log.Println("Get virtual wan client data")
	client := network.NewVirtualWansClient(ClientData.SubscriptionID)
	client.Authorizer = Authorizer
	return client
}

// GetVwan wrapper using the azure go SDK
func GetVwan(vwanName string) (vwan network.VirtualWAN, err error) {
	log.Println("Get virtual wan data")
	vwanClient := GetVirtualWansClient()
	return vwanClient.Get(ctx, ClientData.ResourceGroupName, vwanName)
}

// CreateVwan wrapper using the azure go SDK
func CreateVwan(vwanName string) (vwan network.VirtualWAN, err error) {
	log.Println("Create or update virtual wan data")
	vwanClient := GetVirtualWansClient()

	future, err := vwanClient.CreateOrUpdate(
		ctx,
		ClientData.ResourceGroupName,
		vwanName,
		network.VirtualWAN{
			Location: to.StringPtr(ClientData.ResourceGroupLocation),
			VirtualWanProperties: &network.VirtualWanProperties{
				AllowBranchToBranchTraffic: to.BoolPtr(true),
			},
		},
	)

	if err != nil {
		return vwan, fmt.Errorf("cannot create virtual wan: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, vwanClient.Client)
	if err != nil {
		return vwan, fmt.Errorf("cannot get the vwan create or update future response: %v", err)
	}

	return future.Result(vwanClient)
}

// DeleteVwan wrapper using the azure go SDK
func DeleteVwan(name string) (err error) {
	client := GetVirtualWansClient()

	future, err := client.Delete(ctx, ClientData.ResourceGroupName, name)
	if err != nil {
		return fmt.Errorf("cannot delete vwan: %v", err)

	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("cannot delete vwan - response: %v", err)
	}

	log.Printf("Delete vwan Response: %#v\n", future.Response())

	return nil
}

// GetVirtualHubsClient wrapper using the azure go SDK
func GetVirtualHubsClient() network.VirtualHubsClient {
	client := network.NewVirtualHubsClient(ClientData.SubscriptionID)
	client.Authorizer = Authorizer
	return client
}

// GetVhub wrapper using the azure go SDK
func GetVhub(vhubName string) (vwan network.VirtualHub, err error) {
	vhubClient := GetVirtualHubsClient()
	return vhubClient.Get(ctx, ClientData.ResourceGroupName, vhubName)
}

// CreateVhub wrapper using the azure go SDK
func CreateVhub(vhubName, vwanID, ip, location string) (vwan network.VirtualHub, err error) {
	client := GetVirtualHubsClient()

	future, err := client.CreateOrUpdate(
		ctx,
		ClientData.ResourceGroupName,
		vhubName,
		network.VirtualHub{
			Location: to.StringPtr(location),
			VirtualHubProperties: &network.VirtualHubProperties{
				AddressPrefix: to.StringPtr(ip),
				VirtualWan: &network.SubResource{
					ID: to.StringPtr(vwanID),
				},
			},
		},
	)

	if err != nil {
		return vwan, fmt.Errorf("cannot create virtual hub: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return vwan, fmt.Errorf("cannot get the vhub create or update future response: %v", err)
	}

	return future.Result(client)
}

// DeleteVhub wrapper using the azure go SDK
func DeleteVhub(name string) (err error) {
	client := GetVirtualHubsClient()

	future, err := client.Delete(ctx, ClientData.ResourceGroupName, name)
	if err != nil {
		return fmt.Errorf("cannot delete vhub: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("cannot delete vhub - response: %v", err)
	}

	fmt.Printf("Delete vhub Response: %#v\n", future.Response())

	return nil
}

// GetVpnSitesClient wrapper using the azure go SDK
func GetVpnSitesClient() network.VpnSitesClient {
	client := network.NewVpnSitesClient(ClientData.SubscriptionID)
	client.Authorizer = Authorizer
	return client
}

// GetVpnSite wrapper using the azure go SDK
func GetVpnSite(name string) (vpnsite network.VpnSite, err error) {
	client := GetVpnSitesClient()
	return client.Get(ctx, ClientData.ResourceGroupName, name)
}

// CreateVpnSite wrapper using the azure go SDK
func CreateVpnSite(name, vwanID, location string, nsgConf NsgConfYML) (vpnSite network.VpnSite, err error) {
	client := GetVpnSitesClient()

	var future network.VpnSitesCreateOrUpdateFuture

	if nsgConf.NsgData.BgpEnabled == false {
		future, err = client.CreateOrUpdate(
			ctx,
			ClientData.ResourceGroupName,
			name,
			network.VpnSite{
				Location: to.StringPtr(location),
				VpnSiteProperties: &network.VpnSiteProperties{
					VirtualWan: &network.SubResource{
						ID: to.StringPtr(vwanID),
					},
					DeviceProperties: &network.DeviceProperties{
						DeviceVendor:    to.StringPtr("Nuage Networks"),
						DeviceModel:     to.StringPtr("E306"),
						LinkSpeedInMbps: to.Int32Ptr(100),
					},
					IPAddress: to.StringPtr(nsgConf.NsgData.PublicIP),
					AddressSpace: &network.AddressSpace{
						AddressPrefixes: to.StringSlicePtr(nsgConf.NsgData.LanSubnet),
					},
				},
			},
		)
	} else {
		future, err = client.CreateOrUpdate(
			ctx,
			ClientData.ResourceGroupName,
			name,
			network.VpnSite{
				Location: to.StringPtr(location),
				VpnSiteProperties: &network.VpnSiteProperties{
					VirtualWan: &network.SubResource{
						ID: to.StringPtr(vwanID),
					},
					DeviceProperties: &network.DeviceProperties{
						DeviceVendor:    to.StringPtr("Nuage Networks"),
						DeviceModel:     to.StringPtr("E306"),
						LinkSpeedInMbps: to.Int32Ptr(100),
					},
					IPAddress: to.StringPtr(nsgConf.NsgData.PublicIP),
					AddressSpace: &network.AddressSpace{
						AddressPrefixes: to.StringSlicePtr(nsgConf.NsgData.LanSubnet),
					},
					BgpProperties: &network.BgpSettings{
						BgpPeeringAddress: to.StringPtr("10.0.0.10"),
						Asn:               to.Int64Ptr(nsgConf.NsgData.BgpNsgAsn),
					},
				},
			},
		)
	}

	if err != nil {
		return vpnSite, fmt.Errorf("cannot create vpnSite: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return vpnSite, fmt.Errorf("cannot get the vpnSite create or update future response: %v", err)
	}

	return future.Result(client)
}

// DeleteVpnSite wrapper using the azure go SDK
func DeleteVpnSite(name string) (err error) {
	client := GetVpnSitesClient()

	future, err := client.Delete(ctx, ClientData.ResourceGroupName, name)
	if err != nil {
		return fmt.Errorf("cannot delete vpnSite: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("cannot delete vpnSite - response: %v", err)
	}

	fmt.Printf("Delete vpnSite Response: %#v\n", future.Response())

	return nil
}

// GetVpnGatewaysClient wrapper using the azure go SDK
func GetVpnGatewaysClient() network.VpnGatewaysClient {
	client := network.NewVpnGatewaysClient(ClientData.SubscriptionID)
	client.Authorizer = Authorizer
	return client
}

// GetVpnGateway wrapper using the azure go SDK
func GetVpnGateway(name string) (vpngw network.VpnGateway, err error) {
	client := GetVpnGatewaysClient()
	return client.Get(ctx, ClientData.ResourceGroupName, name)
}

// CreateVpnGateway wrapper using the azure go SDK
func CreateVpnGateway(name, vhubID, vsiteID, location string, nsgConf NsgConfYML) (vpngw network.VpnGateway, err error) {
	client := GetVpnGatewaysClient()

	var future network.VpnGatewaysCreateOrUpdateFuture

	if vsiteID != "" {
		future, err = client.CreateOrUpdate(
			ctx,
			ClientData.ResourceGroupName,
			name,
			network.VpnGateway{
				Location: to.StringPtr(location),
				VpnGatewayProperties: &network.VpnGatewayProperties{
					VirtualHub: &network.SubResource{
						ID: to.StringPtr(vhubID),
					},
					BgpSettings: &network.BgpSettings{
						Asn: to.Int64Ptr(65515),
					},
					Connections: &[]network.VpnConnection{
						{
							Name: to.StringPtr(nsgConf.NsgData.NsgName),
							VpnConnectionProperties: &network.VpnConnectionProperties{
								SharedKey: to.StringPtr("testje"),
								RemoteVpnSite: &network.SubResource{
									ID: to.StringPtr(vsiteID),
								},
							},
						},
					},
				},
			},
		)
	} else {
		future, err = client.CreateOrUpdate(
			ctx,
			ClientData.ResourceGroupName,
			name,
			network.VpnGateway{
				Location: to.StringPtr(location),
				VpnGatewayProperties: &network.VpnGatewayProperties{
					VirtualHub: &network.SubResource{
						ID: to.StringPtr(vhubID),
					},
					BgpSettings: &network.BgpSettings{
						Asn: to.Int64Ptr(65515),
					},
				},
			},
		)
	}

	if err != nil {
		return vpngw, fmt.Errorf("cannot create vpngw: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return vpngw, fmt.Errorf("cannot get the vpngw create or update future response: %v", err)
	}

	return future.Result(client)
}

// UpdateVpnGateway wrapper using the azure go SDK
func UpdateVpnGateway(name, vhubID, location string, conn []network.VpnConnection) (vpngw network.VpnGateway, err error) {
	client := GetVpnGatewaysClient()

	future, err := client.CreateOrUpdate(
		ctx,
		ClientData.ResourceGroupName,
		name,
		network.VpnGateway{
			Location: to.StringPtr(location),
			VpnGatewayProperties: &network.VpnGatewayProperties{
				VirtualHub: &network.SubResource{
					ID: to.StringPtr(vhubID),
				},
				BgpSettings: &network.BgpSettings{
					Asn: to.Int64Ptr(65515),
				},
				Connections: &conn,
			},
		},
	)

	if err != nil {
		return vpngw, fmt.Errorf("cannot update vpngw: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return vpngw, fmt.Errorf("cannot get the vpngw create or update future response: %v", err)
	}

	return future.Result(client)
}

// DeleteVpnGateway wrapper using the azure go SDK
func DeleteVpnGateway(name string) (err error) {
	client := GetVpnGatewaysClient()

	future, err := client.Delete(ctx, ClientData.ResourceGroupName, name)
	if err != nil {
		return fmt.Errorf("cannot delete vpnGw: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("cannot delete vpnGw - response: %v", err)
	}

	fmt.Printf("Delete vpnGw Response: %#v\n", future.Response())

	return nil
}

// GetVpnSitesConfigurationClient wrapper using the azure go SDK
func GetVpnSitesConfigurationClient() network.VpnSitesConfigurationClient {
	client := network.NewVpnSitesConfigurationClient(ClientData.SubscriptionID)
	client.Authorizer = Authorizer
	return client
}

// DownloadVpnSitesConfig wrapper using the azure go SDK
func DownloadVpnSitesConfig(vwanName, vsiteID, url string) (err error) {
	client := GetVpnSitesConfigurationClient()

	future, err := client.Download(
		ctx,
		ClientData.ResourceGroupName,
		vwanName,
		network.GetVpnSitesConfigurationRequest{
			OutputBlobSasURL: to.StringPtr(url),
			VpnSites:         &[]string{vsiteID},
		},
	)

	if err != nil {
		return fmt.Errorf("cannot create vpngw: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("cannot get the vpngw create or update future response: %v", err)
	}

	fmt.Printf("Download VPN Site Config Response: %#v\n", future.Response())

	return nil

}
