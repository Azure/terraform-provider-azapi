package utils_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ms-henglu/terraform-provider-azurermg/utils"
)

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
    "name": "henglulb",
    "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/henglu-rg1111/providers/Microsoft.Network/loadBalancers/henglulb",
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
                "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/henglu-rg1111/providers/Microsoft.Network/loadBalancers/henglulb/frontendIPConfigurations/PublicIPAddress",
                "etag": "W/\"0074bc84-d75b-4496-96c8-c858e0ef0afe\"",
                "type": "Microsoft.Network/loadBalancers/frontendIPConfigurations",
                "properties": {
                    "provisioningState": "Succeeded",
                    "privateIPAllocationMethod": "Dynamic",
                    "publicIPAddress": {
                        "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/henglu-rg1111/providers/Microsoft.Network/publicIPAddresses/hengluIp"
                    },
                    "inboundNatRules": [
                        {
                            "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/henglu-rg1111/providers/Microsoft.Network/loadBalancers/henglulb/inboundNatRules/RDPAccess"
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
                "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/henglu-rg1111/providers/Microsoft.Network/loadBalancers/henglulb/inboundNatRules/RDPAccess",
                "etag": "W/\"0074bc84-d75b-4496-96c8-c858e0ef0afe\"",
                "type": "Microsoft.Network/loadBalancers/inboundNatRules",
                "properties": {
                    "provisioningState": "Succeeded",
                    "frontendIPConfiguration": {
                        "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/henglu-rg1111/providers/Microsoft.Network/loadBalancers/henglulb/frontendIPConfigurations/PublicIPAddress"
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
    "name": "henglulb",
    "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/henglu-rg1111/providers/Microsoft.Network/loadBalancers/henglulb",
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
                "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/henglu-rg1111/providers/Microsoft.Network/loadBalancers/henglulb/frontendIPConfigurations/PublicIPAddress",
                "etag": "W/\"0074bc84-d75b-4496-96c8-c858e0ef0afe\"",
                "type": "Microsoft.Network/loadBalancers/frontendIPConfigurations",
                "properties": {
                    "provisioningState": "Succeeded",
                    "privateIPAllocationMethod": "Dynamic",
                    "publicIPAddress": {
                        "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/henglu-rg1111/providers/Microsoft.Network/publicIPAddresses/hengluIp"
                    },
                    "inboundNatRules": [
                        {
                            "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/henglu-rg1111/providers/Microsoft.Network/loadBalancers/henglulb/inboundNatRules/RDPAccess"
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
                "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/henglu-rg1111/providers/Microsoft.Network/loadBalancers/henglulb/inboundNatRules/RDPAccess",
                "etag": "W/\"0074bc84-d75b-4496-96c8-c858e0ef0afe\"",
                "type": "Microsoft.Network/loadBalancers/inboundNatRules",
                "properties": {
                    "provisioningState": "Succeeded",
                    "frontendIPConfiguration": {
                        "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/henglu-rg1111/providers/Microsoft.Network/loadBalancers/henglulb/frontendIPConfigurations/PublicIPAddress"
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
        "fqdn": "henglu1116-springcloud.azuremicroservices.io",
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
    "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctestRG-henglu1116/providers/Microsoft.AppPlatform/Spring/henglu1116-springcloud/apps/henglu1116-springcloudapp",
    "name": "henglu1116-springcloudapp"
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
        "fqdn": "henglu1116-springcloud.azuremicroservices.io",
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
    "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctestRG-henglu1116/providers/Microsoft.AppPlatform/Spring/henglu1116-springcloud/apps/henglu1116-springcloudapp",
    "name": "henglu1116-springcloudapp"
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
  "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG-henglu924/providers/Microsoft.Kusto/Clusters/acctestkchenglu924/Databases/acctestkd-henglu924/DataConnections/acctestkedc-henglu924",
  "kind": "EventHub",
  "location": "West Europe",
  "name": "acctestkchenglu924/acctestkd-henglu924/acctestkedc-henglu924",
  "properties": {
    "compression": "None",
    "consumerGroup": "acctesteventhubcg-henglu924",
    "dataFormat": "",
    "eventHubResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG-henglu924/providers/Microsoft.EventHub/namespaces/acctesteventhubnamespace-henglu924/eventhubs/acctesteventhub-henglu924",
    "eventSystemProperties": [],
    "managedIdentityResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG-henglu924/providers/Microsoft.ManagedIdentity/userAssignedIdentities/acctesthenglu924",
    "mappingRuleName": "",
    "provisioningState": "Succeeded",
    "tableName": ""
  },
  "type": "Microsoft.Kusto/Clusters/Databases/DataConnections"
}
`
	expectedJson := `
{
  "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG-henglu924/providers/Microsoft.Kusto/Clusters/acctestkchenglu924/Databases/acctestkd-henglu924/DataConnections/acctestkedc-henglu924",
  "kind": "EventHub",
  "location": "West Europe",
  "name": "acctestkchenglu924/acctestkd-henglu924/acctestkedc-henglu924",
  "properties": {
    "compression": "None",
    "consumerGroup": "acctesteventhubcg-henglu924",
    "dataFormat": "",
    "eventHubResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG-henglu924/providers/Microsoft.EventHub/namespaces/acctesteventhubnamespace-henglu924/eventhubs/acctesteventhub-henglu924",
    "eventSystemProperties": [],
    "managedIdentityResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG-henglu924/providers/Microsoft.ManagedIdentity/userAssignedIdentities/acctesthenglu924",
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
  "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG-henglu924/providers/Microsoft.Kusto/Clusters/acctestkchenglu924/Databases/acctestkd-henglu924/DataConnections/acctestkedc-henglu924",
  "kind": "EventHub",
  "location": "West Europe",
  "name": "acctestkchenglu924/acctestkd-henglu924/acctestkedc-henglu924",
  "properties": {
    "compression": "None",
    "consumerGroup": "acctesteventhubcg-henglu924",
    "dataFormat": "",
    "eventHubResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG-henglu924/providers/Microsoft.EventHub/namespaces/acctesteventhubnamespace-henglu924/eventhubs/acctesteventhub-henglu924",
    "eventSystemProperties": [],
    "managedIdentityResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG-henglu924/providers/Microsoft.ManagedIdentity/userAssignedIdentities/acctesthenglu924",
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
    "consumerGroup": "acctesteventhubcg-henglu924"
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
