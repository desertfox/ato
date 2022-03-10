package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type T struct {
	Type       string        `yaml:"type"`
	Properties interface{}   `yaml:"properties,omitempty"`
	Items      []interface{} `yaml:"items,omitempty"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println(errors.New("must pass ansible and crd  file name"))
		os.Exit(1)
	}

	openapischema := make(map[interface{}]interface{})
	for topKey, topValue := range read(os.Args[1]) {
		openapischema[topKey] = turtle(topValue)
	}

	openapiObject := make(map[interface{}]interface{})
	openapiObject["spec"] = T{"object", openapischema, nil}

	schema := make(map[interface{}]interface{})
	schema["openAPIV3Schema"] = T{"object", openapiObject, nil}

	crd := read(os.Args[2])

	spec := crd["spec"].(map[interface{}]interface{})

	delete(spec, "version")
	delete(spec, "subresources")

	versions := spec["versions"].([]interface{})[0].(map[interface{}]interface{})
	versions["schema"] = schema

	crd["apiVersion"] = "apiextensions.k8s.io/v1"

	b, err := yaml.Marshal(crd)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}

//But of course the world is flat and resting on the shell of a giant turtle.
func turtle(node interface{}) interface{} {
	switch value := node.(type) {
	case string:
		return T{"string", nil, nil}
	case int:
		return T{"integer", nil, nil}
	case bool:
		return T{"boolean", nil, nil}
	case []interface{}:
		var arg []interface{}
		for _, arr := range node.([]interface{}) {
			arg = append(arg, turtle(arr))
		}
		return T{"array", nil, arg}
	case map[interface{}]interface{}:
		var collect map[interface{}]interface{} = make(map[interface{}]interface{})
		for k, v := range node.(map[interface{}]interface{}) {
			collect[k] = turtle(v)
		}
		return collect
	default:
		fmt.Printf("unknown type: %T", value)
	}
	return nil
}

func read(file string) map[interface{}]interface{} {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var data map[interface{}]interface{}
	err = yaml.Unmarshal(b, &data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return data
}
