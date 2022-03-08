package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(errors.New("must pass ansible yaml file name"))
		os.Exit(1)
	}

	ansibleYamlFile := os.Args[1]

	b, err := ioutil.ReadFile(ansibleYamlFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	prefixOpenAPI()

	var ansibleData map[interface{}]interface{}
	err = yaml.Unmarshal(b, &ansibleData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	turtle(ansibleData, 10)
}

func prefixOpenAPI() {
	fmt.Println(`schema:
  openAPIV3Schema:
    type: object
    properties:
      spec:
	type: object
	properties:`)
}

//But of course the world is flat and resting on the shell of a giant turtle.
func turtle(i map[interface{}]interface{}, offset int) {
	for k, v := range i {
		switch v.(type) {
		case string:
			String(k.(string), offset)
		case int:
			Int(k.(string), offset)
		case bool:
			Bool(k.(string), offset)
		case []interface{}:
			fmt.Printf("%s%s:\n", pad(offset), k)
			fmt.Printf("%stype: array\n", pad(pad2(offset)))
			fmt.Printf("%sitems:\n", pad(pad2(offset)))
			for _, arr := range v.([]interface{}) {
				array := make(map[interface{}]interface{})
				array[""] = arr
				turtle(array, pad2(offset))
			}
		case map[interface{}]interface{}:
			Object(k.(string), offset)
			turtle(v.(map[interface{}]interface{}), pad2(offset+2))
		default:
			fmt.Printf("unknown type: %T", v)
		}
	}
}

func String(name string, offset int) {
	if name != "" {
		fmt.Printf("%s%s:\n", pad(offset), name)
	}
	fmt.Printf("%stype: string\n", pad(pad2(offset)))
}

func Int(name string, offset int) {
	if name != "" {
		fmt.Printf("%s%s:\n", pad(offset), name)
	}
	fmt.Printf("%stype: integer\n", pad(pad2(offset)))
}

func Bool(name string, offset int) {
	if name != "" {
		fmt.Printf("%s%s:\n", pad(offset), name)
	}
	fmt.Printf("%stype: boolean\n", pad(pad2(offset)))
}

func Object(name string, offset int) {
	if name != "" {
		fmt.Printf("%s%s:\n", pad(offset), name)
	}
	fmt.Printf("%stype: object\n", pad(pad2(offset)))
	fmt.Printf("%sproperties:\n", pad(pad2(offset)))
}

func pad(i int) string {
	return strings.Repeat(" ", i)
}

func pad2(i int) int {
	return i + 2
}

/*
name:
  type: string
*/
