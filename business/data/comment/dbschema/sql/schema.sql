-- Version: 1.1
-- Description: Create table posts
CREATE TABLE posts (
	post_id        UUID,
	title          TEXT,
	description    TEXT,
	user_id        UUID,
	date_created   TIMESTAMP,
	date_updated   TIMESTAMP,

	PRIMARY KEY (post_id)
);

-- Version: 1.2
-- Description: Create table comments
CREATE TABLE comments (
	comment_id     UUID,
	description    TEXT,
	user_id        UUID,
	post_id        UUID,
	date_created   TIMESTAMP,
	date_updated   TIMESTAMP,

	PRIMARY KEY (comment_id),
	FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);