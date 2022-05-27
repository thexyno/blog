package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"math"
	"regexp"
	"time"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type DbConn struct {
	db *sql.DB
}

type Post struct {
	Id         string
	Title      string
	Content    string
	Created    time.Time
	Updated    time.Time
	Tags       []string
	TimeToRead time.Duration
}

const (
	shortPostStmt string = "select id, title, created, updated, tags from posts order by created desc limit ? offset ?"
	postStmt      string = "select id, title, content, created, updated, tags from posts where id = ?"
)

var (
	NotFound error = errors.New("Not Found")
)

func NewDb(uri string) DbConn {
	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		log.Fatal(err)
	}
	return DbConn{db}
}

func (conn *DbConn) Close() error {
	return conn.db.Close()
}

func (conn *DbConn) Seed() error {
	stmt := `
create table if not exists posts (id text primary key not null, title text not null, content text not null, created datetime not null, updated datetime not null, tags text not null);
`
	_, err := conn.db.Exec(stmt)
	if err != nil {
		return err
	}
	err = conn.Add(Post{
		Id:      "i-like-lorem-ipsum",
		Title:   "I like lorem ipsum",
		Content: "Lorem ipsum dolor sit amet",
		Created: time.Now(),
		Updated: time.Now(),
		Tags:    []string{"bread", "enby"},
	})
	return err
}

func (conn *DbConn) Add(post Post) error {
	tags, err := json.Marshal(post.Tags)
	if err != nil {
		return err
	}
	_, err = conn.db.Exec("insert or ignore into posts values (?,?,?,?,?,?)", post.Id, post.Title, post.Content, post.Created, post.Updated, tags)
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
	rows, err := conn.db.Query(postStmt, id)
	if err != nil {
		return Post{}, err
	}
	defer rows.Close()
	var post Post
	hasRow := false
	for rows.Next() {
		hasRow = true
		var id string
		var title string
		var content string
		var created time.Time
		var updated time.Time
		var tmpTags string
		rowErr := rows.Scan(&id, &title, &content, &created, &updated, &tmpTags)
		if rowErr != nil {
			log.Print("Post row Error: ", rowErr)
			return Post{}, err
		}
		var tags []string
		tagsErr := json.Unmarshal([]byte(tmpTags), &tags)
		if tagsErr != nil {
			log.Print("Post tags Error: ", tagsErr)
			return Post{}, err
		}
		duration := time.Duration(math.Ceil(float64(wordCount(content))/239.0) * 60000000000)
		post = Post{id, title, content, created, updated, tags, duration}
	}
	if hasRow {
		return post, nil
	} else {
		return post, NotFound
	}
}

/// returns posts with only 0 chars of post text and no duration
func (conn *DbConn) ShortPosts(limit uint, skip uint) ([]Post, error) {
	rows, err := conn.db.Query(string(shortPostStmt), limit, skip)
	if err != nil {
		return []Post{}, err
	}
	defer rows.Close()
	posts := []Post{}
	for rows.Next() {
		var id string
		var title string
		var created time.Time
		var updated time.Time
		var tmpTags string
		rowErr := rows.Scan(&id, &title, &created, &updated, &tmpTags)
		if rowErr != nil {
			log.Print("ShortPosts row Error: ", rowErr)
		}
		var tags []string
		tagsErr := json.Unmarshal([]byte(tmpTags), &tags)
		if tagsErr != nil {
			log.Print("ShortPosts tags Error: ", tagsErr)
		}
		posts = append(posts, Post{id, title, "", created, updated, tags, time.Duration(0)})
	}
	return posts, nil
}
