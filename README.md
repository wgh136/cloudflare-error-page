# Cloudflare Error Page Generator

## What does this project do?

This project creates customized error pages that mimics the well-known Cloudflare error page. You can also embed it into your website.

## Online Editor

Here's an online editor to create customized error pages. Try it out [here](https://virt.moe/cloudflare-error-page/editor/).

![Editor](https://github.com/donlon/cloudflare-error-page/blob/images/editor.png?raw=true)

## Quickstart for Programmers

### Python

Install `cloudflare-error-page` with pip.

``` Bash
pip install git+https://github.com/donlon/cloudflare-error-page.git
```

Then you can generate an error page with the `render` function. ([example.py](examples/example.py))

``` Python
import webbrowser
from cloudflare_error_page import render as render_cf_error_page

# This function renders an error page based on the input parameters
error_page = render_cf_error_page({
    # Browser status is ok
    'browser_status': {
        "status": 'ok',
    },
    # Cloudflare status is error
    'cloudflare_status': {
        "status": 'error',
        "status_text": 'Error',
    },
    # Host status is also ok
    'host_status': {
        "status": 'ok',
        "location": 'example.com',
    },
    # can be 'browser', 'cloudflare', or 'host'
    'error_source': 'cloudflare',

    # Texts shown in the bottom of the page
    'what_happened': '<p>There is an internal server error on Cloudflare\'s network.</p>',
    'what_can_i_do': '<p>Please try again in a few minutes.</p>',
})

with open('error.html', 'w') as f:
    f.write(error_page)

webbrowser.open('error.html')
```

![Default error page](https://github.com/donlon/cloudflare-error-page/blob/images/default.png?raw=true)

You can also see live demo [here](https://virt.moe/cloudflare-error-page/examples/default).

A demo server using Flask is also available in [flask_demo.py](examples/flask_demo.py).

### Go

Install the package:

``` Bash
go get github.com/donlon/cloudflare-error-page
```

Then you can generate an error page with the `Render` function. ([example.go](examples_go/example/example.go))

``` Go
package main

import (
	"fmt"
	"os"

	errorpage "github.com/donlon/cloudflare-error-page"
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
		"error_source": "cloudflare",

		"what_happened": "<p>There is an internal server error on Cloudflare's network.</p>",
		"what_can_i_do": "<p>Please try again in a few minutes.</p>",
	}

	// Render the error page
	html, err := errorpage.Render(params, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Write to file
	os.WriteFile("error.html", []byte(html), 0644)
}
```

### Node.js

``` JavaScript
// Coming soon!
```

### PHP

``` PHP
/* Coming soon! */
```

## More Examples

### Catastrophic infrastructure failure

``` JavaScript
params = {
    "title": "Catastrophic infrastructure failure",
    "more_information": {
        "for": "no information",
    },
    "browser_status": {
        "status": "error",
        "status_text": "Out of Memory",
    },
    "cloudflare_status": {
        "status": "error",
        "location": "Everywhere",
        "status_text": "Error",
    },
    "host_status": {
        "status": "error",
        "location": "example.com",
        "status_text": "On Fire",
    },
    "error_source": "cloudflare",
    "what_happened": "<p>There is a catastrophic failure.</p>",
    "what_can_i_do": "<p>Please try again in a few years.</p>",
}
```

![Catastrophic infrastructure failure](https://github.com/donlon/cloudflare-error-page/blob/images/example.png?raw=true)

[Demo](https://virt.moe/cloudflare-error-page/examples/catastrophic)

### Web server is working

``` JavaScript
params = {
    "title": "Web server is working",
    "error_code": 200,
    "more_information": {
        "hidden": True,
    },
    "browser_status": {
        "status": "ok",
        "status_text": "Seems Working",
    },
    "cloudflare_status": {
        "status": "ok",
        "status_text": "Often Working",
    },
    "host_status": {
        "status": "ok",
        "location": "example.com",
        "status_text": "Almost Working",
    },
    "error_source": "host",
    "what_happened": "<p>This site is still working. And it looks great.</p>",
    "what_can_i_do": "<p>Visit the site before it crashes someday.</p>",
}
```

![Web server is working](https://github.com/donlon/cloudflare-error-page/blob/images/example2.png?raw=true)

[Demo](https://virt.moe/cloudflare-error-page/examples/working)

## See also

- [cloudflare-error-page-3th.pages.dev](https://cloudflare-error-page-3th.pages.dev/):

    Error page of every HTTP status code (reload to show random page).

- [oftx/cloudflare-error-page](https://github.com/oftx/cloudflare-error-page):

    React reimplementation of the original page, and can be deployed directly to Cloudflare Pages.


## Full Parameter Reference
``` JavaScript
{
    "html_title": "cloudflare.com | 500: Internal server error",
    "title": "Internal server error",
    "error_code": 500,
    "time": "2025-11-18 12:34:56 UTC",  // if not set, current UTC time is shown

    // Configuration for "Visit ... for more information" line
    "more_information": {
        "hidden": false,
        "text": "cloudflare.com", 
        "link": "https://www.cloudflare.com/",
        "for": "more information",
    },

    // Configuration for the Browser/Cloudflare/Host status
    "browser_status": {
        "status": "ok", // "ok" or "error"
        "location": "You",
        "name": "Browser",
        "status_text": "Working",
        "status_text_color": "#9bca3e",
    },
    "cloudflare_status": {
        "status": "error",
        "location": "Cloud",
        "name": "Cloudflare",
        "status_text": "Error",
        "status_text_color": "#bd2426",
    },
    "host_status": {
        "status": "ok",
        "location": "The Site",
        "name": "Host",
        "status_text": "Working",
        "status_text_color": "#9bca3e",
    },
    "error_source": "host", // Position of the error indicator, can be "browser", "cloudflare", or "host"

    "what_happened": "<p>There is an internal server error on Cloudflare's network.</p>",
    "what_can_i_do": "<p>Please try again in a few minutes.</p>",

    "ray_id": '0123456789abcdef',  // if not set, random hex string is shown
    "client_ip": '1.1.1.1',

    // Configuration for 'Performance & security by ...' in the footer
    "perf_sec_by": {
        "text": "Cloudflare",
        "link": "https://www.cloudflare.com/",
    },
}
```
