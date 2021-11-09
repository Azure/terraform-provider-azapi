package types

import (
	"encoding/json"
)

type Schema struct {
	Types []*TypeBase
}

func (s *Schema) UnmarshalJSON(body []byte) error {
	var m []map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	types := make([]*TypeBase, 0)
	for _, v := range m {
		for _, typeBaseKind := range PossibleTypeBaseKindValues() {
			if value := v[string(typeBaseKind)]; value != nil {
				switch typeBaseKind {
				case TypeBaseKindBuiltInType:
					var t BuiltInType
					err = json.Unmarshal(*value, &t)
					if err != nil {
						return err
					}
					types = append(types, t.AsTypeBase())
				case TypeBaseKindObjectType:
					var t ObjectType
					err = json.Unmarshal(*value, &t)
					if err != nil {
						return err
					}
					types = append(types, t.AsTypeBase())
				case TypeBaseKindArrayType:
					var t ArrayType
					err = json.Unmarshal(*value, &t)
					if err != nil {
						return err
					}
					types = append(types, t.AsTypeBase())
				case TypeBaseKindResourceType:
					var t ResourceType
					err = json.Unmarshal(*value, &t)
					if err != nil {
						return err
					}
					types = append(types, t.AsTypeBase())
				case TypeBaseKindUnionType:
					var t UnionType
					err = json.Unmarshal(*value, &t)
					if err != nil {
						return err
					}
					types = append(types, t.AsTypeBase())
				case TypeBaseKindStringLiteralType:
					var t StringLiteralType
					err = json.Unmarshal(*value, &t)
					if err != nil {
						return err
					}
					types = append(types, t.AsTypeBase())
				case TypeBaseKindDiscriminatedObjectType:
					var t DiscriminatedObjectType
					err = json.Unmarshal(*value, &t)
					if err != nil {
						return err
					}
					types = append(types, t.AsTypeBase())
				}
				break
			}
		}

	}
	for index, v := range types {
		if v != nil {
			value := *v
			switch value.(type) {
			case *ObjectType:
				t := value.(*ObjectType)
				t.AdditionalProperties.UpdateType(types)
				for index := range t.Properties {
					reference := t.Properties[index].Type
					reference.UpdateType(types)
				}
				types[index] = t.AsTypeBase()
			case *ArrayType:
				t := value.(*ArrayType)
				t.ItemType.UpdateType(types)
				types[index] = t.AsTypeBase()
			case *ResourceType:
				t := value.(*ResourceType)
				t.Body.UpdateType(types)
				types[index] = t.AsTypeBase()
			case *UnionType:
				t := value.(*UnionType)
				for index := range t.Elements {
					t.Elements[index].UpdateType(types)
				}
				types[index] = t.AsTypeBase()
			case *DiscriminatedObjectType:
				t := value.(*DiscriminatedObjectType)
				for index := range t.Elements {
					reference := t.Elements[index]
					reference.UpdateType(types)
				}
				for index := range t.BaseProperties {
					reference := t.BaseProperties[index].Type
					reference.UpdateType(types)
				}
				types[index] = t.AsTypeBase()
			}
		}
	}
	s.Types = types
	return nil
}
