INSERT INTO fieldkit.user (name, username, email, password, valid, bio)
    VALUES (
        'Jacob Lewallen',
        'jacob',
        'jacob@conservify.org',
        -- asdfasdfasdf
        E'\\x24326124313024747278745357794d54737a5133796b6d3241434b464f48527572616f7969555870644e6d4c577a4275784946346e3739386d4e6247',
        TRUE,
        'A Jacob user.'
    );

INSERT INTO fieldkit.user (name, username, email, password, valid, bio)
	VALUES (
		'Demo User',
		'demo-user',
		'hallway@ocr.nyc',
        -- asdfasdfasdf
        E'\\x24326124313024747278745357794d54737a5133796b6d3241434b464f48527572616f7969555870644e6d4c577a4275784946346e3739386d4e6247',
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
