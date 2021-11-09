package types

type TypeReference struct {
	Type      *TypeBase
	TypeIndex int
}

func (t *TypeReference) UpdateType(types []*TypeBase) {
	if t != nil {
		t.Type = types[t.TypeIndex]
	}
}
