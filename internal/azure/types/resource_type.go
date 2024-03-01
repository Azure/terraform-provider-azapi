package types

import (
	"encoding/json"
	"fmt"
)

var _ TypeBase = &ResourceType{}

type ResourceType struct {
	Type               string
	Name               string
	ScopeTypes         []ScopeType
	ReadOnlyScopeTypes []ScopeType
	Body               *TypeReference
	Flags              []ResourceTypeFlag
}

func (t *ResourceType) GetWriteOnly(body interface{}) interface{} {
	if t == nil || body == nil {
		return nil
	}
	if t.Body != nil && t.Body.Type != nil {
		return (*t.Body.Type).GetWriteOnly(body)
	}
	return nil
}

func (t *ResourceType) Validate(body interface{}, path string) []error {
	if t == nil || body == nil {
		return []error{}
	}
	errors := make([]error, 0)
	if t.Body != nil && t.Body.Type != nil {
		errors = append(errors, (*t.Body.Type).Validate(body, path)...)
	}
	return errors
}

func (t *ResourceType) AsTypeBase() *TypeBase {
	typeBase := TypeBase(t)
	return &typeBase
}

func (t *ResourceType) IsReadOnly() bool {
	for _, value := range t.Flags {
		if value == ResourceTypeFlagReadOnly {
			return true
		}
	}
	return false
}

func (t *ResourceType) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "$type":
			if v != nil {
				var typeRef string
				err := json.Unmarshal(*v, &typeRef)
				if err != nil {
					return err
				}
				t.Type = typeRef
			}
		case "name":
			if v != nil {
				var name string
				err := json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				t.Name = name
			}
		case "scopeType":
			if v != nil {
				var scopeType int
				err := json.Unmarshal(*v, &scopeType)
				if err != nil {
					return err
				}
				scopeTypes := make([]ScopeType, 0)
				for _, f := range PossibleScopeTypeValues() {
					if scopeType&int(f) != 0 {
						scopeTypes = append(scopeTypes, f)
					}
				}
				if scopeType == 0 {
					scopeTypes = append(scopeTypes, Unknown)
				}
				t.ScopeTypes = scopeTypes
			}
		case "readOnlyScopes":
			if v != nil {
				var scopeType int
				err := json.Unmarshal(*v, &scopeType)
				if err != nil {
					return err
				}
				scopeTypes := make([]ScopeType, 0)
				for _, f := range PossibleScopeTypeValues() {
					if scopeType&int(f) != 0 {
						scopeTypes = append(scopeTypes, f)
					}
				}
				if scopeType == 0 {
					scopeTypes = append(scopeTypes, Unknown)
				}
				t.ReadOnlyScopeTypes = scopeTypes
			}
		case "body":
			if v != nil {
				var typeRef TypeReference
				err := json.Unmarshal(*v, &typeRef)
				if err != nil {
					return err
				}
				t.Body = &typeRef
			}
		case "flags":
			if v != nil {
				var flag int
				err := json.Unmarshal(*v, &flag)
				if err != nil {
					return err
				}
				flags := make([]ResourceTypeFlag, 0)
				for _, f := range PossibleResourceTypeFlagValues() {
					if flag&int(f) != 0 {
						flags = append(flags, f)
					}
				}
				t.Flags = flags
			}
		default:
			return fmt.Errorf("unmarshalling resource type, unrecognized key: %s", k)
		}
	}
	return nil
}

type ScopeType int

const (
	Unknown ScopeType = 0

	Tenant ScopeType = 1 << 0

	ManagementGroup ScopeType = 1 << 1

	Subscription ScopeType = 1 << 2

	ResourceGroup ScopeType = 1 << 3

	Extension ScopeType = 1 << 4
)

func (scope ScopeType) String() string {
	switch scope {
	case Unknown:
		return "Unknown"

	case Tenant:
		return "Tenant"

	case ManagementGroup:
		return "ManagementGroup"

	case Subscription:
		return "Subscription"

	case ResourceGroup:
		return "ResourceGroup"

	case Extension:
		return "Extension"
	}
	return ""
}

func PossibleScopeTypeValues() []ScopeType {
	return []ScopeType{Unknown, Tenant, ManagementGroup, Subscription, ResourceGroup, Extension}
}

type ResourceTypeFlag int

const (
	ResourceTypeFlagNone ResourceTypeFlag = 0

	ResourceTypeFlagReadOnly ResourceTypeFlag = 1 << 0
)

func PossibleResourceTypeFlagValues() []ResourceTypeFlag {
	return []ResourceTypeFlag{ResourceTypeFlagNone, ResourceTypeFlagReadOnly}
}
