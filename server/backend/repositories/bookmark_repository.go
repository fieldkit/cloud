package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/conservify/sqlxcache"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type BookmarkRepository struct {
	db *sqlxcache.DB
}

func NewBookmarkRepository(db *sqlxcache.DB) (rr *BookmarkRepository) {
	return &BookmarkRepository{db: db}
}

type SavedBookmark struct {
	ID           int64     `db:"id"`
	UserID       *int32    `db:"user_id"`
	Token        string    `db:"token"`
	Bookmark     string    `db:"bookmark"`
	CreatedAt    time.Time `db:"created_at"`
	ReferencedAt time.Time `db:"referenced_at"`
}

func (c *BookmarkRepository) AddNew(ctx context.Context, userID *int32, bookmark string) (*SavedBookmark, error) {
	token, err := gonanoid.New()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	saving := &SavedBookmark{
		UserID:       userID,
		Token:        token,
		Bookmark:     bookmark,
		CreatedAt:    now,
		ReferencedAt: now,
	}

	if err := c.db.NamedGetContext(ctx, saving, `
		INSERT INTO fieldkit.bookmarks
		(user_id, token, bookmark, created_at, referenced_at) VALUES (:user_id, :token, :bookmark, :created_at, :referenced_at)
		RETURNING id
		`, saving); err != nil {
		return nil, err
	}

	return saving, nil
}

func (c *BookmarkRepository) Resolve(ctx context.Context, token string) (*SavedBookmark, error) {
	saved := &SavedBookmark{}
	if err := c.db.GetContext(ctx, saved, `
		SELECT id, user_id, token, bookmark, created_at, referenced_at FROM fieldkit.bookmarks WHERE token = $1
		`, token); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return saved, nil
}
