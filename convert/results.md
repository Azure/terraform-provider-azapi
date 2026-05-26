# Conversion Results
## [Microsoft.AADIAM_diagnosticSettings@2017-04-01/main.bicep](../examples/Microsoft.AADIAM_diagnosticSettings@2017-04-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP135] Scope "resourceGroup" is not valid for this resource type. Permitted scopes: "tenant".
```
## [Microsoft.Advisor_recommendations_suppressions@2023-01-01/main.bicep](../examples/Microsoft.Advisor_recommendations_suppressions@2023-01-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Warning no-unused-params] Parameter "recommendation_id" is declared but never used.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.AlertsManagement_actionRules@2021-08-08/main.bicep](../examples/Microsoft.AlertsManagement_actionRules@2021-08-08/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Error BCP063] The name "resourceGroup" is not a parameter, variable, resource or module.
```
## [Microsoft.AlertsManagement_prometheusRuleGroups@2023-03-01/main.bicep](../examples/Microsoft.AlertsManagement_prometheusRuleGroups@2023-03-01/main.bicep)
Result: success

## [microsoft.alertsManagement_smartDetectorAlertRules@2019-06-01/main.bicep](../examples/microsoft.alertsManagement_smartDetectorAlertRules@2019-06-01/main.bicep)
Result: success

## [Microsoft.AnalysisServices_servers@2017-08-01/main.bicep](../examples/Microsoft.AnalysisServices_servers@2017-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_apis@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_apis@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_apis_diagnostics@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_apis_diagnostics@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_apis_operations@2022-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_apis_operations@2022-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_apis_policies@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_apis_policies@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_apis_releases@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_apis_releases@2021-08-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "trimsuffix" does not exist in the current context.
[Error BCP057] The name "azapi_resource" does not exist in the current context.
[Error BCP103] The following token is not recognized: """. Strings are defined using single quotes in bicep.
[Error BCP009] Expected a literal value, an array, an object, a parenthesized expression, or a function call at this location.
[Error BCP103] The following token is not recognized: """. Strings are defined using single quotes in bicep.
```
## [Microsoft.ApiManagement_service_apis_schemas@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_apis_schemas@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_apis_tagDescriptions@2022-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_apis_tagDescriptions@2022-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_apis_tags@2022-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_apis_tags@2022-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_apiVersionSets@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_apiVersionSets@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_authorizationServers@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_authorizationServers@2021-08-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'oauth_client_secret' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'clientSecret' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.ApiManagement_service_backends@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_backends@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_caches@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_caches@2021-08-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.ApiManagement_service_certificates@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_certificates@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_diagnostics@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_diagnostics@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_gateways@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_gateways@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_gateways_certificateAuthorities@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_gateways_certificateAuthorities@2021-08-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'certificate_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.ApiManagement_service_gateways_hostnameConfigurations@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_gateways_hostnameConfigurations@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_groups@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_groups@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_identityProviders@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_identityProviders@2021-08-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'identity_provider_client_secret' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'clientSecret' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.ApiManagement_service_loggers@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_loggers@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_namedValues@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_namedValues@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_openidConnectProviders@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_openidConnectProviders@2021-08-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'openid_client_secret' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'clientSecret' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.ApiManagement_service_policyFragments@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_policyFragments@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_portalsettings@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_portalsettings@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_products@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_products@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_products_policies@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_products_policies@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_products_tags@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_products_tags@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_schemas@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_schemas@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_subscriptions@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_subscriptions@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_tags@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_tags@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_templates@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_templates@2021-08-01/main.bicep)
Result: success

## [Microsoft.ApiManagement_service_users@2021-08-01/main.bicep](../examples/Microsoft.ApiManagement_service_users@2021-08-01/main.bicep)
Result: success

## [Microsoft.AppConfiguration_configurationStores@2023-03-01/main.bicep](../examples/Microsoft.AppConfiguration_configurationStores@2023-03-01/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_apiPortals@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_apiPortals@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_apiPortals_domains@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_apiPortals_domains@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_applicationAccelerators@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_applicationAccelerators@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_applicationAccelerators_customizedAccelerators@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_applicationAccelerators_customizedAccelerators@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_applicationLiveViews@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_applicationLiveViews@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_apps@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_apps@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_apps_bindings@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_apps_bindings@2023-05-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.AppPlatform_Spring_apps_deployments@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_apps_deployments@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_buildServices@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_buildServices@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_buildServices_builders@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_buildServices_builders@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_buildServices_builders_buildpackBindings@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_buildServices_builders_buildpackBindings@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_configServers@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_configServers@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_configurationServices@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_configurationServices@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_DevToolPortals@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_DevToolPortals@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_gateways@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_gateways@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_gateways_domains@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_gateways_domains@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_gateways_routeConfigs@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_gateways_routeConfigs@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_monitoringSettings@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_monitoringSettings@2023-05-01-preview/main.bicep)
Result: success

## [Microsoft.AppPlatform_Spring_storages@2023-05-01-preview/main.bicep](../examples/Microsoft.AppPlatform_Spring_storages@2023-05-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.App_containerApps@2022-03-01/main.bicep](../examples/Microsoft.App_containerApps@2022-03-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning BCP073] The property "ephemeralStorage" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.App_jobs@2025-01-01/main.bicep](../examples/Microsoft.App_jobs@2025-01-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.App_managedEnvironments@2022-03-01/main.bicep](../examples/Microsoft.App_managedEnvironments@2022-03-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.App_managedEnvironments_certificates@2022-03-01/main.bicep](../examples/Microsoft.App_managedEnvironments_certificates@2022-03-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'certificate_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'password' expects a secure value, but the value provided may not be secure.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.App_managedEnvironments_daprComponents@2022-03-01/main.bicep](../examples/Microsoft.App_managedEnvironments_daprComponents@2022-03-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.App_managedEnvironments_dotNetComponents@2024-10-02-preview/main.bicep](../examples/Microsoft.App_managedEnvironments_dotNetComponents@2024-10-02-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.App_managedEnvironments_privateEndpointConnections@2024-10-02-preview/main.bicep](../examples/Microsoft.App_managedEnvironments_privateEndpointConnections@2024-10-02-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.App_managedEnvironments_storages@2022-03-01/main.bicep](../examples/Microsoft.App_managedEnvironments_storages@2022-03-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP053] The type "Microsoft.Storage/storageAccounts" does not contain property "output". Available properties include "apiVersion", "asserts", "eTag", "extendedLocation", "id", "identity", "kind", "location", "managedBy", "managedByExtended", "name", "plan", "properties", "scale", "sku", "tags", "type", "zones".
```
## [Microsoft.Attestation_attestationProviders@2020-10-01/main.bicep](../examples/Microsoft.Attestation_attestationProviders@2020-10-01/main.bicep)
Result: success

## [Microsoft.Authorization_locks@2020-05-01/main.bicep](../examples/Microsoft.Authorization_locks@2020-05-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.Authorization/locks". Permissible properties include "asserts", "dependsOn", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Authorization_policyAssignments@2022-06-01/main.bicep](../examples/Microsoft.Authorization_policyAssignments@2022-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP135] Scope "resourceGroup" is not valid for this resource type. Permitted scopes: "tenant", "managementGroup", "subscription".
```
## [Microsoft.Authorization_policyDefinitions@2021-06-01/main.bicep](../examples/Microsoft.Authorization_policyDefinitions@2021-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Error BCP135] Scope "resourceGroup" is not valid for this resource type. Permitted scopes: "tenant", "managementGroup", "subscription".
```
## [Microsoft.Authorization_policyExemptions@2020-07-01-preview/main.bicep](../examples/Microsoft.Authorization_policyExemptions@2020-07-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP135] Scope "resourceGroup" is not valid for this resource type. Permitted scopes: "tenant", "managementGroup", "subscription".
```
## [Microsoft.Authorization_policySetDefinitions@2025-01-01/main.bicep](../examples/Microsoft.Authorization_policySetDefinitions@2025-01-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Error BCP135] Scope "resourceGroup" is not valid for this resource type. Permitted scopes: "tenant", "managementGroup", "subscription".
[Error BCP135] Scope "resourceGroup" is not valid for this resource type. Permitted scopes: "tenant", "managementGroup", "subscription".
```
## [Microsoft.Authorization_resourceManagementPrivateLinks@2020-05-01/main.bicep](../examples/Microsoft.Authorization_resourceManagementPrivateLinks@2020-05-01/main.bicep)
Result: success

## [Microsoft.Authorization_roleAssignments@2022-04-01/main.bicep](../examples/Microsoft.Authorization_roleAssignments@2022-04-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Warning no-unused-params] Parameter "resource_name" is declared but never used.
[Error BCP057] The name "azurerm_user_assigned_identity" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Authorization_roleDefinitions@2018-01-01-preview/main.bicep](../examples/Microsoft.Authorization_roleDefinitions@2018-01-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Automanage_configurationProfiles@2022-05-04/main.bicep](../examples/Microsoft.Automanage_configurationProfiles@2022-05-04/main.bicep)
Result: success

## [Microsoft.Automation_automationAccounts@2021-06-22/main.bicep](../examples/Microsoft.Automation_automationAccounts@2021-06-22/main.bicep)
Result: success

## [Microsoft.Automation_automationAccounts_certificates@2020-01-13-preview/main.bicep](../examples/Microsoft.Automation_automationAccounts_certificates@2020-01-13-preview/main.bicep)
Result: success

## [Microsoft.Automation_automationAccounts_configurations@2022-08-08/main.bicep](../examples/Microsoft.Automation_automationAccounts_configurations@2022-08-08/main.bicep)
Result: success

## [Microsoft.Automation_automationAccounts_connections@2020-01-13-preview/main.bicep](../examples/Microsoft.Automation_automationAccounts_connections@2020-01-13-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Automation_automationAccounts_connectionTypes@2020-01-13-preview/main.bicep](../examples/Microsoft.Automation_automationAccounts_connectionTypes@2020-01-13-preview/main.bicep)
Result: success

## [Microsoft.Automation_automationAccounts_credentials@2020-01-13-preview/main.bicep](../examples/Microsoft.Automation_automationAccounts_credentials@2020-01-13-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'automation_credential_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Automation_automationAccounts_hybridRunbookWorkerGroups@2021-06-22/main.bicep](../examples/Microsoft.Automation_automationAccounts_hybridRunbookWorkerGroups@2021-06-22/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'credential_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Automation_automationAccounts_hybridRunbookWorkerGroups_hybridRunbookWorkers@2021-06-22/main.bicep](../examples/Microsoft.Automation_automationAccounts_hybridRunbookWorkerGroups_hybridRunbookWorkers@2021-06-22/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'automation_worker_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'vm_admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
[Warning adminusername-should-not-be-literal] Property 'adminUserName' should not use a literal value. Use a param instead. Found literal string value "adminuser"
```
## [Microsoft.Automation_automationAccounts_jobSchedules@2020-01-13-preview/main.bicep](../examples/Microsoft.Automation_automationAccounts_jobSchedules@2020-01-13-preview/main.bicep)
Result: success

## [Microsoft.Automation_automationAccounts_modules@2020-01-13-preview/main.bicep](../examples/Microsoft.Automation_automationAccounts_modules@2020-01-13-preview/main.bicep)
Result: success

Diagnostics:
```
[Warning no-hardcoded-env-urls] Environment URLs should not be hardcoded. Use the environment() function to ensure compatibility across clouds. Found this disallowed host: "core.windows.net"
```
## [Microsoft.Automation_automationAccounts_powershell72Modules@2020-01-13-preview/main.bicep](../examples/Microsoft.Automation_automationAccounts_powershell72Modules@2020-01-13-preview/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP081] Resource type "Microsoft.Automation/automationAccounts/powerShell72Modules@2020-01-13-preview" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
[Warning no-hardcoded-env-urls] Environment URLs should not be hardcoded. Use the environment() function to ensure compatibility across clouds. Found this disallowed host: "core.windows.net"
```
## [Microsoft.Automation_automationAccounts_python3Packages@2023-11-01/main.bicep](../examples/Microsoft.Automation_automationAccounts_python3Packages@2023-11-01/main.bicep)
Result: success

## [Microsoft.Automation_automationAccounts_runbooks@2019-06-01/main.bicep](../examples/Microsoft.Automation_automationAccounts_runbooks@2019-06-01/main.bicep)
Result: success

## [Microsoft.Automation_automationAccounts_runbooks_draft@2018-06-30/main.bicep](../examples/Microsoft.Automation_automationAccounts_runbooks_draft@2018-06-30/main.bicep)
Result: success

## [Microsoft.Automation_automationAccounts_schedules@2020-01-13-preview/main.bicep](../examples/Microsoft.Automation_automationAccounts_schedules@2020-01-13-preview/main.bicep)
Result: success

## [Microsoft.Automation_automationAccounts_softwareUpdateConfigurations@2019-06-01/main.bicep](../examples/Microsoft.Automation_automationAccounts_softwareUpdateConfigurations@2019-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP063] The name "resourceGroup" is not a parameter, variable, resource or module.
```
## [Microsoft.Automation_automationAccounts_sourceControls@2023-11-01/main.bicep](../examples/Microsoft.Automation_automationAccounts_sourceControls@2023-11-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning BCP036] The property "identity" expected a value of type "Identity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Automation_automationAccounts_sourceControls_sourceControlSyncJobs@2023-11-01/main.bicep](../examples/Microsoft.Automation_automationAccounts_sourceControls_sourceControlSyncJobs@2023-11-01/main.bicep)
Result: failed (unexpected error)

```
Error: Invalid HCL input or no resources found.
```
## [Microsoft.Automation_automationAccounts_variables@2020-01-13-preview/main.bicep](../examples/Microsoft.Automation_automationAccounts_variables@2020-01-13-preview/main.bicep)
Result: success

## [Microsoft.Automation_automationAccounts_webHooks@2015-10-31/main.bicep](../examples/Microsoft.Automation_automationAccounts_webHooks@2015-10-31/main.bicep)
Result: success

## [Microsoft.AVS_privateClouds@2022-05-01/main.bicep](../examples/Microsoft.AVS_privateClouds@2022-05-01/main.bicep)
Result: success

## [Microsoft.AVS_privateClouds_authorizations@2022-05-01/main.bicep](../examples/Microsoft.AVS_privateClouds_authorizations@2022-05-01/main.bicep)
Result: success

## [Microsoft.AzureActiveDirectory_b2cDirectories@2021-04-01-preview/main.bicep](../examples/Microsoft.AzureActiveDirectory_b2cDirectories@2021-04-01-preview/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Warning no-unused-params] Parameter "resource_name" is declared but never used.
```
## [Microsoft.Batch_batchAccounts@2022-10-01/main.bicep](../examples/Microsoft.Batch_batchAccounts@2022-10-01/main.bicep)
Result: success

## [Microsoft.Batch_batchAccounts_applications@2022-10-01/main.bicep](../examples/Microsoft.Batch_batchAccounts_applications@2022-10-01/main.bicep)
Result: success

## [Microsoft.Batch_batchAccounts_certificates@2022-10-01/main.bicep](../examples/Microsoft.Batch_batchAccounts_certificates@2022-10-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Batch_batchAccounts_pools@2022-10-01/main.bicep](../examples/Microsoft.Batch_batchAccounts_pools@2022-10-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "placement" expected a value of type "'CacheDisk' | null" but the provided value is of type "''". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Warning BCP036] The property "nodeDeallocationOption" expected a value of type "'Requeue' | 'RetainedData' | 'TaskCompletion' | 'Terminate' | null" but the provided value is of type "''". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Billing_billingAccounts_billingProfiles@2024-04-01/main.bicep](../examples/Microsoft.Billing_billingAccounts_billingProfiles@2024-04-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "billing_account_id" is declared but never used.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning no-unused-params] Parameter "payment_method_id" is declared but never used.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning no-unused-params] Parameter "payment_sca_id" is declared but never used.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP135] Scope "resourceGroup" is not valid for this resource type. Permitted scopes: "tenant".
```
## [Microsoft.BotService_botServices@2021-05-01-preview/main.bicep](../examples/Microsoft.BotService_botServices@2021-05-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.BotService_botServices_channels@2021-05-01-preview/main.bicep](../examples/Microsoft.BotService_botServices_channels@2021-05-01-preview/main.bicep)
Result: success

## [Microsoft.BotService_botServices_connections@2021-05-01-preview/main.bicep](../examples/Microsoft.BotService_botServices_connections@2021-05-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'client_secret' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning no-unused-params] Parameter "location" is declared but never used.
[Error BCP057] The name "data" does not exist in the current context.
[Warning use-secure-value-for-secure-inputs] Property 'clientSecret' expects a secure value, but the value provided may not be secure.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Cache_redis@2023-04-01/main.bicep](../examples/Microsoft.Cache_redis@2023-04-01/main.bicep)
Result: success

## [Microsoft.Cache_redisEnterprise@2022-01-01/main.bicep](../examples/Microsoft.Cache_redisEnterprise@2022-01-01/main.bicep)
Result: success

## [Microsoft.Cache_redisEnterprise_databases@2024-10-01/main.bicep](../examples/Microsoft.Cache_redisEnterprise_databases@2024-10-01/main.bicep)
Result: success

## [Microsoft.Cache_redis_accessPolicies@2024-11-01/main.bicep](../examples/Microsoft.Cache_redis_accessPolicies@2024-11-01/main.bicep)
Result: success

## [Microsoft.Cache_redis_accessPolicyAssignments@2023-04-01/main.bicep](../examples/Microsoft.Cache_redis_accessPolicyAssignments@2023-04-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Cache_redis_firewallRules@2024-11-01/main.bicep](../examples/Microsoft.Cache_redis_firewallRules@2024-11-01/main.bicep)
Result: success

## [Microsoft.Cache_redis_linkedServers@2024-11-01/main.bicep](../examples/Microsoft.Cache_redis_linkedServers@2024-11-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP240] The "parent" property only permits direct references to resources. Expressions are not supported.
[Error BCP057] The name "resourceGroup_secondary" does not exist in the current context.
```
## [Microsoft.Cdn_profiles@2020-09-01/main.bicep](../examples/Microsoft.Cdn_profiles@2020-09-01/main.bicep)
Result: success

## [Microsoft.Cdn_profiles@2021-06-01/main.bicep](../examples/Microsoft.Cdn_profiles@2021-06-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Cdn_profiles_afdEndpoints@2021-06-01/main.bicep](../examples/Microsoft.Cdn_profiles_afdEndpoints@2021-06-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Cdn_profiles_afdEndpoints_routes@2021-06-01/main.bicep](../examples/Microsoft.Cdn_profiles_afdEndpoints_routes@2021-06-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Cdn_profiles_customDomains@2021-06-01/main.bicep](../examples/Microsoft.Cdn_profiles_customDomains@2021-06-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Cdn_profiles_endpoints@2020-09-01/main.bicep](../examples/Microsoft.Cdn_profiles_endpoints@2020-09-01/main.bicep)
Result: success

## [Microsoft.Cdn_profiles_originGroups@2021-06-01/main.bicep](../examples/Microsoft.Cdn_profiles_originGroups@2021-06-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Cdn_profiles_originGroups_origins@2021-06-01/main.bicep](../examples/Microsoft.Cdn_profiles_originGroups_origins@2021-06-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Cdn_profiles_ruleSets@2021-06-01/main.bicep](../examples/Microsoft.Cdn_profiles_ruleSets@2021-06-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Cdn_profiles_ruleSets_rules@2024-09-01/main.bicep](../examples/Microsoft.Cdn_profiles_ruleSets_rules@2024-09-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Error BCP057] The name "substr" does not exist in the current context.
[Error BCP057] The name "var" does not exist in the current context.
[Error BCP057] The name "substr" does not exist in the current context.
[Error BCP057] The name "var" does not exist in the current context.
```
## [Microsoft.Cdn_profiles_securityPolicies@2021-06-01/main.bicep](../examples/Microsoft.Cdn_profiles_securityPolicies@2021-06-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.CertificateRegistration_certificateOrders@2021-02-01/main.bicep](../examples/Microsoft.CertificateRegistration_certificateOrders@2021-02-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.CodeSigning_codeSigningAccounts@2024-09-30-preview/main.bicep](../examples/Microsoft.CodeSigning_codeSigningAccounts@2024-09-30-preview/main.bicep)
Result: success

## [Microsoft.CognitiveServices_accounts@2022-10-01/main.bicep](../examples/Microsoft.CognitiveServices_accounts@2022-10-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "Identity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.CognitiveServices_accounts_deployments@2022-10-01/main.bicep](../examples/Microsoft.CognitiveServices_accounts_deployments@2022-10-01/main.bicep)
Result: success

## [Microsoft.CognitiveServices_accounts_raiBlocklists@2024-10-01/main.bicep](../examples/Microsoft.CognitiveServices_accounts_raiBlocklists@2024-10-01/main.bicep)
Result: success

## [Microsoft.CognitiveServices_accounts_raiPolicies@2024-10-01/main.bicep](../examples/Microsoft.CognitiveServices_accounts_raiPolicies@2024-10-01/main.bicep)
Result: success

## [Microsoft.Communication_communicationServices@2023-03-31/main.bicep](../examples/Microsoft.Communication_communicationServices@2023-03-31/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Communication_emailServices@2023-03-31/main.bicep](../examples/Microsoft.Communication_emailServices@2023-03-31/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Communication_emailServices_domains@2023-04-01-preview/main.bicep](../examples/Microsoft.Communication_emailServices_domains@2023-04-01-preview/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Communication_emailServices_domains_senderUsernames@2023-04-01-preview/main.bicep](../examples/Microsoft.Communication_emailServices_domains_senderUsernames@2023-04-01-preview/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Compute_availabilitySets@2021-11-01/main.bicep](../examples/Microsoft.Compute_availabilitySets@2021-11-01/main.bicep)
Result: success

## [Microsoft.Compute_capacityReservationGroups@2022-03-01/main.bicep](../examples/Microsoft.Compute_capacityReservationGroups@2022-03-01/main.bicep)
Result: success

## [Microsoft.Compute_capacityReservationGroups_capacityReservations@2022-03-01/main.bicep](../examples/Microsoft.Compute_capacityReservationGroups_capacityReservations@2022-03-01/main.bicep)
Result: success

## [Microsoft.Compute_diskAccesses@2022-03-02/main.bicep](../examples/Microsoft.Compute_diskAccesses@2022-03-02/main.bicep)
Result: success

## [Microsoft.Compute_diskEncryptionSets@2022-03-02/main.bicep](../examples/Microsoft.Compute_diskEncryptionSets@2022-03-02/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "azapi_resource_action" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Compute_disks@2022-03-02/main.bicep](../examples/Microsoft.Compute_disks@2022-03-02/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "osType" expected a value of type "'Linux' | 'Windows' | null" but the provided value is of type "''". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Compute_galleries@2022-03-03/main.bicep](../examples/Microsoft.Compute_galleries@2022-03-03/main.bicep)
Result: success

## [Microsoft.Compute_galleries_applications@2022-03-03/main.bicep](../examples/Microsoft.Compute_galleries_applications@2022-03-03/main.bicep)
Result: success

## [Microsoft.Compute_galleries_applications_versions@2022-03-03/main.bicep](../examples/Microsoft.Compute_galleries_applications_versions@2022-03-03/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP169] Expected resource name to contain 2 "/" character(s). The number of name segments must match the number of segments in the resource type.
[Warning no-hardcoded-env-urls] Environment URLs should not be hardcoded. Use the environment() function to ensure compatibility across clouds. Found this disallowed host: "core.windows.net"
```
## [Microsoft.Compute_galleries_images@2022-03-03/main.bicep](../examples/Microsoft.Compute_galleries_images@2022-03-03/main.bicep)
Result: success

## [Microsoft.Compute_hostGroups@2021-11-01/main.bicep](../examples/Microsoft.Compute_hostGroups@2021-11-01/main.bicep)
Result: success

## [Microsoft.Compute_hostGroups_hosts@2021-11-01/main.bicep](../examples/Microsoft.Compute_hostGroups_hosts@2021-11-01/main.bicep)
Result: success

## [Microsoft.Compute_proximityPlacementGroups@2022-03-01/main.bicep](../examples/Microsoft.Compute_proximityPlacementGroups@2022-03-01/main.bicep)
Result: success

## [Microsoft.Compute_restorePointCollections@2024-03-01/main.bicep](../examples/Microsoft.Compute_restorePointCollections@2024-03-01/main.bicep)
Result: success

Diagnostics:
```
[Warning adminusername-should-not-be-literal] Property 'adminUserName' should not use a literal value. Use a param instead. Found literal string value "adminuser"
```
## [Microsoft.Compute_restorePointCollections_restorePoints@2024-03-01/main.bicep](../examples/Microsoft.Compute_restorePointCollections_restorePoints@2024-03-01/main.bicep)
Result: success

Diagnostics:
```
[Warning adminusername-should-not-be-literal] Property 'adminUserName' should not use a literal value. Use a param instead. Found literal string value "adminuser"
```
## [Microsoft.Compute_snapshots@2022-03-02/main.bicep](../examples/Microsoft.Compute_snapshots@2022-03-02/main.bicep)
Result: success

## [Microsoft.Compute_sshPublicKeys@2021-11-01/main.bicep](../examples/Microsoft.Compute_sshPublicKeys@2021-11-01/main.bicep)
Result: success

## [Microsoft.Compute_virtualMachineScaleSets@2023-03-01/main.bicep](../examples/Microsoft.Compute_virtualMachineScaleSets@2023-03-01/main.bicep)
Result: success

Diagnostics:
```
[Warning adminusername-should-not-be-literal] Property 'adminUserName' should not use a literal value. Use a param instead. Found literal string value "adminuser"
```
## [Microsoft.Compute_virtualMachineScaleSets_extensions@2023-03-01/main.bicep](../examples/Microsoft.Compute_virtualMachineScaleSets_extensions@2023-03-01/main.bicep)
Result: success

Diagnostics:
```
[Warning adminusername-should-not-be-literal] Property 'adminUserName' should not use a literal value. Use a param instead. Found literal string value "adminuser"
```
## [Microsoft.Compute_virtualMachines_extensions@2023-03-01/main.bicep](../examples/Microsoft.Compute_virtualMachines_extensions@2023-03-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'vm_admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
[Warning adminusername-should-not-be-literal] Property 'adminUserName' should not use a literal value. Use a param instead. Found literal string value "testadmin"
```
## [Microsoft.Compute_virtualMachines_runCommands@2023-03-01/main.bicep](../examples/Microsoft.Compute_virtualMachines_runCommands@2023-03-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "VirtualMachineIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
[Warning adminusername-should-not-be-literal] Property 'adminUserName' should not use a literal value. Use a param instead. Found literal string value "adminuser"
```
## [Microsoft.ConfidentialLedger_ledgers@2022-05-13/main.bicep](../examples/Microsoft.ConfidentialLedger_ledgers@2022-05-13/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Consumption_budgets@2019-10-01/main.bicep](../examples/Microsoft.Consumption_budgets@2019-10-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.ContainerInstance_containerGroups@2023-05-01/main.bicep](../examples/Microsoft.ContainerInstance_containerGroups@2023-05-01/main.bicep)
Result: success

## [Microsoft.ContainerRegistry_registries@2021-08-01-preview/main.bicep](../examples/Microsoft.ContainerRegistry_registries@2021-08-01-preview/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "tier" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.ContainerRegistry_registries_agentPools@2019-06-01-preview/main.bicep](../examples/Microsoft.ContainerRegistry_registries_agentPools@2019-06-01-preview/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "tier" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.ContainerRegistry_registries_cacheRules@2023-07-01/main.bicep](../examples/Microsoft.ContainerRegistry_registries_cacheRules@2023-07-01/main.bicep)
Result: success

## [Microsoft.ContainerRegistry_registries_connectedRegistries@2023-11-01-preview/main.bicep](../examples/Microsoft.ContainerRegistry_registries_connectedRegistries@2023-11-01-preview/main.bicep)
Result: success

## [Microsoft.ContainerRegistry_registries_credentialSets@2023-07-01/main.bicep](../examples/Microsoft.ContainerRegistry_registries_credentialSets@2023-07-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "IdentityProperties | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Warning no-hardcoded-env-urls] Environment URLs should not be hardcoded. Use the environment() function to ensure compatibility across clouds. Found this disallowed host: "vault.azure.net"
[Warning no-hardcoded-env-urls] Environment URLs should not be hardcoded. Use the environment() function to ensure compatibility across clouds. Found this disallowed host: "vault.azure.net"
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.ContainerRegistry_registries_scopeMaps@2021-08-01-preview/main.bicep](../examples/Microsoft.ContainerRegistry_registries_scopeMaps@2021-08-01-preview/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "tier" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.ContainerRegistry_registries_taskRuns@2019-06-01-preview/main.bicep](../examples/Microsoft.ContainerRegistry_registries_taskRuns@2019-06-01-preview/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "tier" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.ContainerRegistry_registries_tasks@2019-06-01-preview/main.bicep](../examples/Microsoft.ContainerRegistry_registries_tasks@2019-06-01-preview/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "tier" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.ContainerRegistry_registries_tokens@2021-08-01-preview/main.bicep](../examples/Microsoft.ContainerRegistry_registries_tokens@2021-08-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning BCP073] The property "tier" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.ContainerRegistry_registries_webHooks@2021-08-01-preview/main.bicep](../examples/Microsoft.ContainerRegistry_registries_webHooks@2021-08-01-preview/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "tier" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Warning use-secure-value-for-secure-inputs] Property 'serviceUri' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.ContainerService_fleets@2024-04-01/main.bicep](../examples/Microsoft.ContainerService_fleets@2024-04-01/main.bicep)
Result: success

## [Microsoft.ContainerService_fleets_members@2024-04-01/main.bicep](../examples/Microsoft.ContainerService_fleets_members@2024-04-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "ManagedClusterIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Warning BCP328] The provided value (which will always be less than or equal to 0) is too small to assign to a target for which the minimum allowable value is 1.
```
## [Microsoft.ContainerService_managedClusters@2023-04-02-preview/main.bicep](../examples/Microsoft.ContainerService_managedClusters@2023-04-02-preview/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "ManagedClusterIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.ContainerService_managedClusters_agentPools@2023-04-02-preview/main.bicep](../examples/Microsoft.ContainerService_managedClusters_agentPools@2023-04-02-preview/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "ManagedClusterIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.CostManagement_scheduledActions@2022-06-01-preview/main.bicep](../examples/Microsoft.CostManagement_scheduledActions@2022-06-01-preview/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.CostManagement_scheduledActions@2022-10-01/main.bicep](../examples/Microsoft.CostManagement_scheduledActions@2022-10-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.CostManagement_views@2022-10-01/main.bicep](../examples/Microsoft.CostManagement_views@2022-10-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.CustomProviders_resourceProviders@2018-09-01-preview/main.bicep](../examples/Microsoft.CustomProviders_resourceProviders@2018-09-01-preview/main.bicep)
Result: success

## [Microsoft.Dashboard_grafana@2022-08-01/main.bicep](../examples/Microsoft.Dashboard_grafana@2022-08-01/main.bicep)
Result: success

## [Microsoft.DataBoxEdge_dataBoxEdgeDevices@2022-03-01/main.bicep](../examples/Microsoft.DataBoxEdge_dataBoxEdgeDevices@2022-03-01/main.bicep)
Result: success

## [Microsoft.Databricks_accessConnectors@2022-10-01-preview/main.bicep](../examples/Microsoft.Databricks_accessConnectors@2022-10-01-preview/main.bicep)
Result: success

## [Microsoft.Databricks_workspaces@2023-02-01/main.bicep](../examples/Microsoft.Databricks_workspaces@2023-02-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Databricks_workspaces_virtualNetworkPeerings@2023-02-01/main.bicep](../examples/Microsoft.Databricks_workspaces_virtualNetworkPeerings@2023-02-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.DataFactory_factories@2018-06-01/main.bicep](../examples/Microsoft.DataFactory_factories@2018-06-01/main.bicep)
Result: success

## [Microsoft.DataFactory_factories_credentials@2018-06-01/main.bicep](../examples/Microsoft.DataFactory_factories_credentials@2018-06-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "FactoryIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.DataFactory_factories_dataflows@2018-06-01/main.bicep](../examples/Microsoft.DataFactory_factories_dataflows@2018-06-01/main.bicep)
Result: success

## [Microsoft.DataFactory_factories_datasets@2018-06-01/main.bicep](../examples/Microsoft.DataFactory_factories_datasets@2018-06-01/main.bicep)
Result: success

## [Microsoft.DataFactory_factories_integrationRuntimes@2018-06-01/main.bicep](../examples/Microsoft.DataFactory_factories_integrationRuntimes@2018-06-01/main.bicep)
Result: success

## [Microsoft.DataFactory_factories_linkedservices@2018-06-01/main.bicep](../examples/Microsoft.DataFactory_factories_linkedservices@2018-06-01/main.bicep)
Result: success

## [Microsoft.DataFactory_factories_managedVirtualNetworks@2018-06-01/main.bicep](../examples/Microsoft.DataFactory_factories_managedVirtualNetworks@2018-06-01/main.bicep)
Result: success

## [Microsoft.DataFactory_factories_managedVirtualNetworks_managedPrivateEndpoints@2018-06-01/main.bicep](../examples/Microsoft.DataFactory_factories_managedVirtualNetworks_managedPrivateEndpoints@2018-06-01/main.bicep)
Result: success

## [Microsoft.DataFactory_factories_pipelines@2018-06-01/main.bicep](../examples/Microsoft.DataFactory_factories_pipelines@2018-06-01/main.bicep)
Result: success

## [Microsoft.DataFactory_factories_triggers@2018-06-01/main.bicep](../examples/Microsoft.DataFactory_factories_triggers@2018-06-01/main.bicep)
Result: success

## [Microsoft.DataMigration_services@2018-04-19/main.bicep](../examples/Microsoft.DataMigration_services@2018-04-19/main.bicep)
Result: success

## [Microsoft.DataMigration_services_projects@2018-04-19/main.bicep](../examples/Microsoft.DataMigration_services_projects@2018-04-19/main.bicep)
Result: success

## [Microsoft.DataProtection_backupVaults@2022-04-01/main.bicep](../examples/Microsoft.DataProtection_backupVaults@2022-04-01/main.bicep)
Result: success

## [Microsoft.DataProtection_backupVaults_backupInstances@2022-04-01/main.bicep](../examples/Microsoft.DataProtection_backupVaults_backupInstances@2022-04-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning BCP187] The property "location" does not exist in the resource or type definition, although it might still be valid. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Warning BCP073] The property "id" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Warning BCP036] The property "identity" expected a value of type "DppIdentityDetails | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DataProtection_backupVaults_backupPolicies@2022-04-01/main.bicep](../examples/Microsoft.DataProtection_backupVaults_backupPolicies@2022-04-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "id" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.DataProtection_resourceGuards@2022-04-01/main.bicep](../examples/Microsoft.DataProtection_resourceGuards@2022-04-01/main.bicep)
Result: success

## [Microsoft.DataShare_accounts@2019-11-01/main.bicep](../examples/Microsoft.DataShare_accounts@2019-11-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "Identity" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.DataShare_accounts_shares@2019-11-01/main.bicep](../examples/Microsoft.DataShare_accounts_shares@2019-11-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "Identity" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.DBforMariaDB_servers@2018-06-01/main.bicep](../examples/Microsoft.DBforMariaDB_servers@2018-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforMariaDB_servers_configurations@2018-06-01/main.bicep](../examples/Microsoft.DBforMariaDB_servers_configurations@2018-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforMariaDB_servers_databases@2018-06-01/main.bicep](../examples/Microsoft.DBforMariaDB_servers_databases@2018-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforMariaDB_servers_firewallRules@2018-06-01/main.bicep](../examples/Microsoft.DBforMariaDB_servers_firewallRules@2018-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforMariaDB_servers_virtualNetworkRules@2018-06-01/main.bicep](../examples/Microsoft.DBforMariaDB_servers_virtualNetworkRules@2018-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforMySQL_flexibleServers@2021-05-01/main.bicep](../examples/Microsoft.DBforMySQL_flexibleServers@2021-05-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforMySQL_flexibleServers_databases@2021-05-01/main.bicep](../examples/Microsoft.DBforMySQL_flexibleServers_databases@2021-05-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforMySQL_flexibleServers_firewallRules@2021-05-01/main.bicep](../examples/Microsoft.DBforMySQL_flexibleServers_firewallRules@2021-05-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'mysql_administrator_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforMySQL_servers@2017-12-01/main.bicep](../examples/Microsoft.DBforMySQL_servers@2017-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforMySQL_servers_administrators@2017-12-01/main.bicep](../examples/Microsoft.DBforMySQL_servers_administrators@2017-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforMySQL_servers_configurations@2017-12-01/main.bicep](../examples/Microsoft.DBforMySQL_servers_configurations@2017-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforMySQL_servers_databases@2017-12-01/main.bicep](../examples/Microsoft.DBforMySQL_servers_databases@2017-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforMySQL_servers_firewallRules@2017-12-01/main.bicep](../examples/Microsoft.DBforMySQL_servers_firewallRules@2017-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforMySQL_servers_virtualNetworkRules@2017-12-01/main.bicep](../examples/Microsoft.DBforMySQL_servers_virtualNetworkRules@2017-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforPostgreSQL_flexibleServers@2023-06-01-preview/main.bicep](../examples/Microsoft.DBforPostgreSQL_flexibleServers@2023-06-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforPostgreSQL_flexibleServers_administrators@2022-12-01/main.bicep](../examples/Microsoft.DBforPostgreSQL_flexibleServers_administrators@2022-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.DBforPostgreSQL_flexibleServers_configurations@2022-12-01/main.bicep](../examples/Microsoft.DBforPostgreSQL_flexibleServers_configurations@2022-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'postgresql_administrator_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforPostgreSQL_flexibleServers_databases@2022-12-01/main.bicep](../examples/Microsoft.DBforPostgreSQL_flexibleServers_databases@2022-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'postgresql_administrator_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforPostgreSQL_flexibleServers_firewallRules@2022-12-01/main.bicep](../examples/Microsoft.DBforPostgreSQL_flexibleServers_firewallRules@2022-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'postgresql_administrator_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforPostgreSQL_serverGroupsv2@2022-11-08/main.bicep](../examples/Microsoft.DBforPostgreSQL_serverGroupsv2@2022-11-08/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforPostgreSQL_servers@2017-12-01/main.bicep](../examples/Microsoft.DBforPostgreSQL_servers@2017-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforPostgreSQL_servers_administrators@2017-12-01/main.bicep](../examples/Microsoft.DBforPostgreSQL_servers_administrators@2017-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforPostgreSQL_servers_configurations@2017-12-01/main.bicep](../examples/Microsoft.DBforPostgreSQL_servers_configurations@2017-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforPostgreSQL_servers_databases@2017-12-01/main.bicep](../examples/Microsoft.DBforPostgreSQL_servers_databases@2017-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforPostgreSQL_servers_firewallRules@2017-12-01/main.bicep](../examples/Microsoft.DBforPostgreSQL_servers_firewallRules@2017-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DBforPostgreSQL_servers_virtualNetworkRules@2017-12-01/main.bicep](../examples/Microsoft.DBforPostgreSQL_servers_virtualNetworkRules@2017-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.DesktopVirtualization_applicationGroups_applications@2023-09-05/main.bicep](../examples/Microsoft.DesktopVirtualization_applicationGroups_applications@2023-09-05/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP187] The property "location" does not exist in the resource or type definition, although it might still be valid. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.DesktopVirtualization_scalingPlans_personalSchedules@2023-11-01-preview/main.bicep](../examples/Microsoft.DesktopVirtualization_scalingPlans_personalSchedules@2023-11-01-preview/main.bicep)
Result: success

## [Microsoft.DevCenter_devcenters_attachednetworks@2023-04-01/main.bicep](../examples/Microsoft.DevCenter_devcenters_attachednetworks@2023-04-01/main.bicep)
Result: success

## [Microsoft.DevCenter_devcenters_devboxdefinitions@2024-10-01-preview/main.bicep](../examples/Microsoft.DevCenter_devcenters_devboxdefinitions@2024-10-01-preview/main.bicep)
Result: success

## [Microsoft.DevCenter_networkConnections@2023-04-01/main.bicep](../examples/Microsoft.DevCenter_networkConnections@2023-04-01/main.bicep)
Result: success

## [Microsoft.Devices_IotHubs@2022-04-30-preview/main.bicep](../examples/Microsoft.Devices_IotHubs@2022-04-30-preview/main.bicep)
Result: success

## [Microsoft.Devices_IotHubs_certificates@2022-04-30-preview/main.bicep](../examples/Microsoft.Devices_IotHubs_certificates@2022-04-30-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Devices_provisioningServices@2022-02-05/main.bicep](../examples/Microsoft.Devices_provisioningServices@2022-02-05/main.bicep)
Result: success

## [Microsoft.Devices_provisioningServices_certificates@2022-02-05/main.bicep](../examples/Microsoft.Devices_provisioningServices_certificates@2022-02-05/main.bicep)
Result: success

## [Microsoft.DeviceUpdate_accounts@2022-10-01/main.bicep](../examples/Microsoft.DeviceUpdate_accounts@2022-10-01/main.bicep)
Result: success

## [Microsoft.DeviceUpdate_accounts_instances@2022-10-01/main.bicep](../examples/Microsoft.DeviceUpdate_accounts_instances@2022-10-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "accountName" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.DevTestLab_labs@2018-09-15/main.bicep](../examples/Microsoft.DevTestLab_labs@2018-09-15/main.bicep)
Result: success

## [Microsoft.DevTestLab_labs_schedules@2018-09-15/main.bicep](../examples/Microsoft.DevTestLab_labs_schedules@2018-09-15/main.bicep)
Result: success

## [Microsoft.DevTestLab_labs_virtualMachines@2018-09-15/main.bicep](../examples/Microsoft.DevTestLab_labs_virtualMachines@2018-09-15/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'vm_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.DevTestLab_labs_virtualNetworks@2018-09-15/main.bicep](../examples/Microsoft.DevTestLab_labs_virtualNetworks@2018-09-15/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.DevTestLab_schedules@2018-09-15/main.bicep](../examples/Microsoft.DevTestLab_schedules@2018-09-15/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
[Warning adminusername-should-not-be-literal] Property 'adminUserName' should not use a literal value. Use a param instead. Found literal string value "testadmin"
```
## [Microsoft.DigitalTwins_digitalTwinsInstances@2020-12-01/main.bicep](../examples/Microsoft.DigitalTwins_digitalTwinsInstances@2020-12-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "DigitalTwinsIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.DigitalTwins_digitalTwinsInstances_endpoints@2020-12-01/main.bicep](../examples/Microsoft.DigitalTwins_digitalTwinsInstances_endpoints@2020-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.DocumentDB_databaseAccounts@2021-10-15/main.bicep](../examples/Microsoft.DocumentDB_databaseAccounts@2021-10-15/main.bicep)
Result: success

## [Microsoft.DocumentDB_databaseAccounts_cassandraKeyspaces@2021-10-15/main.bicep](../examples/Microsoft.DocumentDB_databaseAccounts_cassandraKeyspaces@2021-10-15/main.bicep)
Result: success

## [Microsoft.DocumentDB_databaseAccounts_gremlinDatabases@2023-04-15/main.bicep](../examples/Microsoft.DocumentDB_databaseAccounts_gremlinDatabases@2023-04-15/main.bicep)
Result: success

## [Microsoft.DocumentDB_databaseAccounts_gremlinDatabases_graphs@2023-04-15/main.bicep](../examples/Microsoft.DocumentDB_databaseAccounts_gremlinDatabases_graphs@2023-04-15/main.bicep)
Result: success

## [Microsoft.DocumentDB_databaseAccounts_mongodbDatabases@2021-10-15/main.bicep](../examples/Microsoft.DocumentDB_databaseAccounts_mongodbDatabases@2021-10-15/main.bicep)
Result: success

## [Microsoft.DocumentDB_databaseAccounts_services@2022-05-15/main.bicep](../examples/Microsoft.DocumentDB_databaseAccounts_services@2022-05-15/main.bicep)
Result: success

## [Microsoft.DocumentDB_databaseAccounts_sqlDatabases@2021-10-15/main.bicep](../examples/Microsoft.DocumentDB_databaseAccounts_sqlDatabases@2021-10-15/main.bicep)
Result: success

## [Microsoft.DocumentDB_databaseAccounts_sqlDatabases_containers@2023-04-15/main.bicep](../examples/Microsoft.DocumentDB_databaseAccounts_sqlDatabases_containers@2023-04-15/main.bicep)
Result: success

## [Microsoft.DocumentDB_databaseAccounts_sqlDatabases_containers_storedProcedures@2021-10-15/main.bicep](../examples/Microsoft.DocumentDB_databaseAccounts_sqlDatabases_containers_storedProcedures@2021-10-15/main.bicep)
Result: success

## [Microsoft.DocumentDB_databaseAccounts_sqlDatabases_containers_triggers@2021-10-15/main.bicep](../examples/Microsoft.DocumentDB_databaseAccounts_sqlDatabases_containers_triggers@2021-10-15/main.bicep)
Result: success

## [Microsoft.DocumentDB_databaseAccounts_sqlDatabases_containers_userDefinedFunctions@2021-10-15/main.bicep](../examples/Microsoft.DocumentDB_databaseAccounts_sqlDatabases_containers_userDefinedFunctions@2021-10-15/main.bicep)
Result: success

## [Microsoft.DocumentDB_databaseAccounts_sqlRoleAssignments@2021-10-15/main.bicep](../examples/Microsoft.DocumentDB_databaseAccounts_sqlRoleAssignments@2021-10-15/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "Identity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.DocumentDB_databaseAccounts_sqlRoleDefinitions@2021-10-15/main.bicep](../examples/Microsoft.DocumentDB_databaseAccounts_sqlRoleDefinitions@2021-10-15/main.bicep)
Result: success

## [Microsoft.DocumentDB_databaseAccounts_tables@2021-10-15/main.bicep](../examples/Microsoft.DocumentDB_databaseAccounts_tables@2021-10-15/main.bicep)
Result: success

## [Microsoft.EventGrid_domains@2021-12-01/main.bicep](../examples/Microsoft.EventGrid_domains@2021-12-01/main.bicep)
Result: success

## [Microsoft.EventGrid_domains_topics@2021-12-01/main.bicep](../examples/Microsoft.EventGrid_domains_topics@2021-12-01/main.bicep)
Result: success

## [Microsoft.EventGrid_eventSubscriptions@2021-12-01/main.bicep](../examples/Microsoft.EventGrid_eventSubscriptions@2021-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.EventGrid/eventSubscriptions". Permissible properties include "asserts", "dependsOn", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.EventGrid_systemTopics@2021-12-01/main.bicep](../examples/Microsoft.EventGrid_systemTopics@2021-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Error BCP063] The name "resourceGroup" is not a parameter, variable, resource or module.
```
## [Microsoft.EventGrid_topics@2021-12-01/main.bicep](../examples/Microsoft.EventGrid_topics@2021-12-01/main.bicep)
Result: success

## [Microsoft.EventHub_clusters@2021-11-01/main.bicep](../examples/Microsoft.EventHub_clusters@2021-11-01/main.bicep)
Result: success

## [Microsoft.EventHub_namespaces@2022-01-01-preview/main.bicep](../examples/Microsoft.EventHub_namespaces@2022-01-01-preview/main.bicep)
Result: success

## [Microsoft.EventHub_namespaces_authorizationRules@2021-11-01/main.bicep](../examples/Microsoft.EventHub_namespaces_authorizationRules@2021-11-01/main.bicep)
Result: success

## [Microsoft.EventHub_namespaces_disasterRecoveryConfigs@2021-11-01/main.bicep](../examples/Microsoft.EventHub_namespaces_disasterRecoveryConfigs@2021-11-01/main.bicep)
Result: success

## [Microsoft.EventHub_namespaces_eventhubs@2021-11-01/main.bicep](../examples/Microsoft.EventHub_namespaces_eventhubs@2021-11-01/main.bicep)
Result: success

## [Microsoft.EventHub_namespaces_eventhubs_authorizationRules@2021-11-01/main.bicep](../examples/Microsoft.EventHub_namespaces_eventhubs_authorizationRules@2021-11-01/main.bicep)
Result: success

## [Microsoft.EventHub_namespaces_eventhubs_consumerGroups@2021-11-01/main.bicep](../examples/Microsoft.EventHub_namespaces_eventhubs_consumerGroups@2021-11-01/main.bicep)
Result: success

## [Microsoft.EventHub_namespaces_schemaGroups@2021-11-01/main.bicep](../examples/Microsoft.EventHub_namespaces_schemaGroups@2021-11-01/main.bicep)
Result: success

## [Microsoft.FluidRelay_fluidRelayServers@2022-05-26/main.bicep](../examples/Microsoft.FluidRelay_fluidRelayServers@2022-05-26/main.bicep)
Result: success

## [Microsoft.GuestConfiguration_guestConfigurationAssignments@2020-06-25/main.bicep](../examples/Microsoft.GuestConfiguration_guestConfigurationAssignments@2020-06-25/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.GuestConfiguration/guestConfigurationAssignments". Permissible properties include "asserts", "dependsOn", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
[Warning adminusername-should-not-be-literal] Property 'adminUserName' should not use a literal value. Use a param instead. Found literal string value "adminuser"
```
## [Microsoft.HDInsight_clusters@2018-06-01-preview/main.bicep](../examples/Microsoft.HDInsight_clusters@2018-06-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'rest_credential_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'vm_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "data" does not exist in the current context.
[Warning no-hardcoded-env-urls] Environment URLs should not be hardcoded. Use the environment() function to ensure compatibility across clouds. Found this disallowed host: "core.windows.net"
```
## [Microsoft.HealthBot_healthBots@2022-08-08/main.bicep](../examples/Microsoft.HealthBot_healthBots@2022-08-08/main.bicep)
Result: success

## [Microsoft.HealthcareApis_services@2022-12-01/main.bicep](../examples/Microsoft.HealthcareApis_services@2022-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.HealthcareApis_workspaces@2022-12-01/main.bicep](../examples/Microsoft.HealthcareApis_workspaces@2022-12-01/main.bicep)
Result: success

## [Microsoft.HealthcareApis_workspaces_dicomServices@2022-12-01/main.bicep](../examples/Microsoft.HealthcareApis_workspaces_dicomServices@2022-12-01/main.bicep)
Result: success

## [Microsoft.HealthcareApis_workspaces_fhirServices@2022-12-01/main.bicep](../examples/Microsoft.HealthcareApis_workspaces_fhirServices@2022-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-hardcoded-env-urls] Environment URLs should not be hardcoded. Use the environment() function to ensure compatibility across clouds. Found this disallowed host: "login.microsoftonline.com"
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.HealthcareApis_workspaces_iotConnectors@2022-12-01/main.bicep](../examples/Microsoft.HealthcareApis_workspaces_iotConnectors@2022-12-01/main.bicep)
Result: success

## [Microsoft.HealthcareApis_workspaces_iotConnectors_fhirDestinations@2022-12-01/main.bicep](../examples/Microsoft.HealthcareApis_workspaces_iotConnectors_fhirDestinations@2022-12-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-hardcoded-env-urls] Environment URLs should not be hardcoded. Use the environment() function to ensure compatibility across clouds. Found this disallowed host: "login.microsoftonline.com"
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.HybridCompute_privateLinkScopes@2022-11-10/main.bicep](../examples/Microsoft.HybridCompute_privateLinkScopes@2022-11-10/main.bicep)
Result: success

## [Microsoft.Impact_connectors@2024-05-01-preview/main.bicep](../examples/Microsoft.Impact_connectors@2024-05-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Error BCP135] Scope "resourceGroup" is not valid for this resource type. Permitted scopes: "subscription".
[Warning BCP035] The specified "object" declaration is missing the following required properties: "connectorId", "lastRunTimeStamp", "tenantId". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Impact_workloadImpacts@2023-12-01-preview/main.bicep](../examples/Microsoft.Impact_workloadImpacts@2023-12-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
[Warning adminusername-should-not-be-literal] Property 'adminUserName' should not use a literal value. Use a param instead. Found literal string value "testadmin"
[Warning BCP081] Resource type "Microsoft.Impact/workloadImpacts@2023-12-01-preview" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
```
## [Microsoft.Insights_actionGroups@2023-01-01/main.bicep](../examples/Microsoft.Insights_actionGroups@2023-01-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Insights_activityLogAlerts@2020-10-01/main.bicep](../examples/Microsoft.Insights_activityLogAlerts@2020-10-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP063] The name "resourceGroup" is not a parameter, variable, resource or module.
```
## [Microsoft.Insights_autoScaleSettings@2022-10-01/main.bicep](../examples/Microsoft.Insights_autoScaleSettings@2022-10-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.Insights_components@2020-02-02/main.bicep](../examples/Microsoft.Insights_components@2020-02-02/main.bicep)
Result: success

## [microsoft.insights_components_analyticsItems@2015-05-01/main.bicep](../examples/microsoft.insights_components_analyticsItems@2015-05-01/main.bicep)
Result: success

## [Microsoft.Insights_components_ProactiveDetectionConfigs@2015-05-01/main.bicep](../examples/Microsoft.Insights_components_ProactiveDetectionConfigs@2015-05-01/main.bicep)
Result: success

## [Microsoft.Insights_dataCollectionEndpoints@2022-06-01/main.bicep](../examples/Microsoft.Insights_dataCollectionEndpoints@2022-06-01/main.bicep)
Result: success

## [Microsoft.Insights_dataCollectionRuleAssociations@2022-06-01/main.bicep](../examples/Microsoft.Insights_dataCollectionRuleAssociations@2022-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.Insights/dataCollectionRuleAssociations". Permissible properties include "asserts", "dependsOn", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
[Warning adminusername-should-not-be-literal] Property 'adminUserName' should not use a literal value. Use a param instead. Found literal string value "adminuser"
```
## [Microsoft.Insights_dataCollectionRules@2022-06-01/main.bicep](../examples/Microsoft.Insights_dataCollectionRules@2022-06-01/main.bicep)
Result: success

## [Microsoft.Insights_diagnosticSettings@2021-05-01-preview/main.bicep](../examples/Microsoft.Insights_diagnosticSettings@2021-05-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.Insights/diagnosticSettings". Permissible properties include "asserts", "dependsOn", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Insights_logProfiles@2016-03-01/main.bicep](../examples/Microsoft.Insights_logProfiles@2016-03-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP135] Scope "resourceGroup" is not valid for this resource type. Permitted scopes: "subscription".
```
## [Microsoft.Insights_metricAlerts@2018-03-01/main.bicep](../examples/Microsoft.Insights_metricAlerts@2018-03-01/main.bicep)
Result: success

## [Microsoft.Insights_privateLinkScopes@2019-10-17-preview/main.bicep](../examples/Microsoft.Insights_privateLinkScopes@2019-10-17-preview/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Insights_privateLinkScopes_scopedResources@2019-10-17-preview/main.bicep](../examples/Microsoft.Insights_privateLinkScopes_scopedResources@2019-10-17-preview/main.bicep)
Result: success

## [Microsoft.Insights_scheduledQueryRules@2018-04-16/main.bicep](../examples/Microsoft.Insights_scheduledQueryRules@2018-04-16/main.bicep)
Result: success

## [Microsoft.Insights_scheduledQueryRules@2021-08-01/main.bicep](../examples/Microsoft.Insights_scheduledQueryRules@2021-08-01/main.bicep)
Result: success

## [Microsoft.Insights_webTests@2022-06-15/main.bicep](../examples/Microsoft.Insights_webTests@2022-06-15/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "azapi_resource" does not exist in the current context.
```
## [Microsoft.Insights_workbooks@2022-04-01/main.bicep](../examples/Microsoft.Insights_workbooks@2022-04-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "resource_name" is declared but never used.
```
## [Microsoft.Insights_workbookTemplates@2020-11-20/main.bicep](../examples/Microsoft.Insights_workbookTemplates@2020-11-20/main.bicep)
Result: success

## [Microsoft.IoTCentral_iotApps@2021-11-01-preview/main.bicep](../examples/Microsoft.IoTCentral_iotApps@2021-11-01-preview/main.bicep)
Result: success

## [Microsoft.KeyVault_managedHSMs@2021-10-01/main.bicep](../examples/Microsoft.KeyVault_managedHSMs@2021-10-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "resource_name" is declared but never used.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.KeyVault_vaults@2021-10-01/main.bicep](../examples/Microsoft.KeyVault_vaults@2021-10-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.KeyVault_vaults_accessPolicies@2023-02-01/main.bicep](../examples/Microsoft.KeyVault_vaults_accessPolicies@2023-02-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.KeyVault_vaults_keys@2023-02-01/main.bicep](../examples/Microsoft.KeyVault_vaults_keys@2023-02-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.KeyVault_vaults_secrets@2023-02-01/main.bicep](../examples/Microsoft.KeyVault_vaults_secrets@2023-02-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.KubernetesConfiguration_extensions@2022-11-01/main.bicep](../examples/Microsoft.KubernetesConfiguration_extensions@2022-11-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.KubernetesConfiguration/extensions". Permissible properties include "asserts", "dependsOn", "identity", "plan", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Warning BCP036] The property "identity" expected a value of type "ManagedClusterIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.KubernetesConfiguration_fluxConfigurations@2022-11-01/main.bicep](../examples/Microsoft.KubernetesConfiguration_fluxConfigurations@2022-11-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.KubernetesConfiguration/extensions". Permissible properties include "asserts", "dependsOn", "identity", "plan", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.KubernetesConfiguration/fluxConfigurations". Permissible properties include "asserts", "dependsOn", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Warning BCP036] The property "identity" expected a value of type "ManagedClusterIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Kusto_clusters@2022-12-29/main.bicep](../examples/Microsoft.Kusto_clusters@2022-12-29/main.bicep)
Result: success

## [Microsoft.Kusto_clusters@2023-05-02/main.bicep](../examples/Microsoft.Kusto_clusters@2023-05-02/main.bicep)
Result: success

## [Microsoft.Kusto_clusters_databases@2022-12-29/main.bicep](../examples/Microsoft.Kusto_clusters_databases@2022-12-29/main.bicep)
Result: success

## [Microsoft.Kusto_clusters_databases@2023-05-02/main.bicep](../examples/Microsoft.Kusto_clusters_databases@2023-05-02/main.bicep)
Result: success

## [Microsoft.Kusto_clusters_databases_principalAssignments@2023-05-02/main.bicep](../examples/Microsoft.Kusto_clusters_databases_principalAssignments@2023-05-02/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Kusto_clusters_databases_scripts@2022-12-29/main.bicep](../examples/Microsoft.Kusto_clusters_databases_scripts@2022-12-29/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "Identity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Warning use-secure-value-for-secure-inputs] Property 'scriptContent' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.Kusto_clusters_databases_scripts@2023-05-02/main.bicep](../examples/Microsoft.Kusto_clusters_databases_scripts@2023-05-02/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "Identity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Warning use-secure-value-for-secure-inputs] Property 'scriptContent' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.Kusto_clusters_managedPrivateEndpoints@2023-05-02/main.bicep](../examples/Microsoft.Kusto_clusters_managedPrivateEndpoints@2023-05-02/main.bicep)
Result: success

## [Microsoft.Kusto_clusters_principalAssignments@2023-05-02/main.bicep](../examples/Microsoft.Kusto_clusters_principalAssignments@2023-05-02/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.LabServices_labPlans@2022-08-01/main.bicep](../examples/Microsoft.LabServices_labPlans@2022-08-01/main.bicep)
Result: success

## [Microsoft.LabServices_labs@2022-08-01/main.bicep](../examples/Microsoft.LabServices_labs@2022-08-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'password' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.LabServices_labs_schedules@2022-08-01/main.bicep](../examples/Microsoft.LabServices_labs_schedules@2022-08-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'password' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.LabServices_labs_users@2022-08-01/main.bicep](../examples/Microsoft.LabServices_labs_users@2022-08-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'password' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.LoadTestService_loadTests@2022-12-01/main.bicep](../examples/Microsoft.LoadTestService_loadTests@2022-12-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "ManagedServiceIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Logic_integrationAccounts@2019-05-01/main.bicep](../examples/Microsoft.Logic_integrationAccounts@2019-05-01/main.bicep)
Result: success

## [Microsoft.Logic_integrationAccounts_agreements@2019-05-01/main.bicep](../examples/Microsoft.Logic_integrationAccounts_agreements@2019-05-01/main.bicep)
Result: success

## [Microsoft.Logic_integrationAccounts_batchConfigurations@2019-05-01/main.bicep](../examples/Microsoft.Logic_integrationAccounts_batchConfigurations@2019-05-01/main.bicep)
Result: success

## [Microsoft.Logic_integrationAccounts_maps@2019-05-01/main.bicep](../examples/Microsoft.Logic_integrationAccounts_maps@2019-05-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-hardcoded-env-urls] Environment URLs should not be hardcoded. Use the environment() function to ensure compatibility across clouds. Found this disallowed host: "core.windows.net"
```
## [Microsoft.Logic_integrationAccounts_partners@2019-05-01/main.bicep](../examples/Microsoft.Logic_integrationAccounts_partners@2019-05-01/main.bicep)
Result: success

## [Microsoft.Logic_integrationAccounts_schemas@2019-05-01/main.bicep](../examples/Microsoft.Logic_integrationAccounts_schemas@2019-05-01/main.bicep)
Result: success

## [Microsoft.Logic_integrationAccounts_sessions@2019-05-01/main.bicep](../examples/Microsoft.Logic_integrationAccounts_sessions@2019-05-01/main.bicep)
Result: success

## [Microsoft.Logic_workflows@2019-05-01/main.bicep](../examples/Microsoft.Logic_workflows@2019-05-01/main.bicep)
Result: success

## [Microsoft.MachineLearningServices_workspaces@2022-05-01/main.bicep](../examples/Microsoft.MachineLearningServices_workspaces@2022-05-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Warning BCP036] The property "identity" expected a value of type "ManagedServiceIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.MachineLearningServices_workspaces_codes_versions@2024-10-01/main.bicep](../examples/Microsoft.MachineLearningServices_workspaces_codes_versions@2024-10-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP169] Expected resource name to contain 2 "/" character(s). The number of name segments must match the number of segments in the resource type.
[Error BCP025] The property "properties" is declared multiple times in this object. Remove or rename the duplicate properties.
[Error BCP001] The following token is not recognized: "$".
[Error BCP236] Expected a new line or comma character at this location.
[Error BCP022] Expected a property name at this location.
[Error BCP018] Expected the ":" character at this location.
[Error BCP025] The property "properties" is declared multiple times in this object. Remove or rename the duplicate properties.
[Error BCP007] This declaration type is not recognized. Specify a metadata, parameter, variable, resource, or output declaration.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Warning BCP036] The property "identity" expected a value of type "ManagedServiceIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.MachineLearningServices_workspaces_computes@2022-05-01/main.bicep](../examples/Microsoft.MachineLearningServices_workspaces_computes@2022-05-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
[Warning BCP036] The property "identity" expected a value of type "ManagedServiceIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Maintenance_configurationAssignments@2022-07-01-preview/main.bicep](../examples/Microsoft.Maintenance_configurationAssignments@2022-07-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.Maintenance/configurationAssignments". Permissible properties include "asserts", "dependsOn", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
[Warning adminusername-should-not-be-literal] Property 'adminUserName' should not use a literal value. Use a param instead. Found literal string value "adminuser"
```
## [Microsoft.Maintenance_maintenanceConfigurations@2022-07-01-preview/main.bicep](../examples/Microsoft.Maintenance_maintenanceConfigurations@2022-07-01-preview/main.bicep)
Result: success

## [Microsoft.ManagedIdentity_userAssignedIdentities@2023-01-31/main.bicep](../examples/Microsoft.ManagedIdentity_userAssignedIdentities@2023-01-31/main.bicep)
Result: success

## [Microsoft.ManagedIdentity_userAssignedIdentities_federatedIdentityCredentials@2023-01-31/main.bicep](../examples/Microsoft.ManagedIdentity_userAssignedIdentities_federatedIdentityCredentials@2023-01-31/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP187] The property "location" does not exist in the resource or type definition, although it might still be valid. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Maps_accounts@2021-02-01/main.bicep](../examples/Microsoft.Maps_accounts@2021-02-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Maps_accounts_creators@2021-02-01/main.bicep](../examples/Microsoft.Maps_accounts_creators@2021-02-01/main.bicep)
Result: success

## [Microsoft.Media_mediaServices@2021-11-01/main.bicep](../examples/Microsoft.Media_mediaServices@2021-11-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP081] Resource type "Microsoft.Media/mediaServices@2021-11-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
```
## [Microsoft.Media_mediaServices_accountFilters@2022-08-01/main.bicep](../examples/Microsoft.Media_mediaServices_accountFilters@2022-08-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP081] Resource type "Microsoft.Media/mediaServices/accountFilters@2022-08-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
[Warning BCP081] Resource type "Microsoft.Media/mediaServices@2021-11-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
```
## [Microsoft.Media_mediaServices_assets@2022-08-01/main.bicep](../examples/Microsoft.Media_mediaServices_assets@2022-08-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP081] Resource type "Microsoft.Media/mediaServices/assets@2022-08-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
[Warning BCP081] Resource type "Microsoft.Media/mediaServices@2021-11-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
```
## [Microsoft.Media_mediaServices_assets_assetFilters@2022-08-01/main.bicep](../examples/Microsoft.Media_mediaServices_assets_assetFilters@2022-08-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP081] Resource type "Microsoft.Media/mediaServices/assets@2022-08-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
[Warning BCP081] Resource type "Microsoft.Media/mediaServices/assets/assetFilters@2022-08-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
[Warning BCP081] Resource type "Microsoft.Media/mediaServices@2021-11-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
```
## [Microsoft.Media_mediaServices_contentKeyPolicies@2022-08-01/main.bicep](../examples/Microsoft.Media_mediaServices_contentKeyPolicies@2022-08-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP081] Resource type "Microsoft.Media/mediaServices/contentKeyPolicies@2022-08-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
[Warning BCP081] Resource type "Microsoft.Media/mediaServices@2021-11-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
```
## [Microsoft.Media_mediaServices_liveEvents@2022-08-01/main.bicep](../examples/Microsoft.Media_mediaServices_liveEvents@2022-08-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP081] Resource type "Microsoft.Media/mediaServices/liveEvents@2022-08-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
[Warning BCP081] Resource type "Microsoft.Media/mediaServices@2021-11-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
```
## [Microsoft.Media_mediaServices_streamingEndpoints@2022-08-01/main.bicep](../examples/Microsoft.Media_mediaServices_streamingEndpoints@2022-08-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP081] Resource type "Microsoft.Media/mediaServices@2021-11-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
[Warning BCP081] Resource type "Microsoft.Media/mediaServices/streamingEndpoints@2022-08-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
```
## [Microsoft.Media_mediaServices_streamingLocators@2022-08-01/main.bicep](../examples/Microsoft.Media_mediaServices_streamingLocators@2022-08-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP081] Resource type "Microsoft.Media/mediaServices/assets@2022-08-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
[Warning BCP081] Resource type "Microsoft.Media/mediaServices@2021-11-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
[Warning BCP081] Resource type "Microsoft.Media/mediaServices/streamingLocators@2022-08-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
```
## [Microsoft.Media_mediaServices_streamingPolicies@2022-08-01/main.bicep](../examples/Microsoft.Media_mediaServices_streamingPolicies@2022-08-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP081] Resource type "Microsoft.Media/mediaServices@2021-11-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
[Warning BCP081] Resource type "Microsoft.Media/mediaServices/streamingPolicies@2022-08-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
```
## [Microsoft.Media_mediaServices_transforms@2022-07-01/main.bicep](../examples/Microsoft.Media_mediaServices_transforms@2022-07-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP081] Resource type "Microsoft.Media/mediaServices@2021-11-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
[Warning BCP081] Resource type "Microsoft.Media/mediaServices/transforms@2022-07-01" does not have types available. Bicep is unable to validate resource properties prior to deployment, but this will not block the resource from being deployed.
```
## [Microsoft.Migrate_migrateProjects@2018-09-01-preview/main.bicep](../examples/Microsoft.Migrate_migrateProjects@2018-09-01-preview/main.bicep)
Result: success

## [Microsoft.Migrate_migrateProjects_solutions@2018-09-01-preview/main.bicep](../examples/Microsoft.Migrate_migrateProjects_solutions@2018-09-01-preview/main.bicep)
Result: success

## [Microsoft.MixedReality_spatialAnchorsAccounts@2021-01-01/main.bicep](../examples/Microsoft.MixedReality_spatialAnchorsAccounts@2021-01-01/main.bicep)
Result: success

## [Microsoft.MobileNetwork_mobileNetworks@2022-11-01/main.bicep](../examples/Microsoft.MobileNetwork_mobileNetworks@2022-11-01/main.bicep)
Result: success

## [Microsoft.MobileNetwork_mobileNetworks_dataNetworks@2022-11-01/main.bicep](../examples/Microsoft.MobileNetwork_mobileNetworks_dataNetworks@2022-11-01/main.bicep)
Result: success

## [Microsoft.MobileNetwork_mobileNetworks_services@2022-11-01/main.bicep](../examples/Microsoft.MobileNetwork_mobileNetworks_services@2022-11-01/main.bicep)
Result: success

## [Microsoft.MobileNetwork_mobileNetworks_simPolicies@2022-11-01/main.bicep](../examples/Microsoft.MobileNetwork_mobileNetworks_simPolicies@2022-11-01/main.bicep)
Result: success

## [Microsoft.MobileNetwork_mobileNetworks_sites@2022-11-01/main.bicep](../examples/Microsoft.MobileNetwork_mobileNetworks_sites@2022-11-01/main.bicep)
Result: success

## [Microsoft.MobileNetwork_mobileNetworks_slices@2022-11-01/main.bicep](../examples/Microsoft.MobileNetwork_mobileNetworks_slices@2022-11-01/main.bicep)
Result: success

## [Microsoft.MobileNetwork_packetCoreControlPlanes@2022-11-01/main.bicep](../examples/Microsoft.MobileNetwork_packetCoreControlPlanes@2022-11-01/main.bicep)
Result: success

## [Microsoft.MobileNetwork_packetCoreControlPlanes_packetCoreDataPlanes@2022-11-01/main.bicep](../examples/Microsoft.MobileNetwork_packetCoreControlPlanes_packetCoreDataPlanes@2022-11-01/main.bicep)
Result: success

## [Microsoft.MobileNetwork_simGroups@2022-11-01/main.bicep](../examples/Microsoft.MobileNetwork_simGroups@2022-11-01/main.bicep)
Result: success

## [Microsoft.Monitor_accounts@2023-04-03/main.bicep](../examples/Microsoft.Monitor_accounts@2023-04-03/main.bicep)
Result: success

## [Microsoft.Monitor_accounts_privateEndpointConnections@2023-04-03/main.bicep](../examples/Microsoft.Monitor_accounts_privateEndpointConnections@2023-04-03/main.bicep)
Result: success

## [Microsoft.NetApp_netAppAccounts@2022-05-01/main.bicep](../examples/Microsoft.NetApp_netAppAccounts@2022-05-01/main.bicep)
Result: success

## [Microsoft.NetApp_netAppAccounts_capacityPools@2022-05-01/main.bicep](../examples/Microsoft.NetApp_netAppAccounts_capacityPools@2022-05-01/main.bicep)
Result: success

## [Microsoft.NetApp_netAppAccounts_capacityPools_volumes@2022-05-01/main.bicep](../examples/Microsoft.NetApp_netAppAccounts_capacityPools_volumes@2022-05-01/main.bicep)
Result: success

## [Microsoft.NetApp_netAppAccounts_capacityPools_volumes_snapshots@2022-05-01/main.bicep](../examples/Microsoft.NetApp_netAppAccounts_capacityPools_volumes_snapshots@2022-05-01/main.bicep)
Result: success

## [Microsoft.NetApp_netAppAccounts_snapshotPolicies@2022-05-01/main.bicep](../examples/Microsoft.NetApp_netAppAccounts_snapshotPolicies@2022-05-01/main.bicep)
Result: success

## [Microsoft.Network_applicationGateways@2022-07-01/main.bicep](../examples/Microsoft.Network_applicationGateways@2022-07-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Network_ApplicationGatewayWebApplicationFirewallPolicies@2022-07-01/main.bicep](../examples/Microsoft.Network_ApplicationGatewayWebApplicationFirewallPolicies@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_applicationSecurityGroups@2022-09-01/main.bicep](../examples/Microsoft.Network_applicationSecurityGroups@2022-09-01/main.bicep)
Result: success

## [Microsoft.Network_azureFirewalls@2022-07-01/main.bicep](../examples/Microsoft.Network_azureFirewalls@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_bastionHosts@2022-07-01/main.bicep](../examples/Microsoft.Network_bastionHosts@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_ddosProtectionPlans@2022-07-01/main.bicep](../examples/Microsoft.Network_ddosProtectionPlans@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_dnsForwardingRulesets@2022-07-01/main.bicep](../examples/Microsoft.Network_dnsForwardingRulesets@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_dnsForwardingRulesets_forwardingRules@2022-07-01/main.bicep](../examples/Microsoft.Network_dnsForwardingRulesets_forwardingRules@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_dnsForwardingRulesets_virtualNetworkLinks@2022-07-01/main.bicep](../examples/Microsoft.Network_dnsForwardingRulesets_virtualNetworkLinks@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_dnsResolvers@2022-07-01/main.bicep](../examples/Microsoft.Network_dnsResolvers@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_dnsResolvers_inboundEndpoints@2022-07-01/main.bicep](../examples/Microsoft.Network_dnsResolvers_inboundEndpoints@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_dnsResolvers_outboundEndpoints@2022-07-01/main.bicep](../examples/Microsoft.Network_dnsResolvers_outboundEndpoints@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_dnsZones@2018-05-01/main.bicep](../examples/Microsoft.Network_dnsZones@2018-05-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_dnsZones_A@2018-05-01/main.bicep](../examples/Microsoft.Network_dnsZones_A@2018-05-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_dnsZones_AAAA@2018-05-01/main.bicep](../examples/Microsoft.Network_dnsZones_AAAA@2018-05-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_dnsZones_CAA@2018-05-01/main.bicep](../examples/Microsoft.Network_dnsZones_CAA@2018-05-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_dnsZones_CNAME@2018-05-01/main.bicep](../examples/Microsoft.Network_dnsZones_CNAME@2018-05-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_dnsZones_MX@2018-05-01/main.bicep](../examples/Microsoft.Network_dnsZones_MX@2018-05-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_dnsZones_NS@2018-05-01/main.bicep](../examples/Microsoft.Network_dnsZones_NS@2018-05-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_dnsZones_PTR@2018-05-01/main.bicep](../examples/Microsoft.Network_dnsZones_PTR@2018-05-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_dnsZones_SRV@2018-05-01/main.bicep](../examples/Microsoft.Network_dnsZones_SRV@2018-05-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_dnsZones_TXT@2018-05-01/main.bicep](../examples/Microsoft.Network_dnsZones_TXT@2018-05-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_expressRouteCircuits@2022-07-01/main.bicep](../examples/Microsoft.Network_expressRouteCircuits@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_expressRouteCircuits_authorizations@2022-07-01/main.bicep](../examples/Microsoft.Network_expressRouteCircuits_authorizations@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_expressRouteCircuits_peerings@2022-07-01/main.bicep](../examples/Microsoft.Network_expressRouteCircuits_peerings@2022-07-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Network_expressRouteCircuits_peerings_connections@2022-07-01/main.bicep](../examples/Microsoft.Network_expressRouteCircuits_peerings_connections@2022-07-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Network_expressRouteGateways@2022-07-01/main.bicep](../examples/Microsoft.Network_expressRouteGateways@2022-07-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "office365LocalBreakoutCategory" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Network_expressRouteGateways_expressRouteConnections@2022-07-01/main.bicep](../examples/Microsoft.Network_expressRouteGateways_expressRouteConnections@2022-07-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning BCP073] The property "office365LocalBreakoutCategory" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Network_ExpressRoutePorts@2022-07-01/main.bicep](../examples/Microsoft.Network_ExpressRoutePorts@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_expressRoutePorts_authorizations@2022-07-01/main.bicep](../examples/Microsoft.Network_expressRoutePorts_authorizations@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_firewallPolicies@2022-07-01/main.bicep](../examples/Microsoft.Network_firewallPolicies@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_firewallPolicies_ruleCollectionGroups@2022-07-01/main.bicep](../examples/Microsoft.Network_firewallPolicies_ruleCollectionGroups@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_frontDoors_rulesEngines@2020-05-01/main.bicep](../examples/Microsoft.Network_frontDoors_rulesEngines@2020-05-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_FrontDoorWebApplicationFirewallPolicies@2020-11-01/main.bicep](../examples/Microsoft.Network_FrontDoorWebApplicationFirewallPolicies@2020-11-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_ipGroups@2022-07-01/main.bicep](../examples/Microsoft.Network_ipGroups@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_loadBalancers@2022-07-01/main.bicep](../examples/Microsoft.Network_loadBalancers@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_loadBalancers_backendAddressPools@2022-07-01/main.bicep](../examples/Microsoft.Network_loadBalancers_backendAddressPools@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_localNetworkGateways@2022-07-01/main.bicep](../examples/Microsoft.Network_localNetworkGateways@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_natGateways@2022-07-01/main.bicep](../examples/Microsoft.Network_natGateways@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_networkInterfaces@2022-07-01/main.bicep](../examples/Microsoft.Network_networkInterfaces@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_networkManagerConnections@2022-09-01/main.bicep](../examples/Microsoft.Network_networkManagerConnections@2022-09-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP135] Scope "resourceGroup" is not valid for this resource type. Permitted scopes: "managementGroup", "subscription".
```
## [Microsoft.Network_networkManagers@2022-09-01/main.bicep](../examples/Microsoft.Network_networkManagers@2022-09-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Network_networkManagers_connectivityConfigurations@2022-09-01/main.bicep](../examples/Microsoft.Network_networkManagers_connectivityConfigurations@2022-09-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP053] The type "Microsoft.Network/virtualNetworks" does not contain property "output". Available properties include "apiVersion", "asserts", "etag", "eTag", "extendedLocation", "id", "identity", "kind", "location", "managedBy", "managedByExtended", "name", "plan", "properties", "scale", "sku", "tags", "type", "zones".
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Network_networkManagers_ipamPools@2024-01-01-preview/main.bicep](../examples/Microsoft.Network_networkManagers_ipamPools@2024-01-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Network_networkManagers_ipamPools_staticCidr@2024-01-01-preview/main.bicep](../examples/Microsoft.Network_networkManagers_ipamPools_staticCidr@2024-01-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Network_networkManagers_networkGroups@2022-09-01/main.bicep](../examples/Microsoft.Network_networkManagers_networkGroups@2022-09-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Network_networkManagers_networkGroups_staticMembers@2022-09-01/main.bicep](../examples/Microsoft.Network_networkManagers_networkGroups_staticMembers@2022-09-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Network_networkManagers_routingConfigurations_ruleCollections_rules@2024-05-01/main.bicep](../examples/Microsoft.Network_networkManagers_routingConfigurations_ruleCollections_rules@2024-05-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "deploy_locations" is declared but never used.
[Error BCP031] The parameter type is not valid. Please specify one of the following types: "array", "bool", "int", "object", "resourceInput", "resourceOutput", "string".
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Network_networkManagers_scopeConnections@2022-09-01/main.bicep](../examples/Microsoft.Network_networkManagers_scopeConnections@2022-09-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Network_networkManagers_securityAdminConfigurations@2022-09-01/main.bicep](../examples/Microsoft.Network_networkManagers_securityAdminConfigurations@2022-09-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Network_networkManagers_securityAdminConfigurations_ruleCollections@2022-09-01/main.bicep](../examples/Microsoft.Network_networkManagers_securityAdminConfigurations_ruleCollections@2022-09-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Network_networkManagers_securityAdminConfigurations_ruleCollections_rules@2022-09-01/main.bicep](../examples/Microsoft.Network_networkManagers_securityAdminConfigurations_ruleCollections_rules@2022-09-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Network_networkManagers_verifierWorkspace@2024-01-01-preview/main.bicep](../examples/Microsoft.Network_networkManagers_verifierWorkspace@2024-01-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Network_networkManagers_verifierWorkspace_reachabilityAnalysisIntent@2024-01-01-preview/main.bicep](../examples/Microsoft.Network_networkManagers_verifierWorkspace_reachabilityAnalysisIntent@2024-01-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "data" does not exist in the current context.
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.Network_networkManagers_verifierWorkspace_reachabilityAnalysisIntent_reachabilityAnalysisRun@2024-01-01-preview/main.bicep](../examples/Microsoft.Network_networkManagers_verifierWorkspace_reachabilityAnalysisIntent_reachabilityAnalysisRun@2024-01-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'vm_admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "data" does not exist in the current context.
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
[Warning adminusername-should-not-be-literal] Property 'adminUserName' should not use a literal value. Use a param instead. Found literal string value "testadmin"
```
## [Microsoft.Network_networkProfiles@2022-07-01/main.bicep](../examples/Microsoft.Network_networkProfiles@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_networkSecurityGroups@2022-07-01/main.bicep](../examples/Microsoft.Network_networkSecurityGroups@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_networkSecurityGroups_securityRules@2022-09-01/main.bicep](../examples/Microsoft.Network_networkSecurityGroups_securityRules@2022-09-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "resource_name" is declared but never used.
```
## [Microsoft.Network_networkWatchers_flowLogs@2023-11-01/main.bicep](../examples/Microsoft.Network_networkWatchers_flowLogs@2023-11-01/main.bicep)
Result: success

## [Microsoft.Network_p2svpnGateways@2022-07-01/main.bicep](../examples/Microsoft.Network_p2svpnGateways@2022-07-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "office365LocalBreakoutCategory" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Network_privateDnsZones@2018-09-01/main.bicep](../examples/Microsoft.Network_privateDnsZones@2018-09-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_privateDnsZones_A@2018-09-01/main.bicep](../examples/Microsoft.Network_privateDnsZones_A@2018-09-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_privateDnsZones_AAAA@2018-09-01/main.bicep](../examples/Microsoft.Network_privateDnsZones_AAAA@2018-09-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_privateDnsZones_CNAME@2018-09-01/main.bicep](../examples/Microsoft.Network_privateDnsZones_CNAME@2018-09-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_privateDnsZones_MX@2018-09-01/main.bicep](../examples/Microsoft.Network_privateDnsZones_MX@2018-09-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_privateDnsZones_PTR@2018-09-01/main.bicep](../examples/Microsoft.Network_privateDnsZones_PTR@2018-09-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_privateDnsZones_SRV@2018-09-01/main.bicep](../examples/Microsoft.Network_privateDnsZones_SRV@2018-09-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_privateDnsZones_TXT@2018-09-01/main.bicep](../examples/Microsoft.Network_privateDnsZones_TXT@2018-09-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_privateDnsZones_virtualNetworkLinks@2018-09-01/main.bicep](../examples/Microsoft.Network_privateDnsZones_virtualNetworkLinks@2018-09-01/main.bicep)
Result: success

## [Microsoft.Network_privateEndpoints@2022-01-01/main.bicep](../examples/Microsoft.Network_privateEndpoints@2022-01-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "azapi_resource" does not exist in the current context.
```
## [Microsoft.Network_privateLinkServices@2022-07-01/main.bicep](../examples/Microsoft.Network_privateLinkServices@2022-07-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "azapi_resource" does not exist in the current context.
```
## [Microsoft.Network_publicIPAddresses@2022-07-01/main.bicep](../examples/Microsoft.Network_publicIPAddresses@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_publicIPPrefixes@2022-07-01/main.bicep](../examples/Microsoft.Network_publicIPPrefixes@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_routeFilters@2022-09-01/main.bicep](../examples/Microsoft.Network_routeFilters@2022-09-01/main.bicep)
Result: success

## [Microsoft.Network_routeTables@2022-09-01/main.bicep](../examples/Microsoft.Network_routeTables@2022-09-01/main.bicep)
Result: success

## [Microsoft.Network_routeTables_routes@2022-09-01/main.bicep](../examples/Microsoft.Network_routeTables_routes@2022-09-01/main.bicep)
Result: success

## [Microsoft.Network_securityPartnerProviders@2022-07-01/main.bicep](../examples/Microsoft.Network_securityPartnerProviders@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_serviceEndpointPolicies@2022-07-01/main.bicep](../examples/Microsoft.Network_serviceEndpointPolicies@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_trafficManagerProfiles@2018-08-01/main.bicep](../examples/Microsoft.Network_trafficManagerProfiles@2018-08-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_trafficManagerProfiles_AzureEndpoints@2018-08-01/main.bicep](../examples/Microsoft.Network_trafficManagerProfiles_AzureEndpoints@2018-08-01/main.bicep)
Result: success

## [Microsoft.Network_trafficManagerProfiles_ExternalEndpoints@2018-08-01/main.bicep](../examples/Microsoft.Network_trafficManagerProfiles_ExternalEndpoints@2018-08-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_trafficManagerProfiles_NestedEndpoints@2018-08-01/main.bicep](../examples/Microsoft.Network_trafficManagerProfiles_NestedEndpoints@2018-08-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
```
## [Microsoft.Network_virtualHubs@2022-07-01/main.bicep](../examples/Microsoft.Network_virtualHubs@2022-07-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "office365LocalBreakoutCategory" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Network_virtualHubs_hubVirtualNetworkConnections@2022-07-01/main.bicep](../examples/Microsoft.Network_virtualHubs_hubVirtualNetworkConnections@2022-07-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "office365LocalBreakoutCategory" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Network_virtualHubs_ipConfigurations@2022-07-01/main.bicep](../examples/Microsoft.Network_virtualHubs_ipConfigurations@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_virtualHubs_routingIntent@2022-09-01/main.bicep](../examples/Microsoft.Network_virtualHubs_routingIntent@2022-09-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "azurerm_firewall" does not exist in the current context.
[Error BCP057] The name "azurerm_firewall" does not exist in the current context.
[Warning BCP073] The property "office365LocalBreakoutCategory" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Network_virtualNetworkGateways@2022-07-01/main.bicep](../examples/Microsoft.Network_virtualNetworkGateways@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_virtualNetworkGateways_natRules@2022-07-01/main.bicep](../examples/Microsoft.Network_virtualNetworkGateways_natRules@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_virtualNetworks@2022-07-01/main.bicep](../examples/Microsoft.Network_virtualNetworks@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_virtualNetworks_subnets@2022-07-01/main.bicep](../examples/Microsoft.Network_virtualNetworks_subnets@2022-07-01/main.bicep)
Result: success

## [Microsoft.Network_virtualNetworks_virtualNetworkPeerings@2022-07-01/main.bicep](../examples/Microsoft.Network_virtualNetworks_virtualNetworkPeerings@2022-07-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Network_virtualWans@2022-07-01/main.bicep](../examples/Microsoft.Network_virtualWans@2022-07-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "office365LocalBreakoutCategory" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Network_vpnGateways@2022-07-01/main.bicep](../examples/Microsoft.Network_vpnGateways@2022-07-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "office365LocalBreakoutCategory" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Network_vpnGateways_natRules@2022-07-01/main.bicep](../examples/Microsoft.Network_vpnGateways_natRules@2022-07-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "office365LocalBreakoutCategory" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Network_vpnGateways_vpnConnections@2022-07-01/main.bicep](../examples/Microsoft.Network_vpnGateways_vpnConnections@2022-07-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning BCP073] The property "office365LocalBreakoutCategory" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Network_vpnServerConfigurations@2022-07-01/main.bicep](../examples/Microsoft.Network_vpnServerConfigurations@2022-07-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'radius_server_secret' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Network_vpnServerConfigurations_configurationPolicyGroups@2022-07-01/main.bicep](../examples/Microsoft.Network_vpnServerConfigurations_configurationPolicyGroups@2022-07-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'radius_server_secret' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Network_vpnSites@2022-07-01/main.bicep](../examples/Microsoft.Network_vpnSites@2022-07-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "office365LocalBreakoutCategory" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.NotificationHubs_namespaces@2017-04-01/main.bicep](../examples/Microsoft.NotificationHubs_namespaces@2017-04-01/main.bicep)
Result: success

## [Microsoft.NotificationHubs_namespaces_notificationHubs@2017-04-01/main.bicep](../examples/Microsoft.NotificationHubs_namespaces_notificationHubs@2017-04-01/main.bicep)
Result: success

## [Microsoft.NotificationHubs_namespaces_notificationHubs_authorizationRules@2017-04-01/main.bicep](../examples/Microsoft.NotificationHubs_namespaces_notificationHubs_authorizationRules@2017-04-01/main.bicep)
Result: success

## [Microsoft.OperationalInsights_clusters@2020-08-01/main.bicep](../examples/Microsoft.OperationalInsights_clusters@2020-08-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "Identity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.OperationalInsights_queryPacks@2019-09-01/main.bicep](../examples/Microsoft.OperationalInsights_queryPacks@2019-09-01/main.bicep)
Result: success

## [Microsoft.OperationalInsights_queryPacks_queries@2019-09-01/main.bicep](../examples/Microsoft.OperationalInsights_queryPacks_queries@2019-09-01/main.bicep)
Result: success

## [Microsoft.OperationalInsights_workspaces@2022-10-01/main.bicep](../examples/Microsoft.OperationalInsights_workspaces@2022-10-01/main.bicep)
Result: success

## [Microsoft.OperationalInsights_workspaces_dataExports@2020-08-01/main.bicep](../examples/Microsoft.OperationalInsights_workspaces_dataExports@2020-08-01/main.bicep)
Result: success

## [Microsoft.OperationalInsights_workspaces_dataSources@2020-08-01/main.bicep](../examples/Microsoft.OperationalInsights_workspaces_dataSources@2020-08-01/main.bicep)
Result: success

## [Microsoft.OperationalInsights_workspaces_linkedServices@2020-08-01/main.bicep](../examples/Microsoft.OperationalInsights_workspaces_linkedServices@2020-08-01/main.bicep)
Result: success

## [Microsoft.OperationalInsights_workspaces_linkedStorageAccounts@2020-08-01/main.bicep](../examples/Microsoft.OperationalInsights_workspaces_linkedStorageAccounts@2020-08-01/main.bicep)
Result: success

## [Microsoft.OperationalInsights_workspaces_savedSearches@2020-08-01/main.bicep](../examples/Microsoft.OperationalInsights_workspaces_savedSearches@2020-08-01/main.bicep)
Result: success

## [Microsoft.OperationalInsights_workspaces_storageInsightConfigs@2020-08-01/main.bicep](../examples/Microsoft.OperationalInsights_workspaces_storageInsightConfigs@2020-08-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Orbital_spacecrafts@2022-11-01/main.bicep](../examples/Microsoft.Orbital_spacecrafts@2022-11-01/main.bicep)
Result: success

## [Microsoft.PolicyInsights_policyStates@2019-10-01/main.bicep](../examples/Microsoft.PolicyInsights_policyStates@2019-10-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Warning no-unused-params] Parameter "resource_name" is declared but never used.
```
## [Microsoft.PolicyInsights_remediations@2021-10-01/main.bicep](../examples/Microsoft.PolicyInsights_remediations@2021-10-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Portal_dashboards@2019-01-01-preview/main.bicep](../examples/Microsoft.Portal_dashboards@2019-01-01-preview/main.bicep)
Result: success

## [Microsoft.PowerBIDedicated_capacities@2021-01-01/main.bicep](../examples/Microsoft.PowerBIDedicated_capacities@2021-01-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Purview_accounts@2021-07-01/main.bicep](../examples/Microsoft.Purview_accounts@2021-07-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "Identity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.RecoveryServices_vaults@2022-10-01/main.bicep](../examples/Microsoft.RecoveryServices_vaults@2022-10-01/main.bicep)
Result: success

## [Microsoft.RecoveryServices_vaults_backupconfig@2024-04-01/main.bicep](../examples/Microsoft.RecoveryServices_vaults_backupconfig@2024-04-01/main.bicep)
Result: success

## [Microsoft.RecoveryServices_vaults_backupPolicies@2023-02-01/main.bicep](../examples/Microsoft.RecoveryServices_vaults_backupPolicies@2023-02-01/main.bicep)
Result: success

## [Microsoft.RecoveryServices_vaults_backupResourceGuardProxies@2023-02-01/main.bicep](../examples/Microsoft.RecoveryServices_vaults_backupResourceGuardProxies@2023-02-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP073] The property "type" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.RecoveryServices_vaults_backupStorageConfig@2023-02-01/main.bicep](../examples/Microsoft.RecoveryServices_vaults_backupStorageConfig@2023-02-01/main.bicep)
Result: success

## [Microsoft.RecoveryServices_vaults_replicationFabrics@2022-10-01/main.bicep](../examples/Microsoft.RecoveryServices_vaults_replicationFabrics@2022-10-01/main.bicep)
Result: success

## [Microsoft.RecoveryServices_vaults_replicationFabrics_replicationProtectionContainers@2022-10-01/main.bicep](../examples/Microsoft.RecoveryServices_vaults_replicationFabrics_replicationProtectionContainers@2022-10-01/main.bicep)
Result: success

## [Microsoft.RecoveryServices_vaults_replicationPolicies@2022-10-01/main.bicep](../examples/Microsoft.RecoveryServices_vaults_replicationPolicies@2022-10-01/main.bicep)
Result: success

## [Microsoft.Relay_namespaces@2017-04-01/main.bicep](../examples/Microsoft.Relay_namespaces@2017-04-01/main.bicep)
Result: success

## [Microsoft.Relay_namespaces_authorizationRules@2017-04-01/main.bicep](../examples/Microsoft.Relay_namespaces_authorizationRules@2017-04-01/main.bicep)
Result: success

## [Microsoft.Relay_namespaces_hybridConnections@2017-04-01/main.bicep](../examples/Microsoft.Relay_namespaces_hybridConnections@2017-04-01/main.bicep)
Result: success

## [Microsoft.Relay_namespaces_hybridConnections_authorizationRules@2017-04-01/main.bicep](../examples/Microsoft.Relay_namespaces_hybridConnections_authorizationRules@2017-04-01/main.bicep)
Result: success

## [Microsoft.Resources_deployments@2020-06-01/main.bicep](../examples/Microsoft.Resources_deployments@2020-06-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Warning no-deployments-resources] Resource 'deployment' of type 'Microsoft.Resources/deployments@2020-06-01' should instead be declared as a Bicep module.
```
## [Microsoft.Resources_deploymentScripts@2020-10-01/main.bicep](../examples/Microsoft.Resources_deploymentScripts@2020-10-01/main.bicep)
Result: success

Diagnostics:
```
[Warning use-recent-az-powershell-version] Deployment script is using AzPowerShell version '8.3' which is below the recommended minimum version '11.0'. Consider upgrading to version 11.0 or higher to avoid EOL Ubuntu 20.04 LTS.
```
## [Microsoft.Resources_resourceGroups@2020-06-01/main.bicep](../examples/Microsoft.Resources_resourceGroups@2020-06-01/main.bicep)
Result: success

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Warning no-unused-params] Parameter "resource_name" is declared but never used.
```
## [Microsoft.Search_searchServices@2022-09-01/main.bicep](../examples/Microsoft.Search_searchServices@2022-09-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "publicNetworkAccess" expected a value of type "'disabled' | 'enabled' | null" but the provided value is of type "'Enabled'". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Search_searchServices_sharedPrivateLinkResources@2022-09-01/main.bicep](../examples/Microsoft.Search_searchServices_sharedPrivateLinkResources@2022-09-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "publicNetworkAccess" expected a value of type "'disabled' | 'enabled' | null" but the provided value is of type "'Enabled'". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.SecurityInsights_alertRules@2022-10-01-preview/main.bicep](../examples/Microsoft.SecurityInsights_alertRules@2022-10-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP037] The property "parent" is not allowed on objects of type "NRT". Permissible properties include "dependsOn", "etag", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.SecurityInsights/onboardingStates". Permissible properties include "asserts", "dependsOn", "etag", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.SecurityInsights_automationRules@2022-10-01-preview/main.bicep](../examples/Microsoft.SecurityInsights_automationRules@2022-10-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.SecurityInsights/automationRules". Permissible properties include "asserts", "dependsOn", "etag", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.SecurityInsights/onboardingStates". Permissible properties include "asserts", "dependsOn", "etag", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.SecurityInsights_dataConnectors@2022-10-01-preview/main.bicep](../examples/Microsoft.SecurityInsights_dataConnectors@2022-10-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.SecurityInsights/onboardingStates". Permissible properties include "asserts", "dependsOn", "etag", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.SecurityInsights_metadata@2022-10-01-preview/main.bicep](../examples/Microsoft.SecurityInsights_metadata@2022-10-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP037] The property "parent" is not allowed on objects of type "NRT". Permissible properties include "dependsOn", "etag", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.SecurityInsights/metadata". Permissible properties include "asserts", "dependsOn", "etag", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.SecurityInsights/onboardingStates". Permissible properties include "asserts", "dependsOn", "etag", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.SecurityInsights_onboardingStates@2022-11-01/main.bicep](../examples/Microsoft.SecurityInsights_onboardingStates@2022-11-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.SecurityInsights/onboardingStates". Permissible properties include "asserts", "dependsOn", "etag", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.SecurityInsights_watchlists@2022-11-01/main.bicep](../examples/Microsoft.SecurityInsights_watchlists@2022-11-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.SecurityInsights/onboardingStates". Permissible properties include "asserts", "dependsOn", "etag", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.SecurityInsights/watchlists". Permissible properties include "asserts", "dependsOn", "etag", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.SecurityInsights_watchlists_watchlistItems@2022-11-01/main.bicep](../examples/Microsoft.SecurityInsights_watchlists_watchlistItems@2022-11-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.SecurityInsights/onboardingStates". Permissible properties include "asserts", "dependsOn", "etag", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.SecurityInsights/watchlists". Permissible properties include "asserts", "dependsOn", "etag", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Security_advancedThreatProtectionSettings@2019-01-01/main.bicep](../examples/Microsoft.Security_advancedThreatProtectionSettings@2019-01-01/main.bicep)
Result: success

## [Microsoft.Security_assessmentMetadata@2020-01-01/main.bicep](../examples/Microsoft.Security_assessmentMetadata@2020-01-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Warning no-unused-params] Parameter "resource_name" is declared but never used.
[Error BCP135] Scope "resourceGroup" is not valid for this resource type. Permitted scopes: "tenant", "subscription".
```
## [Microsoft.Security_automations@2019-01-01-preview/main.bicep](../examples/Microsoft.Security_automations@2019-01-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP063] The name "resourceGroup" is not a parameter, variable, resource or module.
```
## [Microsoft.Security_defenderForStorageSettings@2022-12-01-preview/main.bicep](../examples/Microsoft.Security_defenderForStorageSettings@2022-12-01-preview/main.bicep)
Result: success

## [Microsoft.Security_iotSecuritySolutions@2019-08-01/main.bicep](../examples/Microsoft.Security_iotSecuritySolutions@2019-08-01/main.bicep)
Result: success

## [Microsoft.Security_securityContacts@2017-08-01-preview/main.bicep](../examples/Microsoft.Security_securityContacts@2017-08-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Error BCP135] Scope "resourceGroup" is not valid for this resource type. Permitted scopes: "subscription".
```
## [Microsoft.Security_workspaceSettings@2017-08-01-preview/main.bicep](../examples/Microsoft.Security_workspaceSettings@2017-08-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP135] Scope "resourceGroup" is not valid for this resource type. Permitted scopes: "subscription".
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.ServiceBus_namespaces@2022-01-01-preview/main.bicep](../examples/Microsoft.ServiceBus_namespaces@2022-01-01-preview/main.bicep)
Result: success

## [Microsoft.ServiceBus_namespaces_authorizationRules@2021-06-01-preview/main.bicep](../examples/Microsoft.ServiceBus_namespaces_authorizationRules@2021-06-01-preview/main.bicep)
Result: success

## [Microsoft.ServiceBus_namespaces_queues@2021-06-01-preview/main.bicep](../examples/Microsoft.ServiceBus_namespaces_queues@2021-06-01-preview/main.bicep)
Result: success

## [Microsoft.ServiceBus_namespaces_queues_authorizationRules@2021-06-01-preview/main.bicep](../examples/Microsoft.ServiceBus_namespaces_queues_authorizationRules@2021-06-01-preview/main.bicep)
Result: success

## [Microsoft.ServiceBus_namespaces_topics@2021-06-01-preview/main.bicep](../examples/Microsoft.ServiceBus_namespaces_topics@2021-06-01-preview/main.bicep)
Result: success

## [Microsoft.ServiceBus_namespaces_topics_authorizationRules@2021-06-01-preview/main.bicep](../examples/Microsoft.ServiceBus_namespaces_topics_authorizationRules@2021-06-01-preview/main.bicep)
Result: success

## [Microsoft.ServiceBus_namespaces_topics_subscriptions@2021-06-01-preview/main.bicep](../examples/Microsoft.ServiceBus_namespaces_topics_subscriptions@2021-06-01-preview/main.bicep)
Result: success

## [Microsoft.ServiceBus_namespaces_topics_subscriptions_rules@2021-06-01-preview/main.bicep](../examples/Microsoft.ServiceBus_namespaces_topics_subscriptions_rules@2021-06-01-preview/main.bicep)
Result: success

## [Microsoft.ServiceFabric_clusters@2021-06-01/main.bicep](../examples/Microsoft.ServiceFabric_clusters@2021-06-01/main.bicep)
Result: success

## [Microsoft.ServiceFabric_managedClusters@2021-05-01/main.bicep](../examples/Microsoft.ServiceFabric_managedClusters@2021-05-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.ServiceFabric_managedClusters_nodeTypes@2021-05-01/main.bicep](../examples/Microsoft.ServiceFabric_managedClusters_nodeTypes@2021-05-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.ServiceLinker_linkers@2022-05-01/main.bicep](../examples/Microsoft.ServiceLinker_linkers@2022-05-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "ManagedIdentityProperties | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.ServiceLinker/linkers". Permissible properties include "asserts", "dependsOn", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.SignalRService_signalR@2023-02-01/main.bicep](../examples/Microsoft.SignalRService_signalR@2023-02-01/main.bicep)
Result: success

## [Microsoft.SignalRService_signalR_sharedPrivateLinkResources@2023-02-01/main.bicep](../examples/Microsoft.SignalRService_signalR_sharedPrivateLinkResources@2023-02-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.SignalRService_webPubSub@2023-02-01/main.bicep](../examples/Microsoft.SignalRService_webPubSub@2023-02-01/main.bicep)
Result: success

## [Microsoft.SignalRService_webPubSub_hubs@2023-02-01/main.bicep](../examples/Microsoft.SignalRService_webPubSub_hubs@2023-02-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP036] The property "identity" expected a value of type "ManagedIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.SignalRService_webPubSub_sharedPrivateLinkResources@2023-02-01/main.bicep](../examples/Microsoft.SignalRService_webPubSub_sharedPrivateLinkResources@2023-02-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Solutions_applicationDefinitions@2021-07-01/main.bicep](../examples/Microsoft.Solutions_applicationDefinitions@2021-07-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Sql_instancePools@2022-05-01-preview/main.bicep](../examples/Microsoft.Sql_instancePools@2022-05-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Sql_servers@2021-02-01-preview/main.bicep](../examples/Microsoft.Sql_servers@2021-02-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'sql_administrator_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Sql_servers_administrators@2020-11-01-preview/main.bicep](../examples/Microsoft.Sql_servers_administrators@2020-11-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Sql_servers_auditingSettings@2022-05-01-preview/main.bicep](../examples/Microsoft.Sql_servers_auditingSettings@2022-05-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.Sql_servers_automaticTuning@2021-11-01/main.bicep](../examples/Microsoft.Sql_servers_automaticTuning@2021-11-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.Sql_servers_connectionPolicies@2014-04-01/main.bicep](../examples/Microsoft.Sql_servers_connectionPolicies@2014-04-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Sql_servers_databases@2014-04-01/main.bicep](../examples/Microsoft.Sql_servers_databases@2014-04-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Sql_servers_databases@2021-02-01-preview/main.bicep](../examples/Microsoft.Sql_servers_databases@2021-02-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Sql_servers_databases_securityAlertPolicies@2014-04-01/main.bicep](../examples/Microsoft.Sql_servers_databases_securityAlertPolicies@2014-04-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Sql_servers_databases_securityAlertPolicies@2020-11-01-preview/main.bicep](../examples/Microsoft.Sql_servers_databases_securityAlertPolicies@2020-11-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Sql_servers_databases_transparentDataEncryption@2014-04-01/main.bicep](../examples/Microsoft.Sql_servers_databases_transparentDataEncryption@2014-04-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Sql_servers_dnsAliases@2020-11-01-preview/main.bicep](../examples/Microsoft.Sql_servers_dnsAliases@2020-11-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Sql_servers_elasticPools@2014-04-01/main.bicep](../examples/Microsoft.Sql_servers_elasticPools@2014-04-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Sql_servers_elasticPools@2020-11-01-preview/main.bicep](../examples/Microsoft.Sql_servers_elasticPools@2020-11-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Sql_servers_firewallRules@2014-04-01/main.bicep](../examples/Microsoft.Sql_servers_firewallRules@2014-04-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Sql_servers_firewallRules@2020-11-01-preview/main.bicep](../examples/Microsoft.Sql_servers_firewallRules@2020-11-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Sql_servers_jobAgents@2020-11-01-preview/main.bicep](../examples/Microsoft.Sql_servers_jobAgents@2020-11-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'sql_administrator_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Sql_servers_jobAgents_credentials@2020-11-01-preview/main.bicep](../examples/Microsoft.Sql_servers_jobAgents_credentials@2020-11-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'sql_admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Sql_servers_outboundFirewallRules@2021-02-01-preview/main.bicep](../examples/Microsoft.Sql_servers_outboundFirewallRules@2021-02-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning no-hardcoded-env-urls] Environment URLs should not be hardcoded. Use the environment() function to ensure compatibility across clouds. Found this disallowed host: "database.windows.net"
```
## [Microsoft.Sql_servers_securityAlertPolicies@2017-03-01-preview/main.bicep](../examples/Microsoft.Sql_servers_securityAlertPolicies@2017-03-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.Sql_servers_sqlVulnerabilityAssessments@2022-05-01-preview/main.bicep](../examples/Microsoft.Sql_servers_sqlVulnerabilityAssessments@2022-05-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'administratorLoginPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.Sql_servers_virtualNetworkRules@2020-11-01-preview/main.bicep](../examples/Microsoft.Sql_servers_virtualNetworkRules@2020-11-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'sql_administrator_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
```
## [Microsoft.StorageCache_caches@2023-01-01/main.bicep](../examples/Microsoft.StorageCache_caches@2023-01-01/main.bicep)
Result: success

## [Microsoft.StorageMover_storageMovers@2023-03-01/main.bicep](../examples/Microsoft.StorageMover_storageMovers@2023-03-01/main.bicep)
Result: success

## [Microsoft.StorageMover_storageMovers_endpoints@2023-03-01/main.bicep](../examples/Microsoft.StorageMover_storageMovers_endpoints@2023-03-01/main.bicep)
Result: success

## [Microsoft.StorageMover_storageMovers_projects@2023-03-01/main.bicep](../examples/Microsoft.StorageMover_storageMovers_projects@2023-03-01/main.bicep)
Result: success

## [Microsoft.StorageSync_storageSyncServices@2020-03-01/main.bicep](../examples/Microsoft.StorageSync_storageSyncServices@2020-03-01/main.bicep)
Result: success

## [Microsoft.StorageSync_storageSyncServices_syncGroups@2020-03-01/main.bicep](../examples/Microsoft.StorageSync_storageSyncServices_syncGroups@2020-03-01/main.bicep)
Result: success

## [Microsoft.Storage_storageAccounts_blobServices@2021-09-01/main.bicep](../examples/Microsoft.Storage_storageAccounts_blobServices@2021-09-01/main.bicep)
Result: success

## [Microsoft.Storage_storageAccounts_blobServices_containers@2021-09-01/main.bicep](../examples/Microsoft.Storage_storageAccounts_blobServices_containers@2021-09-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP035] The specified "resource" declaration is missing the following required properties: "kind". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Storage_storageAccounts_blobServices_containers_immutabilityPolicies@2023-05-01/main.bicep](../examples/Microsoft.Storage_storageAccounts_blobServices_containers_immutabilityPolicies@2023-05-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP035] The specified "resource" declaration is missing the following required properties: "kind". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Storage_storageAccounts_fileServices_shares@2021-09-01/main.bicep](../examples/Microsoft.Storage_storageAccounts_fileServices_shares@2021-09-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP035] The specified "resource" declaration is missing the following required properties: "kind". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Storage_storageAccounts_managementPolicies@2021-09-01/main.bicep](../examples/Microsoft.Storage_storageAccounts_managementPolicies@2021-09-01/main.bicep)
Result: success

## [Microsoft.Storage_storageAccounts_queueServices_queues@2021-09-01/main.bicep](../examples/Microsoft.Storage_storageAccounts_queueServices_queues@2021-09-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP035] The specified "resource" declaration is missing the following required properties: "kind". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Storage_storageAccounts_tableServices_tables@2021-09-01/main.bicep](../examples/Microsoft.Storage_storageAccounts_tableServices_tables@2021-09-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP035] The specified "resource" declaration is missing the following required properties: "kind". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.StreamAnalytics_clusters@2020-03-01/main.bicep](../examples/Microsoft.StreamAnalytics_clusters@2020-03-01/main.bicep)
Result: success

## [Microsoft.StreamAnalytics_streamingJobs@2020-03-01/main.bicep](../examples/Microsoft.StreamAnalytics_streamingJobs@2020-03-01/main.bicep)
Result: success

## [Microsoft.StreamAnalytics_streamingJobs_functions@2020-03-01/main.bicep](../examples/Microsoft.StreamAnalytics_streamingJobs_functions@2020-03-01/main.bicep)
Result: success

## [Microsoft.StreamAnalytics_streamingJobs_inputs@2020-03-01/main.bicep](../examples/Microsoft.StreamAnalytics_streamingJobs_inputs@2020-03-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.StreamAnalytics_streamingJobs_outputs@2021-10-01-preview/main.bicep](../examples/Microsoft.StreamAnalytics_streamingJobs_outputs@2021-10-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Synapse_privateLinkHubs@2021-06-01/main.bicep](../examples/Microsoft.Synapse_privateLinkHubs@2021-06-01/main.bicep)
Result: success

## [Microsoft.Synapse_workspaces@2021-06-01/main.bicep](../examples/Microsoft.Synapse_workspaces@2021-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'sql_administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning BCP036] The property "identity" expected a value of type "ManagedIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Synapse_workspaces_azureADOnlyAuthentications@2021-06-01-preview/main.bicep](../examples/Microsoft.Synapse_workspaces_azureADOnlyAuthentications@2021-06-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'sql_administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning BCP036] The property "identity" expected a value of type "ManagedIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Synapse_workspaces_bigDataPools@2021-06-01-preview/main.bicep](../examples/Microsoft.Synapse_workspaces_bigDataPools@2021-06-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'sql_administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning BCP036] The property "identity" expected a value of type "ManagedIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Synapse_workspaces_firewallRules@2021-06-01/main.bicep](../examples/Microsoft.Synapse_workspaces_firewallRules@2021-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'sql_administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning BCP036] The property "identity" expected a value of type "ManagedIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Synapse_workspaces_integrationRuntimes@2021-06-01-preview/main.bicep](../examples/Microsoft.Synapse_workspaces_integrationRuntimes@2021-06-01-preview/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'sql_administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning BCP036] The property "identity" expected a value of type "ManagedIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Synapse_workspaces_managedIdentitySqlControlSettings@2021-06-01/main.bicep](../examples/Microsoft.Synapse_workspaces_managedIdentitySqlControlSettings@2021-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'sql_administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning BCP036] The property "identity" expected a value of type "ManagedIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Synapse_workspaces_securityAlertPolicies@2021-06-01/main.bicep](../examples/Microsoft.Synapse_workspaces_securityAlertPolicies@2021-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'sql_administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning BCP036] The property "identity" expected a value of type "ManagedIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Synapse_workspaces_sqlPools@2021-06-01/main.bicep](../examples/Microsoft.Synapse_workspaces_sqlPools@2021-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'sql_administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning BCP036] The property "identity" expected a value of type "ManagedIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Synapse_workspaces_sqlPools_vulnerabilityAssessments@2021-06-01/main.bicep](../examples/Microsoft.Synapse_workspaces_sqlPools_vulnerabilityAssessments@2021-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'sql_administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning BCP036] The property "identity" expected a value of type "ManagedIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Synapse_workspaces_sqlPools_workloadGroups@2021-06-01/main.bicep](../examples/Microsoft.Synapse_workspaces_sqlPools_workloadGroups@2021-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'sql_administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning BCP036] The property "identity" expected a value of type "ManagedIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Synapse_workspaces_sqlPools_workloadGroups_workloadClassifiers@2021-06-01/main.bicep](../examples/Microsoft.Synapse_workspaces_sqlPools_workloadGroups_workloadClassifiers@2021-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning secure-secrets-in-params] Parameter 'sql_administrator_login_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning BCP036] The property "identity" expected a value of type "ManagedIdentity | null" but the provided value is of type "[object]". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.TimeSeriesInsights_environments@2020-05-15/main.bicep](../examples/Microsoft.TimeSeriesInsights_environments@2020-05-15/main.bicep)
Result: success

## [Microsoft.TimeSeriesInsights_environments_accessPolicies@2020-05-15/main.bicep](../examples/Microsoft.TimeSeriesInsights_environments_accessPolicies@2020-05-15/main.bicep)
Result: success

## [Microsoft.TimeSeriesInsights_environments_eventSources@2020-05-15/main.bicep](../examples/Microsoft.TimeSeriesInsights_environments_eventSources@2020-05-15/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.TimeSeriesInsights_environments_referenceDataSets@2020-05-15/main.bicep](../examples/Microsoft.TimeSeriesInsights_environments_referenceDataSets@2020-05-15/main.bicep)
Result: success

## [Microsoft.VoiceServices_communicationsGateways@2023-01-31/main.bicep](../examples/Microsoft.VoiceServices_communicationsGateways@2023-01-31/main.bicep)
Result: success

## [Microsoft.VoiceServices_communicationsGateways_testLines@2023-01-31/main.bicep](../examples/Microsoft.VoiceServices_communicationsGateways_testLines@2023-01-31/main.bicep)
Result: success

## [Microsoft.Web_certificates@2021-02-01/main.bicep](../examples/Microsoft.Web_certificates@2021-02-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'certificate_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "filebase64" does not exist in the current context.
[Error BCP103] The following token is not recognized: """. Strings are defined using single quotes in bicep.
[Error BCP009] Expected a literal value, an array, an object, a parenthesized expression, or a function call at this location.
[Error BCP103] The following token is not recognized: """. Strings are defined using single quotes in bicep.
```
## [Microsoft.Web_connections@2016-06-01/main.bicep](../examples/Microsoft.Web_connections@2016-06-01/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Web_serverfarms@2022-09-01/main.bicep](../examples/Microsoft.Web_serverfarms@2022-09-01/main.bicep)
Result: success

## [Microsoft.Web_sites@2022-09-01/main.bicep](../examples/Microsoft.Web_sites@2022-09-01/main.bicep)
Result: success

## [Microsoft.Web_sites_config@2022-09-01/main.bicep](../examples/Microsoft.Web_sites_config@2022-09-01/main.bicep)
Result: success

## [Microsoft.Web_sites_publicCertificates@2022-09-01/main.bicep](../examples/Microsoft.Web_sites_publicCertificates@2022-09-01/main.bicep)
Result: success

## [Microsoft.Web_sites_siteextensions@2022-09-01/main.bicep](../examples/Microsoft.Web_sites_siteextensions@2022-09-01/main.bicep)
Result: success

Diagnostics:
```
[Warning BCP187] The property "location" does not exist in the resource or type definition, although it might still be valid. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Web_sites_slots@2022-09-01/main.bicep](../examples/Microsoft.Web_sites_slots@2022-09-01/main.bicep)
Result: success

## [Microsoft.Web_sites_slots_config@2022-09-01/main.bicep](../examples/Microsoft.Web_sites_slots_config@2022-09-01/main.bicep)
Result: success

## [Microsoft.Web_sourcecontrols@2021-02-01/main.bicep](../examples/Microsoft.Web_sourcecontrols@2021-02-01/main.bicep)
Result: failed (unexpected error)

```
TypeError: Cannot convert undefined or null to object
```
## [Microsoft.Web_staticSites@2021-02-01/main.bicep](../examples/Microsoft.Web_staticSites@2021-02-01/main.bicep)
Result: success

## [Qumulo.Storage_fileSystems@2024-06-19/main.bicep](../examples/Qumulo.Storage_fileSystems@2024-06-19/main.bicep)
Result: success

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'qumulo_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
[Warning use-secure-value-for-secure-inputs] Property 'email' expects a secure value, but the value provided may not be secure.
[Warning BCP187] The property "location" does not exist in the resource or type definition, although it might still be valid. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Warning BCP073] The property "actions" is read-only. Expressions cannot be assigned to read-only properties. If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
```
## [Microsoft.Compute_virtualMachines@2023-03-01/attach_data_disk/main.bicep](../examples/Microsoft.Compute_virtualMachines@2023-03-01/attach_data_disk/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "local" does not exist in the current context.
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
[Error BCP057] The name "local" does not exist in the current context.
[Error BCP057] The name "local" does not exist in the current context.
```
## [Microsoft.Compute_virtualMachines@2023-03-01/attach_os_disk/main.bicep](../examples/Microsoft.Compute_virtualMachines@2023-03-01/attach_os_disk/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP057] The name "local" does not exist in the current context.
[Error BCP057] The name "local" does not exist in the current context.
[Error BCP057] The name "data" does not exist in the current context.
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
[Error BCP057] The name "local" does not exist in the current context.
```
## [Microsoft.Compute_virtualMachines@2023-03-01/basic/main.bicep](../examples/Microsoft.Compute_virtualMachines@2023-03-01/basic/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
```
## [Microsoft.Compute_virtualMachines@2023-03-01/tag_os_disk/main.bicep](../examples/Microsoft.Compute_virtualMachines@2023-03-01/tag_os_disk/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
[Error BCP057] The name "local" does not exist in the current context.
[Error BCP057] The name "local" does not exist in the current context.
```
## [Microsoft.Network_virtualNetworks@2024-05-01/with_ipam_pool/main.bicep](../examples/Microsoft.Network_virtualNetworks@2024-05-01/with_ipam_pool/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.Network_virtualNetworks_subnets@2024-05-01/with_ipam_pool/main.bicep](../examples/Microsoft.Network_virtualNetworks_subnets@2024-05-01/with_ipam_pool/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "data" does not exist in the current context.
```
## [Microsoft.OperationalInsights_workspaces_tables@2022-10-01/audit_log/main.bicep](../examples/Microsoft.OperationalInsights_workspaces_tables@2022-10-01/audit_log/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "local" does not exist in the current context.
[Error BCP057] The name "local" does not exist in the current context.
[Error BCP057] The name "local" does not exist in the current context.
```
## [Microsoft.OperationalInsights_workspaces_tables@2022-10-01/basic/main.bicep](../examples/Microsoft.OperationalInsights_workspaces_tables@2022-10-01/basic/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "local" does not exist in the current context.
[Error BCP057] The name "local" does not exist in the current context.
[Error BCP057] The name "local" does not exist in the current context.
```
## [Microsoft.OperationalInsights_workspaces_tables@2022-10-01/data_collection_logs/main.bicep](../examples/Microsoft.OperationalInsights_workspaces_tables@2022-10-01/data_collection_logs/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Error BCP057] The name "local" does not exist in the current context.
[Error BCP057] The name "local" does not exist in the current context.
[Error BCP057] The name "local" does not exist in the current context.
```
## [Microsoft.SqlVirtualMachine_sqlVirtualMachines@2023-10-01/basic/main.bicep](../examples/Microsoft.SqlVirtualMachine_sqlVirtualMachines@2023-10-01/basic/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'vm_admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP120] This expression is being used in an assignment to the "location" property of the "Microsoft.SqlVirtualMachine/sqlVirtualMachines" type, which requires a value that can be calculated at the start of the deployment. Properties of virtualMachine which can be calculated at the start include "apiVersion", "id", "name", "type".
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
[Warning adminusername-should-not-be-literal] Property 'adminUserName' should not use a literal value. Use a param instead. Found literal string value "testadmin"
```
## [Microsoft.SqlVirtualMachine_sqlVirtualMachines@2023-10-01/SQL_best_practices_assessment/main.bicep](../examples/Microsoft.SqlVirtualMachine_sqlVirtualMachines@2023-10-01/SQL_best_practices_assessment/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning secure-secrets-in-params] Parameter 'admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Error BCP033] Expected a value of type "string" but the provided value is of type "null".
[Error BCP120] This expression is being used in an assignment to the "name" property of the "Microsoft.Insights/dataCollectionRules" type, which requires a value that can be calculated at the start of the deployment. Properties of workspace which can be calculated at the start include "apiVersion", "id", "name", "type".
[Error BCP037] The property "parent" is not allowed on objects of type "Microsoft.Insights/dataCollectionRuleAssociations". Permissible properties include "asserts", "dependsOn", "scope". If this is a resource type definition inaccuracy, report it using https://aka.ms/bicep-type-issues.
[Error BCP120] This expression is being used in an assignment to the "name" property of the "Microsoft.Insights/dataCollectionRuleAssociations" type, which requires a value that can be calculated at the start of the deployment. Properties of workspace which can be calculated at the start include "apiVersion", "id", "name", "type".
[Error BCP120] This expression is being used in an assignment to the "location" property of the "Microsoft.SqlVirtualMachine/sqlVirtualMachines" type, which requires a value that can be calculated at the start of the deployment. Properties of virtualMachine which can be calculated at the start include "apiVersion", "id", "name", "type".
[Warning use-secure-value-for-secure-inputs] Property 'adminPassword' expects a secure value, but the value provided may not be secure.
[Warning adminusername-should-not-be-literal] Property 'adminUserName' should not use a literal value. Use a param instead. Found literal string value "testadmin"
```
## [Microsoft.StandbyPool_standbyContainerGroupPools@2025-03-01/basic/main.bicep](../examples/Microsoft.StandbyPool_standbyContainerGroupPools@2025-03-01/basic/main.bicep)
Result: success

## [Microsoft.StandbyPool_standbyvirtualmachinepools@2025-03-01/basic/main.bicep](../examples/Microsoft.StandbyPool_standbyvirtualmachinepools@2025-03-01/basic/main.bicep)
Result: success

Diagnostics:
```
[Warning adminusername-should-not-be-literal] Property 'adminUserName' should not use a literal value. Use a param instead. Found literal string value "adminuser"
```
## [Microsoft.Storage_storageAccounts@2021-09-01/basic/main.bicep](../examples/Microsoft.Storage_storageAccounts@2021-09-01/basic/main.bicep)
Result: success

## [Microsoft.Storage_storageAccounts@2021-09-01/with_private_endpoint/main.bicep](../examples/Microsoft.Storage_storageAccounts@2021-09-01/with_private_endpoint/main.bicep)
Result: failed (invalid bicep)

Diagnostics:
```
[Warning no-unused-params] Parameter "location" is declared but never used.
[Warning no-unused-params] Parameter "subscription_id" is declared but never used.
[Warning no-unused-params] Parameter "vm_admin_password" is declared but never used.
[Warning secure-secrets-in-params] Parameter 'vm_admin_password' may represent a secret (according to its name) and must be declared with the '@secure()' attribute.
[Warning no-unused-params] Parameter "vm_admin_username" is declared but never used.
[Error BCP057] The name "azurerm_resource_group" does not exist in the current context.
```
## [Microsoft.Storage_storageAccounts_localUsers@2021-09-01/basic/main.bicep](../examples/Microsoft.Storage_storageAccounts_localUsers@2021-09-01/basic/main.bicep)
Result: success

## [Microsoft.Storage_storageAccounts_localUsers@2021-09-01/generate_password/main.bicep](../examples/Microsoft.Storage_storageAccounts_localUsers@2021-09-01/generate_password/main.bicep)
Result: success

