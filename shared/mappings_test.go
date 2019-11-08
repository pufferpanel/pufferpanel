package shared

import (
	"reflect"
	"testing"
)

var mappingTest = map[string]interface{}{
	mappingStringKey:      mappingStringVal,
	mappingTrueBoolKey:    true,
	mappingFalseBoolKey:   false,
	mappingNullKey:        nil,
	mappingStringArrayKey: mappingStringArrayVal,
	mappingIntArrayKey:    mappingIntArrayVal,
	mappingMixedArrayKey:  mappingMixedArrayVal,
	mappingMapKey:         mappingMapVal,
	mappingObjectArrayKey: mappingObjectArrayVal,
}

var mappingNoKey = "noKey"
var mappingStringKey = "string"
var mappingTrueBoolKey = "trueBool"
var mappingFalseBoolKey = "falseBool"
var mappingNullKey = "null"
var mappingStringArrayKey = "stringArray"
var mappingIntArrayKey = "intArray"
var mappingMixedArrayKey = "mixedArray"
var mappingMapKey = "mapKey"
var mappingObjectArrayKey = "objectKey"

var trueString = "true"
var mappingDefaultString = "someDefault"
var mappingStringVal = "someValue"

var mappingStringArrayVal = []string{
	"val1",
	"2something",
}

var mappingMixedArrayVal = []interface{}{
	"test",
	12345,
	false,
	true,
}
var mappingMixedArrayValExpected = []string{
	"test",
	"12345",
	"false",
	"true",
}

var mappingIntArrayVal = []int{
	64,
	128,
}

var mappingMapVal = map[string]interface{}{
	"key1": "val1",
	"key2": "val2",
	"key3": false,
	"key4": 12345,
}

var mappingObjectArrayVal = []interface{}{
	"val1",
	true,
	false,
	12346,
}

func TestGetBooleanOrDefault(t *testing.T) {
	type args struct {
		data map[string]interface{}
		key  string
		def  bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test for true",
			args: args{
				data: mappingTest,
				key:  mappingTrueBoolKey,
			},
			want: true,
		},
		{
			name: "Test for false",
			args: args{
				data: mappingTest,
				key:  mappingFalseBoolKey,
			},
			want: false,
		},
		{
			name: "Test for true default",
			args: args{
				data: mappingTest,
				key:  mappingNoKey,
				def:  true,
			},
			want: true,
		},
		{
			name: "Test for invalid type",
			args: args{
				data: mappingTest,
				key:  mappingStringKey,
				def:  false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBooleanOrDefault(tt.args.data, tt.args.key, tt.args.def); got != tt.want {
				t.Errorf("GetBooleanOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMapOrNull(t *testing.T) {
	type args struct {
		data map[string]interface{}
		key  string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "Test for valid string array",
			args: args{
				data: mappingTest,
				key:  mappingMapKey,
			},
			want: mappingMapVal,
		},
		{
			name: "Test for missing key",
			args: args{
				data: mappingTest,
				key:  mappingNoKey,
			},
			want: nil,
		},
		{
			name: "Test for invalid type",
			args: args{
				data: mappingTest,
				key:  mappingStringArrayKey,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMapOrNull(tt.args.data, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMapOrNull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetObjectArrayOrNull(t *testing.T) {
	type args struct {
		data map[string]interface{}
		key  string
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "Test for valid object array",
			args: args{
				data: mappingTest,
				key:  mappingObjectArrayKey,
			},
			want: mappingObjectArrayVal,
		},
		{
			name: "Test for missing key",
			args: args{
				data: mappingTest,
				key:  mappingNoKey,
			},
			want: nil,
		},
		{
			name: "Test for invalid type",
			args: args{
				data: mappingTest,
				key:  mappingStringArrayKey,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetObjectArrayOrNull(tt.args.data, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetObjectArrayOrNull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStringArrayOrNull(t *testing.T) {
	type args struct {
		data map[string]interface{}
		key  string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test for valid string array",
			args: args{
				data: mappingTest,
				key:  mappingStringArrayKey,
			},
			want: mappingStringArrayVal,
		},
		{
			name: "Test for missing key",
			args: args{
				data: mappingTest,
				key:  mappingNoKey,
			},
			want: nil,
		},
		{
			name: "Test for invalid type",
			args: args{
				data: mappingTest,
				key:  mappingIntArrayKey,
			},
			want: nil,
		},
		{
			name: "Test for invalid type mix",
			args: args{
				data: mappingTest,
				key:  mappingMixedArrayKey,
			},
			want: mappingMixedArrayValExpected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStringArrayOrNull(tt.args.data, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStringArrayOrNull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStringOrDefault(t *testing.T) {
	type args struct {
		data map[string]interface{}
		key  string
		def  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test for valid string",
			args: args{
				data: mappingTest,
				key:  mappingStringKey,
			},
			want: mappingStringVal,
		},
		{
			name: "Test for string default",
			args: args{
				data: mappingTest,
				key:  mappingNoKey,
				def:  mappingDefaultString,
			},
			want: mappingDefaultString,
		},
		{
			name: "Test for invalid type",
			args: args{
				data: mappingTest,
				key:  mappingTrueBoolKey,
				def:  mappingDefaultString,
			},
			want: trueString,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStringOrDefault(tt.args.data, tt.args.key, tt.args.def); got != tt.want {
				t.Errorf("GetStringOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
