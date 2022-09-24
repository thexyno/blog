CREATE TABLE posts (post_id text primary key not null,
                    title text not null,
                    content text not null,
                    created_at DATETIME not null,
                    updated_at DATETIME not null
);

CREATE TABLE post_tags (post_id text not null,
                        tag text not null,
                        FOREIGN KEY (post_id)
                          REFERENCES posts (post_id),
                        PRIMARY KEY (post_id, tag) -- it's all a primary key
);
