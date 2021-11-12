package types

import "encoding/json"

var _ TypeBase = &ResourceType{}

type ResourceType struct {
	Name       string
	ScopeTypes []ScopeType
	Body       *TypeReference
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

func (t *ResourceType) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "Name":
			if v != nil {
				var name string
				err := json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				t.Name = name
			}
		case "ScopeType":
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
				t.ScopeTypes = scopeTypes
			}
		case "Body":
			if v != nil {
				var index int
				err := json.Unmarshal(*v, &index)
				if err != nil {
					return err
				}
				t.Body = &TypeReference{TypeIndex: index}
			}
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

func PossibleScopeTypeValues() []ScopeType {
	return []ScopeType{Unknown, Tenant, ManagementGroup, Subscription, ResourceGroup, Extension}
}
