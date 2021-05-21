--
-- PostgreSQL database dump
--

-- Dumped from database version 11.2
-- Dumped by pg_dump version 11.12 (Ubuntu 11.12-1.pgdg20.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: migrations; Type: TABLE; Schema: fieldkit; Owner: fieldkit
--

CREATE TABLE fieldkit.migrations (
    id integer NOT NULL,
    name text,
    batch integer,
    completed_at timestamp with time zone
);


ALTER TABLE fieldkit.migrations OWNER TO fieldkit;

--
-- Name: migrations_id_seq; Type: SEQUENCE; Schema: fieldkit; Owner: fieldkit
--

CREATE SEQUENCE fieldkit.migrations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE fieldkit.migrations_id_seq OWNER TO fieldkit;

--
-- Name: migrations_id_seq; Type: SEQUENCE OWNED BY; Schema: fieldkit; Owner: fieldkit
--

ALTER SEQUENCE fieldkit.migrations_id_seq OWNED BY fieldkit.migrations.id;


--
-- Name: migrations id; Type: DEFAULT; Schema: fieldkit; Owner: fieldkit
--

ALTER TABLE ONLY fieldkit.migrations ALTER COLUMN id SET DEFAULT nextval('fieldkit.migrations_id_seq'::regclass);


--
-- Data for Name: migrations; Type: TABLE DATA; Schema: fieldkit; Owner: fieldkit
--

INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (1, '20190101000000_schema.up.sql', 1, '2021-05-21 17:09:39.608226+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (2, '20190101000001_old_funcs.up.sql', 1, '2021-05-21 17:09:39.6121+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (3, '20190101000002_station.up.sql', 1, '2021-05-21 17:09:39.616084+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (4, '20190101000003_ingestion.up.sql', 1, '2021-05-21 17:09:39.620185+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (5, '20190101000004_field_notes.up.sql', 1, '2021-05-21 17:09:39.62412+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (6, '20190102000000_countries_all.up.sql', 1, '2021-05-21 17:09:39.627925+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (7, '20190102000001_countries_srid.up.sql', 1, '2021-05-21 17:09:39.632224+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (8, '20191029102823_test.up.sql', 1, '2021-05-21 17:09:39.636272+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (9, '20191029115022_add_media_columns_to_project.sql.up.sql', 1, '2021-05-21 17:09:39.640229+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (10, '20191029115101_add_media_columns_to_user.sql.up.sql', 1, '2021-05-21 17:09:39.644053+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (11, '20191106130836_create_project_invites_table.up.sql', 1, '2021-05-21 17:09:39.648036+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (12, '20200109103922_add_archive_table.up.sql', 1, '2021-05-21 17:09:39.65192+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (13, '20200122095929_add_project_station.up.sql', 1, '2021-05-21 17:09:39.655633+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (14, '20200122133104_add_more_field_note_categories.up.sql', 1, '2021-05-21 17:09:39.6653+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (15, '20200124122015_add_admin_user_flag.up.sql', 1, '2021-05-21 17:09:39.669196+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (16, '20200303154338_add_ingestion_error_fields.up.sql', 1, '2021-05-21 17:09:39.673664+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (17, '20200326134821_recovery_token.up.sql', 1, '2021-05-21 17:09:39.677342+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (18, '20200413121223_add_project_role.up.sql', 1, '2021-05-21 17:09:39.681235+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (19, '20200420102107_add_firmware_available_flag.up.sql', 1, '2021-05-21 17:09:39.685201+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (20, '20200420125812_add_privacy_flags.up.sql', 1, '2021-05-21 17:09:39.694707+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (21, '20200420125831_add_social_model.up.sql', 1, '2021-05-21 17:09:39.698808+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (22, '20200420130822_add_station_model.up.sql', 1, '2021-05-21 17:09:39.703323+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (23, '20200427105709_add_invite_token.up.sql', 1, '2021-05-21 17:09:39.707399+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (24, '20200428122652_drop_old_tables.up.sql', 1, '2021-05-21 17:09:39.711385+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (25, '20200429102629_add_station_fields.up.sql', 1, '2021-05-21 17:09:39.715184+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (26, '20200504090203_drop_more_old_schema.up.sql', 1, '2021-05-21 17:09:39.7192+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (27, '20200512111445_add_user_tester_flag.up.sql', 1, '2021-05-21 17:09:39.723137+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (28, '20200514101113_add_role_to_invite.up.sql', 1, '2021-05-21 17:09:39.727861+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (29, '20200514105551_add_sensor_reading_columns.up.sql', 1, '2021-05-21 17:09:39.731973+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (30, '20200514130752_add_module_internal_flag.up.sql', 1, '2021-05-21 17:09:39.73606+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (31, '20200514144645_fix_hardware_id_index.up.sql', 1, '2021-05-21 17:09:39.740116+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (32, '20200515075528_add_more_station_fields.up.sql', 1, '2021-05-21 17:09:39.744182+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (33, '20200518092953_add_station_module_index.up.sql', 1, '2021-05-21 17:09:39.748049+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (34, '20200519105300_add_station_module_non_meta_index.up.sql', 1, '2021-05-21 17:09:39.751968+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (35, '20200519160944_add_configuration.up.sql', 1, '2021-05-21 17:09:39.755843+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (36, '20200522152607_rename_meta_column.up.sql', 1, '2021-05-21 17:09:39.759962+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (37, '20200605153934_add_deployed_index.up.sql', 1, '2021-05-21 17:09:39.769936+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (38, '20200628093645_add_record_pb.up.sql', 1, '2021-05-21 17:09:39.774193+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (39, '20200628101742_add_processing_queue.up.sql', 1, '2021-05-21 17:09:39.778439+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (40, '20200628175148_add_aggregated_schema.up.sql', 1, '2021-05-21 17:09:39.782439+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (41, '20200630112726_drop_status_json.up.sql', 1, '2021-05-21 17:09:39.7915+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (42, '20200630115127_add_location_details.up.sql', 1, '2021-05-21 17:09:39.795605+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (43, '20200630202515_add_station_notes.up.sql', 1, '2021-05-21 17:09:39.799682+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (44, '20200710095008_multiple_note_media.up.sql', 1, '2021-05-21 17:09:39.804215+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (45, '20200716213200_add_notes_version.up.sql', 1, '2021-05-21 17:09:39.808624+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (46, '20200720112858_add_user_created.up.sql', 1, '2021-05-21 17:09:39.812778+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (47, '20200720113035_drop_project_slug.up.sql', 1, '2021-05-21 17:09:39.81974+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (48, '20200721120957_drop_old_notes.up.sql', 1, '2021-05-21 17:09:39.823456+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (49, '20200721122406_rename_created.up.sql', 1, '2021-05-21 17:09:39.827398+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (50, '20200725165458_add_number_aggregated.up.sql', 1, '2021-05-21 17:09:39.831328+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (51, '20200803094259_add_station_privacy.up.sql', 1, '2021-05-21 17:09:39.835198+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (52, '20200804083846_add_counties.up.sql', 1, '2021-05-21 17:09:39.839718+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (53, '20200808175614_add_que_jobs.up.sql', 1, '2021-05-21 17:09:39.843363+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (54, '20200810180328_add_exports.up.sql', 1, '2021-05-21 17:09:39.846894+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (55, '20200812081349_add_exports_with_job.up.sql', 1, '2021-05-21 17:09:39.850342+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (56, '20200812133750_add_export_size.up.sql', 1, '2021-05-21 17:09:39.854236+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (57, '20200824122101_fix_export_progress_columns.up.sql', 1, '2021-05-21 17:09:39.857908+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (58, '20200920143420_recreate_activity_view.up.sql', 1, '2021-05-21 17:09:39.861562+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (59, '20201009143536_add_discussion_schema.up.sql', 1, '2021-05-21 17:09:39.86507+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (60, '20201206140344_add_record_range_meta.up.sql', 1, '2021-05-21 17:09:39.874407+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (61, '20210203160935_station_synced_at.up.sql', 1, '2021-05-21 17:09:39.878408+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (62, '20210306122946_add_firmware_logical_address.up.sql', 1, '2021-05-21 17:09:39.882299+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (63, '20210312092507_add_firmware_hidden.up.sql', 1, '2021-05-21 17:09:39.886148+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (64, '20210319102605_add_10s_aggregate.up.sql', 1, '2021-05-21 17:09:39.890127+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (65, '20210331102606_add_project_bounds.up.sql', 1, '2021-05-21 17:09:39.894074+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (66, '20210413193800_add_firmware_version.up.sql', 1, '2021-05-21 17:09:39.898628+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (67, '20210415132305_add_sensor_module_id.up.sql', 1, '2021-05-21 17:09:39.902541+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (68, '20210419111558_enable_module_id_aggregates.up.sql', 1, '2021-05-21 17:09:39.906574+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (69, '20210427104916_add_firmware_regex.up.sql', 1, '2021-05-21 17:09:39.910707+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (70, '20210514202826_add_ttn_messages.up.sql', 1, '2021-05-21 17:09:39.914458+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (71, '20210520162844_add_community_ranking.up.sql', 1, '2021-05-21 17:09:39.918145+00');
INSERT INTO fieldkit.migrations (id, name, batch, completed_at) VALUES (72, '20210521113655_migrate_migrations.up.sql', 1, NULL);


--
-- Name: migrations_id_seq; Type: SEQUENCE SET; Schema: fieldkit; Owner: fieldkit
--

SELECT pg_catalog.setval('fieldkit.migrations_id_seq', 72, true);


--
-- Name: migrations migrations_pkey; Type: CONSTRAINT; Schema: fieldkit; Owner: fieldkit
--

ALTER TABLE ONLY fieldkit.migrations
    ADD CONSTRAINT migrations_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

