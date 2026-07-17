// Package webui embeds the built frontend (web/) SPA and serves it, falling
// back to index.html for client-side routes so react-router splat paths work
// on a full page load.
package webui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist
var distFS embed.FS

func staticFS() fs.FS {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	return sub
}

// Handler serves the embedded frontend for any non-/api/ GET request and
// delegates everything under /api/ to apiHandler.
func Handler(apiHandler http.Handler) http.Handler {
	fileServer := http.FileServer(http.FS(staticFS()))
	static := staticFS()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) >= 5 && r.URL.Path[:5] == "/api/" {
			apiHandler.ServeHTTP(w, r)
			return
		}

		if _, err := fs.Stat(static, trimLeadingSlash(r.URL.Path)); err != nil {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
}

func trimLeadingSlash(p string) string {
	if len(p) > 0 && p[0] == '/' {
		p = p[1:]
	}
	if p == "" {
		return "."
	}
	return p
}
