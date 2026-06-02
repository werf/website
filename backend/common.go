package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const (
	currentDocsMajorEnvVar           = "CURRENT_DOCS_MAJOR"
	supportedDocsMajorVersionsEnvVar = "SUPPORTED_DOCS_MAJOR_VERSIONS"
	docsLatestAliasEnabledEnvVar     = "DOCS_LATEST_ALIAS_ENABLED"
	legacyActiveReleaseEnvVar        = "ACTIVE_RELEASE"
	defaultDocsMajorRoot             = "v1"
)

type docsRoutingConfig struct {
	CurrentMajor   string
	SupportedRoots []string
	LatestAlias    bool
}

type ChannelType struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type ReleaseType struct {
	Group    string        `yaml:"name"`
	Channels []ChannelType `yaml:"channels"`
}

type ReleasesStatusType struct {
	Releases []ReleaseType `yaml:"groups"`
}

type ApiStatusResponseType struct {
	Status         string        `json:"status"`
	Msg            string        `json:"msg"`
	RootVersion    string        `json:"rootVersion"`
	RootVersionURL string        `json:"rootVersionURL"`
	Multiwerf      []ReleaseType `json:"multiwerf"`
}

type versionMenuType struct {
	VersionItems           []versionMenuItems
	HTMLContent            string
	CurrentGroup           string
	CurrentChannel         string
	CurrentVersion         string
	AbsoluteVersion        string // Contains explicit version, used for getting git link to source file
	CurrentVersionURL      string
	CurrentPageURLRelative string // Relative URL, without "/docs/<version>"
	CurrentPageURL         string // Full page URL
	MenuDocumentationLink  string
}

type versionMenuItems struct {
	Group      string
	Channel    string
	Version    string
	VersionURL string // Base URL for corresponding version without a leading /, e.g. 'v1.2.3-plus-fix6'.
	IsCurrent  bool
}

var ReleasesStatus ReleasesStatusType

func (m *versionMenuType) getChannelMenuData(r *http.Request, releases *ReleasesStatusType) (err error) {
	_ = releases
	return m.populateMajorVersionMenu(r)
}

func (m *versionMenuType) getVersionMenuData(r *http.Request, releases *ReleasesStatusType) (err error) {
	_ = releases
	return m.populateMajorVersionMenu(r)
}

func (m *versionMenuType) getGroupMenuData(r *http.Request, releases *ReleasesStatusType) (err error) {
	_ = releases
	return m.populateMajorVersionMenu(r)
}

func (m *versionMenuType) populateMajorVersionMenu(r *http.Request) (err error) {
	err = nil

	m.CurrentPageURLRelative = getDocPageURLRelative(r)
	m.CurrentPageURL = getCurrentPageURL(r)
	m.CurrentVersionURL = getVersionURL(r)
	m.CurrentVersion = URLToVersion(m.CurrentVersionURL)

	requestedVersion := m.CurrentVersion
	if requestedVersion == "" {
		requestedVersion = getCurrentDocsMajorRoot()
	}

	canonicalVersion := getCanonicalDocsVersion(requestedVersion)
	displayVersion := canonicalVersion
	if requestedVersion == "latest" && getDocsLatestAliasEnabled() {
		displayVersion = "latest"
	}
	m.CurrentVersion = displayVersion
	m.CurrentVersionURL = VersionToURL(displayVersion)

	m.MenuDocumentationLink = fmt.Sprintf("/docs/%s/", VersionToURL(canonicalVersion))
	m.AbsoluteVersion = canonicalVersion
	m.CurrentGroup = strings.TrimPrefix(canonicalVersion, "v")
	m.CurrentChannel = ""

	m.VersionItems = append(m.VersionItems, versionMenuItems{
		Version:    m.CurrentVersion,
		VersionURL: VersionToURL(m.CurrentVersion),
		IsCurrent:  true,
	})

	if getDocsLatestAliasEnabled() && requestedVersion != "latest" {
		m.VersionItems = append(m.VersionItems, versionMenuItems{
			Version:    "latest",
			VersionURL: "latest",
			IsCurrent:  false,
		})
	}

	for _, root := range getSupportedDocsMajorRoots() {
		if root == displayVersion {
			continue
		}

		m.VersionItems = append(m.VersionItems, versionMenuItems{
			Version:    root,
			VersionURL: VersionToURL(root),
			IsCurrent:  false,
		})
	}

	return err
}

func getCurrentDocsMajorRoot() string {
	return getDocsRoutingConfig().CurrentMajor
}

func getCanonicalDocsVersion(version string) string {
	version = strings.TrimSpace(URLToVersion(version))
	if version == "" || version == "latest" {
		return getCurrentDocsMajorRoot()
	}

	if strings.HasPrefix(version, "pr-") {
		return version
	}

	re := regexp.MustCompile(`^v([0-9]+)`)
	if res := re.FindStringSubmatch(version); res != nil {
		return fmt.Sprintf("v%s", res[1])
	}

	return getCurrentDocsMajorRoot()
}

func getSupportedDocsMajorRoots() []string {
	return getDocsRoutingConfig().SupportedRoots
}

func getDocsLatestAliasEnabled() bool {
	return getDocsRoutingConfig().LatestAlias
}

func getDocsRoutingConfig() docsRoutingConfig {
	configuredCurrent := normalizeDocsMajorRoot(os.Getenv(currentDocsMajorEnvVar))
	configuredRoots := getConfiguredDocsMajorRoots(os.Getenv(supportedDocsMajorVersionsEnvVar))

	currentRoot := configuredCurrent
	if currentRoot == "" && len(configuredRoots) > 0 {
		currentRoot = configuredRoots[0]
	}
	if currentRoot == "" {
		currentRoot = normalizeDocsMajorRoot(getRootRelease())
	}
	if currentRoot == "" {
		currentRoot = defaultDocsMajorRoot
	}

	roots := []string{currentRoot}
	seen := map[string]struct{}{currentRoot: {}}
	for _, root := range configuredRoots {
		if _, ok := seen[root]; ok {
			continue
		}

		roots = append(roots, root)
		seen[root] = struct{}{}
	}

	latestAliasEnabled := shouldEnableDocsLatestAlias(os.Getenv(docsLatestAliasEnabledEnvVar), currentRoot, roots)

	return docsRoutingConfig{
		CurrentMajor:   currentRoot,
		SupportedRoots: roots,
		LatestAlias:    latestAliasEnabled,
	}
}

func shouldEnableDocsLatestAlias(raw string, currentRoot string, roots []string) bool {
	if value, ok := parseOptionalBool(raw); ok {
		return value
	}

	return !(currentRoot == "v1" && len(roots) == 1 && roots[0] == "v1")
}

func parseOptionalBool(raw string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "on":
		return true, true
	case "0", "false", "no", "off":
		return false, true
	default:
		return false, false
	}
}

func getConfiguredDocsMajorRoots(configuredRoots string) []string {
	roots := []string{}
	seen := map[string]struct{}{}
	for _, rawRoot := range strings.FieldsFunc(configuredRoots, func(r rune) bool {
		return r == ',' || r == ';' || r == ' ' || r == '\n' || r == '\t'
	}) {
		root := normalizeDocsMajorRoot(rawRoot)
		if root == "" {
			continue
		}
		if _, ok := seen[root]; ok {
			continue
		}

		roots = append(roots, root)
		seen[root] = struct{}{}
	}

	return roots
}

func normalizeDocsMajorRoot(raw string) string {
	root := strings.ToLower(strings.TrimSpace(raw))
	root = strings.TrimPrefix(root, "/docs/")
	root = strings.Trim(root, "/")
	root = strings.TrimPrefix(root, "v")
	if root == "" {
		return ""
	}

	if matched, _ := regexp.MatchString(`^[0-9]+$`, root); matched {
		return fmt.Sprintf("v%s", root)
	}

	return ""
}

// Gev version from specified group
// E.g. get v1.2.3+fix6 from v1.2
func getVersionFromGroup(releases *ReleasesStatusType, group string) (err error, version string) {
	_ = releases
	if group == "latest" {
		return nil, getCurrentDocsMajorRoot()
	}
	if strings.HasPrefix(group, "v") {
		return nil, getCanonicalDocsVersion(group)
	}
	if matched, _ := regexp.MatchString(`^[0-9]+$`, group); matched {
		return nil, getCanonicalDocsVersion(fmt.Sprintf("v%s", group))
	}

	return fmt.Errorf("can't get version for %s", group), ""
}

// Add prefix 'v' to a version if it doesn't have yet
func normalizeVersion(version string) string {
	if strings.HasPrefix(version, "v") || version == "latest" {
		return version
	} else {
		return fmt.Sprintf("v%s", version)
	}
}

func getRootReleaseVersion() string {
	return getCurrentDocsMajorRoot()
}

func getRootRelease() (result string) {
	if len(os.Getenv(legacyActiveReleaseEnvVar)) > 0 {
		result = os.Getenv(legacyActiveReleaseEnvVar)
	}

	return
}

// Get the full page URL menu requested for
// E.g /docs/v1.2.3/reference/build_process.html
func getCurrentPageURL(r *http.Request) (result string) {
	originalURI, err := url.Parse(r.Header.Get("x-original-uri"))
	if err != nil {
		return
	}

	if originalURI.Path == "/404.html" {
		return
	} else {
		return originalURI.Path
	}
}

// Get page URL menu requested for without a leading version suffix
// E.g /reference/build_process.html for /docs/v1.2.3/reference/build_process.html
// if useURI == true - use requestURI instead of x-original-uri header value
func getDocPageURLRelative(r *http.Request, useURI ...bool) (result string) {
	var (
		URLtoParse  string
		originalURI *url.URL
		err         error
	)

	if len(useURI) > 0 && useURI[0] {
		originalURI, err = url.Parse(r.RequestURI)
	} else {
		originalURI, err = url.Parse(r.Header.Get("x-original-uri"))
	}

	if err != nil {
		return
	}

	if originalURI.Path == "/404.html" {
		return
	} else {
		URLtoParse = originalURI.Path
	}

	re := regexp.MustCompile(`^/docs/[^/]+(/.+)$`)
	res := re.FindStringSubmatch(URLtoParse)
	if res != nil {
		result = res[1]
	} else {
		result = "/"
	}
	return
}

// Get version URL page belongs to if request came from concrete documentation version, otherwise empty.
// E.g for the /docs/v1.2.3-plus-fix5/reference/build_process.html return "v1.2.3-plus-fix5".
func getVersionURL(r *http.Request) (result string) {
	URLtoParse := ""
	originalURI, err := url.Parse(r.Header.Get("x-original-uri"))
	if err != nil {
		return
	}

	if originalURI.Path == "/404.html" {
		values, err := url.ParseQuery(originalURI.RawQuery)
		if err != nil {
			return
		}
		URLtoParse = values.Get("uri")
	} else {
		URLtoParse = originalURI.Path
	}

	re := regexp.MustCompile(`^/docs/([^/]+)/?.*$`)
	res := re.FindStringSubmatch(URLtoParse)
	if res != nil {
		result = res[1]
	}

	return strings.TrimPrefix(result, "/")
}

func inspectRequestHTML(r *http.Request) string {
	var request []string

	request = append(request, "<p>")
	requestLine := fmt.Sprintf("<b>%v</b> %v %v", r.Method, r.URL, r.Proto)
	request = append(request, requestLine)
	request = append(request, fmt.Sprintf("<b>Host:</b> %v", r.Host))
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("<b>%v:</b> %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		_ = r.ParseForm()
		request = append(request, r.Form.Encode())
	}

	request = append(request, "</p>")
	return strings.Join(request, "<br />")
}

func VersionToURL(version string) string {
	result := strings.ReplaceAll(version, "+", "-plus-")
	result = strings.ReplaceAll(result, "_", "-u-")
	return normalizeVersion(result)
}

func URLToVersion(version string) (result string) {
	result = strings.ReplaceAll(version, "-plus-", "+")
	result = strings.ReplaceAll(result, "-u-", "_")
	return
}

func getRootFilesPath(r *http.Request) (result string) {
	result = "./root/"
	if strings.HasPrefix(r.Host, "ru.") || strings.HasPrefix(r.Host, "ru-") {
		result += "ru"
	} else {
		result += "en"
	}
	return
}
