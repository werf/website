package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Deprecated
func ssiHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "<p>SSI handler (%s).</p>", r.URL.Path[1:])
	_, _ = fmt.Fprintf(w, inspectRequestHTML(r))
}

// Get some status info
func statusHandler(w http.ResponseWriter, r *http.Request) {
	_ = r

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	_ = json.NewEncoder(w).Encode(
		ApiStatusResponseType{
			Status:         "ok",
			Msg:            "",
			RootVersion:    getRootReleaseVersion(),
			RootVersionURL: VersionToURL(getRootReleaseVersion()),
			Multiwerf:      []string{},
		})
}

// X-Redirect to the default root documentation version.
func rootVersionAliasHandler(w http.ResponseWriter, r *http.Request) {
	rawQuery := r.URL.RawQuery
	log.Debugln("Use handler - rootVersionAliasHandler")
	vars := mux.Vars(r)
	if err, version := resolveRootVersionAlias(vars["alias"]); err == nil {
		redirectURL := fmt.Sprintf("/docs/%v%v", VersionToURL(version), getDocPageURLRelative(r, true))
		if rawQuery != "" {
			redirectURL = fmt.Sprintf("%s?%s", redirectURL, rawQuery)
		}
		w.Header().Set("X-Accel-Redirect", redirectURL)
	} else {
		redirectURL := fmt.Sprintf("/docs/%s/", getCurrentDocsRoot())
		if rawQuery != "" {
			redirectURL = fmt.Sprintf("%s?%s", redirectURL, rawQuery)
		}
		http.Redirect(w, r, redirectURL, 302)
	}
}

// Healthcheck handler
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Get HTML content for /includes/topnav.html request
func topnavHandler(w http.ResponseWriter, r *http.Request) {
	versionMenu := versionMenuType{
		VersionItems:           []versionMenuItems{},
		HTMLContent:            "", // not used now
		CurrentVersion:         "",
		CurrentVersionURL:      "",
		CurrentPageURLRelative: "",
		MenuDocumentationLink:  "",
		AbsoluteVersion:        "",
	}

	_ = versionMenu.getVersionMenuData(r)

	tplPath := getRootFilesPath(r) + r.URL.Path
	tpl := template.Must(template.ParseFiles(tplPath))
	err := tpl.Execute(w, versionMenu)
	if err != nil {
		// Log error or maybe make some magic?
		log.Errorf("Internal Server Error (template error), %v ", err.Error())
		http.Error(w, "Internal Server Error (template error)", 500)
	}
}

func serveFilesHandler(fs http.FileSystem) http.Handler {
	fsh := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upath := r.URL.Path

		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
			r.URL.Path = upath
		}

		upath = path.Clean(upath)
		fileInfo, err := os.Stat(fmt.Sprintf("%v%s", fs, upath))
		if err != nil {
			if os.IsNotExist(err) {
				notFoundHandler(w, r)
				return
			}
		}

		if fileInfo.IsDir() {
			indexFile := filepath.Join(upath, "index.html")
			if _, err := os.Stat(fmt.Sprintf("%v%s", fs, indexFile)); err != nil {
				notFoundHandler(w, r)
				return
			}
		}
		fsh.ServeHTTP(w, r)
	})
}

// Redirect to root documentation if request not matches any location (override 404 response)
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	page404File, err := os.Open(getRootFilesPath(r) + "/404.html")
	if err != nil {
		// 404.html built-in stub
		log.Error("404.html file not found")
		http.Error(w, `<html lang="en"><head><meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge"><title>Page Not Found | werf</title><meta name="title" content="Page Not Found | werf">
</head>
<body>
		<h1 class="docs__title">Page Not Found</h1>
		<p>Sorry, the page you were looking for does not exist.<br>
Try searching for it or check the URL to see if it looks correct.</p>
</body></html>`, 404)
		return
	}
	defer page404File.Close()
	io.Copy(w, page404File)
	// w.Header().Set("X-Accel-Redirect", fmt.Sprintf("/404.html"))
}
