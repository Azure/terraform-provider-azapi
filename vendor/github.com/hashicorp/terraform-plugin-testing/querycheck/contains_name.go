// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"
	"strings"
)

var _ QueryResultCheck = contains{}

type contains struct {
	resourceAddress string
	check           string
}

func (c contains) CheckQuery(_ context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	for _, res := range req.Query {
		if strings.EqualFold(c.check, res.DisplayName) {
			return
		}
	}

	resp.Error = fmt.Errorf("expected to find resource with display name %q in results but resource was not found", c.check)

}

// ContainsResourceWithName returns a query check that asserts that a resource with a given display name exists within the returned results of the query.
//
// This query check can only be used with managed resources that support query. Query is only supported in Terraform v1.14+
func ContainsResourceWithName(resourceAddress string, displayName string) QueryResultCheck {
	return contains{
		resourceAddress: resourceAddress,
		check:           displayName,
	}
}
