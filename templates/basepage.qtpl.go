// Code generated by qtc from "basepage.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line ../templates/basepage.qtpl:1
package templates

//line ../templates/basepage.qtpl:1
import (
	t "time"
)

//line ../templates/basepage.qtpl:6
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line ../templates/basepage.qtpl:6
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line ../templates/basepage.qtpl:7
type Page interface {
//line ../templates/basepage.qtpl:7
	Title() string
//line ../templates/basepage.qtpl:7
	StreamTitle(qw422016 *qt422016.Writer)
//line ../templates/basepage.qtpl:7
	WriteTitle(qq422016 qtio422016.Writer)
//line ../templates/basepage.qtpl:7
	Body() string
//line ../templates/basepage.qtpl:7
	StreamBody(qw422016 *qt422016.Writer)
//line ../templates/basepage.qtpl:7
	WriteBody(qq422016 qtio422016.Writer)
//line ../templates/basepage.qtpl:7
	Description() string
//line ../templates/basepage.qtpl:7
	StreamDescription(qw422016 *qt422016.Writer)
//line ../templates/basepage.qtpl:7
	WriteDescription(qq422016 qtio422016.Writer)
//line ../templates/basepage.qtpl:7
	Head() string
//line ../templates/basepage.qtpl:7
	StreamHead(qw422016 *qt422016.Writer)
//line ../templates/basepage.qtpl:7
	WriteHead(qq422016 qtio422016.Writer)
//line ../templates/basepage.qtpl:7
}

//line ../templates/basepage.qtpl:17
func StreamPageTemplate(qw422016 *qt422016.Writer, p Page) {
//line ../templates/basepage.qtpl:17
	qw422016.N().S(`
<!doctype html>
<html>
    <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <!--<meta name="description" content="`)
//line ../templates/basepage.qtpl:23
	p.StreamDescription(qw422016)
//line ../templates/basepage.qtpl:23
	qw422016.N().S(`">-->
      <link href="/statics/dist/main.8a90e2f4.css" rel="stylesheet">
      <title>`)
//line ../templates/basepage.qtpl:25
	p.StreamTitle(qw422016)
//line ../templates/basepage.qtpl:25
	qw422016.N().S(`</title>
      `)
//line ../templates/basepage.qtpl:26
	p.StreamHead(qw422016)
//line ../templates/basepage.qtpl:26
	qw422016.N().S(`
    </head>
    <body class="px-4 mx-2 mx-auto max-w-6xl">
      <header class="top-0 z-40 w-full flex-none mx-auto py-4 relative flex items-center">
        <a class="pr-8 mr-3 text-xl flex-none text-neutral_orange visited:text-neutral_orange hover:text-bright_orange font-semibold overflow-hidden md:w-auto" href="/">xynos space</a>
        <a class="mr-3 text-xl flex-none overflow-hidden md:w-auto" href="/posts">Blog</a>
      </header>
      `)
//line ../templates/basepage.qtpl:33
	p.StreamBody(qw422016)
//line ../templates/basepage.qtpl:33
	qw422016.N().S(`
      <footer class="flex flex-col items-center justify-center bottom-0 pb-2 pt-12 backdrop-blur">
        <p>Copyright (C) `)
//line ../templates/basepage.qtpl:35
	qw422016.N().D(t.Now().Year())
//line ../templates/basepage.qtpl:35
	qw422016.N().S(` xyno (Philipp Hochkamp)</p>
        <br>
        <p>
          <a href="/impressum-de">Impressum / Datenschutzerklärung</a>
        </p>
      </footer>
    </body>
</html>
`)
//line ../templates/basepage.qtpl:43
}

//line ../templates/basepage.qtpl:43
func WritePageTemplate(qq422016 qtio422016.Writer, p Page) {
//line ../templates/basepage.qtpl:43
	qw422016 := qt422016.AcquireWriter(qq422016)
//line ../templates/basepage.qtpl:43
	StreamPageTemplate(qw422016, p)
//line ../templates/basepage.qtpl:43
	qt422016.ReleaseWriter(qw422016)
//line ../templates/basepage.qtpl:43
}

//line ../templates/basepage.qtpl:43
func PageTemplate(p Page) string {
//line ../templates/basepage.qtpl:43
	qb422016 := qt422016.AcquireByteBuffer()
//line ../templates/basepage.qtpl:43
	WritePageTemplate(qb422016, p)
//line ../templates/basepage.qtpl:43
	qs422016 := string(qb422016.B)
//line ../templates/basepage.qtpl:43
	qt422016.ReleaseByteBuffer(qb422016)
//line ../templates/basepage.qtpl:43
	return qs422016
//line ../templates/basepage.qtpl:43
}
