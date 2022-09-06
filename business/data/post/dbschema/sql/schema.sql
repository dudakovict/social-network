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