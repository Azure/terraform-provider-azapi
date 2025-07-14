package customization

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

type Resource interface {
	GetResourceType() string
	CreateFunc() Func
	UpdateFunc() Func
	DeleteFunc() DeleteFunc
	ReadFunc() Func
}

type Func func() error

type DeleteFunc func(ctx context.Context, clients clients.Client, id parse.ResourceId, options clients.RequestOptions) error

var customizations = make(map[string]Resource)

func init() {
	var keyVaultKeyCustomization Resource = KeyVaultKeyCustomization{}
	customizations[keyVaultKeyCustomization.GetResourceType()] = keyVaultKeyCustomization

}

func GetCustomization(resourceType string) *Resource {
	customization, exists := customizations[resourceType]
	if !exists {
		return nil
	}
	return &customization
}
