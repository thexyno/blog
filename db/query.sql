-- name: GetPostsNoContent :many
select post_id, title, created_at, updated_at from posts where draft = FALSE order by created_at desc limit ? offset ?;

-- name: GetPosts :many
select post_id, title, content, created_at, updated_at from posts where draft = FALSE order by created_at desc limit ? offset ?;
-- name: GetPostsWithTags :many
select posts.post_id, title, content, created_at, updated_at, tag from posts join post_tags on posts.post_id = post_tags.post_id where draft = FALSE group by posts.post_id order by created_at desc limit ? offset ?;
-- name: GetPostsByTag :many
select posts.post_id, title, content, created_at, updated_at, tag from posts join post_tags on posts.post_id = post_tags.post_id where tag = ? and draft = FALSE order by created_at desc limit ? offset ?;
-- name: GetPostsByTagNoContent :many
select posts.post_id, title, created_at, updated_at, tag from posts join post_tags on posts.post_id = post_tags.post_id where tag = ? and draft = FALSE order by created_at desc limit ? offset ?;

-- name: GetPostIds :many
select post_id, updated_at from posts where draft = FALSE;

-- name: GetPost :one
select post_id, title, content, created_at, updated_at, draft from posts where post_id = ?;

-- name: GetTags :many
select tag from post_tags where post_id = ?;

-- name: AddPost :one
insert into posts (post_id,title,content,created_at,updated_at,draft) values (?,?,?,?,?,?) returning *;
-- name: UpdatePost :exec
update posts set post_id = ?,title = ?,content = ?,created_at = ?,updated_at = ?, draft = ? where post_id = ?;
-- name: AddTag :one
insert into post_tags (post_id, tag) values (?,?) returning *;
-- name: DeleteTags :exec
delete from post_tags where post_id = ?;
-- name: PublishPost :exec
update posts set draft = FALSE, created_at = datetime(), updated_at = datetime() where post_id = ?;
