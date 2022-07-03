//go:generate go get -d github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=../templates
package server

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
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
		if !gin.IsDebugging() {
			c.Header("Last-Modified", post.Updated.Format(http.TimeFormat))
		}
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

func Mux(db db.DbConn, fontdir string, cssdir string, staticdir string) *gin.Engine {
	r := gin.New()
	log := log.New()

	store := persistence.NewInMemoryStore(300 * time.Second)
	if !gin.IsDebugging() {
		r.Use(cacheControl)
	} else {
		store = persistence.NewInMemoryStore(time.Millisecond)
		r.Use(Logger(log))
	}
	r.Use(gin.Recovery())
	r.GET("/", cache.CachePage(store, 5*time.Minute, renderIndex(db)))
	r.GET("/posts", cache.CachePage(store, 5*time.Minute, renderPosts(db)))
	r.HEAD("/", cache.CachePage(store, 5*time.Minute, renderIndex(db)))
	r.HEAD("/posts", cache.CachePage(store, 5*time.Minute, renderPosts(db)))
	r.GET("/posts.rss", cache.CachePage(store, 5*time.Minute, renderRSS(db)))
	r.GET("/posts.atom", cache.CachePage(store, 5*time.Minute, renderAtom(db)))
	r.GET("/posts.json", cache.CachePage(store, 5*time.Minute, renderJSONFeed(db)))
	r.GET("/sitemap.xml", cache.CachePage(store, 5*time.Minute, renderSitemap(db)))
	r.GET("/post/:id", cache.CachePage(store, 5*time.Minute, renderPost(db)))
	r.HEAD("/post/:id", cache.CachePage(store, 5*time.Minute, renderPost(db)))
	impressumDE, err := ioutil.ReadFile(staticdir + "/data/impressum.de.md")
	if err != nil {
		log.Panic(err)
	}
	datenschutzDE, err := ioutil.ReadFile(staticdir + "/data/datenschutz.de.md")
	if err != nil {
		log.Panic(err)
	}

	r.GET("/impressum-de", renderSimpleMarkdownPage([]byte("Impressum"), impressumDE, false))
	r.GET("/datenschutz-de", renderSimpleMarkdownPage([]byte("Datenschutzerklärung"), datenschutzDE, false))

	r.Static("/css", cssdir)
	r.Static("/fonts", fontdir)
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

func cacheControl(c *gin.Context) {
	path := c.Request.URL.Path
	if hasSuffixes(path, []string{".css", ".txt", ".ico", ".ttf"}) {
		c.Header("Cache-control", "public, max-age=31536000")
	}
	c.Next()
}
