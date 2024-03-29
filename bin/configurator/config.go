package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type config struct {
	Options      []configOption      `yaml:"options"`
	Combinations []configCombination `yaml:"combinations"`
}

type configOption struct {
	Name   string   `yaml:"name"`
	Group  string   `yaml:"group"`
	Values []string `yaml:"values"`
}

type configCombination struct {
	Tabs    []configCombinationTab   `yaml:"tabs"`
	Options configCombinationOptions `yaml:"options"`
}

type configCombinationOptions []configCombinationOption

type configCombinationOption struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type configCombinationSlug string

func (options configCombinationOptions) ToSlug() configCombinationSlug {
	var opts []string
	for _, option := range options {
		opts = append(opts, fmt.Sprintf("%s-%s", option.Name, option.Value))
	}

	return configCombinationSlug(strings.Join(opts, "_"))
}

type configCombinationTab struct {
	Name        string            `yaml:"name"`
	IncludePath string            `yaml:"includePath"`
	Params      map[string]string `yaml:"params,omitempty"`
}

func (c config) mustGetCombinationTree() optionTreeNode {
	rootOption, err := c.getCombinationTree()
	if err != nil {
		panic(err)
	}

	return rootOption
}

type optionTreeNode struct {
	Option string               `json:"option"`
	Values optionTreeNodeValues `json:"values"`
}

type optionTreeNodeValues map[string]*optionTreeNode

func (c config) getCombinationTree() (rootNode optionTreeNode, err error) {
	for _, combination := range c.Combinations {
		var currentNodeList []*optionTreeNode
		var counterCombinationOption int

	optionLoop:
		for _, option := range c.Options {
			for _, combinationOption := range combination.Options {
				if combinationOption.Name == option.Name {
					if counterCombinationOption == 0 { // first combination option
						currentNodeList = append(currentNodeList, &rootNode)
					}

					var newCurrentValueOptions []*optionTreeNode
					for _, currentValueOption := range currentNodeList {
						if currentValueOption.Option == "" { // new branch in the tree
							currentValueOption.Option = combinationOption.Name
						} else if currentValueOption.Option != combinationOption.Name {
							return rootNode, fmt.Errorf("invalid configuration file: conflict between combinations: combination option value cannot be linked to various options (%q != %q)", currentValueOption.Option, combinationOption.Name)
						}

						for _, optionValue := range option.Values {
							if regexp.MustCompile(combinationOption.Value).MatchString(optionValue) {
								if currentValueOption.Values == nil {
									currentValueOption.Values = make(optionTreeNodeValues)
								}

								if valueOpt, ok := currentValueOption.Values[optionValue]; !ok {
									if counterCombinationOption == len(combination.Options)-1 { // last optionTreeNode
										currentValueOption.Values[optionValue] = nil
										newCurrentValueOptions = append(newCurrentValueOptions, nil)
									} else {
										newNode := &optionTreeNode{}
										currentValueOption.Values[optionValue] = newNode
										newCurrentValueOptions = append(newCurrentValueOptions, newNode)
									}
								} else {
									newCurrentValueOptions = append(newCurrentValueOptions, valueOpt)
								}

								currentNodeList = newCurrentValueOptions
							}
						}
					}

					if len(newCurrentValueOptions) == 0 {
						return rootNode, fmt.Errorf("invalid configuration file: combination option %q value %q not supported", combinationOption.Name, combinationOption.Value)
					}

					counterCombinationOption++
					continue optionLoop
				}
			}
		}
	}

	return
}

type optionGroup []configOption

func (c config) getOptionGroupList() []optionGroup {
	var groupList []string
	{
	optionsLoop:
		for _, option := range c.Options {
			for _, group := range groupList {
				if group == option.Group {
					continue optionsLoop
				}
			}

			groupList = append(groupList, option.Group)
		}
	}

	var optionGroupList []optionGroup
	for _, currentGroup := range groupList {
		var currentOptionGroup optionGroup
		for _, option := range c.Options {
			if option.Group != currentGroup {
				continue
			}

			currentOptionGroup = append(currentOptionGroup, option)
		}

		optionGroupList = append(optionGroupList, currentOptionGroup)
	}

	return optionGroupList
}

func (c config) getAllCombinationOptionsList() []configCombinationOptions {
	var combinationOptionsList []configCombinationOptions
	for _, combination := range c.Combinations {
		combinationOptionsList = append(combinationOptionsList, c.getCombinationOptionsList(combination)...)
	}

	return combinationOptionsList
}

func (c config) getCombinationOptionsList(combination configCombination) []configCombinationOptions {
	var combinationOptionsList []configCombinationOptions

optionsLoop:
	for _, option := range c.Options {
		for _, combinationOption := range combination.Options {
			var newCombinationOptionsList []configCombinationOptions
			if option.Name != combinationOption.Name {
				continue
			}

			for _, optionValue := range option.Values {
				if regexp.MustCompile(combinationOption.Value).MatchString(optionValue) {
					combinationOption := configCombinationOption{
						Name:  option.Name,
						Value: optionValue,
					}

					if len(combinationOptionsList) == 0 {
						newCombinationOptionsList = append(newCombinationOptionsList, configCombinationOptions{combinationOption})
					} else {
						for _, combinationOptions := range combinationOptionsList {
							newCombinationOptions := combinationOptions
							newCombinationOptions = append(newCombinationOptions, combinationOption)
							newCombinationOptionsList = append(newCombinationOptionsList, newCombinationOptions)
						}
					}
				}
			}

			combinationOptionsList = newCombinationOptionsList
			continue optionsLoop
		}
	}

	return combinationOptionsList
}

func readConfig() (config, error) {
	var cfg config

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return cfg, fmt.Errorf("unable to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("unable to unmarshal config file: %w", err)
	}

	return cfg, nil
}
