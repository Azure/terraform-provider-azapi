package services

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/require"
)

type privateDataMock struct {
	data map[string][]byte
}

func (m *privateDataMock) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return m.data[key], nil
}

func (m *privateDataMock) SetKey(_ context.Context, key string, value []byte) diag.Diagnostics {
	if value == nil {
		delete(m.data, key)
	} else {
		m.data[key] = value
	}
	return nil
}

func TestIgnoreBodyChangesPrivateMgr(t *testing.T) {
	ctx := context.Background()
	privateData := &privateDataMock{data: make(map[string][]byte)}
	paths := []string{"tags", "properties.sku.name"}

	different, diags := ignoreBodyChangesPrivateMgr.Diff(ctx, privateData, paths)
	require.False(t, diags.HasError())
	require.True(t, different)

	diags = ignoreBodyChangesPrivateMgr.Set(ctx, privateData, paths)
	require.False(t, diags.HasError())

	actual, diags := ignoreBodyChangesPrivateMgr.Get(ctx, privateData)
	require.False(t, diags.HasError())
	require.Equal(t, paths, actual)

	different, diags = ignoreBodyChangesPrivateMgr.Diff(ctx, privateData, paths)
	require.False(t, diags.HasError())
	require.False(t, different)

	different, diags = ignoreBodyChangesPrivateMgr.Diff(ctx, privateData, []string{"properties.sku.name", "tags"})
	require.False(t, diags.HasError())
	require.False(t, different)

	diags = ignoreBodyChangesPrivateMgr.Set(ctx, privateData, nil)
	require.False(t, diags.HasError())
	require.NotContains(t, privateData.data, pkIgnoreBodyChanges)
}

func TestIgnoreBodyChangesPrivateMgrInvalidData(t *testing.T) {
	ctx := context.Background()
	privateData := &privateDataMock{
		data: map[string][]byte{pkIgnoreBodyChanges: []byte("invalid")},
	}

	_, diags := ignoreBodyChangesPrivateMgr.Get(ctx, privateData)
	require.True(t, diags.HasError())
}

func TestOverrideBodyWithPaths(t *testing.T) {
	old := map[string]interface{}{
		"tags": map[string]interface{}{"managed": "terraform"},
		"properties": map[string]interface{}{
			"sku":  map[string]interface{}{"name": "old", "tier": "standard"},
			"name": "configured",
		},
	}
	newValue := map[string]interface{}{
		"tags": map[string]interface{}{"managed": "external"},
		"properties": map[string]interface{}{
			"sku":  map[string]interface{}{"name": "new", "tier": "premium"},
			"name": "remote",
		},
	}

	actual, err := overrideBodyWithPaths(old, newValue, []string{"tags", "properties.sku.name"})
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{
		"tags": map[string]interface{}{"managed": "external"},
		"properties": map[string]interface{}{
			"sku":  map[string]interface{}{"name": "new", "tier": "standard"},
			"name": "configured",
		},
	}, actual)
}

// The update request keeps an ignored property the configuration leaves out or
// sets to null, so a route table attached outside Terraform survives an
// unrelated update (GH-1270). The plan and read paths keep using
// overrideBodyWithPaths, which only touches paths present in the body.
func TestOverrideBodyWithExistingPaths(t *testing.T) {
	existing := map[string]interface{}{
		"properties": map[string]interface{}{
			"addressPrefixes":   []interface{}{"10.250.1.128/26"},
			"routeTable":        map[string]interface{}{"id": "rt"},
			"provisioningState": "Succeeded",
		},
	}
	expected := map[string]interface{}{
		"properties": map[string]interface{}{
			"addressPrefixes":  []interface{}{"10.250.1.128/26"},
			"routeTable":       map[string]interface{}{"id": "rt"},
			"serviceEndpoints": []interface{}{map[string]interface{}{"service": "Microsoft.KeyVault"}},
		},
	}

	for name, routeTable := range map[string]bool{"omitted": false, "null": true} {
		t.Run(name, func(t *testing.T) {
			properties := map[string]interface{}{
				"addressPrefixes":  []interface{}{"10.250.1.128/26"},
				"serviceEndpoints": []interface{}{map[string]interface{}{"service": "Microsoft.KeyVault"}},
			}
			if routeTable {
				properties["routeTable"] = nil
			}
			body := map[string]interface{}{"properties": properties}

			actual, err := overrideBodyWithExistingPaths(body, existing, []string{"properties.routeTable"})
			require.NoError(t, err)
			require.Equal(t, expected, actual)

			// The plan/read helper still leaves a path the body does not set alone.
			unchanged, err := overrideBodyWithPaths(body, existing, []string{"properties.routeTable"})
			require.NoError(t, err)
			require.Equal(t, body, unchanged)
		})
	}
}
