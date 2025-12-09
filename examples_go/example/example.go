package main

import (
	"fmt"
	"os"

	errorpage "github.com/wgh136/cloudflare-error-page"
)

func main() {
	// Create error page parameters
	params := errorpage.Params{
		"browser_status": map[string]interface{}{
			"status": "ok",
		},
		"cloudflare_status": map[string]interface{}{
			"status":      "error",
			"status_text": "Error",
		},
		"host_status": map[string]interface{}{
			"status":   "ok",
			"location": "example.com",
		},
		"error_source": "cloudflare", // 'browser', 'cloudflare', or 'host'

		"what_happened": "<p>There is an internal server error on Cloudflare's network.</p>",
		"what_can_i_do": "<p>Please try again in a few minutes.</p>",
	}

	// Render the error page
	html, err := errorpage.Render(params, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering page: %v\n", err)
		os.Exit(1)
	}

	// Write to file
	err = os.WriteFile("error.html", []byte(html), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Error page generated successfully: error.html")
}
