package acceptance

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type TestResource interface {
	Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error)
}

type TestResourceVerifyingRemoved interface {
	TestResource
	Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error)
}
