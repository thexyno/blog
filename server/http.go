//go:generate go get -d github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=../templates
package server

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/thexyno/xynoblog/db"
	"github.com/thexyno/xynoblog/templates"
	"io/ioutil"
	"net/http"
	"strings"
)

func renderError(c *gin.Context, err error) {
	log.WithField("request", c.FullPath()).Error(err)
	var p *templates.ErrorPage
	if err == db.ErrNotFound {
		c.Status(404)
		p = &templates.ErrorPage{
			Message: "Not Found",
		}
	} else {
		c.Status(500)
		p = &templates.ErrorPage{
			Message: "Internal Server Error",
		}
	}
	templates.WritePageTemplate(c.Writer, p)
}

func renderPosts(db db.DbConn) func(*gin.Context) {
	return func(c *gin.Context) {
		posts, err := db.ShortPosts(-1, 0)
		if err != nil {
			renderError(c, err)
			return
		}
		if len(posts) > 0 {
			c.Header("Last-Modified", posts[0].Created.Format(http.TimeFormat))
		}
		p := &templates.PostsPage{
			Posts: posts,
		}
		templates.WritePageTemplate(c.Writer, p)
	}
}
func renderIndex(db db.DbConn) func(*gin.Context) {
	return func(c *gin.Context) {
		posts, err := db.ShortPosts(5, 0)
		if err != nil {
			renderError(c, err)
			return
		}
		if len(posts) > 0 {
			c.Header("Last-Modified", posts[0].Created.Format(http.TimeFormat))
		}
		p := &templates.IndexPage{
			Posts: posts,
		}
		templates.WritePageTemplate(c.Writer, p)
	}
}
func renderPost(db db.DbConn) func(*gin.Context) {
	return func(c *gin.Context) {
		var id struct {
			Id string `uri:"id" binding:"required"`
		}
		if err := c.ShouldBindUri(&id); err != nil {
			renderError(c, err)
			return
		}
		post, err := db.Post(id.Id)
		if err != nil {
			renderError(c, err)
			return
		}
		if !gin.IsDebugging() {
			c.Header("Last-Modified", post.Updated.Format(http.TimeFormat))
		}
		rendered := Render(post)
		p := &templates.PostPage{
			Post:            post,
			RenderedContent: rendered,
		}
		templates.WritePageTemplate(c.Writer, p)
	}
}

func renderSimpleMarkdownPage(title []byte, content []byte, index bool) func(*gin.Context) {
	return func(c *gin.Context) {
		rendered := RenderSimple(content)
		p := &templates.SimpleMdPage{
			PageTitle:       title,
			RenderedContent: rendered,
		}
		if !index {
			c.Header("X-Robots-Tag", "noindex")
		}
		templates.WritePageTemplate(c.Writer, p)
	}
}

func Mux(db db.DbConn, fontdir string, cssdir string, staticdir string, mediadir string) *gin.Engine {
	r := gin.New()
	log := log.New()

	if !gin.IsDebugging() {
		r.Use(cacheControl)
	} else {
		r.Use(Logger(log))
	}
	r.Use(gin.Recovery())
	r.GET("/", renderIndex(db))
	r.GET("/posts", renderPosts(db))
	r.HEAD("/", renderIndex(db))
	r.HEAD("/posts", renderPosts(db))
	r.GET("/posts.rss", renderRSS(db))
	r.GET("/posts.atom", renderAtom(db))
	r.GET("/posts.json", renderJSONFeed(db))
	r.GET("/sitemap.xml", renderSitemap(db))
	r.GET("/post/:id", renderPost(db))
	r.HEAD("/post/:id", renderPost(db))
	impressumDE, err := ioutil.ReadFile(staticdir + "/data/impressum.de.md")
	if err != nil {
		log.Panic(err)
	}

	r.GET("/impressum-de", renderSimpleMarkdownPage([]byte("Impressum"), impressumDE, false))

	r.Static("/css", cssdir)
	r.Static("/fonts", fontdir)
	r.Static("/media", mediadir)
	r.StaticFile("/favicon.ico", staticdir+"/data/favicon.ico")
	r.StaticFile("/robots.txt", staticdir+"/data/robots.txt")

	return r
}

func hasSuffixes(str string, suff []string) bool {
	for _, v := range suff {
		if strings.HasSuffix(str, v) {
			return true
		}
	}
	return false
}

func hasPrefixes(str string, pref []string) bool {
	for _, v := range pref {
		if strings.HasPrefix(str, v) {
			return true
		}
	}
	return false
}

func cacheControl(c *gin.Context) {
	path := c.Request.URL.Path
	if hasSuffixes(path, []string{".css", ".txt", ".ico", ".ttf"}) || hasPrefixes(path, []string{"/media"}) {
		c.Header("Cache-control", "public, max-age=31536000")
	}
	c.Next()
}
