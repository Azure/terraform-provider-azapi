param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource schema 'Microsoft.ApiManagement/service/schemas@2021-08-01' = {
  parent: service
  name: resource_name
  properties: {
    description: ''
    schemaType: 'xml'
    value: '    <xsd:schema xmlns:xsd="http://www.w3.org/2001/XMLSchema"\n    xmlns:tns="http://tempuri.org/PurchaseOrderSchema.xsd" targetNamespace="http://tempuri.org/PurchaseOrderSchema.xsd" elementFormDefault="qualified">\n    <xsd:element name="PurchaseOrder" type="tns:PurchaseOrderType"/>\n    <xsd:complexType name="PurchaseOrderType">\n        <xsd:sequence>\n            <xsd:element name="ShipTo" type="tns:USAddress" maxOccurs="2"/>\n            <xsd:element name="BillTo" type="tns:USAddress"/>\n        </xsd:sequence>\n        <xsd:attribute name="OrderDate" type="xsd:date"/>\n    </xsd:complexType>\n    <xsd:complexType name="USAddress">\n        <xsd:sequence>\n            <xsd:element name="name" type="xsd:string"/>\n            <xsd:element name="street" type="xsd:string"/>\n            <xsd:element name="city" type="xsd:string"/>\n            <xsd:element name="state" type="xsd:string"/>\n            <xsd:element name="zip" type="xsd:integer"/>\n        </xsd:sequence>\n        <xsd:attribute name="country" type="xsd:NMTOKEN" fixed="US"/>\n    </xsd:complexType>\n</xsd:schema>\n'
  }
}

resource service 'Microsoft.ApiManagement/service@2021-08-01' = {
  location: location
  name: resource_name
  properties: {
    certificates: []
    customProperties: {
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Ssl30': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls10': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls11': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10': 'false'
      'Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls11': 'false'
    }
    disableGateway: false
    publicNetworkAccess: 'Enabled'
    publisherEmail: 'pub1@email.com'
    publisherName: 'pub1'
    virtualNetworkType: 'None'
  }
  sku: {
    capacity: 0
    name: 'Consumption'
  }
}

