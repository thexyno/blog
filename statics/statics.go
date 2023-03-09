package statics

import (
	"embed"
	"io/fs"
)

//go:embed dist
var CSSDir embed.FS

//go:embed data
var dataDir embed.FS
var DataDir, _ = fs.Sub(dataDir, "data")
