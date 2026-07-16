package customization

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

type KeyVaultCertificateCustomization struct{}

func (k KeyVaultCertificateCustomization) CreateFunc() CreateFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
		if _, err := client.DataPlaneClient.Action(ctx, id.AzureResourceId, "create", id.ApiVersion, http.MethodPost, body, options); err != nil {
			return err
		}
		// Certificate creation is a long running process without LRO request header. So we poll the state here.
		return waitForCertificateOperation(ctx, client, id, options)
	}
}

const certificateOperationPollFrequency = 10 * time.Second

// waitForCertificateOperation polls the state of certification creation until it completed or failed.
func waitForCertificateOperation(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) error {
	for {
		result, err := client.DataPlaneClient.Action(ctx, id.AzureResourceId, "pending", id.ApiVersion, http.MethodGet, nil, options)
		if err != nil {
			return fmt.Errorf("polling certificate operation for %s: %+v", id.AzureResourceId, err)
		}

		operation, ok := result.(map[string]interface{})
		if !ok {
			return fmt.Errorf("polling certificate operation for %s: unexpected response, expected a JSON object but got %T", id.AzureResourceId, result)
		}

		status, ok := operation["status"].(string)
		if !ok {
			return fmt.Errorf("polling certificate operation for %s: response is missing a string \"status\" field", id.AzureResourceId)
		}
		status = strings.ToLower(status)
		if status == "completed" {
			return nil
		}
		if status == "failed" || status == "cancelled" || status == "deleted" {
			return fmt.Errorf("certificate operation for %s did not succeed (status %q): %s", id.AzureResourceId, status, certificateOperationError(operation))
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(certificateOperationPollFrequency):
		}
	}
}

func certificateOperationError(operation map[string]interface{}) string {
	errObj, ok := operation["error"].(map[string]interface{})
	if !ok {
		return "no error details provided"
	}
	code, _ := errObj["code"].(string)
	message, _ := errObj["message"].(string)

	if code != "" && message != "" {
		return fmt.Sprintf("error code %s: %s", code, message)
	}
	if message != "" {
		return message
	}
	if code != "" {
		return fmt.Sprintf("error code %s", code)
	}
	return "no error details provided"
}

func (k KeyVaultCertificateCustomization) UpdateFunc() UpdateFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
		_, err := client.DataPlaneClient.Action(ctx, id.AzureResourceId, "", id.ApiVersion, http.MethodPatch, body, options)
		return err
	}
}

func (k KeyVaultCertificateCustomization) ReadFunc() ReadFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) (interface{}, error) {
		return client.DataPlaneClient.Get(ctx, id, options)
	}
}

func (k KeyVaultCertificateCustomization) DeleteFunc() DeleteFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) error {
		_, err := client.DataPlaneClient.DeleteThenPoll(ctx, id, options)
		return err
	}
}

func (k KeyVaultCertificateCustomization) GetResourceType() string {
	return "Microsoft.KeyVault/vaults/certificates"
}

var _ DataPlaneResource = &KeyVaultCertificateCustomization{}
