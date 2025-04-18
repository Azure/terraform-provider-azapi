package services

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const (
	pkEphemeralBody = "sensitive_body"
)

type PrivateData interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

type EphemeralBodyPrivateMgr struct{}

var ephemeralBodyPrivateMgr = EphemeralBodyPrivateMgr{}

func (m EphemeralBodyPrivateMgr) Exists(ctx context.Context, d PrivateData) (bool, diag.Diagnostics) {
	b, diags := d.GetKey(ctx, pkEphemeralBody)
	if diags.HasError() {
		return false, diags
	}
	return b != nil, diags
}

// Set sets the hash of the sensitive_body to the private state.
// If `ebody` is nil, it removes the hash from the private state.
func (m EphemeralBodyPrivateMgr) Set(ctx context.Context, d PrivateData, ebody []byte) (diags diag.Diagnostics) {
	if ebody == nil {
		d.SetKey(ctx, pkEphemeralBody, nil)
		return
	}

	// Calculate the hash of the ephemeral body
	h := sha256.New()
	if _, err := h.Write(ebody); err != nil {
		diags.AddError(
			`Error to hash "sensitive_body"`,
			err.Error(),
		)
		return
	}
	hash := h.Sum(nil)

	b, err := json.Marshal(map[string]interface{}{
		// []byte will be marshaled to base64 encoded string
		"hash": hash,
	})
	if err != nil {
		diags.AddError(
			`Error to marshal "sensitive_body" private data`,
			err.Error(),
		)
		return
	}

	return d.SetKey(ctx, pkEphemeralBody, b)
}

// Diff tells whether the sensitive_body is different than the hash stored in the private state.
// In case private state doesn't have the record, regard the record as "nil" (i.e. will return true if ebody is non-nil).
// In case private state has the record (guaranteed to be non-nil), while ebody is nil, it also returns true.
func (m EphemeralBodyPrivateMgr) Diff(ctx context.Context, d PrivateData, ebody []byte) (bool, diag.Diagnostics) {
	b, diags := d.GetKey(ctx, pkEphemeralBody)
	if diags.HasError() {
		return false, diags
	}
	if b == nil {
		// In case private state doesn't store the key yet, it only diffs when the ebody is not nil.
		return ebody != nil, diags
	}
	var mm map[string]interface{}
	if err := json.Unmarshal(b, &mm); err != nil {
		diags.AddError(
			`Error to unmarshal "sensitive_body" private data`,
			err.Error(),
		)
		return false, diags
	}
	privateHashEnc, ok := mm["hash"]
	if !ok {
		diags.AddError(
			`Invalid "sensitive_body" private data`,
			`Key "hash" not found`,
		)
		return false, diags
	}

	h := sha256.New()
	if _, err := h.Write(ebody); err != nil {
		diags.AddError(
			`Error to hash "sensitive_body"`,
			err.Error(),
		)
		return false, diags
	}
	hash := h.Sum(nil)

	hashEnc := base64.StdEncoding.EncodeToString(hash)

	return hashEnc != privateHashEnc.(string), diags
}
