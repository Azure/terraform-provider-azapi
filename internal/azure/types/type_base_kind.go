package types

type TypeBaseKind string

const (
	TypeBaseKindBuiltInType             TypeBaseKind = "1"
	TypeBaseKindObjectType              TypeBaseKind = "2"
	TypeBaseKindArrayType               TypeBaseKind = "3"
	TypeBaseKindResourceType            TypeBaseKind = "4"
	TypeBaseKindUnionType               TypeBaseKind = "5"
	TypeBaseKindStringLiteralType       TypeBaseKind = "6"
	TypeBaseKindDiscriminatedObjectType TypeBaseKind = "7"
)

func PossibleTypeBaseKindValues() []TypeBaseKind {
	return []TypeBaseKind{
		TypeBaseKindBuiltInType,
		TypeBaseKindObjectType,
		TypeBaseKindArrayType,
		TypeBaseKindResourceType,
		TypeBaseKindUnionType,
		TypeBaseKindStringLiteralType,
		TypeBaseKindDiscriminatedObjectType}
}
