ALTER TABLE fieldkit.project ADD COLUMN community_ranking INTEGER;
UPDATE fieldkit.project SET community_ranking = 0;

CREATE OR REPLACE FUNCTION fk_update_community_ranking()
RETURNS void
LANGUAGE plpgsql
AS
$$
BEGIN

WITH scores AS (
	SELECT q.project_id, SUM(q.score) AS score FROM
	(
	SELECT project_id, COUNT(project_id) AS score FROM fieldkit.discussion_post WHERE project_id IS NOT NULL GROUP BY project_id
	UNION
	SELECT project_id, COUNT(project_id) AS score FROM fieldkit.project_station GROUP BY project_id
	UNION
	SELECT project_id, COUNT(project_id) AS score FROM fieldkit.project_follower GROUP BY project_id
	UNION
	SELECT project_id, COUNT(project_id) AS score FROM fieldkit.project_user GROUP BY project_id
	) AS q
	GROUP BY q.project_id
)
UPDATE fieldkit.project
SET community_ranking = scores.score
FROM scores
WHERE id = scores.project_id;

END
$$;

SELECT fk_update_community_ranking();

ALTER TABLE fieldkit.project ALTER COLUMN community_ranking SET NOT NULL;
