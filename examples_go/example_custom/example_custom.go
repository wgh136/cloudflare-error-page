package main

import (
	"fmt"
	"os"

	errorpage "github.com/wgh136/cloudflare-error-page"
)

func main() {
	// Create a catastrophic failure error page
	params := errorpage.Params{
		"title": "Catastrophic infrastructure failure",
		"more_information": map[string]interface{}{
			"for": "no information",
		},
		"browser_status": map[string]interface{}{
			"status":      "error",
			"status_text": "Out of Memory",
		},
		"cloudflare_status": map[string]interface{}{
			"status":      "error",
			"location":    "Everywhere",
			"status_text": "Error",
		},
		"host_status": map[string]interface{}{
			"status":      "error",
			"location":    "example.com",
			"status_text": "On Fire",
		},
		"error_source":  "cloudflare",
		"what_happened": "<p>There is a catastrophic failure.</p>",
		"what_can_i_do": "<p>Please try again in a few years.</p>",
	}

	// Render with custom options
	opts := &errorpage.RenderOptions{
		AllowHTML: true,
		UseCDN:    true,
	}

	html, err := errorpage.Render(params, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering page: %v\n", err)
		os.Exit(1)
	}

	// Write to file
	err = os.WriteFile("error_catastrophic.html", []byte(html), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Error page generated successfully: error_catastrophic.html")
}
