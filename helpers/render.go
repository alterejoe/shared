package helpers

import (
	"net/http"

	"github.com/a-h/templ"
)

// IsHTMX checks if the request is from HTMX
func IsHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// RenderContent renders either just the content (HTMX) or full page (normal request)
func RenderContent(w http.ResponseWriter, r *http.Request, content templ.Component, fullPage func(templ.Component) templ.Component) error {
	w.Header().Set("Content-Type", "text/html")

	if IsHTMX(r) {
		// HTMX request - render just the content
		return content.Render(r.Context(), w)
	}

	// Normal request - wrap in full page skeleton
	return fullPage(content).Render(r.Context(), w)
}
