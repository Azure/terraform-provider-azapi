package parse_test

import (
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

func Test_NewDataPlaneResourceId(t *testing.T) {
	testData := []struct {
		Name         string
		ParentId     string
		ResourceType string
		Error        bool
		Expected     *parse.DataPlaneResourceId
	}{
		{
			Name:         "test",
			ParentId:     "xxx.xxx.xxx",
			ResourceType: "Microsoft.AppConfiguration/configurationStores/keyValues@api-version",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "xxx.xxx.xxx/kv/test",
				ApiVersion:        "api-version",
				AzureResourceType: "Microsoft.AppConfiguration/configurationStores/keyValues",
			},
		},
		{
			Name:         "test",
			ParentId:     "xxx.xxx.xxx",
			ResourceType: "Microsoft.DeviceUpdate/accounts/v2/groups@8.2",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "xxx.xxx.xxx/v2/management/groups/test",
				ApiVersion:        "8.2",
				AzureResourceType: "Microsoft.DeviceUpdate/accounts/v2/groups",
			},
		},
		{
			Name:         "test",
			ParentId:     "xxx.xxx.xxx",
			ResourceType: "Microsoft.IoTCentral/iotApps/devices/attestation@8.2",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "xxx.xxx.xxx/api/devices/test/attestation",
				ApiVersion:        "8.2",
				AzureResourceType: "Microsoft.IoTCentral/iotApps/devices/attestation",
			},
		},
		{
			Name:         "",
			ParentId:     "foo.keyvault.azure.net",
			ResourceType: "Microsoft.KeyVault/vaults/certificates/contacts@v1.1-preview.2",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "foo.keyvault.azure.net/certificates/contacts",
				ApiVersion:        "v1.1-preview.2",
				AzureResourceType: "Microsoft.KeyVault/vaults/certificates/contacts",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q %q %q", v.Name, v.ParentId, v.ResourceType)

		actual, err := parse.NewDataPlaneResourceId(v.Name, v.ParentId, v.ResourceType)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.AzureResourceId != v.Expected.AzureResourceId {
			t.Fatalf("Expected %q but got %q for AzureResourceId", v.Expected.AzureResourceId, actual.AzureResourceId)
		}
		if actual.ApiVersion != v.Expected.ApiVersion {
			t.Fatalf("Expected %q but got %q for ApiVersion", v.Expected.ApiVersion, actual.ApiVersion)
		}
		if actual.AzureResourceType != v.Expected.AzureResourceType {
			t.Fatalf("Expected %q but got %q for AzureResourceType", v.Expected.AzureResourceType, actual.AzureResourceType)
		}
	}
}

func Test_DataPlaneResourceIDWithResourceType(t *testing.T) {
	testData := []struct {
		ResourceId   string
		ResourceType string
		Error        bool
		Expected     *parse.DataPlaneResourceId
	}{
		{
			ResourceId:   "xxx.xxx.xxx/kv/test",
			ResourceType: "Microsoft.AppConfiguration/configurationStores/keyValues@api-version",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "xxx.xxx.xxx/kv/test",
				ApiVersion:        "api-version",
				AzureResourceType: "Microsoft.AppConfiguration/configurationStores/keyValues",
				ParentId:          "xxx.xxx.xxx",
				Name:              "test",
			},
		},
		{
			ResourceId:   "xxx.xxx.xxx/v2/management/groups/test",
			ResourceType: "Microsoft.DeviceUpdate/accounts/v2/groups@8.2",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "xxx.xxx.xxx/v2/management/groups/test",
				ApiVersion:        "8.2",
				AzureResourceType: "Microsoft.DeviceUpdate/accounts/v2/groups",
				ParentId:          "xxx.xxx.xxx",
				Name:              "test",
			},
		},
		{
			ResourceId:   "xxx.xxx.xxx/api/devices/test/attestation",
			ResourceType: "Microsoft.IoTCentral/iotApps/devices/attestation@8.2",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "xxx.xxx.xxx/api/devices/test/attestation",
				ApiVersion:        "8.2",
				AzureResourceType: "Microsoft.IoTCentral/iotApps/devices/attestation",
				ParentId:          "xxx.xxx.xxx",
				Name:              "test",
			},
		},
		{
			ResourceId:   "foo.keyvault.azure.net/certificates/contacts",
			ResourceType: "Microsoft.KeyVault/vaults/certificates/contacts@v1.1-preview.2",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "foo.keyvault.azure.net/certificates/contacts",
				ApiVersion:        "v1.1-preview.2",
				AzureResourceType: "Microsoft.KeyVault/vaults/certificates/contacts",
				ParentId:          "foo.keyvault.azure.net",
				Name:              "",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q %q", v.ResourceId, v.ResourceType)

		actual, err := parse.DataPlaneResourceIDWithResourceType(v.ResourceId, v.ResourceType)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.AzureResourceId != v.Expected.AzureResourceId {
			t.Fatalf("Expected %q but got %q for AzureResourceId", v.Expected.AzureResourceId, actual.AzureResourceId)
		}
		if actual.ApiVersion != v.Expected.ApiVersion {
			t.Fatalf("Expected %q but got %q for ApiVersion", v.Expected.ApiVersion, actual.ApiVersion)
		}
		if actual.AzureResourceType != v.Expected.AzureResourceType {
			t.Fatalf("Expected %q but got %q for AzureResourceType", v.Expected.AzureResourceType, actual.AzureResourceType)
		}
		if actual.ParentId != v.Expected.ParentId {
			t.Fatalf("Expected %q but got %q for ParentId", v.Expected.ParentId, actual.ParentId)
		}
	}
}
