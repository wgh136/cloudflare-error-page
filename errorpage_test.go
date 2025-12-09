package cloudflare_error_page

import (
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	params := Params{
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
		"error_source":   "cloudflare",
		"what_happened":  "<p>Test error</p>",
		"what_can_i_do": "<p>Try again</p>",
	}

	html, err := Render(params, nil)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Check that HTML was generated
	if len(html) == 0 {
		t.Fatal("Generated HTML is empty")
	}

	// Check that key elements are present
	requiredStrings := []string{
		"<!DOCTYPE html>",
		"Internal server error",
		"Error code 500",
		"Test error",
		"Try again",
		"example.com",
		"cf-error-source",
	}

	for _, s := range requiredStrings {
		if !strings.Contains(html, s) {
			t.Errorf("Generated HTML does not contain expected string: %q", s)
		}
	}
}

func TestRenderWithCustomTitle(t *testing.T) {
	params := Params{
		"title":      "Custom Error",
		"error_code": 404,
	}

	html, err := Render(params, nil)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.Contains(html, "Custom Error") {
		t.Error("Generated HTML does not contain custom title")
	}

	if !strings.Contains(html, "404") {
		t.Error("Generated HTML does not contain custom error code")
	}
}

func TestRenderWithHTMLEscape(t *testing.T) {
	params := Params{
		"what_happened": "<script>alert('xss')</script>",
	}

	opts := &RenderOptions{
		AllowHTML: false,
		UseCDN:    true,
	}

	html, err := Render(params, opts)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Extract just the what_happened section
	start := strings.Index(html, "What happened?")
	var whatHappenedSection string
	if start != -1 {
		end := strings.Index(html[start:], "</div>")
		if end != -1 {
			whatHappenedSection = html[start : start+end]
		}
	}

	// User input should be escaped in the what_happened section
	if strings.Contains(whatHappenedSection, "<script>alert") {
		t.Error("HTML was not properly escaped in what_happened section")
	}

	if !strings.Contains(whatHappenedSection, "&lt;script&gt;") {
		t.Error("HTML escape did not produce expected output")
		t.Logf("what_happened section: %s", whatHappenedSection)
	}
}

func TestRenderWithCDNDisabled(t *testing.T) {
	params := Params{}

	opts := &RenderOptions{
		AllowHTML: true,
		UseCDN:    false,
	}

	html, err := Render(params, opts)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Should not use CDN URL
	if strings.Contains(html, "https://cloudflare.com/cdn-cgi") {
		t.Error("HTML contains CDN URL when UseCDN is false")
	}

	// Should use relative path
	if !strings.Contains(html, `href="/cdn-cgi/styles/main.css"`) {
		t.Error("HTML does not contain expected relative path")
	}
}

func TestGenerateRayID(t *testing.T) {
	id1 := generateRayID()
	id2 := generateRayID()

	if len(id1) != 16 {
		t.Errorf("Ray ID has wrong length: expected 16, got %d", len(id1))
	}

	if id1 == id2 {
		t.Error("generateRayID() generated identical IDs")
	}
}

func TestFillParams(t *testing.T) {
	params := Params{}
	fillParams(params)

	if _, ok := params["time"]; !ok {
		t.Error("fillParams did not set time")
	}

	if _, ok := params["ray_id"]; !ok {
		t.Error("fillParams did not set ray_id")
	}
}

func TestPrepareStatusInfo(t *testing.T) {
	params := Params{
		"browser_status": map[string]interface{}{
			"status":      "error",
			"status_text": "Out of Memory",
			"location":    "Local",
		},
	}

	info := prepareStatusInfo(params, "browser", "browser")

	if info.Icon != "browser" {
		t.Errorf("Expected icon 'browser', got %q", info.Icon)
	}

	if info.Status != "error" {
		t.Errorf("Expected status 'error', got %q", info.Status)
	}

	if info.StatusText != "Out of Memory" {
		t.Errorf("Expected status text 'Out of Memory', got %q", info.StatusText)
	}

	if info.Location != "Local" {
		t.Errorf("Expected location 'Local', got %q", info.Location)
	}

	if info.ErrorClass != "cf-error-source" {
		t.Errorf("Expected error class 'cf-error-source', got %q", info.ErrorClass)
	}
}

