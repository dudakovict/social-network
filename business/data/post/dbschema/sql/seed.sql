INSERT INTO posts (post_id, title, description, date_created, date_updated) VALUES
	('5cf37266-3473-4006-984f-9325122678b7', 'New Song', 'I just released a new album!', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'New Album', 'I just released a new album!', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
	ON CONFLICT DO NOTHING;