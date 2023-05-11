package provider

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type MsiCredential struct {
	cred *azidentity.ManagedIdentityCredential
}

func NewManagedIdentityCredential(options *azidentity.ManagedIdentityCredentialOptions) (*MsiCredential, error) {
	miCred, err := azidentity.NewManagedIdentityCredential(options)
	if err != nil {
		return nil, err
	}

	w := &MsiCredential{
		cred: miCred,
	}

	return w, nil
}

func (w *MsiCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	tk, err := w.cred.GetToken(c, opts)
	if ctxErr := c.Err(); errors.Is(ctxErr, context.DeadlineExceeded) {
		// timeout: signal the chain to try its next credential, if any
		err = azidentity.NewCredentialUnavailableError("managed identity timed out")
	}
	return tk, err
}
