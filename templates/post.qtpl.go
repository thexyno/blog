// Code generated by qtc from "post.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line ../templates/post.qtpl:1
package templates

//line ../templates/post.qtpl:1
import (
	"github.com/thexyno/xynoblog/db"
)

//line ../templates/post.qtpl:6
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line ../templates/post.qtpl:6
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line ../templates/post.qtpl:7
type PostPage struct {
	Post            db.Post
	RenderedContent []byte
}

//line ../templates/post.qtpl:13
func (p *PostPage) StreamTitle(qw422016 *qt422016.Writer) {
//line ../templates/post.qtpl:13
	qw422016.N().S(`
	`)
//line ../templates/post.qtpl:14
	qw422016.E().S(p.Post.Title)
//line ../templates/post.qtpl:14
	qw422016.N().S(` - Xynos Space
`)
//line ../templates/post.qtpl:15
}

//line ../templates/post.qtpl:15
func (p *PostPage) WriteTitle(qq422016 qtio422016.Writer) {
//line ../templates/post.qtpl:15
	qw422016 := qt422016.AcquireWriter(qq422016)
//line ../templates/post.qtpl:15
	p.StreamTitle(qw422016)
//line ../templates/post.qtpl:15
	qt422016.ReleaseWriter(qw422016)
//line ../templates/post.qtpl:15
}

//line ../templates/post.qtpl:15
func (p *PostPage) Title() string {
//line ../templates/post.qtpl:15
	qb422016 := qt422016.AcquireByteBuffer()
//line ../templates/post.qtpl:15
	p.WriteTitle(qb422016)
//line ../templates/post.qtpl:15
	qs422016 := string(qb422016.B)
//line ../templates/post.qtpl:15
	qt422016.ReleaseByteBuffer(qb422016)
//line ../templates/post.qtpl:15
	return qs422016
//line ../templates/post.qtpl:15
}

//line ../templates/post.qtpl:18
func (p *PostPage) StreamBody(qw422016 *qt422016.Writer) {
//line ../templates/post.qtpl:18
	qw422016.N().S(`
	<article class="mx-auto container max-w-7xl">
      <div class="pt-8 py-8">
      	<h1 class="p-0">`)
//line ../templates/post.qtpl:21
	qw422016.E().S(p.Post.Title)
//line ../templates/post.qtpl:21
	qw422016.N().S(`</h1>
        <p class="pl-2 text-xs font-thin text-dark3 dark:text-light3 ">Released On: `)
//line ../templates/post.qtpl:22
	qw422016.E().S(p.Post.Created.Format("2006-01-02 15:04Z07:00"))
//line ../templates/post.qtpl:22
	qw422016.N().S(`</p>
        `)
//line ../templates/post.qtpl:23
	if !p.Post.Created.Equal(p.Post.Updated) {
//line ../templates/post.qtpl:23
		qw422016.N().S(`
        <p class="pl-2 text-xs font-thin text-dark3 dark:text-light3 ">Last Updated On: `)
//line ../templates/post.qtpl:24
		qw422016.E().S(p.Post.Updated.Format("2006-01-02 15:04Z07:00"))
//line ../templates/post.qtpl:24
		qw422016.N().S(`</p>
        `)
//line ../templates/post.qtpl:25
	}
//line ../templates/post.qtpl:25
	qw422016.N().S(`

        <p class="pl-2 text-xs font-thin text-dark3 dark:text-light3 ">A `)
//line ../templates/post.qtpl:27
	qw422016.N().FPrec(p.Post.TimeToRead.Minutes(), 0)
//line ../templates/post.qtpl:27
	qw422016.N().S(` minute read.</p>
      </div>
	  `)
//line ../templates/post.qtpl:29
	qw422016.N().Z(p.RenderedContent)
//line ../templates/post.qtpl:29
	qw422016.N().S(`
	</article>
`)
//line ../templates/post.qtpl:31
}

//line ../templates/post.qtpl:31
func (p *PostPage) WriteBody(qq422016 qtio422016.Writer) {
//line ../templates/post.qtpl:31
	qw422016 := qt422016.AcquireWriter(qq422016)
//line ../templates/post.qtpl:31
	p.StreamBody(qw422016)
//line ../templates/post.qtpl:31
	qt422016.ReleaseWriter(qw422016)
//line ../templates/post.qtpl:31
}

//line ../templates/post.qtpl:31
func (p *PostPage) Body() string {
//line ../templates/post.qtpl:31
	qb422016 := qt422016.AcquireByteBuffer()
//line ../templates/post.qtpl:31
	p.WriteBody(qb422016)
//line ../templates/post.qtpl:31
	qs422016 := string(qb422016.B)
//line ../templates/post.qtpl:31
	qt422016.ReleaseByteBuffer(qb422016)
//line ../templates/post.qtpl:31
	return qs422016
//line ../templates/post.qtpl:31
}
