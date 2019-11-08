package shared

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"reflect"
	"time"
)

//Converts the val parameter to the same type as the target
func Convert(val interface{}, target interface{}) (interface{}, error) {
	switch target.(type) {
	case string:
		if val == nil {
			return "", nil
		}
		return cast.ToStringE(val)
	case int:
		if val == nil {
			return int(0), nil
		}
		return cast.ToIntE(val)
	case int8:
		if val == nil {
			return int8(0), nil
		}
		return cast.ToInt8E(val)
	case int16:
		if val == nil {
			return int16(0), nil
		}
		return cast.ToInt16E(val)
	case int32:
		if val == nil {
			return int32(0), nil
		}
		return cast.ToInt32E(val)
	case int64:
		if val == nil {
			return int64(0), nil
		}
		return cast.ToInt64E(val)
	case uint:
		if val == nil {
			return uint(0), nil
		}
		return cast.ToUintE(val)
	case uint8:
		if val == nil {
			return uint8(0), nil
		}
		return cast.ToUint8E(val)
	case uint16:
		if val == nil {
			return uint16(0), nil
		}
		return cast.ToUint16E(val)
	case uint32:
		if val == nil {
			return uint32(0), nil
		}
		return cast.ToUint32E(val)
	case uint64:
		if val == nil {
			return uint64(0), nil
		}
		return cast.ToUint64E(val)
	case bool:
		if val == nil {
			return false, nil
		}
		return cast.ToBoolE(val)
	case time.Duration:
		if val == nil {
			return time.Duration(0), nil
		}
		return cast.ToDurationE(val)
	case time.Time:
		if val == nil {
			return time.Time{}, nil
		}
		return cast.ToTimeE(val)
	case float32:
		if val == nil {
			return float32(0), nil
		}
		return cast.ToFloat64E(val)
	case float64:
		if val == nil {
			return float64(0), nil
		}
		return cast.ToFloat64E(val)
	case map[string]string:
		if val == nil {
			return map[string]string{}, nil
		}
		return cast.ToStringMapStringE(val)
	case map[string][]string:
		if val == nil {
			return map[string][]string{}, nil
		}
		return cast.ToStringMapStringSliceE(val)
	case map[string]bool:
		if val == nil {
			return map[string]bool{}, nil
		}
		return cast.ToStringMapBoolE(val)
	case map[string]interface{}:
		if val == nil {
			return map[string]interface{}{}, nil
		}
		return cast.ToStringMapE(val)
	case map[string]int:
		if val == nil {
			return map[string]int{}, nil
		}
		return cast.ToStringMapIntE(val)
	case map[string]int64:
		if val == nil {
			return map[string]int64{}, nil
		}
		return cast.ToStringMapInt64E(val)
	case []interface{}:
		if val == nil {
			return []interface{}{}, nil
		}
		return cast.ToSliceE(val)
	case []bool:
		if val == nil {
			return []bool{}, nil
		}
		return cast.ToBoolSliceE(val)
	case []string:
		if val == nil {
			return []string{}, nil
		}
		return cast.ToStringSliceE(val)
	case []int:
		if val == nil {
			return []int{}, nil
		}
		return cast.ToIntSliceE(val)
	case []time.Duration:
		if val == nil {
			return []time.Duration{}, nil
		}
		return cast.ToDurationSliceE(val)
	}

	return nil, errors.New(fmt.Sprintf("cannot convert %s to %s", reflect.TypeOf(val), reflect.TypeOf(target)))
}
