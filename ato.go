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

	var ansibleData map[interface{}]interface{}
	err = yaml.Unmarshal(b, &ansibleData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	turtle(ansibleData, 0)
}

//But of course the world is flat and resting on the shell of a giant turtle.
func turtle(i map[interface{}]interface{}, offset int) {
	for k, v := range i {
		switch v.(type) {
		case string:
			fmt.Printf("%s%s:\n", pad(offset), k)
			fmt.Printf("%stype: string\n", pad(padLevel(offset)))
		case int:
			fmt.Printf("%s%s:\n", pad(offset), k)
			fmt.Printf("%stype: integer\n", pad(padLevel(offset)))
		case bool:
			fmt.Printf("%s%s:\n", pad(offset), k)
			fmt.Printf("%stype: boolean\n", pad(padLevel(offset)))
		case map[interface{}]interface{}:
			fmt.Printf("%s%s:\n", pad(offset), k)
			fmt.Printf("%stype: object\n", pad(padLevel(offset)))
			fmt.Printf("%sproperties:\n", pad(padLevel(offset)))
			turtle(v.(map[interface{}]interface{}), padLevel(offset+2))
		default:
			fmt.Printf("unknown type: %T", v)
		}
	}
}

func pad(i int) string {
	return strings.Repeat(" ", i)
}

func padLevel(i int) int {
	return i + 2
}

/*
name:
  type: string
*/
