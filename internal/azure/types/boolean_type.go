package types

var _ TypeBase = &BooleanType{}

type BooleanType struct {
	Type string `json:"$type"`
}

func (t *BooleanType) AsTypeBase() *TypeBase {
	typeBase := TypeBase(t)
	return &typeBase
}

func (t *BooleanType) Validate(body interface{}, path string) []error {
	return nil
}

func (t *BooleanType) GetWriteOnly(i interface{}) interface{} {
	if t == nil || i == nil {
		return nil
	}
	return i
}
