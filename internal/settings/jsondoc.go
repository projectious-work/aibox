package settings

import (
	"bytes"
	"encoding/json"
	"errors"
	"unicode/utf8"
)

// jsonKind is the kind of a parsed JSON value.
type jsonKind int

const (
	jsonObject jsonKind = iota
	jsonArray
	jsonString
	jsonNumber
	jsonBool
	jsonNull
)

// jsonNode is a JSON value that keeps object members in document order and
// keeps repeated member names, which encoding/json would merge silently.
type jsonNode struct {
	kind    jsonKind
	members []jsonMember
	items   []*jsonNode
	text    string
	number  json.Number
	boolean bool
}

// jsonMember is one name/value pair of an object, in document order.
type jsonMember struct {
	name  string
	value *jsonNode
}

var errNotJSON = errors.New("not JSON")

// parseJSON parses one JSON document. Syntax errors, invalid UTF-8 and
// trailing data are errNotJSON.
func parseJSON(data []byte) (*jsonNode, error) {
	if !utf8.Valid(data) || !json.Valid(data) {
		return nil, errNotJSON
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	return decodeNode(decoder)
}

func decodeNode(decoder *json.Decoder) (*jsonNode, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, errNotJSON
	}
	switch value := token.(type) {
	case json.Delim:
		if value == '{' {
			return decodeObject(decoder)
		}
		if value == '[' {
			return decodeArray(decoder)
		}
		return nil, errNotJSON
	case string:
		return &jsonNode{kind: jsonString, text: value}, nil
	case json.Number:
		return &jsonNode{kind: jsonNumber, number: value}, nil
	case bool:
		return &jsonNode{kind: jsonBool, boolean: value}, nil
	case nil:
		return &jsonNode{kind: jsonNull}, nil
	}
	return nil, errNotJSON
}

func decodeObject(decoder *json.Decoder) (*jsonNode, error) {
	node := &jsonNode{kind: jsonObject}
	for decoder.More() {
		token, err := decoder.Token()
		name, ok := token.(string)
		if err != nil || !ok {
			return nil, errNotJSON
		}
		value, err := decodeNode(decoder)
		if err != nil {
			return nil, err
		}
		node.members = append(node.members, jsonMember{name: name, value: value})
	}
	if _, err := decoder.Token(); err != nil {
		return nil, errNotJSON
	}
	return node, nil
}

func decodeArray(decoder *json.Decoder) (*jsonNode, error) {
	node := &jsonNode{kind: jsonArray}
	for decoder.More() {
		value, err := decodeNode(decoder)
		if err != nil {
			return nil, err
		}
		node.items = append(node.items, value)
	}
	if _, err := decoder.Token(); err != nil {
		return nil, errNotJSON
	}
	return node, nil
}

// member returns the last member with the given name, as JSON parsers do.
func (n *jsonNode) member(name string) *jsonNode {
	var found *jsonNode
	for _, member := range n.members {
		if member.name == name {
			found = member.value
		}
	}
	return found
}
