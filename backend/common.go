package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type ApiStatusResponseType struct {
	Status         string   `json:"status"`
	Msg            string   `json:"msg"`
	RootVersion    string   `json:"rootVersion"`
	RootVersionURL string   `json:"rootVersionURL"`
	Multiwerf      []string `json:"multiwerf"`
}

type versionMenuType struct {
	VersionItems           []versionMenuItems
	HTMLContent            string
	CurrentVersion         string
	AbsoluteVersion        string // Contains explicit version, used for getting git link to source file
	CurrentVersionURL      string
	CurrentPageURLRelative string // Relative URL, without "/docs/<version>"
	CurrentPageURL         string // Full page URL
	MenuDocumentationLink  string
}

type versionMenuItems struct {
	Version    string
	VersionURL string // Base URL for corresponding version without a leading /, e.g. 'v1.2.3-plus-fix6'.
	IsCurrent  bool
}

const (
	defaultCurrentDocsRoot    = "v2"
	currentDocsRootEnv        = "CURRENT_DOCS_MAJOR"
	supportedDocsRootsEnv     = "SUPPORTED_DOCS_MAJOR_VERSIONS"
	latestDocsAliasEnabledEnv = "DOCS_LATEST_ALIAS_ENABLED"
)

var legacyVersionAliasPattern = regexp.MustCompile(`^v(2|1\.2)-(alpha|beta|ea|stable|rock-solid)$`)

func appendVersionItems(items []versionMenuItems, currentVersion string, includeLatest bool) []versionMenuItems {
	if includeLatest && currentVersion != "latest" {
		items = append(items, versionMenuItems{
			Version:    "latest",
			VersionURL: "latest",
			IsCurrent:  false,
		})
	}

	for _, root := range getSupportedDocsRoots() {
		if root == currentVersion {
			continue
		}
		items = append(items, versionMenuItems{
			Version:    root,
			VersionURL: VersionToURL(root),
			IsCurrent:  false,
		})
	}

	return items
}

func (m *versionMenuType) getLegacyVersionMenuData(r *http.Request) (err error) {
	err = nil

	m.CurrentPageURLRelative = getDocPageURLRelative(r)
	m.CurrentPageURL = getCurrentPageURL(r)
	m.CurrentVersion = getCanonicalDocsVersion(getVersionURL(r))
	m.CurrentVersionURL = VersionToURL(m.CurrentVersion)

	m.VersionItems = append(m.VersionItems, versionMenuItems{
		Version:    m.CurrentVersion,
		VersionURL: m.CurrentVersionURL,
		IsCurrent:  true,
	})

	m.VersionItems = appendVersionItems(m.VersionItems, m.CurrentVersion, isLatestAliasEnabled())

	return
}

func (m *versionMenuType) getVersionMenuData(r *http.Request) (err error) {
	err = nil

	m.CurrentPageURLRelative = getDocPageURLRelative(r)
	m.CurrentPageURL = getCurrentPageURL(r)
	m.CurrentVersion = getCanonicalDocsVersion(getVersionURL(r))
	m.CurrentVersionURL = VersionToURL(m.CurrentVersion)
	m.AbsoluteVersion = m.CurrentVersion
	m.MenuDocumentationLink = fmt.Sprintf("/docs/%s/", m.CurrentVersionURL)

	m.VersionItems = append(m.VersionItems, versionMenuItems{
		Version:    m.CurrentVersion,
		VersionURL: m.CurrentVersionURL,
		IsCurrent:  true,
	})

	m.VersionItems = appendVersionItems(m.VersionItems, m.CurrentVersion, isLatestAliasEnabled())

	return
}

func (m *versionMenuType) getRootVersionMenuData(r *http.Request) (err error) {
	err = nil

	m.CurrentPageURLRelative = getDocPageURLRelative(r)
	m.CurrentPageURL = getCurrentPageURL(r)
	m.CurrentVersion = getCanonicalDocsVersion(getVersionURL(r))
	m.CurrentVersionURL = VersionToURL(m.CurrentVersion)

	m.VersionItems = append(m.VersionItems, versionMenuItems{
		Version:    m.CurrentVersion,
		VersionURL: m.CurrentVersionURL,
		IsCurrent:  true,
	})

	m.VersionItems = appendVersionItems(m.VersionItems, m.CurrentVersion, false)

	return
}

// Resolve version alias from legacy path like /docs/v2-stable/.
func resolveLegacyVersionAlias(track, major string) (err error, version string) {
	return nil, getCanonicalDocsVersion(fmt.Sprintf("v%s-%s", strings.TrimPrefix(major, "v"), track))
}

// Resolve docs version from root alias.
func resolveRootVersionAlias(alias string) (err error, version string) {
	candidate := alias
	if candidate != "latest" && !strings.HasPrefix(candidate, "v") {
		candidate = "v" + candidate
	}
	canonical := getCanonicalDocsVersion(candidate)
	if canonical == getCurrentDocsRoot() && candidate != canonical && candidate != "latest" {
		return fmt.Errorf("unsupported docs version alias: %s", alias), ""
	}
	return nil, canonical
}

// Add prefix 'v' to a version if it doesn't have yet
func normalizeVersion(version string) string {
	v := strings.TrimSpace(version)
	if strings.HasPrefix(v, "v") || v == "latest" || strings.HasPrefix(v, "pr-") {
		return v
	}
	return fmt.Sprintf("v%s", v)
}

func getRootReleaseVersion() string {
	return getCurrentDocsRoot()
}

func getRootRelease() (result string) {
	return strings.TrimPrefix(getCurrentDocsRoot(), "v")
}

func isLatestAliasEnabled() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(latestDocsAliasEnabledEnv)))
	if v == "" {
		return true
	}
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

func getCurrentDocsRoot() string {
	candidate := normalizeDocsRoot(os.Getenv(currentDocsRootEnv))
	if candidate != "" {
		return candidate
	}
	return defaultCurrentDocsRoot
}

func getSupportedDocsRoots() []string {
	raw := strings.Split(strings.TrimSpace(os.Getenv(supportedDocsRootsEnv)), ",")
	roots := make([]string, 0, len(raw)+1)
	seen := map[string]struct{}{}

	add := func(v string) {
		if v == "" {
			return
		}
		if _, ok := seen[v]; ok {
			return
		}
		seen[v] = struct{}{}
		roots = append(roots, v)
	}

	add(getCurrentDocsRoot())
	for _, item := range raw {
		add(normalizeDocsRoot(item))
	}

	return roots
}

func normalizeDocsRoot(raw string) string {
	v := strings.TrimSpace(raw)
	if v == "" {
		return ""
	}
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	if matched, _ := regexp.MatchString(`^v[0-9]+(?:\.[0-9]+)?$`, v); !matched {
		return ""
	}
	return v
}

func getCanonicalDocsVersion(raw string) string {
	v := strings.TrimSpace(URLToVersion(raw))
	if v == "" || v == "latest" {
		return getCurrentDocsRoot()
	}
	if strings.HasPrefix(v, "pr-") {
		return v
	}
	if legacyVersionAliasPattern.MatchString(v) {
		parts := strings.SplitN(v[1:], "-", 2)
		if len(parts) == 2 {
			return getCanonicalDocsVersion("v" + parts[0])
		}
	}
	if strings.HasPrefix(v, "v1.2") {
		return "v1.2"
	}
	if strings.HasPrefix(v, "v2") {
		return "v2"
	}
	return getCurrentDocsRoot()
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
	url := fmt.Sprintf("<b>%v</b> %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
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

func validateURL(url string) error {
	_ = url
	return nil
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
