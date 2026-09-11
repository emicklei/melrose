// Package web contains the embedded assets of the Melrōse web frontend.
package web

import "embed"

//go:embed index.html app.js style.css
var Assets embed.FS
