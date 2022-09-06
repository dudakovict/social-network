INSERT INTO comments (comment_id, description, post_id, user_id, date_created, date_updated) VALUES
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'What a great song!', '5cf37266-3473-4006-984f-9325122678b7', '5cf37266-3473-4006-984f-9325122678b7', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('5cf37266-3473-4006-984f-9325122678b7', 'What a great album!', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
	ON CONFLICT DO NOTHING;