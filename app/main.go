package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
)

type Config []map[string][]interface{}

func main() {

	data, err := ioutil.ReadFile("../vinfra.yaml")
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	processRegion(config, "eu-central-1")

	//refreshTerraformOutputs("vpc")
	//TerraformTemplateProcessing("ec2", "main.tf", true)
	//TerraformTemplateProcessing("rds", "main.tf", true)

}

func getConfig(path string) map[string]interface{} {

	obj := make(map[string]interface{})

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, obj)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}

	return obj

}

func getComponents(obj map[string]interface{}) interface{} {

	return obj["components"]

}

func toInterfaceSlice(input interface{}) []interface{} {

	var out []interface{}
	rv := reflect.ValueOf(input)
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			out = append(out, rv.Index(i).Interface())
		}
	}

	return out
}

func toMapFromInterface(input interface{}) {

	v := reflect.ValueOf(input)

	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			strct := v.MapIndex(key)
			fmt.Println(key.Interface(), strct.Interface())
		}
	}

}

func getAllComponents(input []interface{}) {

	for _, component := range input {

		fmt.Println("Iterating over", component)

		switch v := component.(type) {
		case map[interface{}]interface{}:
			fmt.Println("The object type matched, it is map[if]if", " The value is ", v)
		case string:
			fmt.Println(v)
		}

	}

}

func processRegion(config Config, regionName string) {
	for _, region := range config {
		if services, ok := region[regionName]; ok {
			fmt.Println("Found Region:", regionName)
			printServices(services, "  ")
			return // Stop after finding and processing the region
		}
	}
	fmt.Println("Region not found:", regionName)
}

func printServices(services []interface{}, indent string) {
	for i, service := range services {
		prefix := "|-"
		if i == len(services)-1 {
			prefix = "`-"
		}
		switch v := service.(type) {
		case string:
			fmt.Println(indent + prefix + v)
		case map[string]interface{}:
			for key, val := range v {
				fmt.Println(indent + prefix + key + " |")
				if subServices, ok := val.([]interface{}); ok {
					printServices(subServices, indent+"    ")
				} else {
					fmt.Println(indent + "    " + fmt.Sprintf("%v", val))
				}
			}
		case map[interface{}]interface{}:
			for key, val := range v {
				keyStr, ok := key.(string)
				if !ok {
					fmt.Println(indent+"Key type not a string, found:", fmt.Sprintf("%T", key))
					continue
				}
				fmt.Println(indent + prefix + keyStr + " |")
				if subServices, ok := val.([]interface{}); ok {
					printServices(subServices, indent+"    ")
				} else {
					fmt.Println(indent + "    " + fmt.Sprintf("%v", val))
				}
			}
		default:
			fmt.Println(indent + fmt.Sprintf("Unhandled type: %T", v))
		}
	}
}
