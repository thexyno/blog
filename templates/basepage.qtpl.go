// Code generated by qtc from "basepage.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line templates/basepage.qtpl:1
package templates

//line templates/basepage.qtpl:1
import (
	"time"
)

//line templates/basepage.qtpl:6
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line templates/basepage.qtpl:6
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line templates/basepage.qtpl:7
type Page interface {
//line templates/basepage.qtpl:7
	Title() string
//line templates/basepage.qtpl:7
	StreamTitle(qw422016 *qt422016.Writer)
//line templates/basepage.qtpl:7
	WriteTitle(qq422016 qtio422016.Writer)
//line templates/basepage.qtpl:7
	Body() string
//line templates/basepage.qtpl:7
	StreamBody(qw422016 *qt422016.Writer)
//line templates/basepage.qtpl:7
	WriteBody(qq422016 qtio422016.Writer)
//line templates/basepage.qtpl:7
}

//line templates/basepage.qtpl:13
func StreamPageTemplate(qw422016 *qt422016.Writer, p Page) {
//line templates/basepage.qtpl:13
	qw422016.N().S(`
<!doctype html>
<html>
    <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <link href="/css/output.css" rel="stylesheet">
      <title>`)
//line templates/basepage.qtpl:20
	p.StreamTitle(qw422016)
//line templates/basepage.qtpl:20
	qw422016.N().S(`</title>
    </head>
    <body>
      <div class="top-0 z-40 w-full backdrop-blur flex-none">
        <div class="max-w-8xl mx-auto">
          <div class="py-4 lg:px-8 mx-8">
            <div class="relative flex items-center">
              <a class="pr-8 mr-3 text-2xl flex-none text-neutral_orange visited:text-neutral_orange hover:text-bright_orange font-semibold overflow-hidden md:w-auto" href="/">xynos space</a>
              <a class="mr-3 text-xl flex-none font-semibold overflow-hidden md:w-auto" href="/posts">Blog</a>
            </div>
          </div>
        </div>
      </div>
      `)
//line templates/basepage.qtpl:33
	p.StreamBody(qw422016)
//line templates/basepage.qtpl:33
	qw422016.N().S(`
      <footer class="flex justify-center w-screen bottom-0 pt-12 backdrop-blur">
        <p>Copyright (C) `)
//line templates/basepage.qtpl:35
	qw422016.N().D(time.Now().Year())
//line templates/basepage.qtpl:35
	qw422016.N().S(` xyno (Philipp Hochkamp)
      </footer>
    </body>
</html>
`)
//line templates/basepage.qtpl:39
}

//line templates/basepage.qtpl:39
func WritePageTemplate(qq422016 qtio422016.Writer, p Page) {
//line templates/basepage.qtpl:39
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/basepage.qtpl:39
	StreamPageTemplate(qw422016, p)
//line templates/basepage.qtpl:39
	qt422016.ReleaseWriter(qw422016)
//line templates/basepage.qtpl:39
}

//line templates/basepage.qtpl:39
func PageTemplate(p Page) string {
//line templates/basepage.qtpl:39
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/basepage.qtpl:39
	WritePageTemplate(qb422016, p)
//line templates/basepage.qtpl:39
	qs422016 := string(qb422016.B)
//line templates/basepage.qtpl:39
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/basepage.qtpl:39
	return qs422016
//line templates/basepage.qtpl:39
}
