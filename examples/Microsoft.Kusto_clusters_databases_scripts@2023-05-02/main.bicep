param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource cluster 'Microsoft.Kusto/clusters@2023-05-02' = {
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
  location: location
  name: resource_name
  properties: {
    enableAutoStop: true
    enableDiskEncryption: false
    enableDoubleEncryption: false
    enablePurge: false
    enableStreamingIngest: false
    engineType: 'V2'
    publicIPType: 'IPv4'
    publicNetworkAccess: 'Enabled'
    restrictOutboundNetworkAccess: 'Disabled'
    trustedExternalTenants: []
  }
  sku: {
    capacity: 1
    name: 'Dev(No SLA)_Standard_D11_v2'
    tier: 'Basic'
  }
}

resource database 'Microsoft.Kusto/clusters/databases@2023-05-02' = {
  parent: cluster
  location: location
  name: resource_name
  kind: 'ReadWrite'
  properties: {}
}

resource script 'Microsoft.Kusto/clusters/databases/scripts@2023-05-02' = {
  parent: database
  name: 'create-table-script'
  properties: {
    continueOnErrors: false
    forceUpdateTag: '9e2e7874-aa37-7041-81b7-06397f03a37d'
    scriptContent: '.create table TestTable(Id:string, Name:string, _ts:long, _timestamp:datetime)\n.create table TestTable ingestion json mapping "TestMapping"\n\'[\'\n\'    {"column":"Id","path":"$.id"},\'\n\'    {"column":"Name","path":"$.name"},\'\n\'    {"column":"_ts","path":"$._ts"},\'\n\'    {"column":"_timestamp","path":"$._ts", "transform":"DateTimeFromUnixSeconds"}\'\n\']\'\n.alter table TestTable policy ingestionbatching "{\'MaximumBatchingTimeSpan\': \'0:0:10\', \'MaximumNumberOfItems\': 10000}"\n'
  }
}

