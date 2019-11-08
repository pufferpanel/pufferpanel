package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

//designed to be overridden
type MetadataType struct {
	Type     string                 `json:"type,omitempty"`
	Metadata map[string]interface{} `json:"-,omitempty"`
}

//parses a type with this declaration, storing what it needs into metadata and type
func (t *MetadataType) UnmarshalJSON(bs []byte) (err error) {
	err = json.Unmarshal(bs, &t.Metadata)
	if err != nil {
		return
	}
	a := t.Metadata["type"]
	if a == nil {
		return errors.New("no type defined")
	}
	var ok bool
	t.Type, ok = a.(string)
	if !ok {
		return errors.New(fmt.Sprintf("type is of %s instead of string", reflect.TypeOf(a)))
	}

	delete(t.Metadata, "type")
	return
}

func (t *MetadataType) MarshalJSON() ([]byte, error) {
	newMapping := make(map[string]interface{})
	for k, v := range t.Metadata {
		newMapping[k] = v
	}
	newMapping["type"] = t.Type
	return json.Marshal(newMapping)
}

//Parses the metadata into the target interface
func (t *MetadataType) ParseMetadata(target interface{}) (err error) {
	data, err := json.Marshal(t)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &target)
	return
}
