---
layout: "azapi"
page_title: "AzAPI Provider: Understand Azure resource concepts"
description: |-
  This guide will cover some basic Azure Resource concepts
---
# Azure Resource Concepts

## Resource Terminology
1. `Root Resource / Top-Level Resource`: A resource with only a single level of nested types (i.e. there's just a single type after the RP namespace). `Microsoft.Network/networkSecurityGroups` is a top-level resource, whereas `Microsoft.Network/networkSecurityGroups/securityRules` is not.
2. `Child Resource / Nested Resource`: A resource with two or more levels of nested types.
3. `Parent Resource`: The parent to a child resource, identified by removing a level of nesting from the resource type. `Microsoft.Network/networkSecurityGroups` is the parent to `Microsoft.Network/networkSecurityGroups/securityRules`.

## Scopes:
Every resource ID begins with a root scope. There are a few different types of root scopes:
1. Tenant Scope: id starts with `/`, then followed with `/providers/Microsoft.Foo/bar/{name}`
2. Subscription Scope: id starts with `/subscriptions/{subscriptionId}`, then followed with `/providers/Microsoft.Foo/bar/{name}`
3. ResourceGroup Scope: id starts with `/subscriptions/{subscriptionId}/resourceGroups/{groupName}`, then followed with `/providers/Microsoft.Foo/bar/{name}`
4. Extension Scope: id starts with `{resourceID}`, then followed with `/providers/Microsoft.Foo/bar/{name}`
5. ManagementGroup Scope: id starts with `/providers/Microsoft.Management/managementGroups/{managementGroupName}/`, then followed with `/providers/Microsoft.Foo/bar/{name}`

## How does `azapi` provider build resource ID from `type`, `parent_id` and `name`
Let's assume that these values are valid, all the cases including edge cases can be divided into 2 scenarios:
1. `type` is a top level resource, then `resource_id` = `{parent_id}/providers/{type}/{name}`
   1. One special case is `Microsoft.Resources/resourceGroups`.
2. `type` is a child resource, then `resource_id` = `{parent_id}/{last nesting type}/{name}`
 
Then we need to add some validations before building the resourceId.
1. If it's a top level resource, `parent_id` must match with correct scope. There're cases that a resource supports both 
   `Tenant` and `Subscription` scopes, the `parent_id` must match any of them.
2. If it's a child resource, `parent_id`'s type must match with its parent resource's type.
3. For the resources whose scope is unknown or not defined in `bicep-types-az`, validations will be skipped.
