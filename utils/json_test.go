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
	json.Unmarshal([]byte(oldJson), &old)
	json.Unmarshal([]byte(newJson), &new)
	json.Unmarshal([]byte(expectedJson), &expected)

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
	"a":1,
    "b": {
		"b1": "b1",
		"b2": [],
		"b3": "b3"
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
		"b3": null
	}
}
`
	var new, old, expected interface{}
	json.Unmarshal([]byte(oldJson), &old)
	json.Unmarshal([]byte(newJson), &new)
	json.Unmarshal([]byte(expectedJson), &expected)

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
  "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctestRG-henglu924/providers/Microsoft.Kusto/Clusters/acctestkchenglu924/Databases/acctestkd-henglu924/DataConnections/acctestkedc-henglu924",
  "kind": "EventHub",
  "location": "West Europe",
  "name": "acctestkchenglu924/acctestkd-henglu924/acctestkedc-henglu924",
  "properties": {
    "compression": "None",
    "consumerGroup": "acctesteventhubcg-henglu924",
    "dataFormat": "",
    "eventHubResourceId": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctestRG-henglu924/providers/Microsoft.EventHub/namespaces/acctesteventhubnamespace-henglu924/eventhubs/acctesteventhub-henglu924",
    "eventSystemProperties": [],
    "managedIdentityResourceId": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctestRG-henglu924/providers/Microsoft.ManagedIdentity/userAssignedIdentities/acctesthenglu924",
    "mappingRuleName": "",
    "provisioningState": "Succeeded",
    "tableName": ""
  },
  "type": "Microsoft.Kusto/Clusters/Databases/DataConnections"
}
`
	expectedJson := `
{
  "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctestRG-henglu924/providers/Microsoft.Kusto/Clusters/acctestkchenglu924/Databases/acctestkd-henglu924/DataConnections/acctestkedc-henglu924",
  "kind": "EventHub",
  "location": "West Europe",
  "name": "acctestkchenglu924/acctestkd-henglu924/acctestkedc-henglu924",
  "properties": {
    "compression": "None",
    "consumerGroup": "acctesteventhubcg-henglu924",
    "dataFormat": "",
    "eventHubResourceId": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctestRG-henglu924/providers/Microsoft.EventHub/namespaces/acctesteventhubnamespace-henglu924/eventhubs/acctesteventhub-henglu924",
    "eventSystemProperties": [],
    "managedIdentityResourceId": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctestRG-henglu924/providers/Microsoft.ManagedIdentity/userAssignedIdentities/acctesthenglu924",
    "mappingRuleName": "",
    "tableName": ""
  },
  "type": "Microsoft.Kusto/Clusters/Databases/DataConnections"
}
`

	var old, expected interface{}
	json.Unmarshal([]byte(oldJson), &old)
	json.Unmarshal([]byte(expectedJson), &expected)

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
  "id": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctestRG-henglu924/providers/Microsoft.Kusto/Clusters/acctestkchenglu924/Databases/acctestkd-henglu924/DataConnections/acctestkedc-henglu924",
  "kind": "EventHub",
  "location": "West Europe",
  "name": "acctestkchenglu924/acctestkd-henglu924/acctestkedc-henglu924",
  "properties": {
    "compression": "None",
    "consumerGroup": "acctesteventhubcg-henglu924",
    "dataFormat": "",
    "eventHubResourceId": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctestRG-henglu924/providers/Microsoft.EventHub/namespaces/acctesteventhubnamespace-henglu924/eventhubs/acctesteventhub-henglu924",
    "eventSystemProperties": [],
    "managedIdentityResourceId": "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctestRG-henglu924/providers/Microsoft.ManagedIdentity/userAssignedIdentities/acctesthenglu924",
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
	json.Unmarshal([]byte(oldJson), &old)
	json.Unmarshal([]byte(expectedJson), &expected)

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
