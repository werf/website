package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/werf/common-go/pkg/graceful"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()

	staticFileDirectoryMain := http.Dir("./root/en")
	staticFileDirectoryRu := http.Dir("./root/ru")

	var ruHostMatch mux.MatcherFunc = func(r *http.Request, rm *mux.RouteMatch) bool {
		result := false
		result, _ = regexp.MatchString(`^ru\.(localhost|.*(.+\.flant\.com|werf\.io))$`, r.Host)
		return result
	}

	r.PathPrefix("/status").HandlerFunc(statusHandler)
	r.PathPrefix("/backend/").HandlerFunc(ssiHandler)
	r.PathPrefix("/docs/v{group:[0-9]+.[0-9]+}-{channel:alpha|beta|ea|stable|rock-solid}").HandlerFunc(groupChannelHandler)
	r.PathPrefix("/docs/v{version:[0-9]+.[0-9]+.[0-9]+[^/]*}").HandlerFunc(unknownVersionHandler)
	r.PathPrefix("/docs/v{group:[0-9]+}/").HandlerFunc(groupHandler)
	r.PathPrefix("/docs/v{group:[0-9]+.[0-9]+}/").HandlerFunc(groupHandler)
	r.PathPrefix("/docs/{group:latest}/").HandlerFunc(groupHandler)
	r.PathPrefix("/health").HandlerFunc(healthCheckHandler)
	r.Path("/includes/topnav.html").HandlerFunc(topnavHandler)
	r.Path("/includes/version-menu.html").HandlerFunc(topnavHandler)
	r.Path("/includes/group-menu.html").HandlerFunc(groupMenuHandler)
	r.Path("/includes/group-menu-v2.html").HandlerFunc(groupMenuHandler)
	r.Path("/includes/channel-menu.html").HandlerFunc(channelMenuHandler)
	r.Path("/includes/channel-menu-v2.html").HandlerFunc(channelMenuHandler)
	r.Path("/404.html").HandlerFunc(notFoundHandler)
	// Ru static
	r.MatcherFunc(ruHostMatch).Handler(serveFilesHandler(staticFileDirectoryRu))
	// Other (En) static
	r.PathPrefix("/").Handler(serveFilesHandler(staticFileDirectoryMain))

	r.Use(LoggingMiddleware)

	r.NotFoundHandler = r.NewRoute().HandlerFunc(notFoundHandler).GetHandler()

	return r
}

func main() {
	var wait time.Duration
	logLevel := os.Getenv("LOG_LEVEL")
	if strings.ToLower(logLevel) == "debug" {
		log.SetLevel(log.DebugLevel)
	} else if strings.ToLower(logLevel) == "trace" {
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	log.Infoln(fmt.Sprintf("Started with LOG_LEVEL %s", logLevel))
	r := newRouter()

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	terminationCtx := graceful.WithTermination(context.Background())

	defer graceful.Shutdown(terminationCtx, func(err error, exitCode int) {
		shutdownCtx := context.WithoutCancel(terminationCtx)

		if wait > 0 {
			timeoutCtx, cancel := context.WithTimeout(shutdownCtx, wait)
			defer cancel()
			shutdownCtx = timeoutCtx
		}

		if shutdownErr := srv.Shutdown(shutdownCtx); shutdownErr != nil {
			err = errors.Join(err, fmt.Errorf("server shutdown failed: %w", shutdownErr))
		}

		if err != nil {
			log.Errorln("Shutdown errors:", err)
			os.Exit(exitCode)
		} else {
			log.Infoln("Server gracefully stopped")
		}
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorln("Server error:", err)
			graceful.Terminate(terminationCtx, err, 1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case <-c:
		graceful.Terminate(terminationCtx, nil, 0)
	case <-terminationCtx.Done():
	}
}
