package services

import (
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/stretchr/testify/require"
)

func TestBodyRefsChanged(t *testing.T) {
	stateBody, err := dynamic.FromJSONImplied([]byte(`{"policy":{"key_props":{"kty":"RSA","key_size":2048}}}`))
	require.NoError(t, err)
	planBody, err := dynamic.FromJSONImplied([]byte(`{"policy":{"key_props":{"kty":"RSA","key_size":4096}}}`))
	require.NoError(t, err)
	refs := []string{"policy.key_props.kty", "policy.key_props.key_size"}

	changed, err := bodyRefsChanged(stateBody, planBody, refs)
	require.NoError(t, err)
	require.True(t, changed)

	changed, err = bodyRefsChanged(stateBody, stateBody, refs)
	require.NoError(t, err)
	require.False(t, changed)

	_, err = bodyRefsChanged(stateBody, planBody, []string{"policy.["})
	require.Error(t, err)
}
