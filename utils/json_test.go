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

	result := utils.MergeObject(old, new)
	if !reflect.DeepEqual(result, expected) {
		expectedJson, _ := json.Marshal(expected)
		resultJson, _ := json.Marshal(result)
		t.Fatalf("Expected %s but got %s", expectedJson, resultJson)
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

func Test_UpdateObjectPolicyDefinition(t *testing.T) {
	OldJson := `
{
	"properties": {
    "description": "Deploys the diagnostic settings for Database for PostgreSQL to stream to a Log Analytics workspace when any Database for PostgreSQL which is missing this diagnostic settings is created or updated. This policy is superseded by built-in initiative https://www.azadvertizer.net/azpolicyinitiativesadvertizer/0884adba-2312-4468-abeb-5422caed1038.html.",
    "displayName": "[Deprecated]: Deploy Diagnostic Settings for Database for PostgreSQL to Log Analytics workspace",
    "metadata": {
      "alzCloudEnvironments": [
        "AzureCloud",
        "AzureChinaCloud",
        "AzureUSGovernment"
      ],
      "category": "Monitoring",
      "deprecated": true,
      "source": "https://github.com/Azure/Enterprise-Scale/",
      "version": "2.0.0-deprecated"
    },
    "mode": "Indexed",
    "parameters": {
      "effect": {
        "allowedValues": [
          "DeployIfNotExists",
          "Disabled"
        ],
        "defaultValue": "DeployIfNotExists",
        "metadata": {
          "description": "Enable or disable the execution of the policy",
          "displayName": "Effect"
        },
        "type": "String"
      },
      "logAnalytics": {
        "metadata": {
          "description": "Select Log Analytics workspace from dropdown list. If this workspace is outside of the scope of the assignment you must manually grant 'Log Analytics Contributor' permissions (or similar) to the policy assignment's principal ID.",
          "displayName": "Log Analytics workspace",
          "strongType": "omsWorkspace"
        },
        "type": "String"
      },
      "logsEnabled": {
        "allowedValues": [
          "True",
          "False"
        ],
        "defaultValue": "True",
        "metadata": {
          "description": "Whether to enable logs stream to the Log Analytics workspace - True or False",
          "displayName": "Enable logs"
        },
        "type": "String"
      },
      "metricsEnabled": {
        "allowedValues": [
          "True",
          "False"
        ],
        "defaultValue": "True",
        "metadata": {
          "description": "Whether to enable metrics stream to the Log Analytics workspace - True or False",
          "displayName": "Enable metrics"
        },
        "type": "String"
      },
      "profileName": {
        "defaultValue": "setbypolicy",
        "metadata": {
          "description": "The diagnostic settings profile name",
          "displayName": "Profile name"
        },
        "type": "String"
      }
    },
    "policyRule": {
      "if": {
        "anyOf": [
          {
            "equals": "Microsoft.DBforPostgreSQL/flexibleServers",
            "field": "type"
          },
          {
            "equals": "Microsoft.DBforPostgreSQL/servers",
            "field": "type"
          }
        ]
      },
      "then": {
        "details": {
          "deployment": {
            "properties": {
              "mode": "Incremental",
              "parameters": {
                "location": {
                  "value": "[field('location')]"
                },
                "logAnalytics": {
                  "value": "[parameters('logAnalytics')]"
                },
                "logsEnabled": {
                  "value": "[parameters('logsEnabled')]"
                },
                "metricsEnabled": {
                  "value": "[parameters('metricsEnabled')]"
                },
                "profileName": {
                  "value": "[parameters('profileName')]"
                },
                "resourceName": {
                  "value": "[field('name')]"
                },
                "resourceType": {
                  "value": "[field('type')]"
                }
              },
              "template": {
                "$schema": "http://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
                "contentVersion": "1.0.0.0",
                "outputs": {},
                "parameters": {
                  "location": {
                    "type": "String"
                  },
                  "logAnalytics": {
                    "type": "String"
                  },
                  "logsEnabled": {
                    "type": "String"
                  },
                  "metricsEnabled": {
                    "type": "String"
                  },
                  "profileName": {
                    "type": "String"
                  },
                  "resourceName": {
                    "type": "String"
                  },
                  "resourceType": {
                    "type": "String"
                  }
                },
                "resources": [
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
                ],
                "variables": {}
              }
            }
          },
          "existenceCondition": {
            "allOf": [
              {
                "equals": "true",
                "field": "Microsoft.Insights/diagnosticSettings/logs.enabled"
              },
              {
                "equals": "true",
                "field": "Microsoft.Insights/diagnosticSettings/metrics.enabled"
              },
              {
                "equals": "[parameters('logAnalytics')]",
                "field": "Microsoft.Insights/diagnosticSettings/workspaceId"
              }
            ]
          },
          "name": "[parameters('profileName')]",
          "roleDefinitionIds": [
            "/providers/microsoft.authorization/roleDefinitions/749f88d5-cbae-40b8-bcfc-e573ddc772fa",
            "/providers/microsoft.authorization/roleDefinitions/92aaf0da-9dab-42b6-94a3-d43ce8d16293"
          ],
          "type": "Microsoft.Insights/diagnosticSettings"
        },
        "effect": "[parameters('effect')]"
      }
    },
    "policyType": "Custom"
  }
}
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
