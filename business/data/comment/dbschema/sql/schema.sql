-- Version: 1.1
-- Description: Create table comments
CREATE TABLE comments (
	comment_id        UUID,
	description    TEXT,
	post_id        UUID,
	user_id        UUID,
	date_created   TIMESTAMP,
	date_updated   TIMESTAMP,

	PRIMARY KEY (comment_id)
);