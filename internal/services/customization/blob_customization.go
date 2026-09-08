package customization

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

const minimumStorageOAuthVersion = "2017-11-09"

var blobBodyFields = map[string]bool{
	"content_base64":      true,
	"content_type":        true,
	"content_encoding":    true,
	"content_language":    true,
	"content_disposition": true,
	"cache_control":       true,
	"metadata":            true,
	"access_tier":         true,
}

type BlobCustomization struct{}

func (BlobCustomization) GetResourceType() string {
	return "Microsoft.Storage/storageAccounts/blobServices/containers/blobs"
}

func (BlobCustomization) RequiresExclusiveBody() bool {
	return true
}

func (BlobCustomization) CreateFunc() CreateFunc {
	return putBlob
}

func (BlobCustomization) UpdateFunc() UpdateFunc {
	return putBlob
}

func putBlob(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
	content, headers, err := buildBlobRequest(body, id.ApiVersion, options.Headers)
	if err != nil {
		return err
	}
	options.Headers = headers
	resp, err := client.DataPlaneClient.DoRaw(ctx, http.MethodPut, id.AzureResourceId, clients.RawRequestOptions{
		RequestOptions:   options,
		APIVersionHeader: id.ApiVersion,
		Content:          content,
		SuccessCodes:     []int{http.StatusCreated},
	})
	if resp != nil {
		runtime.Drain(resp)
	}
	return err
}

func (BlobCustomization) ReadFunc() ReadFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) (interface{}, error) {
		if err := validateStorageVersion(id.ApiVersion); err != nil {
			return nil, err
		}
		if err := validateReservedBlobHeaders(options.Headers, false); err != nil {
			return nil, err
		}
		resp, err := client.DataPlaneClient.DoRaw(ctx, http.MethodHead, id.AzureResourceId, clients.RawRequestOptions{
			RequestOptions:   options,
			APIVersionHeader: id.ApiVersion,
			SuccessCodes:     []int{http.StatusOK},
		})
		if err != nil {
			return nil, err
		}
		defer runtime.Drain(resp)
		return blobProperties(resp.Header)
	}
}

func (BlobCustomization) DeleteFunc() DeleteFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) error {
		if err := validateStorageVersion(id.ApiVersion); err != nil {
			return err
		}
		if err := validateReservedBlobHeaders(options.Headers, false); err != nil {
			return err
		}
		resp, err := client.DataPlaneClient.DoRaw(ctx, http.MethodDelete, id.AzureResourceId, clients.RawRequestOptions{
			RequestOptions:   options,
			APIVersionHeader: id.ApiVersion,
			SuccessCodes:     []int{http.StatusAccepted, http.StatusNotFound},
		})
		if resp != nil {
			runtime.Drain(resp)
		}
		return err
	}
}

func buildBlobRequest(body interface{}, apiVersion string, operationHeaders map[string]string) ([]byte, map[string]string, error) {
	if err := validateStorageVersion(apiVersion); err != nil {
		return nil, nil, err
	}
	if err := validateReservedBlobHeaders(operationHeaders, true); err != nil {
		return nil, nil, err
	}
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("blob body must be an object")
	}
	for key := range bodyMap {
		if !blobBodyFields[key] {
			return nil, nil, fmt.Errorf("unsupported blob body field %q", key)
		}
	}

	encoded, ok := bodyMap["content_base64"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("blob body field %q must be a Base64-encoded string", "content_base64")
	}
	content, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid content_base64: %w", err)
	}

	headers := make(map[string]string, len(operationHeaders)+8)
	setOptionalStringHeader := func(field, header string) error {
		value, exists := bodyMap[field]
		if !exists || value == nil {
			return nil
		}
		text, ok := value.(string)
		if !ok {
			return fmt.Errorf("blob body field %q must be a string", field)
		}
		headers[header] = text
		return nil
	}
	for field, header := range map[string]string{
		"content_type":        "Content-Type",
		"content_encoding":    "Content-Encoding",
		"content_language":    "Content-Language",
		"content_disposition": "x-ms-blob-content-disposition",
		"cache_control":       "Cache-Control",
		"access_tier":         "x-ms-access-tier",
	} {
		if err := setOptionalStringHeader(field, header); err != nil {
			return nil, nil, err
		}
	}
	if value, exists := bodyMap["access_tier"]; exists && value != nil {
		tier, ok := value.(string)
		if !ok {
			return nil, nil, fmt.Errorf("blob body field %q must be a string", "access_tier")
		}
		switch strings.ToLower(tier) {
		case "hot", "cool", "cold", "archive":
		default:
			return nil, nil, fmt.Errorf("blob body field %q must be one of Hot, Cool, Cold, or Archive", "access_tier")
		}
	}
	if value, exists := bodyMap["metadata"]; exists && value != nil {
		metadata, ok := value.(map[string]interface{})
		if !ok {
			return nil, nil, fmt.Errorf("blob body field %q must be an object of string values", "metadata")
		}
		for key, value := range metadata {
			text, ok := value.(string)
			if !ok {
				return nil, nil, fmt.Errorf("blob metadata %q must be a string", key)
			}
			headers["x-ms-meta-"+key] = text
		}
	}
	for key, value := range operationHeaders {
		headers[key] = value
	}
	headers["x-ms-blob-type"] = "BlockBlob"
	return content, headers, nil
}

func validateStorageVersion(apiVersion string) error {
	version, err := time.Parse("2006-01-02", apiVersion)
	if err != nil {
		return fmt.Errorf("invalid Blob Storage service version %q: expected YYYY-MM-DD", apiVersion)
	}
	minimum, _ := time.Parse("2006-01-02", minimumStorageOAuthVersion)
	if version.Before(minimum) {
		return fmt.Errorf("Blob Storage service version %q does not support Microsoft Entra authentication; use %s or later", apiVersion, minimumStorageOAuthVersion)
	}
	return nil
}

func validateReservedBlobHeaders(headers map[string]string, writing bool) error {
	reserved := map[string]bool{"x-ms-version": true}
	if writing {
		reserved["x-ms-blob-type"] = true
		reserved["content-length"] = true
	}
	for key := range headers {
		if reserved[strings.ToLower(key)] {
			return fmt.Errorf("header %q is managed by the Blob Storage customization and cannot be overridden", key)
		}
	}
	return nil
}

func blobProperties(headers http.Header) (map[string]interface{}, error) {
	contentLength, err := strconv.ParseInt(headers.Get("Content-Length"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid Blob Storage Content-Length response header %q: %w", headers.Get("Content-Length"), err)
	}
	metadata := make(map[string]interface{})
	for key, values := range headers {
		if strings.HasPrefix(strings.ToLower(key), "x-ms-meta-") && len(values) != 0 {
			metadata[strings.ToLower(key[len("x-ms-meta-"):])] = values[0]
		}
	}
	return map[string]interface{}{
		"content_type":        nullableHeader(headers, "Content-Type"),
		"content_encoding":    nullableHeader(headers, "Content-Encoding"),
		"content_language":    nullableHeader(headers, "Content-Language"),
		"content_disposition": nullableHeader(headers, "Content-Disposition"),
		"cache_control":       nullableHeader(headers, "Cache-Control"),
		"metadata":            metadata,
		"access_tier":         nullableHeader(headers, "x-ms-access-tier"),
		"etag":                nullableHeader(headers, "ETag"),
		"content_length":      contentLength,
		"content_md5":         nullableHeader(headers, "Content-MD5"),
		"blob_type":           nullableHeader(headers, "x-ms-blob-type"),
		"last_modified":       nullableHeader(headers, "Last-Modified"),
	}, nil
}

func nullableHeader(headers http.Header, name string) interface{} {
	if value := headers.Get(name); value != "" {
		return value
	}
	return nil
}

var _ DataPlaneResource = &BlobCustomization{}
var _ DataPlaneResourceWithExclusiveBody = &BlobCustomization{}
