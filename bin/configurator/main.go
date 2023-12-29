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
	generatedPagePathFormat             = "generated/pages_%s/configurator/pages/%s.html"
	generatedTitlePathFormat            = "generated/pages_%s/configurator/titles/%s.html"
	generatedPathCombinationTreeConfig  = "generated/configurator-options-list.json"
	generatedPathCombinationTabsConfig  = "generated/configurator-data.json"

	pagePermalinkConfigurator = "/includes/configurator.html"
	includePathFormatTabs     = "/configurator/tabs/%s.md"
	pagePermalinkFormatTabs   = "/configurator/tabs/%s.html"
	pagePermalinkFormatPages  = "/configurator/content/%s.html"
	pagePermalinkFormatTitles = "/configurator/title/%s.html"

	templatePathPagesConfigurator  = "templates/pages/configurator.html"
	templatePathTitlesConfigurator = "templates/pages/title.html"
	templatePathPagesTabs          = "templates/pages/tabs.html"
	templatePathIncludesTabs       = "templates/includes/tabs.html"

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

				// generate pages
				{
					for _, pageOptions := range configuratorPageOptions {
						var pageData configurationPageData
						pageData.Title = pageOptions.Title
						pageData.Groups = conf.getOptionGroupList()
						pageData.DefaultCombinationOptions = conf.getAllCombinationOptionsList()[0]
						pageData.DefaultIncludePath = getTabsIncludePath(pageData.DefaultCombinationOptions.ToSlug())

						pageData.IncludePath = getTabsIncludePath(combinationOptions.ToSlug())
						pageData.Permalink = getPagesPagePermalink(combinationOptions.ToUrlPath())

						pageData.PageTab = false

						pageFilePath := getPagesGeneratedPagePath(lang, combinationOptions.ToSlug())
						if err := createFileByTemplate(templatePathPagesConfigurator, pageFilePath, pageData); err != nil {
							return fmt.Errorf("unable to generate file %q: %s", pageFilePath, err)
						}
					}
				}

				// generate titles
				{
					for _, pageOptions := range configuratorPageOptions {
						var pageData configurationPageData
						pageData.Title = pageOptions.Title
						pageData.Groups = conf.getOptionGroupList()
						pageData.DefaultCombinationOptions = conf.getAllCombinationOptionsList()[0]
						pageData.DefaultIncludePath = getTabsIncludePath(pageData.DefaultCombinationOptions.ToSlug())

						pageData.IncludePath = getTabsIncludePath(combinationOptions.ToSlug())
						pageData.Permalink = getTitlesPagePermalink(combinationOptions.ToUrlPath())

						pageData.TitleParts, pageData.EndingTitle = combinationOptions.GetTitle()
						pageData.PageTab = true

						pageFilePath := getTitlesGeneratedPagePath(lang, combinationOptions.ToSlug())
						if err := createFileByTemplate(templatePathTitlesConfigurator, pageFilePath, pageData); err != nil {
							return fmt.Errorf("unable to generate file %q: %s", pageFilePath, err)
						}
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

func getPagesGeneratedPagePath(lang string, slug configCombinationSlug) string {
	return fmt.Sprintf(generatedPagePathFormat, lang, slug)
}

func getTitlesGeneratedPagePath(lang string, slug configCombinationSlug) string {
	return fmt.Sprintf(generatedTitlePathFormat, lang, slug)
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

func getPagesPagePermalink(slug configCombinationSlug) string {
	return fmt.Sprintf(pagePermalinkFormatPages, slug)
}

func getTitlesPagePermalink(slug configCombinationSlug) string {
	return fmt.Sprintf(pagePermalinkFormatTitles, slug)
}

type configurationPageData struct {
	Title                     string
	EndingTitle               string
	TitleParts                []string
	Permalink                 string
	IncludePath               string
	Groups                    []optionGroup
	DefaultCombinationOptions configCombinationOptions
	DefaultIncludePath        string
	PageTab                   bool
}

func generateConfiguratorPage(conf config) error {
	for _, pageOptions := range configuratorPageOptions {
		var pageData configurationPageData
		pageData.Title = pageOptions.Title
		pageData.Permalink = pageOptions.Permalink
		pageData.Groups = conf.getOptionGroupList()
		pageData.DefaultCombinationOptions = conf.getAllCombinationOptionsList()[0]
		pageData.DefaultIncludePath = getTabsIncludePath(pageData.DefaultCombinationOptions.ToSlug())
		pageData.PageTab = false

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
