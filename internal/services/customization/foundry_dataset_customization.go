package customization

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

// FoundryDatasetCustomization manages a Foundry dataset version and its
// provider-side source upload.
//
// The request body uses Foundry API field names:
//
//	name
//	version
//	description
//	type
//	format
//	dataUri
//
// The provider-only fields are:
//
//	source_url
//	source_sha256
//	computed_sha256
//	pending_upload_id
//	connection_name
//
// The downloaded file is streamed through a temporary file, verified, and
// uploaded to the container SAS returned by startPendingUpload. The file
// contents are never placed in the Foundry API request body.
type FoundryDatasetCustomization struct{}

func (c FoundryDatasetCustomization) GetResourceType() string {
	return "Microsoft.Foundry/datasets/versions"
}

func datasetVersionAction(
	id parse.DataPlaneResourceId,
	action string,
) string {
	return strings.TrimSuffix(id.AzureResourceId, "/") + "/" + action
}

func datasetMap(value interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshalling dataset value: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling dataset value: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("dataset body must be an object")
	}

	return result, nil
}

func datasetField(
	values map[string]interface{},
	names ...string,
) (interface{}, string, bool) {
	for _, name := range names {
		if value, ok := values[name]; ok {
			return value, name, true
		}
	}

	for key, value := range values {
		for _, name := range names {
			if strings.EqualFold(key, name) {
				return value, key, true
			}
		}
	}

	return nil, "", false
}

func datasetStringField(
	values map[string]interface{},
	names ...string,
) (string, string, bool, error) {
	value, name, exists := datasetField(values, names...)
	if !exists || value == nil {
		return "", name, false, nil
	}

	switch value := value.(type) {
	case string:
		return value, name, true, nil

	case float64:
		if value != float64(int64(value)) {
			return "", name, true, fmt.Errorf(
				"dataset field %q must be a string",
				name,
			)
		}

		return strconv.FormatInt(int64(value), 10), name, true, nil

	default:
		return "", name, true, fmt.Errorf(
			"dataset field %q must be a string",
			name,
		)
	}
}

func datasetRequiredString(
	values map[string]interface{},
	name string,
) (string, error) {
	value, _, exists, err := datasetStringField(values, name)
	if err != nil {
		return "", err
	}

	value = strings.TrimSpace(value)
	if !exists || value == "" {
		return "", fmt.Errorf(
			`dataset body field %q is required`,
			name,
		)
	}

	return value, nil
}

func datasetSourceInfo(
	body interface{},
) (string, string, bool, error) {
	values, err := datasetMap(body)
	if err != nil {
		return "", "", false, err
	}

	sourceURL, sourceURLField, exists, err := datasetStringField(
		values,
		"source_url",
		"sourceUrl",
	)
	if err != nil {
		return "", "", false, err
	}

	sourceURL = strings.TrimSpace(sourceURL)
	if !exists || sourceURL == "" {
		return "", "", false, fmt.Errorf(
			`dataset body field "source_url" is required`,
		)
	}

	if err := validateDatasetSourceURL(sourceURL); err != nil {
		return "", "", false, fmt.Errorf(
			`dataset body field %q: %w`,
			sourceURLField,
			err,
		)
	}

	expectedSHA256, _, hasSHA256, err := datasetStringField(
		values,
		"source_sha256",
		"sourceSha256",
	)
	if err != nil {
		return "", "", false, err
	}

	// An omitted or empty checksum disables checksum verification.
	if !hasSHA256 || strings.TrimSpace(expectedSHA256) == "" {
		return sourceURL, "", false, nil
	}

	expectedSHA256 = strings.ToLower(strings.TrimSpace(expectedSHA256))

	if len(expectedSHA256) != sha256.Size*2 {
		return "", "", false, fmt.Errorf(
			`dataset body field "source_sha256" must contain 64 hexadecimal characters`,
		)
	}

	if _, err := hex.DecodeString(expectedSHA256); err != nil {
		return "", "", false, fmt.Errorf(
			`dataset body field "source_sha256" must be hexadecimal: %w`,
			err,
		)
	}

	return sourceURL, expectedSHA256, true, nil
}

func validateDatasetSourceURL(rawURL string) error {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("must be an HTTPS URL")
	}

	if !strings.EqualFold(parsedURL.Scheme, "https") ||
		parsedURL.Hostname() == "" ||
		parsedURL.User != nil ||
		parsedURL.Fragment != "" ||
		(parsedURL.Port() != "" && parsedURL.Port() != "443") {
		return fmt.Errorf("must be an HTTPS URL")
	}

	host := strings.ToLower(parsedURL.Hostname())

	if host == "raw.github.kp.org" || isAzureBlobHost(host) {
		return nil
	}

	return fmt.Errorf(
		"must use raw.github.kp.org or an Azure Blob Storage host",
	)
}

func isAzureBlobHost(host string) bool {
	for _, suffix := range []string{
		".blob.core.windows.net",
	} {
		if strings.HasSuffix(host, suffix) &&
			len(host) > len(suffix) {
			return true
		}
	}

	return false
}

func datasetPendingUploadBody(
	body interface{},
) (map[string]interface{}, error) {
	values, err := datasetMap(body)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"pendingUploadType": "BlobReference",
	}

	if value, _, exists, err := datasetStringField(
		values,
		"pending_upload_id",
		"pendingUploadId",
	); err != nil {
		return nil, err
	} else if exists && strings.TrimSpace(value) != "" {
		result["pendingUploadId"] = strings.TrimSpace(value)
	}

	if value, _, exists, err := datasetStringField(
		values,
		"connection_name",
		"connectionName",
	); err != nil {
		return nil, err
	} else if exists && strings.TrimSpace(value) != "" {
		result["connectionName"] = strings.TrimSpace(value)
	}

	return result, nil
}

func setDatasetDefaults(
	body interface{},
	versionFallback string,
) error {
	values, ok := body.(map[string]interface{})
	if !ok || values == nil {
		return fmt.Errorf("dataset body must be a mutable object")
	}

	// The Foundry API field is "version". The resource path version is used
	// as the fallback because the generic AzAPI resource uses "name" for the
	// final path segment.
	version, _, exists, err := datasetStringField(values, "version")
	if err != nil {
		return err
	}

	version = strings.TrimSpace(version)

	if !exists || version == "" {
		version = strings.TrimSpace(versionFallback)

		if version == "" || version == "__generated__" {
			version = "1"
		}

		values["version"] = version
	}

	datasetType, _, exists, err := datasetStringField(values, "type")
	if err != nil {
		return err
	}

	if !exists || strings.TrimSpace(datasetType) == "" {
		values["type"] = "uri_file"
	}

	// Format is a required Foundry API field.
	format, _, exists, err := datasetStringField(values, "format")
	if err != nil {
		return err
	}

	if !exists || strings.TrimSpace(format) == "" {
		return fmt.Errorf(
			`dataset body field "format" is required`,
		)
	}

	return nil
}

func datasetVersionRequestBody(
	body interface{},
	versionFallback string,
	dataURI string,
) (map[string]interface{}, string, error) {
	values, err := datasetMap(body)
	if err != nil {
		return nil, "", err
	}

	// Foundry API field: name
	name, err := datasetRequiredString(values, "name")
	if err != nil {
		return nil, "", err
	}

	// Foundry API field: version
	version, _, exists, err := datasetStringField(values, "version")
	if err != nil {
		return nil, "", err
	}

	version = strings.TrimSpace(version)

	if !exists || version == "" {
		version = strings.TrimSpace(versionFallback)

		if version == "" || version == "__generated__" {
			version = "1"
		}

		values["version"] = version
	}

	// Foundry API field: description
	description, err := datasetRequiredString(values, "description")
	if err != nil {
		return nil, "", err
	}

	// Foundry API field: type
	datasetType, _, exists, err := datasetStringField(values, "type")
	if err != nil {
		return nil, "", err
	}

	datasetType = strings.TrimSpace(datasetType)

	if !exists || datasetType == "" {
		datasetType = "uri_file"
	}

	if datasetType != "uri_file" && datasetType != "uri_folder" {
		return nil, "", fmt.Errorf(
			`dataset body field "type" must be "uri_file" or "uri_folder"`,
		)
	}

	// Foundry API field: format
	format, err := datasetRequiredString(values, "format")
	if err != nil {
		return nil, "", err
	}

	// Provider-only fields such as source_url and source_sha256 are
	// intentionally excluded from the Foundry API request body.
	return map[string]interface{}{
		"name":        name,
		"version":     version,
		"description": description,
		"type":        datasetType,
		"dataUri":     dataURI,
		"format":      format,
	}, version, nil
}

func datasetResponseString(
	values map[string]interface{},
	name string,
) (string, error) {
	value, exists := values[name]
	if !exists || value == nil {
		return "", fmt.Errorf(
			"dataset response field %q is missing",
			name,
		)
	}

	result, ok := value.(string)
	if !ok || strings.TrimSpace(result) == "" {
		return "", fmt.Errorf(
			"dataset response field %q must be a non-empty string",
			name,
		)
	}

	return strings.TrimSpace(result), nil
}

func datasetNestedMap(
	values map[string]interface{},
	name string,
) (map[string]interface{}, error) {
	value, exists := values[name]
	if !exists || value == nil {
		return nil, fmt.Errorf(
			"dataset response field %q is missing",
			name,
		)
	}

	result, err := datasetMap(value)
	if err != nil {
		return nil, fmt.Errorf(
			"dataset response field %q must be an object: %w",
			name,
			err,
		)
	}

	return result, nil
}

func datasetUploadDetails(
	response interface{},
) (string, string, error) {
	values, err := datasetMap(response)
	if err != nil {
		return "", "", fmt.Errorf(
			"invalid startPendingUpload response: %w",
			err,
		)
	}

	blobReference, err := datasetNestedMap(values, "blobReference")
	if err != nil {
		return "", "", err
	}

	credential, err := datasetNestedMap(blobReference, "credential")
	if err != nil {
		return "", "", err
	}

	uploadURL, err := datasetResponseString(credential, "sasUri")
	if err != nil {
		return "", "", err
	}

	if err := validateDatasetSASURL(uploadURL); err != nil {
		return "", "", fmt.Errorf(
			"dataset response field blobReference.credential.sasUri: %w",
			err,
		)
	}

	dataURIValues := values

	if consumption, ok := values["blobReferenceForConsumption"]; ok &&
		consumption != nil {
		dataURIValues, err = datasetMap(consumption)
		if err != nil {
			return "", "", fmt.Errorf(
				"invalid blobReferenceForConsumption: %w",
				err,
			)
		}
	}

	dataURI, err := datasetResponseString(dataURIValues, "blobUri")
	if err != nil {
		return "", "", err
	}

	return uploadURL, dataURI, nil
}

func validateDatasetSASURL(rawURL string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil ||
		!strings.EqualFold(parsedURL.Scheme, "https") ||
		parsedURL.Hostname() == "" ||
		parsedURL.RawQuery == "" {
		return fmt.Errorf("must be an HTTPS SAS URL")
	}

	return nil
}

func datasetSourceFilename(sourceURL string) (string, error) {
	parsedURL, err := url.Parse(sourceURL)
	if err != nil {
		return "", fmt.Errorf("parsing source_url: %w", err)
	}

	if parsedURL.Path == "" || strings.HasSuffix(parsedURL.Path, "/") {
		return "", fmt.Errorf("source_url must identify a file")
	}

	filename := path.Base(parsedURL.Path)
	if filename == "" || filename == "." || filename == "/" {
		return "", fmt.Errorf("source_url must identify a file")
	}

	return filename, nil
}

func datasetBlobUploadURL(
	containerSASURL string,
	filename string,
) (string, error) {
	parsedURL, err := url.Parse(containerSASURL)
	if err != nil {
		return "", fmt.Errorf(
			"parsing upload SAS URL: %w",
			err,
		)
	}

	query := parsedURL.Query()

	// A blob-level SAS URL can be used directly.
	if strings.EqualFold(query.Get("sr"), "b") {
		return parsedURL.String(), nil
	}

	// Otherwise, the returned SAS URL identifies a container. Append the
	// source filename to create the blob-level upload URL.
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/") +
		"/" +
		filename
	parsedURL.RawPath = ""

	return parsedURL.String(), nil
}

func datasetSourceHTTPClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(
			request *http.Request,
			_ []*http.Request,
		) error {
			if err := validateDatasetSourceURL(request.URL.String()); err != nil {
				return fmt.Errorf(
					"refusing redirect for source_url: %w",
					err,
				)
			}

			return nil
		},
	}
}

func datasetUploadHTTPClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(
			*http.Request,
			[]*http.Request,
		) error {
			return fmt.Errorf(
				"refusing redirect from dataset upload SAS URL",
			)
		},
	}
}

func streamDatasetToUpload(
	ctx context.Context,
	sourceURL string,
	containerSASURL string,
	expectedSHA256 string,
	verifySHA256 bool,
) (string, error) {
	filename, err := datasetSourceFilename(sourceURL)
	if err != nil {
		return "", err
	}

	uploadURL, err := datasetBlobUploadURL(containerSASURL, filename)
	if err != nil {
		return "", err
	}

	sourceRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		sourceURL,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf(
			"creating dataset download request: %w",
			err,
		)
	}

	sourceRequest.Header.Set("Accept-Encoding", "identity")

	sourceResponse, err := datasetSourceHTTPClient().Do(sourceRequest)
	if err != nil {
		return "", fmt.Errorf(
			"downloading dataset: %w",
			err,
		)
	}

	defer sourceResponse.Body.Close()

	if sourceResponse.StatusCode < http.StatusOK ||
		sourceResponse.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf(
			"downloading dataset returned HTTP %s",
			sourceResponse.Status,
		)
	}

	tempFile, err := os.CreateTemp("", "azapi-dataset-*")
	if err != nil {
		return "", fmt.Errorf(
			"creating temporary dataset file: %w",
			err,
		)
	}

	tempFileName := tempFile.Name()
	defer os.Remove(tempFileName)
	defer tempFile.Close()

	hasher := sha256.New()

	if _, err := io.Copy(
		io.MultiWriter(tempFile, hasher),
		sourceResponse.Body,
	); err != nil {
		return "", fmt.Errorf(
			"downloading dataset: %w",
			err,
		)
	}

	actualSHA256 := hex.EncodeToString(hasher.Sum(nil))

	if verifySHA256 &&
		!strings.EqualFold(actualSHA256, expectedSHA256) {
		return "", fmt.Errorf(
			"dataset SHA-256 mismatch: expected %s, got %s",
			expectedSHA256,
			actualSHA256,
		)
	}

	if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf(
			"rewinding temporary dataset file: %w",
			err,
		)
	}

	uploadRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		uploadURL,
		tempFile,
	)
	if err != nil {
		return "", fmt.Errorf(
			"creating dataset upload request: %w",
			err,
		)
	}

	fileInfo, err := tempFile.Stat()
	if err != nil {
		return "", fmt.Errorf(
			"reading temporary dataset file metadata: %w",
			err,
		)
	}

	uploadRequest.ContentLength = fileInfo.Size()
	uploadRequest.Header.Set(
		"Content-Type",
		"application/octet-stream",
	)
	uploadRequest.Header.Set(
		"x-ms-blob-type",
		"BlockBlob",
	)

	uploadResponse, err := datasetUploadHTTPClient().Do(uploadRequest)
	if err != nil {
		return "", fmt.Errorf(
			"uploading dataset: %w",
			err,
		)
	}

	defer uploadResponse.Body.Close()

	if uploadResponse.StatusCode < http.StatusOK ||
		uploadResponse.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf(
			"uploading dataset returned HTTP %s",
			uploadResponse.Status,
		)
	}

	return actualSHA256, nil
}

func (c FoundryDatasetCustomization) createOrUpdate(
	ctx context.Context,
	client clients.Client,
	id parse.DataPlaneResourceId,
	body interface{},
	options clients.RequestOptions,
) (interface{}, error) {
	resourceVersion := strings.TrimSpace(id.Name)

	if resourceVersion == "" ||
		resourceVersion == "__generated__" {
		return nil, fmt.Errorf(
			`resource-level "name" is required for Foundry dataset versions; set it to the dataset version, for example name = "1"`,
		)
	}

	// Never mutate the body supplied by Terraform. The body is a dynamic
	// Terraform object whose attribute set is determined by configuration.
	// Adding fields to it during apply causes:
	//
	//   .body: wrong final value type: incorrect object attributes
	//
	// Use a detached request copy for defaults and provider-side processing.
	requestBody, err := datasetMap(body)
	if err != nil {
		return nil, err
	}

	if err := setDatasetDefaults(requestBody, resourceVersion); err != nil {
		return nil, err
	}

	version, _, exists, err := datasetStringField(
		requestBody,
		"version",
	)
	if err != nil {
		return nil, err
	}

	if !exists || strings.TrimSpace(version) == "" {
		return nil, fmt.Errorf(
			`dataset body field "version" must match resource-level "name" %q`,
			resourceVersion,
		)
	}

	if strings.TrimSpace(version) != resourceVersion {
		return nil, fmt.Errorf(
			`dataset body field "version" must match resource-level "name": got %q, expected %q`,
			version,
			resourceVersion,
		)
	}

	sourceURL, expectedSHA256, verifySHA256, err := datasetSourceInfo(
		requestBody,
	)
	if err != nil {
		return nil, err
	}

	_, _, err = datasetVersionRequestBody(
		requestBody,
		resourceVersion,
		"",
	)
	if err != nil {
		return nil, err
	}

	pendingBody, err := datasetPendingUploadBody(requestBody)
	if err != nil {
		return nil, err
	}

	pendingResponse, err := client.DataPlaneClient.Action(
		ctx,
		datasetVersionAction(id, "startPendingUpload"),
		"",
		id.ApiVersion,
		http.MethodPost,
		pendingBody,
		options,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"starting pending dataset upload: %w",
			err,
		)
	}

	uploadSASURL, dataURI, err := datasetUploadDetails(pendingResponse)
	if err != nil {
		return nil, err
	}

	computedSHA256, err := streamDatasetToUpload(
		ctx,
		sourceURL,
		uploadSASURL,
		expectedSHA256,
		verifySHA256,
	)
	if err != nil {
		return nil, err
	}

	versionBody, _, err := datasetVersionRequestBody(
		requestBody,
		resourceVersion,
		dataURI,
	)
	if err != nil {
		return nil, err
	}

	responseBody, err := client.DataPlaneClient.ActionWithContentType(
		ctx,
		id.AzureResourceId,
		"",
		id.ApiVersion,
		http.MethodPatch,
		versionBody,
		options,
		"application/merge-patch+json",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"creating dataset version: %w",
			err,
		)
	}

	// Add the checksum only to the response/output object. Do not add it to
	// the Terraform body object because computed_sha256 was not part of the
	// configured body's object type.
	if responseValues, err := datasetMap(responseBody); err == nil {
		responseValues["computed_sha256"] = computedSHA256
		responseBody = responseValues
	}

	return responseBody, nil
}

func (c FoundryDatasetCustomization) CreateFunc() CreateFunc {
	return func(
		ctx context.Context,
		client clients.Client,
		id parse.DataPlaneResourceId,
		body interface{},
		options clients.RequestOptions,
	) error {
		_, err := c.createOrUpdate(
			ctx,
			client,
			id,
			body,
			options,
		)

		return err
	}
}

func (c FoundryDatasetCustomization) ReadFunc() ReadFunc {
	return func(
		ctx context.Context,
		client clients.Client,
		id parse.DataPlaneResourceId,
		options clients.RequestOptions,
	) (interface{}, error) {
		return client.DataPlaneClient.Get(ctx, id, options)
	}
}

func (c FoundryDatasetCustomization) UpdateFunc() UpdateFunc {
	return func(
		_ context.Context,
		_ clients.Client,
		id parse.DataPlaneResourceId,
		_ interface{},
		_ clients.RequestOptions,
	) error {
		return fmt.Errorf(
			"Foundry dataset version %q is immutable; create a new dataset version instead",
			id.Name,
		)
	}
}

func (c FoundryDatasetCustomization) DeleteFunc() DeleteFunc {
	return nil
}

func (c FoundryDatasetCustomization) StateBodyFunc() StateBodyFunc {
	return func(body interface{}) (interface{}, error) {
		return body, nil
	}
}

func (c FoundryDatasetCustomization) PlanBodyFunc() PlanBodyFunc {
	return func(
		planBody interface{},
		stateBody interface{},
	) (interface{}, error) {
		planValues, err := datasetMap(planBody)
		if err != nil {
			return nil, err
		}

		stateValues, err := datasetMap(stateBody)
		if err != nil {
			return planValues, nil
		}

		// Only copy fields that are defaults or API fields. Do not copy
		// computed_sha256 into the configured body object.
		for _, field := range []string{
			"version",
			"type",
			"format",
		} {
			if _, _, exists := datasetField(planValues, field); exists {
				continue
			}

			if value, _, exists := datasetField(stateValues, field); exists {
				planValues[field] = value
			}
		}

		return planValues, nil
	}
}

func (c FoundryDatasetCustomization) PreserveBodyStateOnRead() bool {
	return true
}

func (c FoundryDatasetCustomization) UseResponseBodyAsOutput() bool {
	return true
}

func (c FoundryDatasetCustomization) AugmentReadOutput(
	responseBody interface{},
	stateBody interface{},
) (interface{}, error) {
	stateValues, err := datasetMap(stateBody)
	if err != nil {
		return responseBody, nil
	}

	computedSHA256, _, exists, err := datasetStringField(
		stateValues,
		"computed_sha256",
	)
	if err != nil ||
		!exists ||
		strings.TrimSpace(computedSHA256) == "" {
		return responseBody, nil
	}

	outputValues, err := datasetMap(responseBody)
	if err != nil {
		return responseBody, nil
	}

	outputValues["computed_sha256"] = computedSHA256
	return outputValues, nil
}

var _ DataPlaneResource = &FoundryDatasetCustomization{}
var _ DataPlaneResourceWithPlanBody = &FoundryDatasetCustomization{}
var _ DataPlaneResourceWithStateBody = &FoundryDatasetCustomization{}
var _ DataPlaneResourceWithReadOptions = &FoundryDatasetCustomization{}
