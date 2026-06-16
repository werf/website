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

func TestLegacyDocsVersionHandler(t *testing.T) {
	t.Setenv(currentDocsRootEnv, "v2")
	t.Setenv(supportedDocsRootsEnv, "v2,v1.2")

	r := newRouter()

	cases := []struct {
		name       string
		path       string
		wantStatus int
		wantLoc    string
	}{
		{name: "legacy v2 patch", path: "/docs/v2.3.0/usage/", wantStatus: http.StatusFound, wantLoc: "/docs/v2/usage/"},
		{name: "legacy v1.2 patch", path: "/docs/v1.2.9-plus-fix1/reference/", wantStatus: http.StatusFound, wantLoc: "/docs/v1.2/reference/"},
		{name: "legacy v2 alias", path: "/docs/v2-stable/reference/", wantStatus: http.StatusFound, wantLoc: "/docs/v2/reference/"},
		{name: "legacy v1.2 alias", path: "/docs/v1.2-beta/reference/", wantStatus: http.StatusFound, wantLoc: "/docs/v1.2/reference/"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("unexpected status: got %d want %d", rec.Code, tc.wantStatus)
			}
			if rec.Header().Get("Location") != tc.wantLoc {
				t.Fatalf("unexpected location: got %q want %q", rec.Header().Get("Location"), tc.wantLoc)
			}
		})
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
