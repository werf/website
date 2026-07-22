package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCanonicalDocsVersion(t *testing.T) {
	cases := []struct {
		name           string
		currentRoot    string
		supportedRoots string
		in             string
		want           string
	}{
		{name: "empty", currentRoot: "v2", supportedRoots: "v2,v1.2", in: "", want: "v2"},
		{name: "latest", currentRoot: "v2", supportedRoots: "v2,v1.2", in: "latest", want: "v2"},
		{name: "pr", currentRoot: "v2", supportedRoots: "v2,v1.2", in: "pr-123", want: "pr-123"},
		{name: "v2 patch", currentRoot: "v2", supportedRoots: "v2,v1.2", in: "v2.10.3", want: "v2"},
		{name: "v1.2 patch", currentRoot: "v2", supportedRoots: "v2,v1.2", in: "v1.2.4-plus-fix1", want: "v1.2"},
		{name: "v3 patch from supported roots", currentRoot: "v2", supportedRoots: "v2,v1.2,v3", in: "v3.1.4", want: "v3"},
		{name: "v3 root from supported roots", currentRoot: "v2", supportedRoots: "v2,v1.2,v3", in: "v3", want: "v3"},
		{name: "unknown version fallback", currentRoot: "v2", supportedRoots: "v2,v1.2,v3", in: "v9.1", want: "v2"},
		{name: "prefer longest root match", currentRoot: "v1", supportedRoots: "v1,v1.2,v2", in: "v1.2.5", want: "v1.2"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(currentDocsRootEnv, tc.currentRoot)
			t.Setenv(supportedDocsRootsEnv, tc.supportedRoots)
			if got := getCanonicalDocsVersion(tc.in); got != tc.want {
				t.Fatalf("unexpected canonical version: got %q want %q", got, tc.want)
			}
		})
	}
}

func TestResolveRootVersionAlias(t *testing.T) {
	t.Setenv(currentDocsRootEnv, "v2")

	err, version := resolveRootVersionAlias("latest")
	if err != nil {
		t.Fatalf("unexpected error for latest alias: %v", err)
	}
	if version != "v2" {
		t.Fatalf("unexpected version for latest alias: got %q want %q", version, "v2")
	}

	err, _ = resolveRootVersionAlias("v2")
	if err == nil {
		t.Fatal("expected unsupported alias error for canonical version")
	}
}

func TestLatestAliasHandler(t *testing.T) {
	t.Setenv(currentDocsRootEnv, "v2")
	t.Setenv(supportedDocsRootsEnv, "v2,v1.2")
	t.Setenv(latestDocsAliasEnabledEnv, "true")

	r := newRouter()

	req, err := http.NewRequest(http.MethodGet, "/docs/latest/reference/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RequestURI = "/docs/latest/reference/"
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("X-Accel-Redirect"); got != "/docs/v2/reference/" {
		t.Fatalf("unexpected redirect target: got %q want %q", got, "/docs/v2/reference/")
	}

	reqQuery, err := http.NewRequest(http.MethodGet, "/docs/latest/reference/?foo=bar", nil)
	if err != nil {
		t.Fatal(err)
	}
	reqQuery.RequestURI = "/docs/latest/reference/?foo=bar"
	recQuery := httptest.NewRecorder()
	r.ServeHTTP(recQuery, reqQuery)

	if recQuery.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", recQuery.Code, http.StatusOK)
	}
	if got := recQuery.Header().Get("X-Accel-Redirect"); got != "/docs/v2/reference/?foo=bar" {
		t.Fatalf("unexpected redirect target: got %q want %q", got, "/docs/v2/reference/?foo=bar")
	}
}

func TestStaticFileServer(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/")
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Status should be 404, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("Wrong content type, expected %s, got %s", expectedContentType, contentType)
	}
}
