package main

import (
	"fmt"
	"log"
	"net/http"

	errorpage "github.com/wgh136/cloudflare-error-page"
)

func errorHandler(w http.ResponseWriter, r *http.Request) {
	params := errorpage.Params{
		"error_code": 500,
		"title":      "Internal server error",
		"browser_status": map[string]interface{}{
			"status": "ok",
		},
		"cloudflare_status": map[string]interface{}{
			"status":      "error",
			"status_text": "Error",
		},
		"host_status": map[string]interface{}{
			"status":   "ok",
			"location": r.Host,
		},
		"error_source": "cloudflare",

		"what_happened": "<p>There is an internal server error on Cloudflare's network.</p>",
		"what_can_i_do": "<p>Please try again in a few minutes.</p>",
		"client_ip":     r.RemoteAddr,
	}

	html, err := errorpage.Render(params, nil)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(html))
}

func main() {
	http.HandleFunc("/error", errorHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Visit <a href=\"/error\">/error</a> to see the error page")
	})

	addr := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", addr)
	fmt.Printf("Visit http://localhost%s/error to see the error page\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
