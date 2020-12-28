package parser

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func ParseFile(data interface{}, path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(file, data)
	if err != nil {
		log.Fatal(err)
	}
}
