package schema_test

import (
	"errors"
	"testing"

	"github.com/Azure/terraform-provider-azapi/pkg/schema"
)

func TestTopLevelPropertyStatus(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		resourceType string
		property     string
		want         schema.PropertyStatus
	}{
		"writable tags": {
			resourceType: "Microsoft.Resources/resourceGroups@2021-04-01",
			property:     "tags",
			want:         schema.PropertyStatusWritable,
		},
		"unsupported tags": {
			resourceType: "Microsoft.Authorization/roleAssignments@2022-04-01",
			property:     "tags",
			want:         schema.PropertyStatusUnsupported,
		},
		"read-only system data": {
			resourceType: "Microsoft.Authorization/roleAssignments@2022-04-01",
			property:     "systemData",
			want:         schema.PropertyStatusReadOnly,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := schema.TopLevelPropertyStatus(test.resourceType, test.property)
			if err != nil {
				t.Fatalf("TopLevelPropertyStatus() error = %v", err)
			}
			if got != test.want {
				t.Errorf("TopLevelPropertyStatus() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestTopLevelPropertyStatusUnavailableDefinition(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"unknown resource type": "Microsoft.Contoso/notAResource@2021-01-01",
		"unknown API version":   "Microsoft.Resources/resourceGroups@1900-01-01",
	}

	for name, resourceType := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := schema.TopLevelPropertyStatus(resourceType, "tags")
			if got != schema.PropertyStatusUnknown {
				t.Errorf("TopLevelPropertyStatus() = %q, want %q", got, schema.PropertyStatusUnknown)
			}
			if !errors.Is(err, schema.ErrDefinitionUnavailable) {
				t.Errorf("TopLevelPropertyStatus() error = %v, want ErrDefinitionUnavailable", err)
			}
		})
	}
}

func TestTopLevelPropertyStatusInvalidInput(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		resourceType string
		property     string
		wantErr      error
	}{
		"empty resource type": {
			wantErr: schema.ErrInvalidResourceType,
		},
		"missing API version separator": {
			resourceType: "Microsoft.Resources/resourceGroups",
			wantErr:      schema.ErrInvalidResourceType,
		},
		"empty resource name": {
			resourceType: "@2021-04-01",
			wantErr:      schema.ErrInvalidResourceType,
		},
		"empty API version": {
			resourceType: "Microsoft.Resources/resourceGroups@",
			wantErr:      schema.ErrInvalidAPIVersion,
		},
		"multiple API version separators": {
			resourceType: "Microsoft.Resources/resourceGroups@2021-04-01@preview",
			wantErr:      schema.ErrInvalidResourceType,
		},
		"empty property": {
			resourceType: "Microsoft.Resources/resourceGroups@2021-04-01",
			wantErr:      schema.ErrInvalidProperty,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := schema.TopLevelPropertyStatus(test.resourceType, test.property)
			if got != schema.PropertyStatusUnknown {
				t.Errorf("TopLevelPropertyStatus() = %q, want %q", got, schema.PropertyStatusUnknown)
			}
			if !errors.Is(err, test.wantErr) {
				t.Errorf("TopLevelPropertyStatus() error = %v, want %v", err, test.wantErr)
			}
		})
	}
}
