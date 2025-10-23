package provider

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"

const (
	AppConfiguration cloud.ServiceName = "AppConfiguration"
	DeviceUpdate     cloud.ServiceName = "DeviceUpdate"
	DigitalTwins     cloud.ServiceName = "DigitalTwins"
	IoTCentral       cloud.ServiceName = "IoTCentral"
	KeyVault         cloud.ServiceName = "KeyVault"
	Purview          cloud.ServiceName = "Purview"
	Synapse          cloud.ServiceName = "Synapse"
	Search           cloud.ServiceName = "Search"
)

func init() {
	cloud.AzurePublic.Services[AppConfiguration] = cloud.ServiceConfiguration{
		Audience: "https://azconfig.io",
		Endpoint: "https://azconfig.io",
	}
	cloud.AzurePublic.Services[DeviceUpdate] = cloud.ServiceConfiguration{
		Audience: "https://api.adu.microsoft.com",
		Endpoint: "https://api.adu.microsoft.com",
	}
	cloud.AzurePublic.Services[DigitalTwins] = cloud.ServiceConfiguration{
		Audience: "https://digitaltwins.azure.net",
		Endpoint: "https://digitaltwins.azure.net",
	}
	cloud.AzurePublic.Services[IoTCentral] = cloud.ServiceConfiguration{
		Audience: "https://apps.azureiotcentral.com",
		Endpoint: "https://azureiotcentral.com",
	}
	cloud.AzurePublic.Services[KeyVault] = cloud.ServiceConfiguration{
		Audience: "https://vault.azure.net",
		Endpoint: "https://vault.azure.net",
	}
	cloud.AzurePublic.Services[Purview] = cloud.ServiceConfiguration{
		Audience: "https://purview.azure.net",
		Endpoint: "https://purview.azure.com",
	}
	cloud.AzurePublic.Services[Synapse] = cloud.ServiceConfiguration{
		Audience: "https://dev.azuresynapse.net",
		Endpoint: "https://dev.azuresynapse.net",
	}
	cloud.AzurePublic.Services[Search] = cloud.ServiceConfiguration{
		Audience: "https://search.azure.com",
		Endpoint: "https://search.windows.net",
	}

}
