package types

var _ TypeBase = &AnyType{}

type AnyType struct {
	Type string `json:"$type"`
}

func (t *AnyType) AsTypeBase() *TypeBase {
	typeBase := TypeBase(t)
	return &typeBase
}

func (t *AnyType) Validate(body interface{}, path string) []error {
	return nil
}

func (t *AnyType) GetWriteOnly(i interface{}) interface{} {
	if t == nil || i == nil {
		return nil
	}
	return i
}
