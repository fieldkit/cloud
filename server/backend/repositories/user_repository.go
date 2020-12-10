package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type UserRepository struct {
	db *sqlxcache.DB
}

func NewUserRepository(db *sqlxcache.DB) (r *UserRepository) {
	return &UserRepository{db: db}
}

func (r *UserRepository) QueryByID(ctx context.Context, id int32) (*data.User, error) {
	user := &data.User{}
	if err := r.db.GetContext(ctx, user, `SELECT * FROM fieldkit.user WHERE id = $1`, id); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Search(ctx context.Context, query string) ([]*data.User, error) {
	likeQuery := "%" + query + "%"
	users := make([]*data.User, 0)
	if err := r.db.SelectContext(ctx, &users, `SELECT * FROM fieldkit.user WHERE LOWER(name) LIKE LOWER($1) OR LOWER(email) LIKE LOWER($1)`, likeQuery); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Delete(outerCtx context.Context, id int32) (err error) {
	return r.db.WithNewTransaction(outerCtx, func(ctx context.Context) error {
		if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.project_follower WHERE follower_id = $1`, id); err != nil {
			return err
		}
		if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.project_invite WHERE user_id = $1`, id); err != nil {
			return err
		}
		if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.project_update WHERE author_id = $1`, id); err != nil {
			return err
		}

		projects := []*data.Project{}
		if err := r.db.SelectContext(ctx, &projects, `SELECT * FROM fieldkit.project WHERE id IN (SELECT project_id FROM fieldkit.project_user WHERE user_id = $1)`, id); err != nil {
			return err
		}

		for _, p := range projects {
			memberships := []*data.ProjectUser{}
			if err := r.db.SelectContext(ctx, &memberships, `SELECT * FROM fieldkit.project_user WHERE project_id = $1`, p.ID); err != nil {
				return err
			}

			// We delete the whole project if they're the only member.
			delete := true
			for _, m := range memberships {
				if m.UserID != id {
					delete = false
				}
			}

			// Eventually we'll delete all of these rows, we only do this project here.
			if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.project_user WHERE user_id = $1 AND project_id = $2`, id, p.ID); err != nil {
				return err
			}

			if delete {
				if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.project_station WHERE project_id = $1`, p.ID); err != nil {
					return err
				}
				if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.project_activity WHERE project_id = $1`, p.ID); err != nil {
					return err
				}
				if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.project WHERE id = $1`, p.ID); err != nil {
					return err
				}
			}
		}

		if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.refresh_token WHERE user_id = $1`, id); err != nil {
			return err
		}
		if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.recovery_token WHERE user_id = $1`, id); err != nil {
			return err
		}
		if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.validation_token WHERE user_id = $1`, id); err != nil {
			return err
		}
		if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.notes_media WHERE user_id = $1`, id); err != nil {
			return err
		}
		if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.notes WHERE user_id = $1`, id); err != nil {
			return err
		}
		if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.station_ingestion WHERE uploader_id = $1`, id); err != nil {
			return err
		}
		if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.project_station WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)`, id); err != nil {
			return err
		}
		if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.station WHERE owner_id = $1`, id); err != nil {
			return err
		}
		if _, err := r.db.ExecContext(ctx, `UPDATE fieldkit.ingestion SET user_id = 2 WHERE user_id = $1`, id); err != nil {
			return err
		}
		if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.user WHERE id = $1`, id); err != nil {
			return err
		}
		return nil
	})
}
