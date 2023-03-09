package statics

import (
	"embed"
	"io/fs"

	"github.com/thexyno/xynoblog/hashFS"
)

//go:embed dist
var cssDir embed.FS
var CSSDir, _ = fs.Sub(cssDir, "dist")
var CSSHashDir = hashFS.GenHashFS(CSSDir)

//go:embed data
var dataDir embed.FS
var DataDir, _ = fs.Sub(dataDir, "data")
