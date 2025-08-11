package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	pathToConfiguratorYaml = "static/_data/_common/configurator.yaml"
	startPhraseRu          = "Установка и запуск на "
	startPhraseEn          = "Install and run on "
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

type configNames struct {
	Groups []struct {
		Name  string `yaml:"name"`
		Title struct {
			En string `yaml:"en"`
			Ru string `yaml:"ru"`
		} `yaml:"title"`
		Tooltip struct {
			En string `yaml:"en"`
			Ru string `yaml:"ru"`
		} `yaml:"tooltip,omitempty"`
		Buttons []struct {
			Name  string `yaml:"name"`
			Title struct {
				En       string `yaml:"en"`
				Ru       string `yaml:"ru"`
				PageName struct {
					En string `yaml:"en"`
					Ru string `yaml:"ru"`
				} `yaml:"page-name"`
			} `yaml:"title"`
		} `yaml:"buttons"`
	} `yaml:"groups"`
	Tabs []struct {
		Name  string `yaml:"name"`
		Title struct {
			En string `yaml:"en"`
			Ru string `yaml:"ru"`
		} `yaml:"title"`
	} `yaml:"tabs"`
}

type titlesStruct struct {
	Ru string
	En string
}

func (options configCombinationOptions) ToSlug() configCombinationSlug {
	var opts []string
	for _, option := range options {
		opts = append(opts, fmt.Sprintf("%s-%s", option.Name, option.Value))
	}

	return configCombinationSlug(strings.Join(opts, "_"))
}

func (options configCombinationOptions) ToUrlPath() configCombinationSlug {
	var opts []string
	var usage string
	for _, option := range options {
		if option.Value == "localDev" {
			usage = "local/"
		} else if option.Value == "ci" {
			usage = "cicd/"
		} else {
			opts = append(opts, strings.ToLower(option.Value))
		}
	}

	return configCombinationSlug(usage + strings.Join(opts, "-"))
}

func (options configCombinationOptions) GetTitle(lang string) string {
	var opts []configCombinationOption
	for _, option := range options {
		if option.Name != "usage" &&
			option.Name != "repoType" &&
			option.Name != "sharedCICD" &&
			option.Name != "projectType" {
			opts = append(opts, option)
		}
	}

	file, err := os.Open(pathToConfiguratorYaml)
	if err != nil {
		log.Fatalf("Failed to open configurator.yaml file: %v", err)
	}
	defer file.Close()
	confData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to open configurator.yaml file: %v", err)
	}
	var configNames configNames
	err = yaml.Unmarshal(confData, &configNames)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML data: %v", err)
	}

	var titles titlesStruct
	titles.Ru = startPhraseRu
	titles.En = startPhraseEn
	count := 0
	for _, opt := range opts {
		count++
		for _, group := range configNames.Groups {
			if group.Name == opt.Name {
				for _, btn := range group.Buttons {
					if opt.Value == btn.Name {
						switch lang {
						case "ru":
							if len(btn.Title.PageName.Ru) != 0 {
								titles.Ru += btn.Title.PageName.Ru
							} else {
								titles.Ru += btn.Title.Ru
							}
							if count < len(opts) {
								if count < len(opts)-1 {
									titles.Ru += ", "
								} else {
									titles.Ru += " и "
								}
							}
						case "en":
							if len(btn.Title.PageName.Ru) != 0 {
								titles.En += btn.Title.PageName.En
							} else {
								titles.En += btn.Title.En
							}
							if count < len(opts) {
								if count < len(opts)-1 {
									titles.En += ", "
								} else {
									titles.En += " and "
								}
							}
						}
					}
				}
			}
		}
	}

	if lang == "ru" {
		return titles.Ru
	} else if lang == "en" {
		return titles.En
	}

	return ""
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
