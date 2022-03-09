package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Type struct {
	Type string `yaml:"type"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println(errors.New("must pass crd and ansible yaml file name"))
		os.Exit(1)
	}

	oldCRD := read(os.Args[1])

	ansiblePlaybookYaml := read(os.Args[2])
	openapischema := turtle(ansiblePlaybookYaml)

	schema := updateCRD(openapischema)

	newCRD := migrateCRD(oldCRD, schema)

	fmt.Printf("%v", newCRD)
}

func migrateCRD(crd map[interface{}]interface{}, openAPIschema map[interface{}]interface{}) map[interface{}]interface{} {
	var spec map[interface{}]interface{} = crd["spec"].(map[interface{}]interface{})

	delete(spec, "version")
	delete(spec, "subresources")

	versions := spec["versions"].([]interface{})[0].(map[interface{}]interface{})
	versions["schema"] = openAPIschema

	crd["apiVersion"] = "apiextensions.k8s.io/v1"

	return crd
}

type O struct {
	Type       string      `yaml:"type"`
	Properties interface{} `yaml:"properties"`
}

func updateCRD(oas map[interface{}]interface{}) map[interface{}]interface{} {
	var spec map[interface{}]interface{} = make(map[interface{}]interface{})
	spec["type"] = "object"
	spec["properties"] = oas

	var openAPIschema map[interface{}]interface{} = make(map[interface{}]interface{})
	openAPIschema["type"] = "object"
	openAPIschema["properties"] = spec

	var schema map[interface{}]interface{} = make(map[interface{}]interface{})
	schema["openAPIV3Schema"] = openAPIschema

	return schema
}

//But of course the world is flat and resting on the shell of a giant turtle.
func turtle(i map[interface{}]interface{}) map[interface{}]interface{} {
	var collect map[interface{}]interface{} = make(map[interface{}]interface{})

	for k, v := range i {

		switch v.(type) {
		case string:
			collect[k.(string)] = Type{"string"}
		case int:
			collect[k.(string)] = Type{"integer"}
		case bool:
			collect[k.(string)] = Type{"boolean"}
		case []interface{}:
			/*
				fmt.Printf("%s%s:\n", pad(offset), k)
				fmt.Printf("%stype: array\n", pad(pad2(offset)))
				fmt.Printf("%sitems:\n", pad(pad2(offset)))
				for _, arr := range v.([]interface{}) {
					array := make(map[interface{}]interface{})
					array[""] = arr
					turtle(array, pad2(offset))
				}
			*/
		case map[interface{}]interface{}:
			collect[k.(string)] = Object(turtle(v.(map[interface{}]interface{})))
		default:
			fmt.Printf("unknown type: %T", v)
		}
	}

	return collect
}

func String(name string, d map[interface{}]interface{}) map[interface{}]interface{} {
	d[name] = Type{"string"}
	return d
}

func Int(name string, d map[interface{}]interface{}) map[interface{}]interface{} {
	d[name] = Type{"integer"}
	return d
}

func Bool(name string, d map[interface{}]interface{}) map[interface{}]interface{} {
	d[name] = Type{"boolean"}
	return d
}

func Object(nested map[interface{}]interface{}) map[interface{}]interface{} {
	var obj map[interface{}]interface{} = make(map[interface{}]interface{})
	obj["type"] = "object"
	obj["properties"] = nested
	return obj
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
