package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DbConn struct {
	db *sql.DB
}

type Post struct {
	Id      string
	Title   string
	Content string
	Created time.Time
	Updated time.Time
	Tags    []string
}

const (
	shortPostStmt string = "select id, title, substr(content,0,100), created, updated, tags from posts limit ? offset ?"
	postStmt      string = "select id, title, content, created, updated, tags from posts where id = ?"
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
		Id:      "i-like-bread2",
		Title:   "I like Bread",
		Content: "# Bread is life\n## Bread is Love\n\nBread will consume us all\n\n- That's a very enby thing to say",
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

/// returns post  with id = id
func (conn *DbConn) Post(id string) (Post, error) {
	rows, err := conn.db.Query(postStmt, id)
	if err != nil {
		return Post{}, err
	}
	defer rows.Close()
	var post Post
	for rows.Next() {
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
		post = Post{id, title, content, created, updated, tags}
	}
	return post, nil
}

/// returns posts with only 100 chars of post text
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
		var content string
		var created time.Time
		var updated time.Time
		var tmpTags string
		rowErr := rows.Scan(&id, &title, &content, &created, &updated, &tmpTags)
		if rowErr != nil {
			log.Print("ShortPosts row Error: ", rowErr)
		}
		var tags []string
		tagsErr := json.Unmarshal([]byte(tmpTags), &tags)
		if tagsErr != nil {
			log.Print("ShortPosts tags Error: ", tagsErr)
		}
		posts = append(posts, Post{id, title, content, created, updated, tags})
	}
	return posts, nil
}
