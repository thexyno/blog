package server

import (
	"io"

	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	log "github.com/sirupsen/logrus"
)

func renderCodeBlock(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	n, ok := node.(*ast.CodeBlock)
	if !ok {
		return ast.GoToNext, false
	}

	lexer := lexers.Get(string(n.Info))
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
	opts := html.RendererOptions{Flags: html.CommonFlags, RenderNodeHook: renderCodeBlock}
	renderer := html.NewRenderer(opts)
	return markdown.ToHTML([]byte(text), parser, renderer)
}
