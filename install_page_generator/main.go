package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

var c config

// Configuration file structure
type config struct {
	Options []struct {
		Name    string `yaml:"name"`
		GroupID string `yaml:"groupId"`
		Help    string `yaml:"help"`
		Values  []struct {
			Name    string `yaml:"name"`
			TabName string `yaml:"tabName"`
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

type pageData struct {
	Groups []group
}

type group struct {
	Name        string
	GroupName   string
	Help        string
	Description string
	Buttons     []struct {
		Name    string
		TabName string
	}
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

	var pD pageData

	if _, err := os.Stat("generated"); os.IsNotExist(err) {
		err := os.Mkdir("generated", 0777)
		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat("templates/configurator_buttons.html"); !os.IsNotExist(err) {
		t, err := template.ParseFiles("templates/configurator_buttons.html")
		if err != nil {
			log.Print(err)
			return
		}

		pD.Groups = getGroups()

		if _, err := os.Stat("generated/configurator.html"); os.IsNotExist(err) {
			f, err := os.Create("generated/configurator.html")
			if err != nil {
				log.Println("Create file: ", err)
				return
			}
			err = t.Execute(f, pD)
			if err != nil {
				log.Print("Execute: ", err)
				return
			}
		}
	}
}

func getGroups() []group {
	var groups []group
	for i := range c.Options {
		var g group
		g.Name = c.Options[i].Name
		g.GroupName = c.Options[i].GroupID
		g.Help = c.Options[i].Help
		for y := range c.Options[i].Values {
			g.Buttons = append(g.Buttons, struct {
				Name    string
				TabName string
			}{Name: c.Options[i].Values[y].Name, TabName: c.Options[i].Values[y].TabName})
		}
		groups = append(groups, g)
	}
	return groups
}
