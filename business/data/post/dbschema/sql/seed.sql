INSERT INTO posts (post_id, title, description, user_id, date_created, date_updated) VALUES
	('3dc0a440-2e05-11ed-a261-0242ac120002', 'New Song', 'I just released a new song!', '5cf37266-3473-4006-984f-9325122678b7', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('47d0e86e-2e05-11ed-a261-0242ac120002', 'New Album', 'I just released a new album!', '5cf37266-3473-4006-984f-9325122678b7', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
	ON CONFLICT DO NOTHING;