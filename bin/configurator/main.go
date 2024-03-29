package main

import (
	"fmt"
	"os"
)

const (
	configFile = "config.yml"

	generatedDirectory                  = "generated"
	generatedIncludePathFormatTabs      = "generated/_includes/%s/configurator/tabs/%s.md"
	generatedPagePathFormatConfigurator = "generated/pages_%s/configurator.md"
	generatedPagePathFormatTabs         = "generated/pages_%s/configurator/tabs/%s.md"
	generatedPathCombinationTreeConfig  = "generated/configurator-options-list.json"
	generatedPathCombinationTabsConfig  = "generated/configurator-data.json"

	pagePermalinkConfigurator = "/includes/configurator.html"
	includePathFormatTabs     = "/configurator/tabs/%s.md"
	pagePermalinkFormatTabs   = "/configurator/tabs/%s.html"

	templatePathPagesConfigurator = "templates/pages/configurator.html"
	templatePathPagesTabs         = "templates/pages/tabs.html"
	templatePathIncludesTabs      = "templates/includes/tabs.html"

	pageTitleEnConfigurator = "Configurator"
	pageTitleRuConfigurator = "Конфигуратор"
)

var (
	languages = []string{
		"en",
		"ru",
	}
	configuratorPageOptions = []struct {
		Title     string
		Permalink string
		Language  string
	}{
		{pageTitleEnConfigurator, pagePermalinkConfigurator, "en"},
		{pageTitleRuConfigurator, pagePermalinkConfigurator, "ru"},
	}
)

type tabsIncludeData configCombination

type tabsPageData struct {
	Permalink   string
	IncludePath string
}

func main() {
	conf, err := readConfig()
	if err != nil {
		panic(fmt.Sprintf("Error while reading config: %s", err))
	}

	if err := cleanup(); err != nil {
		panic(fmt.Sprintf("Error while cleaning up: %s", err))
	}

	if err := generateCombinationTreeConfig(conf); err != nil {
		panic(fmt.Sprintf("Error while generating combination tree config: %s", err))
	}

	if err := generateCombinationTabsConfig(conf); err != nil {
		panic(fmt.Sprintf("Error while generating combination tabs config: %s", err))
	}

	if err := generateTabsIncludesAndPages(conf); err != nil {
		panic(fmt.Sprintf("Error while generating tabs includes and pages: %s", err))
	}

	if err := generateConfiguratorPage(conf); err != nil {
		panic(fmt.Sprintf("Error while generating configurator page: %s", err))
	}
}

func cleanup() error {
	if err := os.RemoveAll(generatedDirectory); err != nil {
		return fmt.Errorf("unable to remove directory %s: %w", generatedDirectory, err)
	}

	return nil
}

func generateCombinationTreeConfig(conf config) error {
	rootOption, err := conf.getCombinationTree()
	if err != nil {
		return fmt.Errorf("unable to generate combination tree: %s", err)
	}

	if err := createJsonFile(generatedPathCombinationTreeConfig, rootOption); err != nil {
		return fmt.Errorf("unable to create file %q: %s", generatedPathCombinationTreeConfig, err)
	}

	return nil
}

func generateTabsIncludesAndPages(conf config) error {
	for _, lang := range languages {
		for _, combination := range conf.Combinations {
			for _, combinationOptions := range conf.getCombinationOptionsList(combination) {
				// generate tabs includes
				{
					var pageData tabsIncludeData
					pageData.Options = combinationOptions
					for _, cCombTab := range combination.Tabs {
						pageData.Tabs = append(pageData.Tabs, cCombTab)
					}

					tabsIncludePath := getTabsGeneratedIncludePath(lang, combinationOptions.ToSlug())
					if err := createFileByTemplate(templatePathIncludesTabs, tabsIncludePath, pageData); err != nil {
						return fmt.Errorf("unable to create file %q: %s", tabsIncludePath, err)
					}
				}

				// generate tabs pages
				{
					tabsPagePath := getTabsGeneratedPagePath(lang, combinationOptions.ToSlug())
					tabsIncludePath := getTabsIncludePath(combinationOptions.ToSlug())
					permalink := getTabsPagePermalink(combinationOptions.ToSlug())
					if err := createFileByTemplate(templatePathPagesTabs, tabsPagePath, tabsPageData{
						Permalink:   permalink,
						IncludePath: tabsIncludePath,
					}); err != nil {
						return fmt.Errorf("unable to create file %q: %s", tabsPagePath, err)
					}
				}
			}
		}
	}

	return nil
}

func getTabsGeneratedPagePath(lang string, slug configCombinationSlug) string {
	return fmt.Sprintf(generatedPagePathFormatTabs, lang, slug)
}

func getTabsGeneratedIncludePath(lang string, slug configCombinationSlug) string {
	return fmt.Sprintf(generatedIncludePathFormatTabs, lang, slug)
}

func getTabsIncludePath(slug configCombinationSlug) string {
	return fmt.Sprintf(includePathFormatTabs, slug)
}

func getTabsPagePermalink(slug configCombinationSlug) string {
	return fmt.Sprintf(pagePermalinkFormatTabs, slug)
}

type configurationPageData struct {
	Title                     string
	Permalink                 string
	Groups                    []optionGroup
	DefaultCombinationOptions configCombinationOptions
	DefaultIncludePath        string
}

func generateConfiguratorPage(conf config) error {
	for _, pageOptions := range configuratorPageOptions {
		var pageData configurationPageData
		pageData.Title = pageOptions.Title
		pageData.Permalink = pageOptions.Permalink
		pageData.Groups = conf.getOptionGroupList()
		pageData.DefaultCombinationOptions = conf.getAllCombinationOptionsList()[0]
		pageData.DefaultIncludePath = getTabsIncludePath(pageData.DefaultCombinationOptions.ToSlug())

		pageFilePath := getConfiguratorGeneratedPagePath(pageOptions.Language)
		if err := createFileByTemplate(templatePathPagesConfigurator, pageFilePath, pageData); err != nil {
			return fmt.Errorf("unable to generate file %q: %s", pageFilePath, err)
		}
	}

	return nil
}

func getConfiguratorGeneratedPagePath(lang string) string {
	return fmt.Sprintf(generatedPagePathFormatConfigurator, lang)
}

func generateCombinationTabsConfig(conf config) error {
	combinationPermalinkMap := make(map[configCombinationSlug]string)
	for _, combinationOptions := range conf.getAllCombinationOptionsList() {
		combinationPermalinkMap[combinationOptions.ToSlug()] = getTabsPagePermalink(combinationOptions.ToSlug())
	}

	if err := createJsonFile(generatedPathCombinationTabsConfig, combinationPermalinkMap); err != nil {
		return fmt.Errorf("unable to create file %q: %s", generatedPathCombinationTabsConfig, err)
	}

	return nil
}
