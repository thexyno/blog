//go:generate go get -d github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=templates
package server

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/thexyno/xynoblog/db"
	"github.com/thexyno/xynoblog/templates"
)

func renderError(c *gin.Context, err error) {
	log.WithField("request", c.FullPath()).Error(err)
	var p *templates.ErrorPage
	if err == db.NotFound {
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
		rendered := Render([]byte(post.Content))
		p := &templates.PostPage{
			Post:            post,
			RenderedContent: rendered,
		}
		templates.WritePageTemplate(c.Writer, p)
	}
}

func Mux(db db.DbConn, fontdir string, cssdir string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/", renderIndex(db))
	r.GET("/posts", renderPosts(db))
	r.GET("/posts.rss", renderRSS(db))
	r.GET("/posts.atom", renderAtom(db))
	r.GET("/posts.json", renderJSONFeed(db))
	r.GET("/sitemap.xml", renderSitemap(db))
	r.GET("/post/:id", renderPost(db))
	r.Static("/css", cssdir)
	r.Static("/fonts", fontdir)

	return r
}
