param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource integrationAccount 'Microsoft.Logic/integrationAccounts@2019-05-01' = {
  location: location
  name: resource_name
  properties: {}
  sku: {
    name: 'Basic'
  }
}

resource map 'Microsoft.Logic/integrationAccounts/maps@2019-05-01' = {
  parent: integrationAccount
  name: resource_name
  properties: {
    content: '<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"\n                xmlns:msxsl="urn:schemas-microsoft-com:xslt"\n                xmlns:var="http://schemas.microsoft.com/BizTalk/2003/var"\n                exclude-result-prefixes="msxsl var s0 userCSharp"\n                version="1.0"\n                xmlns:ns0="http://BizTalk_Server_Project4.StringFunctoidsDestinationSchema"\n                xmlns:s0="http://BizTalk_Server_Project4.StringFunctoidsSourceSchema"\n                xmlns:userCSharp="http://schemas.microsoft.com/BizTalk/2003/userCSharp">\n<xsl:import href="http://btsfunctoids.blob.core.windows.net/functoids/functoids.xslt" />\n<xsl:output omit-xml-declaration="yes"\n            method="xml"\n            version="1.0" />\n<xsl:template match="/">\n<xsl:apply-templates select="/s0:Root" />\n</xsl:template>\n<xsl:template match="/s0:Root">\n<xsl:variable name="var:v1"\n              select="userCSharp:StringFind(string(StringFindSource/text()) , &quot;SearchString&quot;)" />\n<xsl:variable name="var:v2"\n              select="userCSharp:StringLeft(string(StringLeftSource/text()) , &quot;2&quot;)" />\n<xsl:variable name="var:v3"\n              select="userCSharp:StringRight(string(StringRightSource/text()) , &quot;2&quot;)" />\n<xsl:variable name="var:v4"\n              select="userCSharp:StringUpperCase(string(UppercaseSource/text()))" />\n<xsl:variable name="var:v5"\n              select="userCSharp:StringLowerCase(string(LowercaseSource/text()))" />\n<xsl:variable name="var:v6"\n              select="userCSharp:StringSize(string(SizeSource/text()))" />\n<xsl:variable name="var:v7"\n              select="userCSharp:StringSubstring(string(StringExtractSource/text()) , &quot;0&quot; , &quot;2&quot;)" />\n<xsl:variable name="var:v8"\n              select="userCSharp:StringConcat(string(StringConcatSource/text()))" />\n<xsl:variable name="var:v9"\n              select="userCSharp:StringTrimLeft(string(StringLeftTrimSource/text()))" />\n<xsl:variable name="var:v10"\n              select="userCSharp:StringTrimRight(string(StringRightTrimSource/text()))" />\n<ns0:Root>\n<StringFindDestination>\n<xsl:value-of select="$var:v1" />\n</StringFindDestination>\n<StringLeftDestination>\n<xsl:value-of select="$var:v2" />\n</StringLeftDestination>\n<StringRightDestination>\n<xsl:value-of select="$var:v3" />\n</StringRightDestination>\n<UppercaseDestination>\n<xsl:value-of select="$var:v4" />\n</UppercaseDestination>\n<LowercaseDestination>\n<xsl:value-of select="$var:v5" />\n</LowercaseDestination>\n<SizeDestination>\n<xsl:value-of select="$var:v6" />\n</SizeDestination>\n<StringExtractDestination>\n<xsl:value-of select="$var:v7" />\n</StringExtractDestination>\n<StringConcatDestination>\n<xsl:value-of select="$var:v8" />\n</StringConcatDestination>\n<StringLeftTrimDestination>\n<xsl:value-of select="$var:v9" />\n</StringLeftTrimDestination>\n<StringRightTrimDestination>\n<xsl:value-of select="$var:v10" />\n</StringRightTrimDestination>\n</ns0:Root>\n</xsl:template>\n</xsl:stylesheet>\n'
    contentType: 'application/xml'
    mapType: 'Xslt'
  }
}

