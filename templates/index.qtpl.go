// Code generated by qtc from "index.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line templates/index.qtpl:1
package templates

//line templates/index.qtpl:1
import (
	"github.com/thexyno/xynoblog/db"
)

//line templates/index.qtpl:6
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line templates/index.qtpl:6
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line templates/index.qtpl:7
type IndexPage struct {
	Posts []db.Post
}

//line templates/index.qtpl:13
func (p *IndexPage) StreamTitle(qw422016 *qt422016.Writer) {
//line templates/index.qtpl:13
	qw422016.N().S(`
	xynos space
`)
//line templates/index.qtpl:15
}

//line templates/index.qtpl:15
func (p *IndexPage) WriteTitle(qq422016 qtio422016.Writer) {
//line templates/index.qtpl:15
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/index.qtpl:15
	p.StreamTitle(qw422016)
//line templates/index.qtpl:15
	qt422016.ReleaseWriter(qw422016)
//line templates/index.qtpl:15
}

//line templates/index.qtpl:15
func (p *IndexPage) Title() string {
//line templates/index.qtpl:15
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/index.qtpl:15
	p.WriteTitle(qb422016)
//line templates/index.qtpl:15
	qs422016 := string(qb422016.B)
//line templates/index.qtpl:15
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/index.qtpl:15
	return qs422016
//line templates/index.qtpl:15
}

//line templates/index.qtpl:18
func (p *IndexPage) StreamBody(qw422016 *qt422016.Writer) {
//line templates/index.qtpl:18
	qw422016.N().S(`
	<div class="mx-auto container max-w-6xl">
    <h1 class="pt-8 font-bold">xyno</h1>
    <h3 class="px-4 text-xl text-dark1 dark:text-light1">Full Stack Engineer</h3>
    <h3 class="text-xl">Skills:</h3>
    <div class="px-8">
      <ul>
        <li>Go</li>
        <li>Typescript/Javascript</li>
      </ul>
      <ul class="px-4">
        <li>Angular</li>
        <li>Nest.JS</li>
      </ul>
      <ul>
        <li>Kubernetes</li>
        <li>Terraform</li>
        <li>Nix/NixOS</li>
        <li>GitLab CI</li>
      </ul>
    </div>
	<h3 class="pt-4 text-xl font-semibold">Links:</h3>
      <div class="px-8">
      <ul>
        <li><a target="_blank" href="https://github.com/thexyno">GitHub - thexyno</a></li>
        <li><a rel="me" target="_blank" href="https://matrix.to/#/@me:ragon.xyz">Matrix - @me:ragon.xyz</a></li>
        <li><a rel="me" target="_blank" href="https://chaos.social/@xyno">Mastodon - @xyno@chaos.social</a></li>
      </ul>
	  </div>
	<h3 class="pt-4 text-xl font-semibold">Latest Posts:</h3>
      <div class="px-8 flex flex-col">
      `)
//line templates/index.qtpl:49
	if len(p.Posts) == 0 {
//line templates/index.qtpl:49
		qw422016.N().S(`
	  	No posts.
	  `)
//line templates/index.qtpl:51
	} else {
//line templates/index.qtpl:51
		qw422016.N().S(`
	  		`)
//line templates/index.qtpl:52
		streamemitPosts(qw422016, p.Posts)
//line templates/index.qtpl:52
		qw422016.N().S(`
	  `)
//line templates/index.qtpl:53
	}
//line templates/index.qtpl:53
	qw422016.N().S(`
	  </div>
	</div>
`)
//line templates/index.qtpl:56
}

//line templates/index.qtpl:56
func (p *IndexPage) WriteBody(qq422016 qtio422016.Writer) {
//line templates/index.qtpl:56
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/index.qtpl:56
	p.StreamBody(qw422016)
//line templates/index.qtpl:56
	qt422016.ReleaseWriter(qw422016)
//line templates/index.qtpl:56
}

//line templates/index.qtpl:56
func (p *IndexPage) Body() string {
//line templates/index.qtpl:56
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/index.qtpl:56
	p.WriteBody(qb422016)
//line templates/index.qtpl:56
	qs422016 := string(qb422016.B)
//line templates/index.qtpl:56
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/index.qtpl:56
	return qs422016
//line templates/index.qtpl:56
}

//line templates/index.qtpl:58
func streamemitPosts(qw422016 *qt422016.Writer, posts []db.Post) {
//line templates/index.qtpl:58
	qw422016.N().S(`
   `)
//line templates/index.qtpl:59
	for _, v := range posts {
//line templates/index.qtpl:59
		qw422016.N().S(`
     <a href="/post/`)
//line templates/index.qtpl:60
		qw422016.E().S(string(v.Id))
//line templates/index.qtpl:60
		qw422016.N().S(`">
      <span>`)
//line templates/index.qtpl:61
		qw422016.E().S(v.Title)
//line templates/index.qtpl:61
		qw422016.N().S(`</span>
      <span class="text-xs font-thin text-dark3 dark:text-light3 ">(`)
//line templates/index.qtpl:62
		qw422016.E().S(v.Created.Format("2006-01-02"))
//line templates/index.qtpl:62
		qw422016.N().S(`)</span>
     </a>
   `)
//line templates/index.qtpl:64
	}
//line templates/index.qtpl:64
	qw422016.N().S(`
`)
//line templates/index.qtpl:65
}

//line templates/index.qtpl:65
func writeemitPosts(qq422016 qtio422016.Writer, posts []db.Post) {
//line templates/index.qtpl:65
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/index.qtpl:65
	streamemitPosts(qw422016, posts)
//line templates/index.qtpl:65
	qt422016.ReleaseWriter(qw422016)
//line templates/index.qtpl:65
}

//line templates/index.qtpl:65
func emitPosts(posts []db.Post) string {
//line templates/index.qtpl:65
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/index.qtpl:65
	writeemitPosts(qb422016, posts)
//line templates/index.qtpl:65
	qs422016 := string(qb422016.B)
//line templates/index.qtpl:65
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/index.qtpl:65
	return qs422016
//line templates/index.qtpl:65
}
