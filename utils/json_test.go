package utils_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Azure/terraform-provider-azapi/utils"
)

func Test_UpdateObject(t *testing.T) {
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
    "addressSpace": {
      "addressPrefixes": [
        "10.0.0.0/16"
      ]
    },
    "dhcpOptions": {
      "dnsServers": null
    },
    "subnets": [
      {
        "name": "default",
        "properties": {
          "addressPrefix": "10.0.3.0/24"
        }
      }
    ]
  }
}
`,
			NewJson: `
{
  "etag": "W/\"5633f421-aaa4-4cf3-b9a6-40625807bd75\"",
  "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/henglu422/providers/Microsoft.Network/virtualNetworks/henglu422",
  "location": "westus",
  "name": "henglu422",
  "properties": {
    "addressSpace": {
      "addressPrefixes": [
        "10.0.0.0/16"
      ]
    },
    "dhcpOptions": {
      "dnsServers": []
    },
    "enableDdosProtection": false,
    "provisioningState": "Succeeded",
    "resourceGuid": "86ec5186-a958-4b83-b5e7-3b88a728c34c",
    "subnets": [
      {
        "etag": "W/\"5633f421-aaa4-4cf3-b9a6-40625807bd75\"",
        "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/henglu422/providers/Microsoft.Network/virtualNetworks/henglu422/subnets/henglu422",
        "name": "henglu422",
        "properties": {
          "addressPrefix": "10.0.2.0/24",
          "delegations": [],
          "networkSecurityGroup": {
            "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/henglu422/providers/Microsoft.Network/networkSecurityGroups/NRMS-zpedr2ai6dglmhenglu422"
          },
          "privateEndpointNetworkPolicies": "Disabled",
          "privateLinkServiceNetworkPolicies": "Enabled",
          "provisioningState": "Succeeded"
        },
        "type": "Microsoft.Network/virtualNetworks/subnets"
      },
      {
        "etag": "W/\"5633f421-aaa4-4cf3-b9a6-40625807bd75\"",
        "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/henglu422/providers/Microsoft.Network/virtualNetworks/henglu422/subnets/default",
        "name": "default",
        "properties": {
          "addressPrefix": "10.0.3.0/24",
          "delegations": [],
          "networkSecurityGroup": {
            "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/henglu422/providers/Microsoft.Network/networkSecurityGroups/NRMS-zpedr2ai6dglmhenglu422"
          },
          "privateEndpointNetworkPolicies": "Disabled",
          "privateLinkServiceNetworkPolicies": "Enabled",
          "provisioningState": "Succeeded"
        },
        "type": "Microsoft.Network/virtualNetworks/subnets"
      }
    ],
    "virtualNetworkPeerings": []
  },
  "type": "Microsoft.Network/virtualNetworks"
}
`,
			ExpectJson: `
{
  "properties": {
    "addressSpace": {
      "addressPrefixes": [
        "10.0.0.0/16"
      ]
    },
    "dhcpOptions": {
      "dnsServers": []
    },
    "subnets": [
      {
        "name": "default",
        "properties": {
          "addressPrefix": "10.0.3.0/24"
        }
      },
      {
        "etag": "W/\"5633f421-aaa4-4cf3-b9a6-40625807bd75\"",
        "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/henglu422/providers/Microsoft.Network/virtualNetworks/henglu422/subnets/henglu422",
        "name": "henglu422",
        "properties": {
          "addressPrefix": "10.0.2.0/24",
          "delegations": [],
          "networkSecurityGroup": {
            "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/henglu422/providers/Microsoft.Network/networkSecurityGroups/NRMS-zpedr2ai6dglmhenglu422"
          },
          "privateEndpointNetworkPolicies": "Disabled",
          "privateLinkServiceNetworkPolicies": "Enabled",
          "provisioningState": "Succeeded"
        },
        "type": "Microsoft.Network/virtualNetworks/subnets"
      }
    ]
  }
}
`,
		},
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
    "mode": "Deferred",
    "basePolicyName": "Microsoft.Default",
    "contentFilters": [
      {
        "allowedContentLevel": "Low",
        "blocking": true,
        "enabled": true,
        "name": "selfharm",
        "source": "Completion"
      },
      {
        "allowedContentLevel": "Low",
        "blocking": true,
        "enabled": true,
        "name": "violence",
        "source": "Completion"
      },
      {
        "allowedContentLevel": "Low",
        "blocking": true,
        "enabled": true,
        "name": "selfharm",
        "source": "Prompt"
      },
      {
        "allowedContentLevel": "Low",
        "blocking": true,
        "enabled": true,
        "name": "violence",
        "source": "Prompt"
      },
      {
        "allowedContentLevel": "Medium",
        "blocking": true,
        "enabled": true,
        "name": "hate",
        "source": "Prompt"
      },
      {
        "allowedContentLevel": "Medium",
        "blocking": true,
        "enabled": true,
        "name": "sexual",
        "source": "Prompt"
      },
      {
        "blocking": true,
        "enabled": true,
        "name": "jailbreak",
        "source": "Prompt"
      },
      {
        "blocking": true,
        "enabled": true,
        "name": "indirect_attack",
        "source": "Prompt"
      }
    ]
  }
}
`,
			NewJson: `
{
  "properties": {
    "mode": "Deferred",
    "basePolicyName": "Microsoft.Default",
    "contentFilters": [
      {
        "allowedContentLevel": "Low",
        "blocking": true,
        "enabled": true,
        "name": "selfharm",
        "source": "Completion"
      },
      {
        "allowedContentLevel": "Low",
        "blocking": true,
        "enabled": true,
        "name": "violence",
        "source": "Completion"
      },
      {
        "allowedContentLevel": "Low",
        "blocking": true,
        "enabled": true,
        "name": "selfharm",
        "source": "Prompt"
      },
      {
        "allowedContentLevel": "Low",
        "blocking": true,
        "enabled": true,
        "name": "violence",
        "source": "Prompt"
      },
      {
        "allowedContentLevel": "Medium",
        "blocking": true,
        "enabled": true,
        "name": "hate",
        "source": "Prompt"
      },
      {
        "allowedContentLevel": "Medium",
        "blocking": true,
        "enabled": true,
        "name": "sexual",
        "source": "Prompt"
      },
      {
        "blocking": true,
        "enabled": true,
        "name": "jailbreak",
        "source": "Prompt"
      },
      {
        "blocking": true,
        "enabled": true,
        "name": "indirect_attack",
        "source": "Prompt"
      }
    ]
  }
}`,
			ExpectJson: `
{
  "properties": {
    "mode": "Deferred",
    "basePolicyName": "Microsoft.Default",
    "contentFilters": [
      {
        "allowedContentLevel": "Low",
        "blocking": true,
        "enabled": true,
        "name": "selfharm",
        "source": "Completion"
      },
      {
        "allowedContentLevel": "Low",
        "blocking": true,
        "enabled": true,
        "name": "violence",
        "source": "Completion"
      },
      {
        "allowedContentLevel": "Low",
        "blocking": true,
        "enabled": true,
        "name": "selfharm",
        "source": "Prompt"
      },
      {
        "allowedContentLevel": "Low",
        "blocking": true,
        "enabled": true,
        "name": "violence",
        "source": "Prompt"
      },
      {
        "allowedContentLevel": "Medium",
        "blocking": true,
        "enabled": true,
        "name": "hate",
        "source": "Prompt"
      },
      {
        "allowedContentLevel": "Medium",
        "blocking": true,
        "enabled": true,
        "name": "sexual",
        "source": "Prompt"
      },
      {
        "blocking": true,
        "enabled": true,
        "name": "jailbreak",
        "source": "Prompt"
      },
      {
        "blocking": true,
        "enabled": true,
        "name": "indirect_attack",
        "source": "Prompt"
      }
    ]
  }
}`,
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

		result := utils.UpdateObject(old, new, testcase.Option)
		if !reflect.DeepEqual(result, expected) {
			expectedJson, _ := json.Marshal(expected)
			resultJson, _ := json.Marshal(result)
			t.Fatalf("Expected %s but got %s", expectedJson, resultJson)
		}
	}
}

func Test_MergeObject(t *testing.T) {
	testcases := []struct {
		Name       string
		OldJson    string
		NewJson    string
		ExpectJson string
	}{
		{
			Name: "Merge simple objects",
			OldJson: ` 
{
	"a":1,
    "b": {
		"b1": "b1",
		"b2": []
	}
}
`,
			NewJson: `
{
	"b": {
		"b3": "b3"
	}
}
`,
			ExpectJson: `
{
	"a":1,
    "b": {
		"b1": "b1",
		"b2": [],
		"b3": "b3"
	}
}`,
		},
		{
			Name: "Merge objects with arrays",
			OldJson: `
{
  "array": [
    {
      "name": "item1",
      "key1": "value1"
    },
    {
      "name": "item2",
      "key2": "value2"
    }
  ]
}`,
			NewJson: `
{
  "array": [
    {
      "name": "item2",
      "key2": "value3"
    },
    {
      "name": "item1",
      "key1": "value2"
    }
  ]
}`,
			ExpectJson: `
{
  "array": [
    {
      "name": "item1",
      "key1": "value2"
    },
    {
      "name": "item2",
      "key2": "value3"
    }
  ]
}
`,
		},
	}

	for _, testcase := range testcases {
		t.Logf("Running test case: %s", testcase.Name)
		var new, old, expected interface{}
		_ = json.Unmarshal([]byte(testcase.OldJson), &old)
		_ = json.Unmarshal([]byte(testcase.NewJson), &new)
		_ = json.Unmarshal([]byte(testcase.ExpectJson), &expected)

		result := utils.MergeObject(old, new)
		if !reflect.DeepEqual(result, expected) {
			expectedJson, _ := json.Marshal(expected)
			resultJson, _ := json.Marshal(result)
			t.Fatalf("Test %s: Expected %s but got %s", testcase.Name, expectedJson, resultJson)
		}
	}
}

func Test_MergeObjectWithArray(t *testing.T) {
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

	result := utils.MergeObject(old, new)
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

func Test_OverrideWithPaths(t *testing.T) {
	testcases := []struct {
		OldJson       string
		NewJson       string
		ExpectJson    string
		IgnoreChanges []string
	}{
		{
			OldJson: `
{
    "properties": {
        "provisioningState": "Foo",
        "subnets": [
            {
                "name": "default",
                "value": "foo"
            }
        ]
    }
}`,
			NewJson: `
{
    "properties": {
        "provisioningState": "Bar",
        "subnets": [
            {
                "name": "default",
                "value": "foo"
            },
            {
                "name": "spin-app-subnet",
                "id": "bar"
            }
        ]
    }
}`,
			ExpectJson: `
{
    "properties": {
        "provisioningState": "Foo",
        "subnets": [
            {
                "name": "default",
                "value": "foo"
            },
            {
                "name": "spin-app-subnet",
                "id": "bar"
            }
        ]
    }
}`,
			IgnoreChanges: []string{"properties.subnets"},
		},
		{
			OldJson: `
{
    "properties": {
        "provisioningState": "Foo"
    }
}`,
			NewJson: `
{
    "properties": {
        "provisioningState": "Bar"
    }
}`,
			ExpectJson: `
{
    "properties": {
        "provisioningState": "Bar"
    }
}`,
			IgnoreChanges: []string{"properties.provisioningState"},
		},
	}

	for _, testcase := range testcases {
		var new, old, expected interface{}
		_ = json.Unmarshal([]byte(testcase.OldJson), &old)
		_ = json.Unmarshal([]byte(testcase.NewJson), &new)
		_ = json.Unmarshal([]byte(testcase.ExpectJson), &expected)

		pathSet := make(map[string]bool)
		for _, path := range testcase.IgnoreChanges {
			pathSet[path] = true
		}
		result, err := utils.OverrideWithPaths(old, new, "", pathSet)
		if err != nil {
			t.Fatalf("Expected no error but got %s", err)
		}
		if !reflect.DeepEqual(result, expected) {
			expectedJson, _ := json.Marshal(expected)
			resultJson, _ := json.Marshal(result)
			t.Fatalf("Expected %s but got %s", expectedJson, resultJson)
		}
	}
}

func Test_ExtractObjectJMES(t *testing.T) {
	testcases := []struct {
		InputJson  string
		PathKey    string
		Path       string
		ExpectJson string
	}{
		{
			InputJson: `
{
  "values": [
    {
      "id": "1",
      "name": "test1",
      "properties": {
        "status": "active"
      }
    },
    {
      "id": "2",
      "name": "test2",
      "properties": {
        "status": "inactive"
      }
    }
  ]
}
`,
			PathKey: "values[*].name",
			Path:    "values[*].name",
			ExpectJson: `
{
  "values[*].name": ["test1", "test2"]
}
`,
		},
		{
			InputJson: `
{
  "values": [
    {
      "id": "1",
      "name": "test1",
      "properties": {
        "status": "active"
      }
    },
    {
      "id": "2",
      "name": "test2",
      "properties": {
        "status": "inactive"
      }
    }
  ]
}
`,
			PathKey: "values[*].status",
			Path:    "values[*].properties.status",
			ExpectJson: `
{
  "values[*].status": ["active", "inactive"]
}
`,
		},
		{
			InputJson: `
{
  "values": [
    {
      "id": "1",
      "name": "test1",
      "properties": {
        "status": "active"
      }
    },
    {
      "id": "2",
      "name": "test2",
      "properties": {
        "status": "inactive"
      }
    }
  ]
}
`,
			PathKey: "values[*].nonexistent",
			Path:    "values[*].nonexistent",
			ExpectJson: `
{
  "values[*].nonexistent": []
}
`,
		},
	}

	for _, testcase := range testcases {
		var input, expected interface{}
		_ = json.Unmarshal([]byte(testcase.InputJson), &input)
		_ = json.Unmarshal([]byte(testcase.ExpectJson), &expected)

		result := utils.ExtractObjectJMES(input, testcase.PathKey, testcase.Path)
		if !reflect.DeepEqual(result, expected) {
			expectedJson, _ := json.Marshal(expected)
			resultJson, _ := json.Marshal(result)
			t.Fatalf("Expected %s but got %s", string(expectedJson), string(resultJson))
		}
	}
}

func Test_UpdateObjectDuplicateIdentifiers(t *testing.T) {
	OldJson := `
[
	{
		"apiVersion": "2021-05-01-preview",
		"condition": "[startsWith(parameters('resourceType'),'Microsoft.DBforPostgreSQL/flexibleServers')]",
		"dependsOn": [],
		"location": "[parameters('location')]",
		"name": "[concat(parameters('resourceName'), '/', 'Microsoft.Insights/', parameters('profileName'))]",
		"properties": {
			"logs": [
				{
					"category": "PostgreSQLLogs",
					"enabled": "[parameters('logsEnabled')]"
				}
			],
			"metrics": [
				{
					"category": "AllMetrics",
					"enabled": "[parameters('metricsEnabled')]",
					"retentionPolicy": {
						"days": 0,
						"enabled": false
					},
					"timeGrain": null
				}
			],
			"workspaceId": "[parameters('logAnalytics')]"
		},
		"type": "Microsoft.DBforPostgreSQL/flexibleServers/providers/diagnosticSettings"
	},
	{
		"apiVersion": "2021-05-01-preview",
		"condition": "[startsWith(parameters('resourceType'),'Microsoft.DBforPostgreSQL/servers')]",
		"dependsOn": [],
		"location": "[parameters('location')]",
		"name": "[concat(parameters('resourceName'), '/', 'Microsoft.Insights/', parameters('profileName'))]",
		"properties": {
			"logs": [
				{
					"category": "PostgreSQLLogs",
					"enabled": "[parameters('logsEnabled')]"
				},
				{
					"category": "QueryStoreRuntimeStatistics",
					"enabled": "[parameters('logsEnabled')]"
				},
				{
					"category": "QueryStoreWaitStatistics",
					"enabled": "[parameters('logsEnabled')]"
				}
			],
			"metrics": [
				{
					"category": "AllMetrics",
					"enabled": "[parameters('metricsEnabled')]",
					"retentionPolicy": {
						"days": 0,
						"enabled": false
					},
					"timeGrain": null
				}
			],
			"workspaceId": "[parameters('logAnalytics')]"
		},
		"type": "Microsoft.DBforPostgreSQL/servers/providers/diagnosticSettings"
	}
]
`
	var old, new, expected any
	_ = json.Unmarshal([]byte(OldJson), &old)
	_ = json.Unmarshal([]byte(OldJson), &new)
	_ = json.Unmarshal([]byte(OldJson), &expected)

	got := utils.UpdateObject(old, new, utils.UpdateJsonOption{
		IgnoreCasing:          false,
		IgnoreMissingProperty: true,
	})
	if !reflect.DeepEqual(got, expected) {
		expectedJson, _ := json.MarshalIndent(expected, "", "  ")
		gotJson, _ := json.MarshalIndent(got, "", "  ")
		t.Fatalf("Expected:\n%s\n\n but got\n%s", expectedJson, gotJson)
	}
}

func Test_UpdateObjectDuplicateIdentifiersWithInconsistentOrdering(t *testing.T) {
	OldJson := `
{
	"contentFilters": [
		{
			"name": "Hate",
			"allowedContentLevel": "Medium",
			"blocking": true,
			"enabled": true,
			"source": "Prompt"
		},
		{
			"name": "Hate",
			"allowedContentLevel": "Medium",
			"blocking": true,
			"enabled": true,
			"source": "Completion"
		}
	]
}
`

	NewJson := `
{
	"contentFilters": [
		{
			"name": "Hate",
			"allowedContentLevel": "Medium",
			"blocking": true,
			"enabled": true,
			"source": "Completion"
		},
		{
			"name": "Hate",
			"allowedContentLevel": "Medium",
			"blocking": true,
			"enabled": true,
			"source": "Prompt"
		}
	]
}
`
	var old, new, expected any
	_ = json.Unmarshal([]byte(OldJson), &old)
	_ = json.Unmarshal([]byte(NewJson), &new)
	_ = json.Unmarshal([]byte(OldJson), &expected)

	got := utils.UpdateObject(old, new, utils.UpdateJsonOption{
		IgnoreCasing:          false,
		IgnoreMissingProperty: true,
	})
	if !reflect.DeepEqual(got, expected) {
		expectedJson, _ := json.MarshalIndent(expected, "", "  ")
		gotJson, _ := json.MarshalIndent(got, "", "  ")
		t.Fatalf("Expected:\n%s\n\n but got\n%s", expectedJson, gotJson)
	}
}

func Test_UpdateObject_ListUniqueIdProperty_IgnoreOtherItemsInList(t *testing.T) {
	oldJson := `
{
  "properties": {
    "serviceEndpoints": [
      {
        "service": "Microsoft.Web"
      }
    ]
  }
}
`
	newJson := `
{
  "properties": {
    "serviceEndpoints": [
      {
        "service": "Microsoft.KeyVault",
        "locations": ["westus"]
      },
      {
        "service": "Microsoft.Web",
        "locations": ["westus2"]
      }
    ]
  }
}
`

	var oldObj interface{}
	var newObj interface{}
	if err := json.Unmarshal([]byte(oldJson), &oldObj); err != nil {
		t.Fatalf("failed to unmarshal oldJson: %v", err)
	}
	if err := json.Unmarshal([]byte(newJson), &newObj); err != nil {
		t.Fatalf("failed to unmarshal newJson: %v", err)
	}

	got := utils.UpdateObject(oldObj, newObj, utils.UpdateJsonOption{
		ListUniqueIdProperty: map[string]string{
			"properties.serviceEndpoints": "service",
		},
		IgnoreOtherItemsInList: map[string]bool{
			"properties.serviceEndpoints": true,
		},
	})

	if !reflect.DeepEqual(got, oldObj) {
		gotBytes, _ := json.MarshalIndent(got, "", "  ")
		expBytes, _ := json.MarshalIndent(oldObj, "", "  ")
		t.Fatalf("unexpected result:\n got: %s\nwant: %s", string(gotBytes), string(expBytes))
	}
}

func Test_MergeObjectWithOption_ListUniqueIdProperty_AppendsNewItems(t *testing.T) {
	oldJson := `
{
  "properties": {
    "serviceEndpoints": [
      {
        "service": "Microsoft.KeyVault",
        "x": 1
      }
    ]
  }
}
`
	newJson := `
{
  "properties": {
    "serviceEndpoints": [
      {
        "service": "Microsoft.Web"
      },
      {
        "service": "Microsoft.KeyVault",
        "x": 2
      }
    ]
  }
}
`

	var oldObj interface{}
	var newObj interface{}
	if err := json.Unmarshal([]byte(oldJson), &oldObj); err != nil {
		t.Fatalf("failed to unmarshal oldJson: %v", err)
	}
	if err := json.Unmarshal([]byte(newJson), &newObj); err != nil {
		t.Fatalf("failed to unmarshal newJson: %v", err)
	}

	got := utils.MergeObjectWithOption(oldObj, newObj, utils.UpdateJsonOption{
		ListUniqueIdProperty: map[string]string{
			"properties.serviceEndpoints": "service",
		},
	})

	gotMap := got.(map[string]interface{})
	props := gotMap["properties"].(map[string]interface{})
	endpoints := props["serviceEndpoints"].([]interface{})
	if len(endpoints) != 2 {
		t.Fatalf("expected 2 endpoints, got %d", len(endpoints))
	}
	first := endpoints[0].(map[string]interface{})
	second := endpoints[1].(map[string]interface{})
	if first["service"] != "Microsoft.KeyVault" {
		t.Fatalf("expected first endpoint to be Microsoft.KeyVault, got %v", first["service"])
	}
	if first["x"] != float64(2) {
		t.Fatalf("expected KeyVault.x to be overwritten to 2, got %v", first["x"])
	}
	if second["service"] != "Microsoft.Web" {
		t.Fatalf("expected second endpoint to be Microsoft.Web, got %v", second["service"])
	}
}

func Test_RemoveFields(t *testing.T) {
	testcases := []struct {
		OldJson    string
		Fields     []string
		ExpectJson string
	}{
		{
			OldJson: `
{
    "apiVersion": "2024-05-01",
    "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/acctestheng125/providers/Microsoft.Network/routeTables/acctestheng125",
    "name": "acctestheng125",
    "type": "microsoft.network/routetables",
    "location": "westus",
    "properties": {
        "provisioningState": "Succeeded",
        "resourceGuid": "c7e6268d-eef3-4aa0-86a2-9a2cdedd59a8",
        "disableBgpRoutePropagation": false,
        "routes": [
            {
                "name": "route1",
                "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/acctestheng125/providers/Microsoft.Network/routeTables/acctestheng125/routes/route1",
                "etag": "W/\"8cbb63f5-3125-4a0b-86d0-a4818311b154\"",
                "properties": {
                    "provisioningState": "Succeeded",
                    "addressPrefix": "10.1.0.0/16",
                    "nextHopType": "VnetLocal",
                    "nextHopIpAddress": "",
                    "hasBgpOverride": false
                },
                "type": "Microsoft.Network/routeTables/routes"
            }
        ]
    },
    "etag": "W/\"8cbb63f5-3125-4a0b-86d0-a4818311b154\""
}
`,
			Fields: []string{"etag"},
			ExpectJson: `
{
    "apiVersion": "2024-05-01",
    "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/acctestheng125/providers/Microsoft.Network/routeTables/acctestheng125",
    "name": "acctestheng125",
    "type": "microsoft.network/routetables",
    "location": "westus",
    "properties": {
        "provisioningState": "Succeeded",
        "resourceGuid": "c7e6268d-eef3-4aa0-86a2-9a2cdedd59a8",
        "disableBgpRoutePropagation": false,
        "routes": [
            {
                "name": "route1",
                "id": "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/acctestheng125/providers/Microsoft.Network/routeTables/acctestheng125/routes/route1",
                "properties": {
                    "provisioningState": "Succeeded",
                    "addressPrefix": "10.1.0.0/16",
                    "nextHopType": "VnetLocal",
                    "nextHopIpAddress": "",
                    "hasBgpOverride": false
                },
                "type": "Microsoft.Network/routeTables/routes"
            }
        ]
    }
}
`,
		},
	}

	for _, testcase := range testcases {
		var old, expected interface{}
		_ = json.Unmarshal([]byte(testcase.OldJson), &old)
		_ = json.Unmarshal([]byte(testcase.ExpectJson), &expected)

		result := utils.RemoveFields(old, testcase.Fields)
		if !reflect.DeepEqual(result, expected) {
			expectedJson, _ := json.Marshal(expected)
			resultJson, _ := json.Marshal(result)
			t.Fatalf("Expected %s but got %s", expectedJson, resultJson)
		}
	}
}

func Test_FilterFields(t *testing.T) {
	testcases := []struct {
		OldJson    string
		Fields     []string
		ExpectJson string
	}{
		{
			OldJson: `
{
	"map": {
		"key1": "value1",
		"key2": "value2"
	},
	"list": [
		{
			"name": "item1",
			"value": "value1"
		},
		{
			"name": "item2",
			"value": "value2"
		}
	],
	"key1": "value1",
	"key2": "value2"
}`,
			Fields: []string{"map.key1", "list[1].name", "key1"},
			ExpectJson: `
{
	"map": {
		"key1": "value1"
	},
	"list": [
		{
			"name": "item2"
		}
	],
	"key1": "value1"
}`,
		},
	}

	for _, testcase := range testcases {
		var old, expected interface{}
		_ = json.Unmarshal([]byte(testcase.OldJson), &old)
		_ = json.Unmarshal([]byte(testcase.ExpectJson), &expected)

		fieldsToKeep := make(map[string]bool)
		for _, field := range testcase.Fields {
			fieldsToKeep[field] = true
		}

		result := utils.FilterFields(old, fieldsToKeep, "")
		if !reflect.DeepEqual(result, expected) {
			expectedJson, _ := json.Marshal(expected)
			resultJson, _ := json.Marshal(result)
			t.Fatalf("Expected %s but got %s", expectedJson, resultJson)
		}
	}
}

// Test_UpdateObject_CompositeKeyIdentifier tests composite keys for list item matching.
func Test_UpdateObject_CompositeKeyIdentifier(t *testing.T) {
	// Simulates diagnostic settings logs where items are identified by
	// either 'category' OR 'categoryGroup' (mutually exclusive).
	oldJson := `{
  "properties": {
    "logs": [
      { "category": "AuditEvent", "enabled": true },
      { "categoryGroup": "allLogs", "enabled": true }
    ]
  }
}`
	newJson := `{
  "properties": {
    "logs": [
      { "category": "Alert", "enabled": false },
      { "categoryGroup": "allLogs", "enabled": true },
      { "category": "AuditEvent", "enabled": true }
    ]
  }
}`
	expectedJson := `{
  "properties": {
    "logs": [
      { "category": "AuditEvent", "enabled": true },
      { "categoryGroup": "allLogs", "enabled": true }
    ]
  }
}`

	var oldObj, newObj, expected interface{}
	if err := json.Unmarshal([]byte(oldJson), &oldObj); err != nil {
		t.Fatalf("failed to unmarshal oldJson: %v", err)
	}
	if err := json.Unmarshal([]byte(newJson), &newObj); err != nil {
		t.Fatalf("failed to unmarshal newJson: %v", err)
	}
	if err := json.Unmarshal([]byte(expectedJson), &expected); err != nil {
		t.Fatalf("failed to unmarshal expectedJson: %v", err)
	}

	// Use composite key "category, categoryGroup" to match on either field
	got := utils.UpdateObject(oldObj, newObj, utils.UpdateJsonOption{
		ListUniqueIdProperty: map[string]string{
			"properties.logs": "category, categoryGroup",
		},
		IgnoreOtherItemsInList: map[string]bool{
			"properties.logs": true,
		},
	})

	if !reflect.DeepEqual(got, expected) {
		gotBytes, _ := json.MarshalIndent(got, "", "  ")
		expBytes, _ := json.MarshalIndent(expected, "", "  ")
		t.Fatalf("composite key identifier failed:\n got: %s\nwant: %s", string(gotBytes), string(expBytes))
	}
}
