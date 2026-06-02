package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestGetCurrentDocsMajorRootUsesSupportedRootsWhenCurrentIsUnset(t *testing.T) {
	t.Setenv("CURRENT_DOCS_MAJOR", "")
	t.Setenv("SUPPORTED_DOCS_MAJOR_VERSIONS", "v3, v2, v1")
	t.Setenv("ACTIVE_RELEASE", "")

	if actual := getCurrentDocsMajorRoot(); actual != "v3" {
		t.Fatalf("expected current docs major v3, got %s", actual)
	}

	expected := []string{"v3", "v2", "v1"}
	if actual := getSupportedDocsMajorRoots(); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected supported roots %v, got %v", expected, actual)
	}
}

func TestGetCurrentDocsMajorRootPrefersExplicitConfigOverLegacyFallback(t *testing.T) {
	t.Setenv("CURRENT_DOCS_MAJOR", "v3")
	t.Setenv("SUPPORTED_DOCS_MAJOR_VERSIONS", "v2,v1")
	t.Setenv("ACTIVE_RELEASE", "2")

	if actual := getCurrentDocsMajorRoot(); actual != "v3" {
		t.Fatalf("expected current docs major v3, got %s", actual)
	}

	expected := []string{"v3", "v2", "v1"}
	if actual := getSupportedDocsMajorRoots(); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected supported roots %v, got %v", expected, actual)
	}
}

func TestGetCurrentDocsMajorRootUsesLegacyFallbackOnlyWhenModernConfigMissing(t *testing.T) {
	t.Setenv("CURRENT_DOCS_MAJOR", "")
	t.Setenv("SUPPORTED_DOCS_MAJOR_VERSIONS", "")
	t.Setenv("ACTIVE_RELEASE", "3")

	if actual := getCurrentDocsMajorRoot(); actual != "v3" {
		t.Fatalf("expected current docs major v3, got %s", actual)
	}

	expected := []string{"v3"}
	if actual := getSupportedDocsMajorRoots(); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected supported roots %v, got %v", expected, actual)
	}
}

func TestGetCurrentDocsMajorRootFallsBackToNeutralDefaultWithoutImplicitV2(t *testing.T) {
	t.Setenv("CURRENT_DOCS_MAJOR", "")
	t.Setenv("SUPPORTED_DOCS_MAJOR_VERSIONS", "")
	t.Setenv("ACTIVE_RELEASE", "")

	for _, key := range []string{"CURRENT_DOCS_MAJOR", "SUPPORTED_DOCS_MAJOR_VERSIONS", "ACTIVE_RELEASE"} {
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("failed to unset %s: %v", key, err)
		}
	}

	if actual := getCurrentDocsMajorRoot(); actual != "v1" {
		t.Fatalf("expected fallback current docs major v1, got %s", actual)
	}

	expected := []string{"v1"}
	if actual := getSupportedDocsMajorRoots(); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected supported roots %v, got %v", expected, actual)
	}
}

func TestGetCanonicalDocsVersion(t *testing.T) {
	t.Setenv("CURRENT_DOCS_MAJOR", "v3")
	t.Setenv("SUPPORTED_DOCS_MAJOR_VERSIONS", "v3,v2,v1")

	testCases := map[string]string{
		"":                   "v3",
		"latest":             "v3",
		"v3":                 "v3",
		"v3.0.1":             "v3",
		"v2":                 "v2",
		"v2.68.2":            "v2",
		"v2-stable":          "v2",
		"v1":                 "v1",
		"v1.2.294-plus-fix1": "v1",
	}

	for input, expected := range testCases {
		if actual := getCanonicalDocsVersion(input); actual != expected {
			t.Fatalf("expected canonical version for %q to be %q, got %q", input, expected, actual)
		}
	}
}

func TestGetSupportedDocsMajorRoots(t *testing.T) {
	t.Setenv("CURRENT_DOCS_MAJOR", "v2")
	t.Setenv("SUPPORTED_DOCS_MAJOR_VERSIONS", "v3, v2, v1")

	expected := []string{"v2", "v3", "v1"}
	if actual := getSupportedDocsMajorRoots(); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected supported roots %v, got %v", expected, actual)
	}
}

func TestGetSupportedDocsMajorRootsDoesNotInjectImplicitLegacyRoots(t *testing.T) {
	t.Setenv("CURRENT_DOCS_MAJOR", "v3")
	t.Setenv("SUPPORTED_DOCS_MAJOR_VERSIONS", "v3, v2")

	expected := []string{"v3", "v2"}
	if actual := getSupportedDocsMajorRoots(); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected supported roots %v, got %v", expected, actual)
	}
}

func TestDocsCompatibilityRoutesRedirectToMajorRoots(t *testing.T) {
	t.Setenv("CURRENT_DOCS_MAJOR", "v3")
	t.Setenv("SUPPORTED_DOCS_MAJOR_VERSIONS", "v3,v2,v1")

	r := newRouter()
	testCases := []struct {
		name     string
		path     string
		location string
	}{
		{name: "latest alias", path: "/docs/latest/reference/cli/overview.html", location: "/docs/v3/reference/cli/overview.html"},
		{name: "legacy patch version", path: "/docs/v1.2.294-plus-fix1/usage/build/overview.html", location: "/docs/v1/usage/build/overview.html"},
		{name: "legacy channel", path: "/docs/v2-stable/reference/werf_yaml.html", location: "/docs/v2/reference/werf_yaml.html"},
		{name: "current major channel alias", path: "/docs/v3-stable/reference/werf_yaml.html", location: "/docs/v3/reference/werf_yaml.html"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != http.StatusFound {
				t.Fatalf("expected %d for %s, got %d", http.StatusFound, tc.path, rec.Code)
			}

			if location := rec.Header().Get("Location"); location != tc.location {
				t.Fatalf("expected redirect location %q, got %q", tc.location, location)
			}
		})
	}
}

func TestVersionMenuDataUsesMajorRoots(t *testing.T) {
	t.Setenv("CURRENT_DOCS_MAJOR", "v3")
	t.Setenv("SUPPORTED_DOCS_MAJOR_VERSIONS", "v3,v2,v1")

	req := httptest.NewRequest(http.MethodGet, "/includes/channel-menu-v2.html", nil)
	req.Header.Set("x-original-uri", "/docs/v1.2/reference/cli/overview.html")

	menu := versionMenuType{VersionItems: []versionMenuItems{}}
	if err := menu.getVersionMenuData(req, nil); err != nil {
		t.Fatalf("unexpected error populating version menu: %v", err)
	}

	if menu.AbsoluteVersion != "v1" {
		t.Fatalf("expected absolute version v1, got %s", menu.AbsoluteVersion)
	}

	if menu.MenuDocumentationLink != "/docs/v1/" {
		t.Fatalf("expected menu documentation link /docs/v1/, got %s", menu.MenuDocumentationLink)
	}

	gotVersions := []string{}
	for _, item := range menu.VersionItems {
		gotVersions = append(gotVersions, item.Version)
	}

	expectedVersions := []string{"v1", "latest", "v3", "v2"}
	if !reflect.DeepEqual(gotVersions, expectedVersions) {
		t.Fatalf("expected menu versions %v, got %v", expectedVersions, gotVersions)
	}
}

func TestVersionMenuDataKeepsCanonicalMajorWhenCurrentIsLatest(t *testing.T) {
	t.Setenv("CURRENT_DOCS_MAJOR", "v3")
	t.Setenv("SUPPORTED_DOCS_MAJOR_VERSIONS", "v3,v2,v1")

	req := httptest.NewRequest(http.MethodGet, "/includes/channel-menu-v2.html", nil)
	req.Header.Set("x-original-uri", "/docs/latest/reference/cli/overview.html")

	menu := versionMenuType{VersionItems: []versionMenuItems{}}
	if err := menu.getVersionMenuData(req, nil); err != nil {
		t.Fatalf("unexpected error populating version menu: %v", err)
	}

	gotVersions := []string{}
	for _, item := range menu.VersionItems {
		gotVersions = append(gotVersions, item.Version)
	}

	expectedVersions := []string{"latest", "v3", "v2", "v1"}
	if !reflect.DeepEqual(gotVersions, expectedVersions) {
		t.Fatalf("expected menu versions %v, got %v", expectedVersions, gotVersions)
	}
}
