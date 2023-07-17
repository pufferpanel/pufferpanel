package main

import (
	"github.com/pufferpanel/pufferpanel/v3/conditions"
	"syscall/js"
)

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("CalculateIf", js.FuncOf(func(this js.Value, args []js.Value) any {
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
			switch value.Type().String() {
			case js.TypeBoolean.String():
				inputData[fieldName] = value.Bool()
			}
		}

		result, err := conditions.CalculateIf(condition, inputData, nil)

		if err != nil {
			return err.Error()
		}
		return result
	}))
	<-done
}
