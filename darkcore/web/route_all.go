package web

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/darklab8/fl-darkcore/darkcore/core_types"
	"github.com/darklab8/fl-darkcore/darkcore/settings/logus"
	"github.com/darklab8/fl-darkcore/darkcore/web/registry"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/utils_types"
)

const UrlStatic core_types.Url = "/"

func (w *Web) NewEndpointStatic() *registry.Endpoint {
	return &registry.Endpoint{
		Url: UrlStatic,
		Handler: func(resp http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodOptions:
			case http.MethodGet:

				requested := req.URL.Path[1:]

				requested = strings.ReplaceAll(requested, "/", PATH_SEPARATOR)

				for _, acceptor := range w.site_root_acceptors {
					requested = strings.ReplaceAll(requested, acceptor, w.site_root[1:])
				}

				if requested == "" {
					requested = "index.html"
				}

				logger := logus.Log.WithFields(
					typelog.String("requested_path", requested),
					typelog.Int("files_count", len(w.filesystem.Files)),
				)

				logger.Info("having get request")

				content, ok := w.filesystem.Files[utils_types.FilePath(requested)]

				if strings.Contains(requested, ".css") {
					resp.Header().Set("Content-Type", "text/css; charset=utf-8")
				} else if strings.Contains(requested, ".html") {
					resp.Header().Set("Content-Type", "text/html; charset=utf-8")
				} else if strings.Contains(requested, ".js") {
					resp.Header().Set("Content-Type", "application/javascript; charset=utf-8")
				}

				if ok {
					fmt.Fprint(resp, string(content.Render()))
				} else {
					resp.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(resp, "content is not found at %s!, %q", req.URL, html.EscapeString(requested))
					logus.Log.Error("content is not found")
				}

			default:
				http.Error(resp, "Method not allowed", http.StatusMethodNotAllowed)
			}
		},
	}
}
