package pufferpanel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

// MetadataType designed to be overridden
type MetadataType struct {
	Type     string                 `json:"type,omitempty"`
	Metadata map[string]interface{} `json:"-"`
} //@name Metadata

type ConditionalMetadataType struct {
	If string `json:"if,omitempty"`
	MetadataType
} //@name MetadataWithIf

// UnmarshalJSON parses a type with this declaration, storing what it needs into metadata and type
func (t *MetadataType) UnmarshalJSON(bs []byte) error {
	err := json.Unmarshal(bs, &t.Metadata)
	if err != nil {
		return err
	}

	a := t.Metadata["type"]
	var ok bool
	t.Type, ok = a.(string)
	if !ok && reflect.TypeOf(a) != reflect.TypeOf(nil) {
		return fmt.Errorf("type is of %s instead of string", reflect.TypeOf(a))
	}

	delete(t.Metadata, "type")
	return nil
}

func (t *MetadataType) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{ ")

	err := encode(&buf, "type", t.Type)
	if err != nil {
		return nil, err
	}

	for k, v := range t.Metadata {
		buf.WriteString(", ")
		err = encode(&buf, k, v)
		if err != nil {
			return nil, err
		}
	}

	buf.WriteString(" }")
	return buf.Bytes(), nil
}

// ParseMetadata Parses the metadata into the target interface
func (t *MetadataType) ParseMetadata(target interface{}) error {
	return UnmarshalTo(t, target)
}

// UnmarshalJSON parses a type with this declaration, storing what it needs into metadata and type
func (t *ConditionalMetadataType) UnmarshalJSON(bs []byte) error {
	err := json.Unmarshal(bs, &t.Metadata)
	if err != nil {
		return err
	}

	a := t.Metadata["type"]
	var ok bool
	t.Type, ok = a.(string)
	if !ok && reflect.TypeOf(a) != reflect.TypeOf(nil) {
		return fmt.Errorf("type is of %s instead of string", reflect.TypeOf(a))
	}

	a, exists := t.Metadata["if"]
	if exists {
		t.If, ok = a.(string)
		if !ok && reflect.TypeOf(a) != reflect.TypeOf(nil) {
			return fmt.Errorf("if is of %s instead of string", reflect.TypeOf(a))
		}
	}

	delete(t.Metadata, "type")
	delete(t.Metadata, "if")
	return nil
}

func (t *ConditionalMetadataType) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{ ")

	var err error
	if t.If != "" {
		err = encode(&buf, "if", t.If)
		if err != nil {
			return nil, err
		}
		buf.WriteString(", ")
	}

	err = encode(&buf, "type", t.Type)
	if err != nil {
		return nil, err
	}

	for k, v := range t.Metadata {
		buf.WriteString(", ")
		err = encode(&buf, k, v)
		if err != nil {
			return nil, err
		}
	}

	buf.WriteString(" }")
	return buf.Bytes(), nil
}

// ParseMetadata Parses the metadata into the target interface
func (t *ConditionalMetadataType) ParseMetadata(target interface{}) error {
	return UnmarshalTo(t, target)
}

func encode(buf *bytes.Buffer, k string, v any) error {
	d, err := json.Marshal(v)
	if err != nil {
		return err
	}
	buf.WriteString("\"" + k + "\": " + string(d) + "")
	return nil
}
