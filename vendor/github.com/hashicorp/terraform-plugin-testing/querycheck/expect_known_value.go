// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ QueryResultCheck = expectKnownValue{}

type expectKnownValue struct {
	listResourceAddress string
	resourceName        string
	attributePath       tfjsonpath.Path
	knownValue          knownvalue.Check
}

func (e expectKnownValue) CheckQuery(_ context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	for _, res := range req.Query {
		var diags []error

		if e.listResourceAddress == strings.TrimPrefix(res.Address, "list.") && e.resourceName == res.DisplayName {
			if res.ResourceObject == nil {
				resp.Error = fmt.Errorf("%s - no resource object was returned, ensure `include_resource` has been set to `true` in the list resource config`", e.listResourceAddress)
				return
			}

			resource, err := tfjsonpath.Traverse(res.ResourceObject, e.attributePath)
			if err != nil {
				resp.Error = err
				return
			}

			if err := e.knownValue.CheckValue(resource); err != nil {
				diags = append(diags, fmt.Errorf("error checking value for attribute at path: %s for resource %s, err: %s", e.attributePath.String(), e.resourceName, err))
			}

			if diags == nil {
				return
			}
		}

		if diags != nil {
			var diagsStr string
			for _, diag := range diags {
				diagsStr += diag.Error() + "; "
			}
			resp.Error = fmt.Errorf("the following errors were found while checking values: %s", diagsStr)
			return
		}
	}

	resp.Error = fmt.Errorf("%s - the resource %s was not found", e.listResourceAddress, e.resourceName)
}

// ExpectKnownValue returns a query check that asserts the specified attribute values are present for a given resource object
// returned by a list query. The resource object can only be identified by providing the list resource address as well as the
// resource name (display name).
//
// This query check can only be used with managed resources that support resource identity and query. Query is only supported in Terraform v1.14+
func ExpectKnownValue(listResourceAddress string, resourceName string, attributePath tfjsonpath.Path, knownValue knownvalue.Check) QueryResultCheck {
	return expectKnownValue{
		listResourceAddress: listResourceAddress,
		resourceName:        resourceName,
		attributePath:       attributePath,
		knownValue:          knownValue,
	}
}
