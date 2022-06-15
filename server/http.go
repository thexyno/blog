//go:generate go get -d github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=../templates
package server

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/thexyno/xynoblog/db"
	"github.com/thexyno/xynoblog/templates"
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
		posts, err := db.ShortPosts(0, 1000, 0)
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
		posts, err := db.ShortPosts(0, 5, 0)
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
		c.Header("Last-Modified", post.Updated.Format(http.TimeFormat))
		rendered := Render([]byte(post.Content))
		p := &templates.PostPage{
			Post:            post,
			RenderedContent: rendered,
		}
		templates.WritePageTemplate(c.Writer, p)
	}
}

func renderSimpleMarkdownPage(title []byte, content []byte, index bool) func(*gin.Context) {
	return func(c *gin.Context) {
		rendered := Render(content)
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

func Mux(db db.DbConn, fontdir string, cssdir string) *gin.Engine {
	r := gin.New()
	log := log.New()
	r.Use(gin.Recovery())
	r.Use(Logger(log))
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

	impressumDE, err := ioutil.ReadFile("./data/impressum.de.md")
	if err != nil {
		log.Panic(err)
	}
	datenschutzDE, err := ioutil.ReadFile("./data/datenschutz.de.md")
	if err != nil {
		log.Panic(err)
	}

	r.GET("/impressum-de", renderSimpleMarkdownPage([]byte("Impressum"), impressumDE, false))
	r.GET("/datenschutz-de", renderSimpleMarkdownPage([]byte("Datenschutzerklärung"), datenschutzDE, false))

	r.Static("/css", cssdir)
	r.Static("/fonts", fontdir)
	r.StaticFile("/favicon.ico", "./data/favicon.ico")
	r.StaticFile("/robots.txt", "./data/robots.txt")

	return r
}
