//go:generate go get -d github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=templates
package xynoblog

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

func renderError(w http.ResponseWriter, r *http.Request, err string) {
	w.WriteHeader(500)
	fmt.Fprint(w, "Uff: ", err)
}

func renderIndex(db db.DbConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := db.ShortPosts(5, 0)
		if err != nil {
			renderError(w, r, err.Error())
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
			renderError(w, r, err.Error())
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

func Mux(db db.DbConn) *mux.Router {
	mux := mux.NewRouter()
	mux.Use(RequestLoggerMiddleware(mux))
	mux.HandleFunc("/", renderIndex(db))
	mux.HandleFunc("/post/{id}", renderPost(db))
	fileServer := http.FileServer(http.Dir("./cssdist"))
	mux.PathPrefix("/css/").Handler(http.StripPrefix("/css", fileServer))

	return mux
}
