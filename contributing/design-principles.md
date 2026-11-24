# Design Principles

- The AzAPI provider is a very thin layer on top of the [Azure ARM REST APIs](https://learn.microsoft.com/rest/api/azure) and data plane APIs
- A single `azapi_resource` should correspond to CRUD of a single REST API domain object, avoid shoving complex logic such as payload manipulation or polling only for specific API
- API problems and bugs should be fixed upstream
- Prefer builtin Terraform features and constructs over implementing custom logic in the provider
- Past code that has been merged in the past is not necessarily a guideline / pattern for future changes, they might have to be kept solely for backward compatibility
