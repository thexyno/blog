-- name: GetPostsNoContent :many
select post_id, title, created_at, updated_at from posts order by created_at desc limit ? offset ?;

-- name: GetPosts :many
select post_id, title, content, created_at, updated_at from posts order by created_at desc limit ? offset ?;
-- name: GetPostsWithTags :many
select posts.post_id, title, content, created_at, updated_at, tag from posts join post_tags on posts.post_id = post_tags.post_id group by posts.post_id order by created_at desc limit ? offset ?;

-- name: GetPostIds :many
select post_id, updated_at from posts;

-- name: GetPost :one
select post_id, title, content, created_at, updated_at from posts where post_id = ?;

-- name: GetTags :many
select tag from post_tags where post_id = ?;

-- name: AddPost :one
insert into posts (post_id,title,content,created_at,updated_at) values (?,?,?,?,?) returning *;
-- name: UpdatePost :exec
update posts set post_id = ?,title = ?,content = ?,created_at = ?,updated_at = ? where post_id = ?;
-- name: AddTag :one
insert into post_tags (post_id, tag) values (?,?) returning *;
-- name: DeleteTags :exec
delete from post_tags where post_id = ?;
