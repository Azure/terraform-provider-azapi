package retry

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestValueType(t *testing.T) {
	ctx := context.Background()
	v := RetryValue{}
	ty := v.Type(ctx)
	_, ok := ty.(basetypes.ObjectType)
	assert.True(t, ok)
}
