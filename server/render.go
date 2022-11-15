package server

import (
	"io"
	"regexp"

	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/thexyno/xynoblog/db"

	log "github.com/sirupsen/logrus"
)

var renderCache = make(map[db.IdUpdated]([]byte))

func renderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	_, ok := node.(*ast.CodeBlock)
	if ok {
		return renderCodeBlock(w, node, entering)
	}
	_, ok = node.(*ast.HTMLSpan)
	if ok {
		return renderHTMLBlock(w, node, entering)
	}
	_, ok = node.(*ast.Image)
	if ok {
		return renderIMGBlock(w, node, entering)
	}
	return ast.GoToNext, false
}

func renderIMGBlock(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	n, ok := node.(*ast.Image)

	if !ok {
		return ast.GoToNext, false
	}
	tmp := append([]byte("/media/"), n.Destination...)
	n.Destination = tmp
	return ast.GoToNext, false
}

func renderHTMLBlock(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	n, ok := node.(*ast.HTMLSpan)

	if !ok {
		return ast.GoToNext, false
	}

	b := n.Literal
	r, err := regexp.Compile("< *(/?)(.*) *>")
	if err != nil {
		log.Error(err)
		return ast.GoToNext, false
	}
	b2 := r.FindSubmatch(b)
	closing := len(b2[1]) > 0
	tag := b2[2]
	allowedTags := [][]byte{[]byte("tangent"), []byte("box")}
	if contains(allowedTags, tag) {
		if !closing {
			w.Write([]byte("<div class=\""))
			w.Write(tag)
			w.Write([]byte("\">"))
		} else {
			w.Write([]byte("</div>"))
		}
	}

	return ast.GoToNext, true
}

func equals(arr []byte, x []byte) bool {
	if len(arr) != len(x) {
		return false
	}
	for i := range arr {
		if x[i] != arr[i] {
			return false
		}
	}
	return true
}

func contains(arr [][]byte, x []byte) bool {
	for _, v := range arr {
		if equals(v, x) {
			return true
		}
	}
	return false
}

func renderCodeBlock(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	n, ok := node.(*ast.CodeBlock)

	if !ok {
		return ast.GoToNext, false
	}

	var language string
	if n.Info == nil {
		language = "bash"
	} else {
		language = string(n.Info)
	}

	lexer := lexers.Get(language)
	formatter := chromahtml.New(chromahtml.WithClasses(true))
	tokens, err := lexer.Tokenise(nil, string(n.Literal))
	style := styles.Get("gruvbox")
	if style == nil {
		style = styles.Fallback
	}
	if err != nil {
		log.Error(err)
	}
	formatter.Format(w, style, tokens)

	return ast.GoToNext, true
}

func Render(post db.Post) []byte {
	idu := post.ToIdUpdated()
	if cached, exists := renderCache[idu]; exists {
		return cached
	}

	parser := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs)
	opts := html.RendererOptions{Flags: html.CommonFlags, RenderNodeHook: renderHook}
	renderer := html.NewRenderer(opts)
	bytes := markdown.ToHTML([]byte(post.Content), parser, renderer)
	renderCache[idu] = bytes
	return bytes
}

func RenderSimple(text []byte) []byte {
	parser := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.Footnotes | parser.LaxHTMLBlocks)
	opts := html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank | html.LazyLoadImages, RenderNodeHook: renderHook}
	renderer := html.NewRenderer(opts)
	return markdown.ToHTML(text, parser, renderer)
}
