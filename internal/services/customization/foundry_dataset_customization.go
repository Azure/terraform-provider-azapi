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
// The request body uses Foundry API fields:
//
//	name
//	version
//	description
//	type
//	format
//	dataUri
//
// The following fields are provider-side fields:
//
//	source_url
//	source_sha256
//
// computed_sha256 is exposed in the resource output. It is intentionally not
// added to body because body is a Terraform dynamic object whose attribute
// type is determined by configuration.
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

	sourceURL, _, exists, err := datasetStringField(
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

	expectedSHA256, _, hasSHA256, err := datasetStringField(
		values,
		"source_sha256",
		"sourceSha256",
	)
	if err != nil {
		return "", "", false, err
	}

	// An empty checksum disables verification. The provider still calculates
	// and exposes the actual checksum in output.computed_sha256.
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

	version, _, exists, err := datasetStringField(values, "version")
	if err != nil {
		return err
	}

	version = strings.TrimSpace(version)

	if !exists || version == "" {
		version = strings.TrimSpace(versionFallback)

		if version == "" || version == "__generated__" {
			return fmt.Errorf(
				`resource-level "name" must contain the dataset version`,
			)
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

	name, err := datasetRequiredString(values, "name")
	if err != nil {
		return nil, "", err
	}

	version, _, exists, err := datasetStringField(values, "version")
	if err != nil {
		return nil, "", err
	}

	version = strings.TrimSpace(version)

	if !exists || version == "" {
		version = strings.TrimSpace(versionFallback)

		if version == "" || version == "__generated__" {
			return nil, "", fmt.Errorf(
				`resource-level "name" must contain the dataset version`,
			)
		}

		values["version"] = version
	}

	description, err := datasetRequiredString(values, "description")
	if err != nil {
		return nil, "", err
	}

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

	format, err := datasetRequiredString(values, "format")
	if err != nil {
		return nil, "", err
	}

	// source_url and source_sha256 are provider-only fields and are not sent
	// to the Foundry API.
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

func datasetSourceFilename(sourceURL string) (string, error) {
	parsedURL, err := url.Parse(sourceURL)
	// err is always non-nil for invalid URLs.
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

	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/") +
		"/" +
		filename
	parsedURL.RawPath = ""

	return parsedURL.String(), nil
}

func datasetSourceHTTPClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(
			*http.Request,
			[]*http.Request,
		) error {
			return fmt.Errorf(
				"refusing redirect for source_url; validate and use the final URL",
			)
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

func downloadDatasetSHA256(sourceURL string) (string, error) {
	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		sourceURL,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf(
			"creating dataset checksum request: %w",
			err,
		)
	}

	request.Header.Set("Accept-Encoding", "identity")

	response, err := datasetSourceHTTPClient().Do(request)
	if err != nil {
		return "", fmt.Errorf(
			"downloading dataset for checksum: %w",
			err,
		)
	}

	defer response.Body.Close()

	if response.StatusCode < http.StatusOK ||
		response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf(
			"downloading dataset for checksum returned HTTP %s",
			response.Status,
		)
	}

	hasher := sha256.New()

	if _, err := io.Copy(hasher, response.Body); err != nil {
		return "", fmt.Errorf(
			"calculating dataset checksum: %w",
			err,
		)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
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

	// Work on a detached copy. The original Terraform body must not be
	// mutated because its object type is determined by the configuration.
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

	if _, _, err := datasetVersionRequestBody(
		requestBody,
		resourceVersion,
		"",
	); err != nil {
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

	// computed_sha256 belongs in output, not body.
	responseValues, err := datasetMap(responseBody)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid dataset version response: %w",
			err,
		)
	}

	responseValues["computed_sha256"] = computedSHA256

	return responseValues, nil
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
		// Preserve the configured body, including source_sha256.
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

		// Do not copy computed_sha256 into body. It belongs in output.
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
	outputValues, err := datasetMap(responseBody)
	if err != nil {
		return responseBody, nil
	}

	// Preserve a checksum already returned by another provider operation.
	if _, _, exists := datasetField(
		outputValues,
		"computed_sha256",
	); exists {
		return outputValues, nil
	}

	sourceURL, _, _, err := datasetSourceInfo(stateBody)
	if err != nil {
		// Imported resources or older state may not contain source_url.
		// The missing checksum must not prevent refresh or deletion.
		outputValues["computed_sha256"] = nil
		return outputValues, nil
	}

	// Calculate the checksum for output. This is best effort during reads:
	// source_url can expire, be deleted, or become unavailable after the
	// Azure AI asset has already been created.
	computedSHA256, err := downloadDatasetSHA256(sourceURL)
	if err != nil {
		outputValues["computed_sha256"] = nil
		return outputValues, nil
	}

	outputValues["computed_sha256"] = computedSHA256

	return outputValues, nil
}

var _ DataPlaneResource = &FoundryDatasetCustomization{}
var _ DataPlaneResourceWithPlanBody = &FoundryDatasetCustomization{}
var _ DataPlaneResourceWithStateBody = &FoundryDatasetCustomization{}
var _ DataPlaneResourceWithReadOptions = &FoundryDatasetCustomization{}
