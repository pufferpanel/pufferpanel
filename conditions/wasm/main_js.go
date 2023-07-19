package main

import (
	"github.com/pufferpanel/pufferpanel/v3/conditions"
	"syscall/js"
)

func main() {
	js.Global().Set("resolveIf", resolveIf)
	<-make(chan bool)
}

func resolveIf(_ js.Value, args []js.Value) any {
	if len(args) != 2 {
		return "invalid arguments"
	}

	condition := args[0].String()
	data := args[1]

	inputData := map[string]interface{}{}

	for i := 0; i < data.Length(); i++ {
		item := data.Index(i)
		fieldName := item.Get("name").String()
		value := item.Get("value")
		switch value.Type() {
		case js.TypeBoolean:
			inputData[fieldName] = value.Bool()
		case js.TypeString:
			inputData[fieldName] = value.String()
		case js.TypeNumber:
			inputData[fieldName] = value.Int()
		}
	}

	result, err := conditions.ResolveIf(condition, inputData, nil)

	if err != nil {
		return err.Error()
	}
	return result
}
