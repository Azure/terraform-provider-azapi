package types

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Schema struct {
	Types []*TypeBase
}

type typeItem struct {
	Type string `json:"$type"`
}

func (s *Schema) UnmarshalJSON(body []byte) error {
	var m []*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	types := make([]*TypeBase, 0)
	for _, v := range m {
		var t typeItem
		err = json.Unmarshal(*v, &t)
		if err != nil {
			return err
		}
		dec := json.NewDecoder(bytes.NewReader(*v))
		dec.DisallowUnknownFields()
		switch t.Type {
		case "StringType":
			var t StringType
			if err = dec.Decode(&t); err != nil {
				return err
			}
			types = append(types, t.AsTypeBase())
		case "IntegerType":
			var t IntegerType
			if err = dec.Decode(&t); err != nil {
				return err
			}
			types = append(types, t.AsTypeBase())
		case "BooleanType":
			var t BooleanType
			if err = dec.Decode(&t); err != nil {
				return err
			}
			types = append(types, t.AsTypeBase())
		case "AnyType":
			var t AnyType
			if err = dec.Decode(&t); err != nil {
				return err
			}
			types = append(types, t.AsTypeBase())
		case "StringLiteralType":
			var t StringLiteralType
			if err = dec.Decode(&t); err != nil {
				return err
			}
			types = append(types, t.AsTypeBase())
		case "ObjectType":
			var t ObjectType
			if err = dec.Decode(&t); err != nil {
				return err
			}
			types = append(types, t.AsTypeBase())
		case "ArrayType":
			var t ArrayType
			if err = dec.Decode(&t); err != nil {
				return err
			}
			types = append(types, t.AsTypeBase())
		case "ResourceType":
			var t ResourceType
			// ResourceType has a custom UnmarshalJSON method
			err = json.Unmarshal(*v, &t)
			if err != nil {
				return err
			}
			types = append(types, t.AsTypeBase())
		case "UnionType":
			var t UnionType
			if err = dec.Decode(&t); err != nil {
				return err
			}
			types = append(types, t.AsTypeBase())
		case "DiscriminatedObjectType":
			var t DiscriminatedObjectType
			if err = dec.Decode(&t); err != nil {
				return err
			}
			types = append(types, t.AsTypeBase())
		case "ResourceFunctionType":
			var t ResourceFunctionType
			if err = dec.Decode(&t); err != nil {
				return err
			}
			types = append(types, t.AsTypeBase())
		default:
			return fmt.Errorf("unknown type %s", t.Type)
		}
	}
	for index, v := range types {
		if v != nil {
			value := *v
			switch t := value.(type) {
			case *ObjectType:
				t.AdditionalProperties.UpdateType(types)
				for index := range t.Properties {
					reference := t.Properties[index].Type
					reference.UpdateType(types)
				}
				types[index] = t.AsTypeBase()
			case *ArrayType:
				t.ItemType.UpdateType(types)
				types[index] = t.AsTypeBase()
			case *ResourceType:
				t.Body.UpdateType(types)
				types[index] = t.AsTypeBase()
			case *UnionType:
				for index := range t.Elements {
					t.Elements[index].UpdateType(types)
				}
				types[index] = t.AsTypeBase()
			case *DiscriminatedObjectType:
				for index := range t.Elements {
					reference := t.Elements[index]
					reference.UpdateType(types)
				}
				for index := range t.BaseProperties {
					reference := t.BaseProperties[index].Type
					reference.UpdateType(types)
				}
				types[index] = t.AsTypeBase()
			case *ResourceFunctionType:
				t.Input.UpdateType(types)
				t.Output.UpdateType(types)
				types[index] = t.AsTypeBase()
			}
		}
	}
	s.Types = types
	return nil
}
