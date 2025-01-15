package function

import (
	"fmt"
	"geminiapiclient/function/esp32"
	"geminiapiclient/function/spotify"
	"log"
	"reflect"
)

var Map = map[string]interface{}{
	"LEDControl": esp32.LEDControl,
	"SearchSong": spotify.SearchSong,
	"PlayAlbum":  spotify.PlayAlbum,
}

func CallFunctionByName(name string, args ...interface{}) ([]interface{}, error) {
	// Get the function from the map
	fn, exists := Map[name]

	// Replace error handling with useful prompt for better experience
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
		log.Println("reflectArgs:", reflectArgs[i])
	}

	// Call the function
	results := fnValue.Call(reflectArgs)

	// Convert the reflect Value results to a slice of interface{}
	output := make([]interface{}, len(results))
	for i, result := range results {
		output[i] = result.Interface()
	}

	return output, nil
}
