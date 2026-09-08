package provider

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
)

func TestStorageCloudConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		config   cloud.Configuration
		endpoint string
		audience string
	}{
		{"public", cloud.AzurePublic, "https://.blob.core.windows.net", "https://storage.azure.com"},
		{"government", cloud.AzureGovernment, "https://.blob.core.usgovcloudapi.net", "https://storage.azure.us"},
		{"china", cloud.AzureChina, "https://.blob.core.chinacloudapi.cn", "https://storage.azure.cn"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := test.config.Services[Storage]
			if service.Endpoint != test.endpoint || service.Audience != test.audience {
				t.Fatalf("unexpected storage service configuration %#v", service)
			}
		})
	}
}
