package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var c config

// Configuration file structure
type config struct {
	Options []struct {
		Name   string `yaml:"name"`
		Row    int    `yaml:"row"`
		Values []struct {
			Name string `yaml:"name"`
		} `yaml:"values"`
	} `yaml:"options"`
	Combinations []struct {
		Tabs []struct {
			Name         string `yaml:"name"`
			TemplateName string `yaml:"templateName"`
			Params       []struct {
				Type string `yaml:"type"`
			} `yaml:"params,omitempty"`
		} `yaml:"tabs"`
		Options []struct {
			Name  string `yaml:"name"`
			Value string `yaml:"value"`
		} `yaml:"options"`
	} `yaml:"combinations"`
}

type jsonElem struct {
	Option string
	Values []jsonElem
}

// Main function
func main() {
	c.getConf()

	var elems jsonElem

	elems.Option = c.Options[0].Name
	elems.Values = getElems(0)

	file, err := json.Marshal(elems)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	err = ioutil.WriteFile("js_conf.json", file, 0777)

	if err != nil {
		panic(err)
	}

	genSelectorTemplate()
}

// Recursive tree traversal function
func getElems(index int) []jsonElem {
	var elems []jsonElem

	for i := range c.Options[index].Values {
		var elem jsonElem
		elem.Option = c.Options[index].Values[i].Name
		if index < len(c.Options[index].Values)-1 {
			elem.Values = getElems(index + 1)
		}
		elems = append(elems, elem)
	}

	return elems
}

// Configuration file upload function
func (c *config) getConf() *config {

	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func genSelectorTemplate() {

	err := os.Mkdir("generated", 0777)
	if err != nil {
		panic(err)
	}
}
