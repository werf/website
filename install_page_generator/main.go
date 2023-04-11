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
			Name         string              `yaml:"name"`
			LinkName     string              `yaml:"linkName"`
			TemplateName string              `yaml:"templateName"`
			Params       []map[string]string `yaml:"params,omitempty"`
		} `yaml:"tabs"`
		Options []struct {
			Name  string `yaml:"name"`
			Value string `yaml:"value"`
		} `yaml:"options"`
	} `yaml:"combinations"`
}

type jsonElem struct {
	Option string
	Values []map[string]interface{}
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

type tab struct {
	Name      string
	Permalink string
	data      []map[string]string
}

// Main function
func main() {
	c.getConf()

	var elems jsonElem

	elems.Option = c.Options[0].GroupID
	elems.Values = getElems(0)

	file, err := json.Marshal(elems)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	if _, err = os.Stat("generated"); os.IsNotExist(err) {
		err = os.Mkdir("generated", 0777)
		if err != nil {
			panic(err)
		}
	}

	err = ioutil.WriteFile("generated/js_conf.json", file, 0777)

	if err != nil {
		panic(err)
	}

	genSelectorTemplate()
	genPartialsForTabs()
}

// Recursive tree traversal function
func getElems(index int) []map[string]interface{} {
	var elems []map[string]interface{}

	if index < len(c.Options)-1 {
		for i := range c.Options[index].Values {
			var elem map[string]interface{}
			elem = make(map[string]interface{})
			if index < len(c.Options[index].Values)-1 {
				elem[c.Options[index].Values[i].TabName] = getValues(index)
			}
			elems = append(elems, elem)
		}
	}

	return elems
}

func getValues(index int) []jsonElem {
	var values []jsonElem
	var elem jsonElem

	for i := index; i <= len(c.Options[index].Values); i++ {
		elem.Option = c.Options[index+1].GroupID
		if index < len(c.Options[index].Values) {
			elem.Values = getElems(index + 1)
		}
		values = append(values, elem)
	}
	return values
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

func genPartialsForTabs() {
	for i := range c.Combinations {

		for y := range c.Combinations[i].Tabs {
			link := "/"
			t := tab{}
			t.Name = c.Combinations[i].Tabs[y].Name

			for o := range c.Combinations[i].Options {
				for opt := range c.Options {
					if c.Options[opt].Name == c.Combinations[i].Options[o].Name {
						link = link + c.Options[opt].GroupID + "_"
					}
				}
			}
			link = link + c.Combinations[i].Tabs[y].LinkName + ".html"
			t.Permalink = link

			t.data = c.Combinations[i].Tabs[y].Params

			if _, err := os.Stat("templates/tabs/" + c.Combinations[i].Tabs[y].TemplateName + ".html"); !os.IsNotExist(err) {
				tpl, err := template.ParseFiles("templates/tabs/" + c.Combinations[i].Tabs[y].TemplateName + ".html")
				if err != nil {
					log.Print(err)
					return
				}

				if _, err := os.Stat("generated/" + t.Permalink); os.IsNotExist(err) {
					f, err := os.Create("generated/" + t.Permalink)
					if err != nil {
						log.Println("Create file: ", err)
						return
					}
					err = tpl.Execute(f, t)
					if err != nil {
						log.Print("Execute: ", err)
						return
					}
				}
			}
		}
	}
}
