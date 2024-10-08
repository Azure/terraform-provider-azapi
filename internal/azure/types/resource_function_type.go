package types

var _ TypeBase = &ResourceFunctionType{}

type ResourceFunctionType struct {
	Type         string         `json:"$type"`
	Name         string         `json:"name"`
	ResourceType string         `json:"resourceType"`
	ApiVersion   string         `json:"apiVersion"`
	Input        *TypeReference `json:"input"`
	Output       *TypeReference `json:"output"`
}

func (t ResourceFunctionType) GetReadOnly(i interface{}) interface{} {
	return i
}

func (t ResourceFunctionType) AsTypeBase() *TypeBase {
	typeBase := TypeBase(t)
	return &typeBase
}

func (t ResourceFunctionType) Validate(body interface{}, path string) []error {
	return []error{}
}

func (t ResourceFunctionType) GetWriteOnly(body interface{}) interface{} {
	return body
}
