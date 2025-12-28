package web

import (
	"embed"
	"io/fs"
)

//go:embed all:static/*
var staticFS embed.FS

var Static, _ = fs.Sub(staticFS, "static")

