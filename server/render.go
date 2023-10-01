package server

import (
	"bytes"
	"io"
	"strings"
	"sync"

	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/thexyno/xynoblog/db"
	"github.com/thexyno/xynoblog/templates"
	xhtml "golang.org/x/net/html"

	log "github.com/sirupsen/logrus"
)

var renderCache sync.Map

type component struct {
	start func(w io.Writer, props map[string]string)
	end   func(w io.Writer, props map[string]string)
}

var components = map[string]component{
	"tangent": {
		start: templates.WriteTangentStart,
		end:   templates.WriteTangentEnd,
	},
	"box": {
		start: templates.WriteBoxStart,
		end:   templates.WriteBoxEnd,
	},
	"greybox": {
		start: templates.WriteGreyBoxStart,
		end:   templates.WriteGreyBoxEnd,
	},
}

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
	tkn := xhtml.NewTokenizer(bytes.NewReader(b))
	tkn.Next()
	token := tkn.Token()
	for k, v := range components {
		if strings.ToLower(token.Data) == strings.ToLower(k) {
			if token.Type == xhtml.StartTagToken {
				v.start(w, nil)
			} else {
				v.end(w, nil)
			}
			return ast.GoToNext, true
		}
	}

	return ast.GoToNext, false
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
	if cached, exists := renderCache.Load(idu); exists {
		return cached.([]byte)
	}

	parser := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs)
	opts := html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank | html.FootnoteReturnLinks, RenderNodeHook: renderHook}
	renderer := html.NewRenderer(opts)
	bytes := markdown.ToHTML([]byte(post.Content), parser, renderer)
	renderCache.Store(idu, bytes)
	return bytes
}

func RenderSimple(text []byte) []byte {
	parser := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.Footnotes | parser.LaxHTMLBlocks)
	opts := html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank | html.LazyLoadImages, RenderNodeHook: renderHook}
	renderer := html.NewRenderer(opts)
	return markdown.ToHTML(text, parser, renderer)
}
