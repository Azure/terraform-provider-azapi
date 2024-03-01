package types

import (
	"log"
	"strconv"
	"strings"
)

type TypeReference struct {
	Type *TypeBase
	Ref  string `json:"$ref"`
}

func (t *TypeReference) UpdateType(types []*TypeBase) {
	if t == nil {
		return
	}
	if t.Type != nil {
		return
	}
	if t.Ref == "" {
		log.Printf("[WARN] invalid, the ref is empty")
		return
	}
	if !strings.HasPrefix(t.Ref, "#/") {
		log.Printf("[WARN] invalid, the ref is invalid: %s", t.Ref)
		return
	}
	index, err := strconv.ParseInt(t.Ref[2:], 10, 64)
	if err != nil {
		log.Printf("[WARN] invalid, the ref is invalid: %s: %s", t.Ref, err)
		return
	}
	t.Type = types[int(index)]
}
