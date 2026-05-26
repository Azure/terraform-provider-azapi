param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource api 'Microsoft.ApiManagement/service/apis@2021-08-01' = {
  parent: service
  name: '${resource_name};rev=1'
  properties: {
    apiRevisionDescription: ''
    apiType: 'http'
    apiVersion: ''
    apiVersionDescription: ''
    authenticationSettings: {}
    description: ''
    displayName: 'api1'
    path: 'api1'
    protocols: [
      'https'
    ]
    serviceUrl: ''
    subscriptionRequired: true
    type: 'http'
  }
}

resource schema 'Microsoft.ApiManagement/service/apis/schemas@2021-08-01' = {
  parent: api
  name: resource_name
  properties: {
    contentType: 'application/vnd.ms-azure-apim.xsd+xml'
    document: {
      value: '<!--\n Copyright (c) HashiCorp, Inc.\n SPDX-License-Identifier: MPL-2.0\n-->\n\n<s:schema elementFormDefault="qualified" targetNamespace="http://ws.cdyne.com/WeatherWS/" xmlns:tns="http://ws.cdyne.com/WeatherWS/" xmlns:s="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://schemas.xmlsoap.org/wsdl/soap12/" xmlns:mime="http://schemas.xmlsoap.org/wsdl/mime/" xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/" xmlns:tm="http://microsoft.com/wsdl/mime/textMatching/" xmlns:http="http://schemas.xmlsoap.org/wsdl/http/" xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/" xmlns:apim-wsdltns="http://ws.cdyne.com/WeatherWS/"> <s:element name="GetWeatherInformation"> <s:complexType /> </s:element> <s:element name="GetWeatherInformationResponse"> <s:complexType> <s:sequence> <s:element minOccurs="0" maxOccurs="1" name="GetWeatherInformationResult" type="tns:ArrayOfWeatherDescription" /> </s:sequence> </s:complexType> </s:element> <s:complexType name="ArrayOfWeatherDescription"> <s:sequence> <s:element minOccurs="0" maxOccurs="unbounded" name="WeatherDescription" type="tns:WeatherDescription" /> </s:sequence> </s:complexType> <s:complexType name="WeatherDescription"> <s:sequence> <s:element minOccurs="1" maxOccurs="1" name="WeatherID" type="s:short" /> <s:element minOccurs="0" maxOccurs="1" name="Description" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="PictureURL" type="s:string" /> </s:sequence> </s:complexType> <s:element name="GetCityForecastByZIP"> <s:complexType> <s:sequence> <s:element minOccurs="0" maxOccurs="1" name="ZIP" type="s:string" /> </s:sequence> </s:complexType> </s:element> <s:element name="GetCityForecastByZIPResponse"> <s:complexType> <s:sequence> <s:element minOccurs="0" maxOccurs="1" name="GetCityForecastByZIPResult" type="tns:ForecastReturn" /> </s:sequence> </s:complexType> </s:element> <s:complexType name="ForecastReturn"> <s:sequence> <s:element minOccurs="1" maxOccurs="1" name="Success" type="s:boolean" /> <s:element minOccurs="0" maxOccurs="1" name="ResponseText" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="State" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="City" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="WeatherStationCity" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="ForecastResult" type="tns:ArrayOfForecast" /> </s:sequence> </s:complexType> <s:complexType name="ArrayOfForecast"> <s:sequence> <s:element minOccurs="0" maxOccurs="unbounded" name="Forecast" nillable="true" type="tns:Forecast" /> </s:sequence> </s:complexType> <s:complexType name="Forecast"> <s:sequence> <s:element minOccurs="1" maxOccurs="1" name="Date" type="s:dateTime" /> <s:element minOccurs="1" maxOccurs="1" name="WeatherID" type="s:short" /> <s:element minOccurs="0" maxOccurs="1" name="Desciption" type="s:string" /> <s:element minOccurs="1" maxOccurs="1" name="Temperatures" type="tns:temp" /> <s:element minOccurs="1" maxOccurs="1" name="ProbabilityOfPrecipiation" type="tns:POP" /> </s:sequence> </s:complexType> <s:complexType name="temp"> <s:sequence> <s:element minOccurs="0" maxOccurs="1" name="MorningLow" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="DaytimeHigh" type="s:string" /> </s:sequence> </s:complexType> <s:complexType name="POP"> <s:sequence> <s:element minOccurs="0" maxOccurs="1" name="Nighttime" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="Daytime" type="s:string" /> </s:sequence> </s:complexType> <s:element name="GetCityWeatherByZIP"> <s:complexType> <s:sequence> <s:element minOccurs="0" maxOccurs="1" name="ZIP" type="s:string" /> </s:sequence> </s:complexType> </s:element> <s:element name="GetCityWeatherByZIPResponse"> <s:complexType> <s:sequence> <s:element minOccurs="1" maxOccurs="1" name="GetCityWeatherByZIPResult" type="tns:WeatherReturn" /> </s:sequence> </s:complexType> </s:element> <s:complexType name="WeatherReturn"> <s:sequence> <s:element minOccurs="1" maxOccurs="1" name="Success" type="s:boolean" /> <s:element minOccurs="0" maxOccurs="1" name="ResponseText" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="State" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="City" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="WeatherStationCity" type="s:string" /> <s:element minOccurs="1" maxOccurs="1" name="WeatherID" type="s:short" /> <s:element minOccurs="0" maxOccurs="1" name="Description" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="Temperature" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="RelativeHumidity" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="Wind" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="Pressure" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="Visibility" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="WindChill" type="s:string" /> <s:element minOccurs="0" maxOccurs="1" name="Remarks" type="s:string" /> </s:sequence> </s:complexType> <s:element name="ArrayOfWeatherDescription" nillable="true" type="tns:ArrayOfWeatherDescription" /> <s:element name="ForecastReturn" nillable="true" type="tns:ForecastReturn" /> <s:element name="WeatherReturn" type="tns:WeatherReturn" /> </s:schema>'
    }
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

