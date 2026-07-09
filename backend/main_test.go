package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCanonicalDocsVersion(t *testing.T) {
	t.Setenv(currentDocsRootEnv, "v2")
	t.Setenv(supportedDocsRootsEnv, "v2,v1.2")

	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: "v2"},
		{name: "latest", in: "latest", want: "v2"},
		{name: "pr", in: "pr-123", want: "pr-123"},
		{name: "v2 patch", in: "v2.10.3", want: "v2"},
		{name: "v1.2 patch", in: "v1.2.4-plus-fix1", want: "v1.2"},
		{name: "v1.1 fallback", in: "v1.1.9", want: "v2"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
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
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusOK)
	}
	if rec.Header().Get("X-Accel-Redirect") != "/docs/v2/" {
		t.Fatalf("unexpected redirect target: got %q want %q", rec.Header().Get("X-Accel-Redirect"), "/docs/v2/")
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
