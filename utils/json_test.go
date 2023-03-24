package utils_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Azure/terraform-provider-azapi/utils"
)

func Test_GetUpdatedJson(t *testing.T) {
	testcases := []struct {
		OldJson    string
		NewJson    string
		ExpectJson string
		Option     utils.UpdateJsonOption
	}{
		{
			OldJson: `
{
  "properties": {
    "compression": "None",
    "consumerGroup": "acctesteventhubcg",
    "dataFormat": "",
    "eventHubResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG924/providers/Microsoft.EventHub/namespaces/acctesteventhubnamespace/eventhubs/acctesteventhub",
    "eventSystemProperties": [],
    "managedIdentityResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG924/providers/Microsoft.ManagedIdentity/userAssignedIdentities/acctest924",
    "mappingRuleName": "",
    "provisioningState": "Succeeded",
    "tableName": ""
  }
}`,
			NewJson: `
{
  "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG924/providers/Microsoft.Kusto/Clusters/acctestkc924/Databases/acctestkd/DataConnections/acctestkedc",
  "kind": "EventHub",
  "location": "West Europe",
  "name": "acctestkc924/acctestkd/acctestkedc",
  "properties": {
    "compression": "None",
    "dataFormat": "",
    "eventHubResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/acctestRG924/providers/Microsoft.EventHub/namespaces/acctesteventhubnamespace/eventhubs/acctesteventhub",
    "eventSystemProperties": [],
    "managedIdentityResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG924/providers/Microsoft.ManagedIdentity/userAssignedIdentities/acctest924",
    "mappingRuleName": "",
    "provisioningState": "Succeeded",
    "tableName": ""
  },
  "type": "Microsoft.Kusto/Clusters/Databases/DataConnections"
}
`,
			ExpectJson: `
{
  "properties": {
    "compression": "None",
    "consumerGroup": "acctesteventhubcg",
    "dataFormat": "",
    "eventHubResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG924/providers/Microsoft.EventHub/namespaces/acctesteventhubnamespace/eventhubs/acctesteventhub",
    "eventSystemProperties": [],
    "managedIdentityResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG924/providers/Microsoft.ManagedIdentity/userAssignedIdentities/acctest924",
    "mappingRuleName": "",
    "provisioningState": "Succeeded",
    "tableName": ""
  }
}
`,
			Option: utils.UpdateJsonOption{
				IgnoreMissingProperty: true,
				IgnoreCasing:          true,
			},
		},
		{
			OldJson: `
{
  "properties": {
    "sshPrivateKey": "asdf"
  }
}`,
			NewJson: `
{
  "properties": {
    "sshPrivateKey": "<redacted>"
  }
}
`,
			ExpectJson: `
{
  "properties": {
    "sshPrivateKey": "asdf"
  }
}
`,
			Option: utils.UpdateJsonOption{
				IgnoreMissingProperty: true,
				IgnoreCasing:          true,
			},
		},
	}

	for _, testcase := range testcases {
		var new, old, expected interface{}
		_ = json.Unmarshal([]byte(testcase.OldJson), &old)
		_ = json.Unmarshal([]byte(testcase.NewJson), &new)
		_ = json.Unmarshal([]byte(testcase.ExpectJson), &expected)

		result := utils.GetUpdatedJson(old, new, testcase.Option)
		if !reflect.DeepEqual(result, expected) {
			expectedJson, _ := json.Marshal(expected)
			resultJson, _ := json.Marshal(result)
			t.Fatalf("Expected %s but got %s", expectedJson, resultJson)
		}
	}
}

func Test_GetMergedJson(t *testing.T) {
	oldJson := `
 {
	"a":1,
    "b": {
		"b1": "b1",
		"b2": []
	}
}
`

	newJson := `
{
	"b": {
		"b3": "b3"
	}
}
`
	expectedJson := `
{
	"a":1,
    "b": {
		"b1": "b1",
		"b2": [],
		"b3": "b3"
	}
}
`
	var new, old, expected interface{}
	_ = json.Unmarshal([]byte(oldJson), &old)
	_ = json.Unmarshal([]byte(newJson), &new)
	_ = json.Unmarshal([]byte(expectedJson), &expected)

	result := utils.GetMergedJson(old, new)
	if !reflect.DeepEqual(result, expected) {
		expectedJson, _ := json.Marshal(expected)
		resultJson, _ := json.Marshal(result)
		t.Fatalf("Expected %s but got %s", expectedJson, resultJson)
	}
}

func Test_GetMergedJsonWithArray(t *testing.T) {
	oldJson := `
{
    "name": "mylb",
    "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/rg/providers/Microsoft.Network/loadBalancers/mylb",
    "etag": "W/\"0074bc84-d75b-4496-96c8-c858e0ef0afe\"",
    "type": "Microsoft.Network/loadBalancers",
    "location": "westus",
    "tags": {},
    "properties": {
        "provisioningState": "Succeeded",
        "resourceGuid": "0a6d40b0-4f2e-4582-9bb8-cda0e3f8099a",
        "frontendIPConfigurations": [
            {
                "name": "PublicIPAddress",
                "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/rg/providers/Microsoft.Network/loadBalancers/mylb/frontendIPConfigurations/PublicIPAddress",
                "etag": "W/\"0074bc84-d75b-4496-96c8-c858e0ef0afe\"",
                "type": "Microsoft.Network/loadBalancers/frontendIPConfigurations",
                "properties": {
                    "provisioningState": "Succeeded",
                    "privateIPAllocationMethod": "Dynamic",
                    "publicIPAddress": {
                        "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/rg/providers/Microsoft.Network/publicIPAddresses/myip"
                    },
                    "inboundNatRules": [
                        {
                            "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/rg/providers/Microsoft.Network/loadBalancers/mylb/inboundNatRules/RDPAccess"
                        }
                    ]
                }
            }
        ],
        "backendAddressPools": [],
        "loadBalancingRules": [],
        "probes": [],
        "inboundNatRules": [
            {
                "name": "RDPAccess",
                "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/rg/providers/Microsoft.Network/loadBalancers/mylb/inboundNatRules/RDPAccess",
                "etag": "W/\"0074bc84-d75b-4496-96c8-c858e0ef0afe\"",
                "type": "Microsoft.Network/loadBalancers/inboundNatRules",
                "properties": {
                    "provisioningState": "Succeeded",
                    "frontendIPConfiguration": {
                        "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/rg/providers/Microsoft.Network/loadBalancers/mylb/frontendIPConfigurations/PublicIPAddress"
                    },
                    "frontendPort": 3389,
                    "backendPort": 3389,
                    "enableFloatingIP": false,
                    "protocol": "Tcp",
                    "enableDestinationServiceEndpoint": false,
                    "enableTcpReset": false,
                    "allowBackendPortConflict": false
                }
            }
        ],
        "inboundNatPools": []
    },
    "sku": {
        "name": "Basic",
        "tier": "Regional"
    }
}
`

	newJson := `
  {
      "properties": {
        "inboundNatRules": [
          {
            "properties": {
               "idleTimeoutInMinutes": 15
            }
          }
        ]
      }
    }
`

	expectedJson := `
{
    "name": "mylb",
    "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/rg/providers/Microsoft.Network/loadBalancers/mylb",
    "etag": "W/\"0074bc84-d75b-4496-96c8-c858e0ef0afe\"",
    "type": "Microsoft.Network/loadBalancers",
    "location": "westus",
    "tags": {},
    "properties": {
        "provisioningState": "Succeeded",
        "resourceGuid": "0a6d40b0-4f2e-4582-9bb8-cda0e3f8099a",
        "frontendIPConfigurations": [
            {
                "name": "PublicIPAddress",
                "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/rg/providers/Microsoft.Network/loadBalancers/mylb/frontendIPConfigurations/PublicIPAddress",
                "etag": "W/\"0074bc84-d75b-4496-96c8-c858e0ef0afe\"",
                "type": "Microsoft.Network/loadBalancers/frontendIPConfigurations",
                "properties": {
                    "provisioningState": "Succeeded",
                    "privateIPAllocationMethod": "Dynamic",
                    "publicIPAddress": {
                        "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/rg/providers/Microsoft.Network/publicIPAddresses/myip"
                    },
                    "inboundNatRules": [
                        {
                            "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/rg/providers/Microsoft.Network/loadBalancers/mylb/inboundNatRules/RDPAccess"
                        }
                    ]
                }
            }
        ],
        "backendAddressPools": [],
        "loadBalancingRules": [],
        "probes": [],
        "inboundNatRules": [
            {
                "name": "RDPAccess",
                "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/rg/providers/Microsoft.Network/loadBalancers/mylb/inboundNatRules/RDPAccess",
                "etag": "W/\"0074bc84-d75b-4496-96c8-c858e0ef0afe\"",
                "type": "Microsoft.Network/loadBalancers/inboundNatRules",
                "properties": {
                    "provisioningState": "Succeeded",
                    "frontendIPConfiguration": {
                        "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/rg/providers/Microsoft.Network/loadBalancers/mylb/frontendIPConfigurations/PublicIPAddress"
                    },
                    "frontendPort": 3389,
                    "backendPort": 3389,
                    "enableFloatingIP": false,
                    "idleTimeoutInMinutes": 15,
                    "protocol": "Tcp",
                    "enableDestinationServiceEndpoint": false,
                    "enableTcpReset": false,
                    "allowBackendPortConflict": false
                }
            }
        ],
        "inboundNatPools": []
    },
    "sku": {
        "name": "Basic",
        "tier": "Regional"
    }
}
`
	var new, old, expected interface{}
	_ = json.Unmarshal([]byte(oldJson), &old)
	_ = json.Unmarshal([]byte(newJson), &new)
	_ = json.Unmarshal([]byte(expectedJson), &expected)

	result := utils.GetMergedJson(old, new)
	if !reflect.DeepEqual(result, expected) {
		expectedJson, _ := json.Marshal(expected)
		resultJson, _ := json.Marshal(result)
		t.Fatalf("Expected %s but got %s", expectedJson, resultJson)
	}
}

func Test_GetRemovedJson(t *testing.T) {
	oldJson := `
{
    "properties": {
        "public": false,
        "provisioningState": "Succeeded",
        "fqdn": "springcloud.azuremicroservices.io",
        "httpsOnly": false,
        "createdTime": "2021-11-16T08:49:54.966Z",
        "temporaryDisk": {
            "sizeInGB": 4,
            "mountPath": "/temp"
        },
        "persistentDisk": {
            "sizeInGB": 0,
            "mountPath": "/persistent"
        },
        "enableEndToEndTLS": false
    },
    "type": "Microsoft.AppPlatform/Spring/apps",
    "identity": {
        "type": "SystemAssigned",
        "principalId": "d44e42c2-173f-456b-883a-7433aa870a18",
        "tenantId": "72f988bf-86f1-41af-91ab-2d7cd011db47"
    },
    "location": "westeurope",
    "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/acctestRG1116/providers/Microsoft.AppPlatform/Spring/springcloud/apps/springcloudapp",
    "name": "springcloudapp"
}
`

	newJson := `
  {
    "properties": {
      "temporaryDisk": {
        "mountPath": "/temp",
        "sizeInGB": 4
      }
    }
  }
`
	expectedJson := `
{
    "properties": {
        "public": false,
        "provisioningState": "Succeeded",
        "fqdn": "springcloud.azuremicroservices.io",
        "httpsOnly": false,
        "createdTime": "2021-11-16T08:49:54.966Z",
        "persistentDisk": {
            "sizeInGB": 0,
            "mountPath": "/persistent"
        },
        "enableEndToEndTLS": false
    },
    "type": "Microsoft.AppPlatform/Spring/apps",
    "identity": {
        "type": "SystemAssigned",
        "principalId": "d44e42c2-173f-456b-883a-7433aa870a18",
        "tenantId": "72f988bf-86f1-41af-91ab-2d7cd011db47"
    },
    "location": "westeurope",
    "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/acctestRG1116/providers/Microsoft.AppPlatform/Spring/springcloud/apps/springcloudapp",
    "name": "springcloudapp"
}
`
	var new, old, expected interface{}
	_ = json.Unmarshal([]byte(oldJson), &old)
	_ = json.Unmarshal([]byte(newJson), &new)
	_ = json.Unmarshal([]byte(expectedJson), &expected)

	result := utils.GetRemovedJson(old, new)
	if !reflect.DeepEqual(result, expected) {
		expectedJson, _ := json.Marshal(expected)
		resultJson, _ := json.Marshal(result)
		t.Fatalf("Expected %s but got %s", expectedJson, resultJson)
	}
}

func Test_GetIgnoredJson(t *testing.T) {
	oldJson := `
{
  "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG924/providers/Microsoft.Kusto/Clusters/acctestkc924/Databases/acctestkd/DataConnections/acctestkedc",
  "kind": "EventHub",
  "location": "West Europe",
  "name": "acctestkc924/acctestkd/acctestkedc",
  "properties": {
    "compression": "None",
    "consumerGroup": "acctesteventhubcg",
    "dataFormat": "",
    "eventHubResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG924/providers/Microsoft.EventHub/namespaces/acctesteventhubnamespace/eventhubs/acctesteventhub",
    "eventSystemProperties": [],
    "managedIdentityResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG924/providers/Microsoft.ManagedIdentity/userAssignedIdentities/acctest924",
    "mappingRuleName": "",
    "provisioningState": "Succeeded",
    "tableName": ""
  },
  "type": "Microsoft.Kusto/Clusters/Databases/DataConnections"
}
`
	expectedJson := `
{
  "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG924/providers/Microsoft.Kusto/Clusters/acctestkc924/Databases/acctestkd/DataConnections/acctestkedc",
  "kind": "EventHub",
  "location": "West Europe",
  "name": "acctestkc924/acctestkd/acctestkedc",
  "properties": {
    "compression": "None",
    "consumerGroup": "acctesteventhubcg",
    "dataFormat": "",
    "eventHubResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG924/providers/Microsoft.EventHub/namespaces/acctesteventhubnamespace/eventhubs/acctesteventhub",
    "eventSystemProperties": [],
    "managedIdentityResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG924/providers/Microsoft.ManagedIdentity/userAssignedIdentities/acctest924",
    "mappingRuleName": "",
    "tableName": ""
  },
  "type": "Microsoft.Kusto/Clusters/Databases/DataConnections"
}
`

	var old, expected interface{}
	_ = json.Unmarshal([]byte(oldJson), &old)
	_ = json.Unmarshal([]byte(expectedJson), &expected)

	result := utils.GetIgnoredJson(old, []string{"provisioningState"})
	if !reflect.DeepEqual(result, expected) {
		expectedJson, _ := json.Marshal(expected)
		resultJson, _ := json.Marshal(result)
		t.Fatalf("Expected %s but got %s", expectedJson, resultJson)
	}
}

func Test_ExtractObject(t *testing.T) {
	oldJson := `
{
  "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG924/providers/Microsoft.Kusto/Clusters/acctestkc924/Databases/acctestkd/DataConnections/acctestkedc",
  "kind": "EventHub",
  "location": "West Europe",
  "name": "acctestkc924/acctestkd/acctestkedc",
  "properties": {
    "compression": "None",
    "consumerGroup": "acctesteventhubcg",
    "dataFormat": "",
    "eventHubResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG924/providers/Microsoft.EventHub/namespaces/acctesteventhubnamespace/eventhubs/acctesteventhub",
    "eventSystemProperties": [],
    "managedIdentityResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG924/providers/Microsoft.ManagedIdentity/userAssignedIdentities/acctest924",
    "mappingRuleName": "",
    "provisioningState": "Succeeded",
    "tableName": ""
  },
  "type": "Microsoft.Kusto/Clusters/Databases/DataConnections"
}
`
	expectedJson := `
{
  "properties": {
    "consumerGroup": "acctesteventhubcg"
  }
}
`

	var old, expected interface{}
	_ = json.Unmarshal([]byte(oldJson), &old)
	_ = json.Unmarshal([]byte(expectedJson), &expected)

	result := utils.ExtractObject(old, "properties.consumerGroup")
	if !reflect.DeepEqual(result, expected) {
		expectedJson, _ := json.Marshal(expected)
		resultJson, _ := json.Marshal(result)
		t.Fatalf("Expected %s but got %s", expectedJson, resultJson)
	}

	// invalid path
	result = utils.ExtractObject(old, "properties.consumerGroup1")
	if result != nil {
		resultJson, _ := json.Marshal(result)
		t.Fatalf("Expected nil but got %s", resultJson)
	}
}
