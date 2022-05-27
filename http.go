//go:generate go get -d github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=templates
package xynoblog

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/thexyno/xynoblog/db"
	"github.com/thexyno/xynoblog/templates"
)

func RequestLoggerMiddleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				log.Printf(
					"[%s] %s %s %s",
					req.Method,
					req.Host,
					req.URL.Path,
					req.URL.RawQuery,
				)
			}()

			next.ServeHTTP(w, req)
		})
	}
}

func renderError(w http.ResponseWriter, r *http.Request, err error) {

	log.WithField("request", r.URL.Path).Error(err)
	var p *templates.ErrorPage
	if err == db.NotFound {
		w.WriteHeader(404)
		p = &templates.ErrorPage{
			Message: "Not Found",
		}
	} else {
		w.WriteHeader(500)
		p = &templates.ErrorPage{
			Message: "Internal Server Error",
		}
	}
	templates.WritePageTemplate(w, p)
}

func renderPosts(db db.DbConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := db.ShortPosts(0, 1000, 0)
		if err != nil {
			renderError(w, r, err)
			return
		}
		p := &templates.PostsPage{
			Posts: posts,
		}
		templates.WritePageTemplate(w, p)
	}
}
func renderIndex(db db.DbConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := db.ShortPosts(0, 5, 0)
		if err != nil {
			renderError(w, r, err)
			return
		}
		p := &templates.IndexPage{
			Posts: posts,
		}
		templates.WritePageTemplate(w, p)
	}
}
func renderPost(db db.DbConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["id"]
		post, err := db.Post(key)
		if err != nil {
			renderError(w, r, err)
			return
		}
		rendered := Render([]byte(post.Content))
		p := &templates.PostPage{
			Post:            post,
			RenderedContent: rendered,
		}
		templates.WritePageTemplate(w, p)
	}
}

func Mux(db db.DbConn, fontdir string, cssdir string) *mux.Router {
	mux := mux.NewRouter()
	mux.Use(RequestLoggerMiddleware(mux))
	mux.HandleFunc("/", renderIndex(db))
	mux.HandleFunc("/posts", renderPosts(db))
	mux.HandleFunc("/posts.rss", renderRSS(db))
	mux.HandleFunc("/posts.atom", renderAtom(db))
	mux.HandleFunc("/posts.json", renderJSONFeed(db))
	mux.HandleFunc("/sitemap.xml", renderSitemap(db))
	mux.HandleFunc("/post/{id}", renderPost(db))
	CSSFileServer := http.FileServer(http.Dir(cssdir))
	mux.PathPrefix("/css/").Handler(http.StripPrefix("/css", CSSFileServer))
	FontFileServer := http.FileServer(http.Dir(fontdir))
	mux.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts", FontFileServer))

	return mux
}
