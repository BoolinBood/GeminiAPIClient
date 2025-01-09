package function

import (
	"fmt"
	"geminiapiclient/function/lights"
	"reflect"
)

var Map = map[string]interface{}{
	"LivingRoomLight": lights.LivingRoomLight,
}

func CallFunctionByName(name string, args ...interface{}) ([]interface{}, error) {
	// Get the function from the map
	fn, exists := Map[name]
	if !exists {
		return nil, fmt.Errorf("function %s not found", name)
	}

	// Use reflection to call the function with arguments
	fnValue := reflect.ValueOf(fn)
	if len(args) != fnValue.Type().NumIn() {
		return nil, fmt.Errorf("function %s expects %d arguments, but got %d", name, fnValue.Type().NumIn(), len(args))
	}

	// Convert args to reflect.Value slice
	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		reflectArgs[i] = reflect.ValueOf(arg)
	}

	// Call the function
	results := fnValue.Call(reflectArgs)

	// Convert the reflect.Value results to a slice of interface{}
	output := make([]interface{}, len(results))
	for i, result := range results {
		output[i] = result.Interface()
	}

	return output, nil
}
