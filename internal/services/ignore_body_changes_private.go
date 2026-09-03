package services

import (
	"context"
	"encoding/json"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const pkIgnoreBodyChanges = "ignore_body_changes"

type IgnoreBodyChangesPrivateMgr struct{}

var ignoreBodyChangesPrivateMgr = IgnoreBodyChangesPrivateMgr{}

func (m IgnoreBodyChangesPrivateMgr) Get(ctx context.Context, d PrivateData) ([]string, diag.Diagnostics) {
	b, diags := d.GetKey(ctx, pkIgnoreBodyChanges)
	if diags.HasError() || b == nil {
		return nil, diags
	}

	var paths []string
	if err := json.Unmarshal(b, &paths); err != nil {
		diags.AddError(`Invalid "ignore_body_changes" private data`, err.Error())
		return nil, diags
	}
	return paths, diags
}

func (m IgnoreBodyChangesPrivateMgr) Set(ctx context.Context, d PrivateData, paths []string) diag.Diagnostics {
	if len(paths) == 0 {
		return d.SetKey(ctx, pkIgnoreBodyChanges, nil)
	}

	b, err := json.Marshal(paths)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(`Error marshaling "ignore_body_changes" private data`, err.Error())
		return diags
	}
	return d.SetKey(ctx, pkIgnoreBodyChanges, b)
}

func (m IgnoreBodyChangesPrivateMgr) Diff(ctx context.Context, d PrivateData, paths []string) (bool, diag.Diagnostics) {
	storedPaths, diags := m.Get(ctx, d)
	if diags.HasError() {
		return false, diags
	}
	storedPaths = slices.Clone(storedPaths)
	paths = slices.Clone(paths)
	slices.Sort(storedPaths)
	slices.Sort(paths)
	return !slices.Equal(storedPaths, paths), diags
}
