package db

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"math"
	"regexp"
	"time"

	_ "github.com/mattn/go-sqlite3"
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
}
type PostNoContent struct {
	Id      PostId
	Title   string
	Created time.Time
	Updated time.Time
}

func (post *Post) ToIdUpdated() IdUpdated {
	return IdUpdated{post.Id, post.Updated}
}

var ErrNotFound error = errors.New("not found")

//go:embed schema.sql
var ddl string

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
	// create tables
	if _, err := conn.db.ExecContext(conn.ctx, ddl); err != nil {
		return err
	}
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

/// returns post  with id = id
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
	}, nil
}

/// returns post  with id = id
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
			}
	}
	return toReturn, nil
}

func (conn *DbConn) PostIds() ([]dbconn.GetPostIdsRow, error) {
	return conn.queries.GetPostIds(conn.ctx)
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
		}
	}
	return toReturn, nil
}
