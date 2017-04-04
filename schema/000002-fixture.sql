-- password: testpasword
INSERT INTO fieldkit.user VALUES (1, 'test name', 'test', 'test@test.com', E'\\x2432612430342439322e6e30632f70794739452e456247395477454b4f33724e4c6b735a464d536c787368626b555667665568536c4e6b2f71623771', true, 'test bio');
INSERT INTO fieldkit.user VALUES (2, 'test02 name', 'test02', 'test02@test.com', E'\\x2432612430342439322e6e30632f70794739452e456247395477454b4f33724e4c6b735a464d536c787368626b555667665568536c4e6b2f71623771', true, 'test02 bio');

INSERT INTO fieldkit.project VALUES (1, 'project01', 'project01');
INSERT INTO fieldkit.project_user VALUES (1, 1);
INSERT INTO fieldkit.expedition VALUES (1, 1, 'expedition01', 'expedition01');

INSERT INTO fieldkit.team VALUES (1, 1, 'team01', 'team01');
INSERT INTO fieldkit.team_user VALUES (1, 1, 'role01');
INSERT INTO fieldkit.team_user VALUES (1, 2, 'role02');

INSERT INTO fieldkit.input VALUES (1, 1, 'input01-twitter', 1, NULL, true);

INSERT INTO fieldkit.twitter_oauth VALUES (1, 'AFDN:ASDN', 'jnasdnjkaf');
INSERT INTO fieldkit.twitter_account VALUES (1, 'ocr_test', 'lknadslknad', 'njkladskjnads');

INSERT INTO fieldkit.input_twitter_account VALUES (1, 1);

INSERT INTO fieldkit.schema VALUES (1, 1, '{ "id": 1 }');

INSERT INTO fieldkit.document VALUES (1, 1, 1, NULL, 1, current_timestamp - interval '0 hour', ST_SetSRID(ST_MakePoint(38.99, -76.81), 4326), '{}');
INSERT INTO fieldkit.document VALUES (2, 1, 1, NULL, 1, current_timestamp - interval '1 hour', ST_SetSRID(ST_MakePoint(39.01, -76.81), 4326), '{}');
INSERT INTO fieldkit.document VALUES (3, 1, 1, NULL, 1, current_timestamp - interval '2 hour', ST_SetSRID(ST_MakePoint(39.03, -76.81), 4326), '{}');
INSERT INTO fieldkit.document VALUES (4, 1, 1, NULL, 1, current_timestamp - interval '3 hour', ST_SetSRID(ST_MakePoint(39.05, -76.81), 4326), '{}');
INSERT INTO fieldkit.document VALUES (5, 1, 1, NULL, 1, current_timestamp - interval '4 hour', ST_SetSRID(ST_MakePoint(39.07, -76.81), 4326), '{}');
