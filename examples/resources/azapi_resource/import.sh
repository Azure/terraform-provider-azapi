# Azure resource can be imported using the resource id, e.g.
terraform import azapi_resource.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/cluster1

# It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
terraform import azapi_resource.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/cluster1?api-version=2021-07-01
