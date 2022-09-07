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
	log "github.com/sirupsen/logrus"
)

func renderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	_, ok := node.(*ast.CodeBlock)
	if ok {
		return renderCodeBlock(w, node, entering)
	}
	_, ok = node.(*ast.HTMLSpan)
	if ok {
		return renderHTMLBlock(w, node, entering)
	}
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

func Render(text []byte) []byte {
	parser := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs)
	opts := html.RendererOptions{Flags: html.CommonFlags, RenderNodeHook: renderHook}
	renderer := html.NewRenderer(opts)
	return markdown.ToHTML([]byte(text), parser, renderer)
}
