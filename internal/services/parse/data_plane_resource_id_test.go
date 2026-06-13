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
		{
			Name:         "myindex",
			ParentId:     "mysearchservice.search.windows.net",
			ResourceType: "Microsoft.Search/searchServices/indexes@2024-07-01",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "mysearchservice.search.windows.net/indexes('myindex')",
				ApiVersion:        "2024-07-01",
				AzureResourceType: "Microsoft.Search/searchServices/indexes",
			},
		},
		{
			Name:         "mydatasource",
			ParentId:     "mysearchservice.search.windows.net",
			ResourceType: "Microsoft.Search/searchServices/datasources@2024-07-01",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "mysearchservice.search.windows.net/datasources('mydatasource')",
				ApiVersion:        "2024-07-01",
				AzureResourceType: "Microsoft.Search/searchServices/datasources",
			},
		},
		{
			Name:         "myindexer",
			ParentId:     "mysearchservice.search.windows.net",
			ResourceType: "Microsoft.Search/searchServices/indexers@2024-07-01",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "mysearchservice.search.windows.net/indexers('myindexer')",
				ApiVersion:        "2024-07-01",
				AzureResourceType: "Microsoft.Search/searchServices/indexers",
			},
		},
		{
			Name:         "myskillset",
			ParentId:     "mysearchservice.search.windows.net",
			ResourceType: "Microsoft.Search/searchServices/skillsets@2024-07-01",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "mysearchservice.search.windows.net/skillsets('myskillset')",
				ApiVersion:        "2024-07-01",
				AzureResourceType: "Microsoft.Search/searchServices/skillsets",
			},
		},
		{
			Name:         "mysynonymmap",
			ParentId:     "mysearchservice.search.windows.net",
			ResourceType: "Microsoft.Search/searchServices/synonymmaps@2024-07-01",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "mysearchservice.search.windows.net/synonymmaps('mysynonymmap')",
				ApiVersion:        "2024-07-01",
				AzureResourceType: "Microsoft.Search/searchServices/synonymmaps",
			},
		},
		{
			Name:         "acctesttable",
			ParentId:     "mystorage.table.core.windows.net",
			ResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables@2026-04-06",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "mystorage.table.core.windows.net/Tables('acctesttable')",
				ApiVersion:        "2026-04-06",
				AzureResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables",
			},
		},
		{
			Name:         "",
			ParentId:     "mystorage.table.core.windows.net/mytable",
			ResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables/entities@2026-04-06",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "mystorage.table.core.windows.net/mytable(PartitionKey='pk',RowKey='rk')",
				ApiVersion:        "2026-04-06",
				AzureResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables/entities",
				Identifiers: map[string]string{
					"partitionKey": "pk",
					"rowKey":       "rk",
				},
			},
		},
		// {tokenId} — non-name placeholder; previously broken (silently left as literal text)
		{
			Name:         "",
			ParentId:     "myapp.azureiotcentral.com",
			ResourceType: "Microsoft.IoTCentral/iotApps/apiTokens@2022-07-31",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "myapp.azureiotcentral.com/api/apiTokens/mytoken",
				ApiVersion:        "2022-07-31",
				AzureResourceType: "Microsoft.IoTCentral/iotApps/apiTokens",
				Identifiers: map[string]string{
					"tokenId": "mytoken",
				},
			},
		},
		// {name=defaultValue} — singleton with a fixed name enforced by the service
		{
			Name:         "defaultResourceSetRuleConfig",
			ParentId:     "mypurview.purview.azure.com",
			ResourceType: "Microsoft.Purview/accounts/Account/resourceSetRuleConfigs@2021-07-01",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "mypurview.purview.azure.com/resourceSetRuleConfigs/defaultResourceSetRuleConfig",
				ApiVersion:        "2021-07-01",
				AzureResourceType: "Microsoft.Purview/accounts/Account/resourceSetRuleConfigs",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q %q %q", v.Name, v.ParentId, v.ResourceType)

		actual, err := parse.NewDataPlaneResourceIdWithIdentifiers(v.Name, v.ParentId, v.ResourceType, v.Expected.Identifiers)
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
		if len(actual.Identifiers) != len(v.Expected.Identifiers) {
			t.Fatalf("Expected identifiers %#v but got %#v", v.Expected.Identifiers, actual.Identifiers)
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
		{
			ResourceId:   "mysearchservice.search.windows.net/indexes('myindex')",
			ResourceType: "Microsoft.Search/searchServices/indexes@2024-07-01",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "mysearchservice.search.windows.net/indexes('myindex')",
				ApiVersion:        "2024-07-01",
				AzureResourceType: "Microsoft.Search/searchServices/indexes",
				ParentId:          "mysearchservice.search.windows.net",
				Name:              "myindex",
			},
		},
		{
			ResourceId:   "mysearchservice.search.windows.net/datasources('mydatasource')",
			ResourceType: "Microsoft.Search/searchServices/datasources@2024-07-01",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "mysearchservice.search.windows.net/datasources('mydatasource')",
				ApiVersion:        "2024-07-01",
				AzureResourceType: "Microsoft.Search/searchServices/datasources",
				ParentId:          "mysearchservice.search.windows.net",
				Name:              "mydatasource",
			},
		},
		{
			ResourceId:   "mysearchservice.search.windows.net/indexers('myindexer')",
			ResourceType: "Microsoft.Search/searchServices/indexers@2024-07-01",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "mysearchservice.search.windows.net/indexers('myindexer')",
				ApiVersion:        "2024-07-01",
				AzureResourceType: "Microsoft.Search/searchServices/indexers",
				ParentId:          "mysearchservice.search.windows.net",
				Name:              "myindexer",
			},
		},
		{
			ResourceId:   "mysearchservice.search.windows.net/skillsets('myskillset')",
			ResourceType: "Microsoft.Search/searchServices/skillsets@2024-07-01",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "mysearchservice.search.windows.net/skillsets('myskillset')",
				ApiVersion:        "2024-07-01",
				AzureResourceType: "Microsoft.Search/searchServices/skillsets",
				ParentId:          "mysearchservice.search.windows.net",
				Name:              "myskillset",
			},
		},
		{
			ResourceId:   "mysearchservice.search.windows.net/synonymmaps('mysynonymmap')",
			ResourceType: "Microsoft.Search/searchServices/synonymmaps@2024-07-01",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "mysearchservice.search.windows.net/synonymmaps('mysynonymmap')",
				ApiVersion:        "2024-07-01",
				AzureResourceType: "Microsoft.Search/searchServices/synonymmaps",
				ParentId:          "mysearchservice.search.windows.net",
				Name:              "mysynonymmap",
			},
		},
		{
			ResourceId:   "mystorage.table.core.windows.net/Tables('acctesttable')",
			ResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables@2026-04-06",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "mystorage.table.core.windows.net/Tables('acctesttable')",
				ApiVersion:        "2026-04-06",
				AzureResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables",
				ParentId:          "mystorage.table.core.windows.net",
				Name:              "acctesttable",
			},
		},
		{
			ResourceId:   "mystorage.table.core.windows.net/mytable(PartitionKey='pk',RowKey='rk')",
			ResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables/entities@2026-04-06",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "mystorage.table.core.windows.net/mytable(PartitionKey='pk',RowKey='rk')",
				ApiVersion:        "2026-04-06",
				AzureResourceType: "Microsoft.Storage/storageAccounts/tableServices/tables/entities",
				ParentId:          "mystorage.table.core.windows.net/mytable",
				Identifiers: map[string]string{
					"partitionKey": "pk",
					"rowKey":       "rk",
				},
			},
		},
		// {tokenId} round-trip — verify the non-name placeholder is parsed correctly
		{
			ResourceId:   "myapp.azureiotcentral.com/api/apiTokens/mytoken",
			ResourceType: "Microsoft.IoTCentral/iotApps/apiTokens@2022-07-31",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "myapp.azureiotcentral.com/api/apiTokens/mytoken",
				ApiVersion:        "2022-07-31",
				AzureResourceType: "Microsoft.IoTCentral/iotApps/apiTokens",
				ParentId:          "myapp.azureiotcentral.com",
				Name:              "",
				Identifiers: map[string]string{
					"tokenId": "mytoken",
				},
			},
		},
		// {name=defaultValue} round-trip — verify fixed-default singletons parse correctly
		{
			ResourceId:   "mypurview.purview.azure.com/resourceSetRuleConfigs/defaultResourceSetRuleConfig",
			ResourceType: "Microsoft.Purview/accounts/Account/resourceSetRuleConfigs@2021-07-01",
			Error:        false,
			Expected: &parse.DataPlaneResourceId{
				AzureResourceId:   "mypurview.purview.azure.com/resourceSetRuleConfigs/defaultResourceSetRuleConfig",
				ApiVersion:        "2021-07-01",
				AzureResourceType: "Microsoft.Purview/accounts/Account/resourceSetRuleConfigs",
				ParentId:          "mypurview.purview.azure.com",
				Name:              "defaultResourceSetRuleConfig",
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
		if len(actual.Identifiers) != len(v.Expected.Identifiers) {
			t.Fatalf("Expected identifiers %#v but got %#v", v.Expected.Identifiers, actual.Identifiers)
		}
	}
}
