package pufferpanel

import "github.com/spf13/cast"

func GetStringOrDefault(data map[string]interface{}, key string, def string) string {
	if data == nil {
		return def
	}
	var section = data[key]
	if section == nil {
		return def
	} else {
		val, err := cast.ToStringE(section)
		if err != nil {
			return def
		}
		return val
	}
}

func GetBooleanOrDefault(data map[string]interface{}, key string, def bool) bool {
	if data == nil {
		return def
	}
	var section = data[key]
	if section == nil {
		return def
	} else {
		val, err := cast.ToBoolE(section)
		if err != nil {
			return def
		}
		return val
	}
}

func GetMapOrNull(data map[string]interface{}, key string) map[string]interface{} {
	if data == nil {
		return nil
	}
	var section = data[key]
	if section == nil {
		return nil
	} else {
		val, err := cast.ToStringMapE(section)
		if err != nil {
			return nil
		}
		return val
	}
}

func GetObjectArrayOrNull(data map[string]interface{}, key string) []interface{} {
	if data == nil {
		return nil
	}
	var section = data[key]
	if section == nil {
		return nil
	} else {
		val, err := cast.ToSliceE(section)
		if err != nil {
			return nil
		}
		return val
	}
}

func GetStringArrayOrNull(data map[string]interface{}, key string) []string {
	if data == nil {
		return nil
	}
	var section = data[key]
	if section == nil {
		return nil
	} else {
		val, err := cast.ToStringSliceE(section)
		if err != nil {
			return nil
		}
		return val
	}
}
