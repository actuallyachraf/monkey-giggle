package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
)

// HashKey represents keys for hashmaps
type HashKey struct {
	Type  Type
	Value uint64
}

// Each type we implement Integers,Boolean will implement a hashkey method

// HashKey for boolean types is either 1 for true or 0 for false
func (b *Boolean) HashKey() HashKey {
	var val uint64

	if b.Value {
		val = 1
	} else {
		val = 0
	}

	return HashKey{Type: b.Type(), Value: val}
}

// HashKey for Integers uses the integer value as key
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// HashKey for strings is the fnv hash of the string literal
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// HashPair represents a key value entry in the hashmap
type HashPair struct {
	Key   Object
	Value Object
}

// HashMap represents the actual hashmap
type HashMap struct {
	Pairs map[HashKey]HashPair
}

// Type implements the object interface
func (hm *HashMap) Type() Type {
	return HASH
}

// Inspect implements the object interface
func (hm *HashMap) Inspect() string {

	var out bytes.Buffer

	pairs := []string{}

	for _, pair := range hm.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
