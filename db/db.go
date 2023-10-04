package db

import (
	"context"
	"database/sql"
	"embed"
	_ "embed"
	"errors"
	"math"
	"net/http"
	"regexp"
	"time"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"github.com/thexyno/xynoblog/dbconn"
)

type DbConn struct {
	db      *sql.DB
	ctx     context.Context
	queries *dbconn.Queries
}

type PostId string

type IdUpdated struct {
	Id      PostId
	Updated time.Time
}

type Post struct {
	Id         PostId
	Title      string
	Content    string
	Created    time.Time
	Updated    time.Time
	Tags       []string
	TimeToRead time.Duration
	Draft      bool
}
type PostNoContent struct {
	Id      PostId
	Title   string
	Created time.Time
	Updated time.Time
	Draft   bool
}

func (post *Post) ToIdUpdated() IdUpdated {
	return IdUpdated{post.Id, post.Updated}
}

var ErrNotFound error = errors.New("not found")

//go:embed migrations
var ddl embed.FS

func NewDb(uri string) DbConn {
	log.Infof("dbpath is %v", uri)
	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		log.Fatal(err)
	}
	return DbConn{db, context.Background(), dbconn.New(db)}
}

func (conn *DbConn) Close() error {
	return conn.db.Close()
}

func (conn *DbConn) Seed() error {
	httpFs := http.FS(ddl)
	migrationSource := &migrate.HttpFileSystemMigrationSource{
		FileSystem: httpFs,
	}
	n, err := migrate.Exec(conn.db, "sqlite3", migrationSource, migrate.Up)
	if err != nil {
		return err
	}
	log.Infof("Applied %d migrations!", n)
	return nil
}

func (conn *DbConn) updatePost(post Post) error {
	tx, err := conn.db.Begin()
	if err != nil {
		return err
	}
	qtx := conn.queries.WithTx(tx)
	err = qtx.UpdatePost(conn.ctx,
		dbconn.UpdatePostParams{
			PostID:    string(post.Id),
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.Created,
			UpdatedAt: time.Now(),
			PostID_2:  string(post.Id),
		},
	)
	if err != nil {
		return err
	}
	err = qtx.DeleteTags(conn.ctx, string(post.Id))
	if err != nil {
		return err
	}
	for _, tag := range post.Tags {
		_, err = qtx.AddTag(conn.ctx, dbconn.AddTagParams{
			PostID: string(post.Id),
			Tag:    tag,
		})
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (conn *DbConn) insertPost(post Post) error {
	tx, err := conn.db.Begin()
	if err != nil {
		return err
	}
	qtx := conn.queries.WithTx(tx)
	_, err = qtx.AddPost(conn.ctx,
		dbconn.AddPostParams{
			PostID:    string(post.Id),
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.Created,
			UpdatedAt: post.Created,
			Draft:     post.Draft,
		},
	)
	if err != nil {
		return err
	}
	for _, tag := range post.Tags {
		_, err = qtx.AddTag(conn.ctx, dbconn.AddTagParams{
			PostID: string(post.Id),
			Tag:    tag,
		})
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (conn *DbConn) Add(post Post) error {
	dbpost, err := conn.Post(string(post.Id))
	log.Infof("adding %s to posts", post.Id)
	if err != nil || dbpost.Id == "" {
		if err := conn.insertPost(post); err != nil {
			return err
		}
		log.Info("added via insert")
	} else if dbpost.Updated.Before(post.Updated) || dbpost.Content != post.Content {
		if err := conn.updatePost(post); err != nil {
			return err
		}
		log.Info("added via update")
	} else {
		log.Info("post didnt change")
	}
	return err
}

func wordCount(value string) int {
	// Match non-space character sequences.
	re := regexp.MustCompile(`[\S]+`)

	// Find all matches and return count.
	results := re.FindAllString(value, -1)
	return len(results)
}

// / returns post  with id = id
func (conn *DbConn) Post(id string) (Post, error) {
	post, err := conn.queries.GetPost(conn.ctx, id)
	if err != nil {
		return Post{}, err
	}
	tags, err := conn.queries.GetTags(conn.ctx, id)
	if err != nil {
		return Post{}, err
	}
	duration := time.Duration(math.Ceil(float64(wordCount(post.Content))/239.0) * 60000000000)
	return Post{
		PostId(post.PostID),
		post.Title,
		post.Content,
		post.CreatedAt,
		post.UpdatedAt,
		tags,
		duration,
		post.Draft,
	}, nil
}

// / returns posts limited by limit and offset order by updated desc
func (conn *DbConn) Posts(limit int64, offset int64) ([]Post, error) {
	posts, err := conn.queries.GetPosts(conn.ctx, dbconn.GetPostsParams{Limit: limit, Offset: offset})
	if err != nil {
		return []Post{}, err
	}
	var toReturn = make([]Post, len(posts))
	for i, post := range posts {
		toReturn[i] =
			Post{
				PostId(post.PostID),
				post.Title,
				post.Content,
				post.CreatedAt,
				post.UpdatedAt,
				[]string{},
				0,
				false,
			}
	}
	return toReturn, nil
}

// / returns shortPosts with tag = tag
func (conn *DbConn) ShortPostsByTag(tag string, limit, offset int64) ([]PostNoContent, error) {
	posts, err := conn.queries.GetPostsByTagNoContent(conn.ctx, dbconn.GetPostsByTagNoContentParams{Tag: tag, Limit: limit, Offset: offset})
	if err != nil {
		return []PostNoContent{}, err
	}
	var toReturn = make([]PostNoContent, len(posts))
	for i, post := range posts {
		toReturn[i] = PostNoContent{
			PostId(post.PostID),
			post.Title,
			post.CreatedAt,
			post.UpdatedAt,
			false,
		}
	}
	return toReturn, nil
}

func (conn *DbConn) PostIds() ([]dbconn.GetPostIdsRow, error) {
	return conn.queries.GetPostIds(conn.ctx)
}

func (conn *DbConn) Publish(id string) error {
	return conn.queries.PublishPost(conn.ctx, id)
}

// returns posts with only textLength chars of post text and no duration
// limit = -1 returns all
func (conn *DbConn) ShortPosts(limit int64, skip int64) ([]PostNoContent, error) {
	posts, err := conn.queries.GetPostsNoContent(conn.ctx, dbconn.GetPostsNoContentParams{Limit: limit, Offset: skip})

	if err != nil {
		return []PostNoContent{}, err
	}
	var toReturn = make([]PostNoContent, len(posts))
	for i, post := range posts {
		toReturn[i] = PostNoContent{
			PostId(post.PostID),
			post.Title,
			post.CreatedAt,
			post.UpdatedAt,
			false,
		}
	}
	return toReturn, nil
}
