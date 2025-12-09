package cloudflare_error_page

import (
	"bytes"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"fmt"
	"html/template"
	"path/filepath"
	"time"
)

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed resources
var resourcesFS embed.FS

var defaultTemplate *template.Template

func init() {
	// Parse the default template
	tmplContent, err := templatesFS.ReadFile("templates/error.html")
	if err != nil {
		panic("failed to read default template: " + err.Error())
	}

	// Create template with custom functions
	funcMap := template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	defaultTemplate = template.Must(template.New("error.html").Funcs(funcMap).Parse(string(tmplContent)))
}

// GetResourcesFolder returns the embedded resources filesystem
// This can be used to serve static files (CSS, images)
func GetResourcesFolder() embed.FS {
	return resourcesFS
}

// GetResourcePath returns the path to a resource file
func GetResourcePath(filename string) string {
	return filepath.Join("resources", filename)
}

// Params represents the parameters for rendering the error page
type Params map[string]interface{}

type statusInfo struct {
	Icon            string
	DefaultLocation string
	DefaultName     string
	Location        string
	Name            string
	Status          string
	StatusText      string
	StatusTextColor string
	ErrorClass      string
}

// fillParams fills in default values for missing parameters
func fillParams(params Params) {
	if _, ok := params["time"]; !ok {
		params["time"] = time.Now().UTC().Format("2006-01-02 15:04:05 UTC")
	}
	if _, ok := params["ray_id"]; !ok {
		params["ray_id"] = generateRayID()
	}
}

func getMapValue(m map[string]interface{}, key string, defaultVal string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultVal
}

func getBoolValue(m map[string]interface{}, key string, defaultVal bool) bool {
	if val, ok := m[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return defaultVal
}

// prepareStatusInfo prepares status information for the template
func prepareStatusInfo(params Params, itemID string, errorSource string) statusInfo {
	info := statusInfo{}
	
	// Set icon and defaults based on item ID
	switch itemID {
	case "browser":
		info.Icon = "browser"
		info.DefaultLocation = "You"
		info.DefaultName = "Browser"
	case "cloudflare":
		info.Icon = "cloud"
		info.DefaultLocation = "San Francisco"
		info.DefaultName = "Cloudflare"
	case "host":
		info.Icon = "server"
		info.DefaultLocation = "Website"
		info.DefaultName = "Host"
	}
	
	// Get item status from params
	itemKey := itemID + "_status"
	var item map[string]interface{}
	if val, ok := params[itemKey]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			item = m
		}
	}
	if item == nil {
		item = make(map[string]interface{})
	}
	
	// Set status info
	info.Status = getMapValue(item, "status", "ok")
	info.Location = getMapValue(item, "location", info.DefaultLocation)
	info.Name = getMapValue(item, "name", info.DefaultName)
	
	// Set status text
	if val := getMapValue(item, "status_text", ""); val != "" {
		info.StatusText = val
	} else if info.Status == "ok" {
		info.StatusText = "Working"
	} else {
		info.StatusText = "Error"
	}
	
	// Set status text color
	if val := getMapValue(item, "status_text_color", ""); val != "" {
		info.StatusTextColor = val
	} else if info.Status == "ok" {
		info.StatusTextColor = "#9bca3e"
	} else if info.Status == "error" {
		info.StatusTextColor = "#bd2426"
	}
	
	// Set error class
	if errorSource == itemID {
		info.ErrorClass = "cf-error-source"
	}
	
	return info
}

// generateRayID generates a random 16-character hex string
func generateRayID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "0000000000000000"
	}
	return hex.EncodeToString(b)
}

// RenderOptions contains options for rendering the error page
type RenderOptions struct {
	AllowHTML bool
	UseCDN    bool
}

// Render renders a customized Cloudflare error page
func Render(params Params, opts *RenderOptions) (string, error) {
	if opts == nil {
		opts = &RenderOptions{
			AllowHTML: true,
			UseCDN:    true,
		}
	}

	// Make a copy of params to avoid modifying the original
	paramsCopy := make(Params)
	for k, v := range params {
		paramsCopy[k] = v
	}

	fillParams(paramsCopy)

	// Note: We don't manually escape HTML here when AllowHTML is false.
	// Instead, we rely on the template engine's automatic escaping when
	// we don't use the 'safe' function in the template.
	allowHTMLContent := opts.AllowHTML

	// Set defaults for params
	if _, ok := paramsCopy["error_code"]; !ok {
		paramsCopy["error_code"] = 500
	}
	if _, ok := paramsCopy["title"]; !ok {
		paramsCopy["title"] = "Internal server error"
	}
	if _, ok := paramsCopy["html_title"]; !ok {
		paramsCopy["html_title"] = fmt.Sprintf("%d: %s", paramsCopy["error_code"], paramsCopy["title"])
	}
	if _, ok := paramsCopy["what_happened"]; !ok {
		paramsCopy["what_happened"] = "<p>There is an internal server error on Cloudflare's network.</p>"
	}
	if _, ok := paramsCopy["what_can_i_do"]; !ok {
		paramsCopy["what_can_i_do"] = "<p>Please try again in a few minutes.</p>"
	}
	if _, ok := paramsCopy["client_ip"]; !ok {
		paramsCopy["client_ip"] = "1.1.1.1"
	}

	// Prepare more_information
	moreInfo := make(map[string]interface{})
	if val, ok := paramsCopy["more_information"]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			moreInfo = m
		}
	}
	if _, ok := moreInfo["hidden"]; !ok {
		moreInfo["hidden"] = false
	}
	if _, ok := moreInfo["link"]; !ok {
		moreInfo["link"] = "https://www.cloudflare.com/"
	}
	if _, ok := moreInfo["text"]; !ok {
		moreInfo["text"] = "cloudflare.com"
	}
	if _, ok := moreInfo["for"]; !ok {
		moreInfo["for"] = "more information"
	}
	paramsCopy["more_information"] = moreInfo

	// Prepare perf_sec_by
	perfSecBy := make(map[string]interface{})
	if val, ok := paramsCopy["perf_sec_by"]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			perfSecBy = m
		}
	}
	if _, ok := perfSecBy["link"]; !ok {
		perfSecBy["link"] = "https://www.cloudflare.com/"
	}
	if _, ok := perfSecBy["text"]; !ok {
		perfSecBy["text"] = "Cloudflare"
	}
	paramsCopy["perf_sec_by"] = perfSecBy

	// Prepare creator_info
	creatorInfo := make(map[string]interface{})
	if val, ok := paramsCopy["creator_info"]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			creatorInfo = m
		}
	}
	if _, ok := creatorInfo["hidden"]; !ok {
		creatorInfo["hidden"] = true
	}
	paramsCopy["creator_info"] = creatorInfo

	// Get error source
	errorSource := ""
	if val, ok := paramsCopy["error_source"]; ok {
		if str, ok := val.(string); ok {
			errorSource = str
		}
	}

	// Prepare status info for each section
	browserStatus := prepareStatusInfo(paramsCopy, "browser", errorSource)
	cloudflareStatus := prepareStatusInfo(paramsCopy, "cloudflare", errorSource)
	hostStatus := prepareStatusInfo(paramsCopy, "host", errorSource)

	// Prepare template data
	data := map[string]interface{}{
		"params":             paramsCopy,
		"resources_use_cdn":  opts.UseCDN,
		"resources_cdn":      "",
		"browser_status":     browserStatus,
		"cloudflare_status":  cloudflareStatus,
		"host_status":        hostStatus,
		"allow_html":         allowHTMLContent,
	}

	if opts.UseCDN {
		data["resources_cdn"] = "https://cloudflare.com"
	}

	// Render the template
	var buf bytes.Buffer
	err := defaultTemplate.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	return buf.String(), nil
}

