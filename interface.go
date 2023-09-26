package pufferpanel

import (
	"encoding/json"
	"fmt"
)

func ToString(v interface{}) string {
	if t, ok := v.(string); ok {
		return t
	}
	if stringer, ok := v.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%v", v)
}

func UnmarshalTo(source, target interface{}) error {
	data, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
