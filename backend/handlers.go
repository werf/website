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
	_, _ = fmt.Fprint(w, inspectRequestHTML(r))
}

// Get some status info
func statusHandler(w http.ResponseWriter, _ *http.Request) {
	var msg []string
	status := "ok"

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	_ = json.NewEncoder(w).Encode(
		ApiStatusResponseType{
			Status:         status,
			Msg:            strings.Join(msg, " "),
			RootVersion:    getRootReleaseVersion(),
			RootVersionURL: VersionToURL(getRootReleaseVersion()),
			Multiwerf:      ReleasesStatus.Releases,
		})
}

// X-Redirect to the stablest documentation version for specific group
func groupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debugln("Use handler - groupHandler")

	if err, version := getVersionFromGroup(&ReleasesStatus, vars["group"]); err == nil {
		redirectToDocsRoot(w, r, version)
		return
	}

	redirectToDocsRoot(w, r, getCurrentDocsMajorRoot())
}

// Handler for compatibility URLs which should resolve to supported docs major roots.
func legacyDocsVersionHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugln("Use handler - legacyDocsVersionHandler")
	vars := mux.Vars(r)
	redirectToDocsRoot(w, r, getCanonicalDocsVersion(fmt.Sprintf("v%s", vars["legacy"])))
}

// Handles request to legacy channel URLs like /docs/v2-stable/ and redirects them to the corresponding major root.
func groupChannelHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugln("Use handler - groupChannelHandler")
	vars := mux.Vars(r)
	redirectToDocsRoot(w, r, getCanonicalDocsVersion(fmt.Sprintf("v%s-%s", vars["group"], vars["channel"])))
}

func redirectToDocsRoot(w http.ResponseWriter, r *http.Request, root string) {
	redirectURL := fmt.Sprintf("/docs/%s%s", VersionToURL(root), getDocPageURLRelative(r, true))
	if rawQuery := r.URL.RawQuery; rawQuery != "" {
		redirectURL = fmt.Sprintf("%s?%s", redirectURL, rawQuery)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// Healthcheck handler
func healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Get HTML content for /includes/topnav.html request
func topnavHandler(w http.ResponseWriter, r *http.Request) {
	versionMenu := versionMenuType{
		VersionItems:           []versionMenuItems{},
		HTMLContent:            "", // not used now
		CurrentGroup:           "", // not used now
		CurrentChannel:         "",
		CurrentVersion:         "",
		CurrentVersionURL:      "",
		CurrentPageURLRelative: "",
		MenuDocumentationLink:  "",
		AbsoluteVersion:        "",
	}

	_ = versionMenu.getVersionMenuData(r, nil)

	tplPath := getRootFilesPath(r) + r.URL.Path
	tpl := template.Must(template.ParseFiles(tplPath))
	err := tpl.Execute(w, versionMenu)
	if err != nil {
		// Log error or maybe make some magic?
		log.Errorf("Internal Server Error (template error), %v ", err.Error())
		http.Error(w, "Internal Server Error (template error)", 500)
	}
}

func groupMenuHandler(w http.ResponseWriter, r *http.Request) {
	versionMenu := versionMenuType{
		VersionItems:           []versionMenuItems{},
		HTMLContent:            "", // not used now
		CurrentGroup:           "", // not used now
		CurrentChannel:         "",
		CurrentVersion:         "",
		CurrentVersionURL:      "",
		CurrentPageURLRelative: "",
		MenuDocumentationLink:  "",
	}

	_ = versionMenu.getGroupMenuData(r, nil)

	tplPath := getRootFilesPath(r) + r.RequestURI
	tpl := template.Must(template.ParseFiles(tplPath))
	err := tpl.Execute(w, versionMenu)
	if err != nil {
		// Log error or maybe make some magic?
		log.Errorf("Internal Server Error (template error), %v ", err.Error())
		http.Error(w, "Internal Server Error (template error)", 500)
	}
}

func channelMenuHandler(w http.ResponseWriter, r *http.Request) {
	versionMenu := versionMenuType{
		VersionItems:           []versionMenuItems{},
		HTMLContent:            "", // not used now
		CurrentGroup:           "", // not used now
		CurrentChannel:         "",
		CurrentVersion:         "",
		CurrentVersionURL:      "",
		CurrentPageURLRelative: "",
		MenuDocumentationLink:  "",
	}

	_ = versionMenu.getChannelMenuData(r, nil)

	tplPath := getRootFilesPath(r) + r.RequestURI
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
	defer func() {
		if closeErr := page404File.Close(); closeErr != nil {
			log.Error(closeErr)
		}
	}()
	if _, err := io.Copy(w, page404File); err != nil {
		log.Error(err)
	}
	// w.Header().Set("X-Accel-Redirect", fmt.Sprintf("/404.html"))
}
