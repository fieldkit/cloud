/* I checked and we don't have any invites right now so, hehe I'm
 * going NOT NULL and gonna get this in fast! */
ALTER TABLE fieldkit.project_invite ADD COLUMN role_id INTEGER NOT NULL;
