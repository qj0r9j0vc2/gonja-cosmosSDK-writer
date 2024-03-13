package main

import (
	"fmt"
	"github.com/noirbizarre/gonja"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
)

var (
	configTomlJ2 = gonja.Must(gonja.FromFile("templates/config.toml.j2"))
	appTomlJ2    = gonja.Must(gonja.FromFile("templates/app.toml.j2"))
)

type CustomConfig struct {
	// ConfigToml specifies parameter(s) in config.toml
	ConfigToml []Entry `yaml:"configToml"`
	// ConfigToml specifies parameter(s) in config.toml
	AppToml []Entry `yaml:"appToml"`
}

type Entry interface{}

func ConvertToMapStringInterface(entryList []Entry) map[string]interface{} {
	mapInterface := make(map[string]interface{})
	for _, entry := range entryList {
		mapInterface = mergeMaps(mapInterface, interfaceConverter(entry))
	}
	return mapInterface
}

func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for key, value := range m {
			result[key] = value
		}
	}
	return result
}

func interfaceConverter(entry interface{}) map[string]interface{} {
	mapInterface := make(map[string]interface{})
	var convertValue interface{}

	convert, ok := entry.(map[string]interface{})
	if ok {
		for key, value := range convert {
			if convertValue, ok = value.(string); ok {
				mapInterface[key] = convertValue
			} else if convertValue, ok = value.(int); ok {
				mapInterface[key] = convertValue
			} else if convertValue, ok = value.(bool); ok {
				mapInterface[key] = convertValue
			} else {
				mapInterface[key] = interfaceConverter(value)
			}
		}
	} else if reflect.TypeOf(entry).Kind() == reflect.Slice {
		s := reflect.ValueOf(entry)
		for i := 0; i < s.Len(); i++ {
			mapInterface = mergeMaps(mapInterface, interfaceConverter(s.Index(i).Interface()))
		}
	} else {
		mapInterface = mergeMaps(mapInterface, interfaceConverter(entry))
	}

	return mapInterface
}

func parse() {
	var config CustomConfig
	f, err := os.ReadFile("config.yaml")

	err = yaml.Unmarshal(f, &config)
	if err != nil {
		panic(err)
	}

	if len(config.ConfigToml) <= 0 {
		return
	}
	fmt.Printf("======== custom config.toml ========\n")
	for key, value := range ConvertToMapStringInterface(config.ConfigToml) {
		fmt.Printf("%s, %s\n", key, value)
	}

	configTomlOut, err := configTomlJ2.Execute(ConvertToMapStringInterface(config.ConfigToml))
	if err != nil {
		panic(err)
	}

	configTomlOutFile, err := os.Create("out/config.toml")
	if err != nil {
		panic(err)
	}
	_, err = configTomlOutFile.Write([]byte(configTomlOut))
	if err != nil {
		panic(err)
	}

	if len(config.AppToml) <= 0 {
		return
	}
	fmt.Printf("======== custom app.toml ========\n")
	for key, value := range ConvertToMapStringInterface(config.AppToml) {
		fmt.Printf("%s, %s\n", key, value)
	}

	appTomlOut, err := appTomlJ2.Execute(ConvertToMapStringInterface(config.AppToml))
	if err != nil {
		panic(err)
	}

	appTomlOutFile, err := os.Create("out/app.toml")
	if err != nil {
		panic(err)
	}
	_, err = appTomlOutFile.Write([]byte(appTomlOut))
	if err != nil {
		panic(err)
	}

}

func main() {
	parse()
}
