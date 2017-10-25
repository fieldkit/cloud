INSERT INTO fieldkit.user (name, username, email, password, valid, bio)
	VALUES (
		'Demo User',
		'demo-user',
		'hallway@ocr.nyc',
		E'\\x2432612430342439322e6e30632f70794739452e456247395477454b4f33724e4c6b735a464d536c787368626b555667665568536c4e6b2f71623771',
		TRUE,
		'A demo user.'
	);

INSERT INTO fieldkit.project (name, slug, description)
	VALUES (
		'Demo Project',
		'demo-project',
		'A demo project.'
	);

INSERT INTO fieldkit.project_user (user_id, project_id)
	VALUES (
		(SELECT id FROM fieldkit.user WHERE username = 'demo-user'),
		(SELECT id FROM fieldkit.project WHERE slug = 'demo-project')
	);

INSERT INTO fieldkit.expedition (project_id, name, slug, description)
	VALUES (
		(SELECT id FROM fieldkit.project WHERE slug = 'demo-project'),
		'Demo Expedition',
		'demo-expedition',
		'A demo expedition.'
	);

INSERT INTO fieldkit.project (name, slug, description)
VALUES (
'Www Project',
'www',
'A www project.'
);

INSERT INTO fieldkit.project_user (user_id, project_id)
VALUES (
(SELECT id FROM fieldkit.user WHERE username = 'demo-user'),
(SELECT id FROM fieldkit.project WHERE slug = 'www')
);

INSERT INTO fieldkit.expedition (project_id, name, slug, description)
VALUES (
(SELECT id FROM fieldkit.project WHERE slug = 'www'),
'Www Expedition',
'www-expedition',
'A www expedition.'
);
