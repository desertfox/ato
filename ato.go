package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

var crdPath string

type T struct {
	Type       string      `yaml:"type"`
	Properties interface{} `yaml:"properties,omitempty"`
	Items      interface{} `yaml:"items,omitempty"`
}

func init() {
	flag.StringVar(&crdPath, "crd", "", "--crd=path/to/crd.yaml")
}

func main() {
	flag.Parse()

	if crdPath == "" {
		log.Fatal(errors.New("must pass crd file path"))
	}

	log.Printf("Loading CustomResourceDefinition at %s\n", crdPath)
	crd := read(crdPath)

	spec := crd["spec"].(map[interface{}]interface{})

	ansiblePlaybookPath := fmt.Sprintf("./roles/%s/defaults/main.yml", spec["names"].(map[interface{}]interface{})["singular"])

	log.Printf("Loading Ansible Playbook defaults at %s\n", ansiblePlaybookPath)
	ansiblePlaybook := read(ansiblePlaybookPath)

	log.Println("Translate variables into OpenAPIv3Schema")
	openapischema := trutles(ansiblePlaybook)
	openapiObject := make(map[interface{}]interface{})
	openapiObject["spec"] = T{"object", openapischema, nil}
	schema := make(map[interface{}]interface{})
	schema["openAPIV3Schema"] = T{"object", openapiObject, nil}

	log.Println("Update CustomResourceDefinition")
	delete(spec, "version")
	delete(spec, "subresources")

	versions := spec["versions"].([]interface{})[0].(map[interface{}]interface{})
	versions["schema"] = schema

	crd["apiVersion"] = "apiextensions.k8s.io/v1"

	log.Printf("Write CustomResourceDefinition at %s\n", crdPath)
	b, err := yaml.Marshal(crd)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(crdPath, b, os.FileMode(int(0777)))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Finished")
}

func trutles(node map[interface{}]interface{}) map[interface{}]interface{} {
	t := make(map[interface{}]interface{})
	for k, v := range node {
		t[k] = turtle(v)
	}
	return t
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
		for _, arr := range node.([]interface{}) {
			return T{"array", nil, turtle(arr)}
		}
	case map[interface{}]interface{}:
		return T{"object", trutles(node.(map[interface{}]interface{})), nil}
	default:
		fmt.Printf("unknown type: %T", value)
	}
	return nil
}

func read(file string) map[interface{}]interface{} {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	var data map[interface{}]interface{}
	err = yaml.Unmarshal(b, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
