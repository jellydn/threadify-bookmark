package bookmark

import (
	"context"
	"fmt"
	"time"

	"encore.dev/storage/sqldb"
	"encore.dev/types/uuid"
)

// insert inserts a URL into the database.
func insert(ctx context.Context, id uuid.UUID, url string, owner string, note string) error {
	createdAt := time.Now().UTC()
	_, err := sqldb.Exec(ctx, `
		INSERT INTO bookmark (id, url, owner, note, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, id, url, owner, note, createdAt)
	return err
}

type Bookmark struct {
	ID         uuid.UUID // unique ID
	OWNER      string    // owner of the bookmark
	URL        string    // url of the bookmark
	NOTE       string    // optional note
	CREATED_AT time.Time // date time of creation
}

type BookmarkParams struct {
	URL         string // the URL to bookmark
	Owner       string // the owner of the bookmark
	Description string // optional description of the bookmark
}

// Bookmark a URL.
//
//encore:api public method=POST path=/bookmark
func CreateBookmark(ctx context.Context, p *BookmarkParams) (*Bookmark, error) {
	id, err := uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", id)
	if err != nil {
		return nil, err
	}

	if err := insert(ctx, id, p.URL, p.Owner, p.Description); err != nil {
		return nil, err
	}

	return &Bookmark{ID: id, URL: p.URL, OWNER: p.Owner, NOTE: p.Description}, nil
}

type GetResponse struct {
	Bookmarks []*Bookmark
}

// Get retrieves the all bookmark URLs for the owner id.
// encore:api public method=GET path=/bookmark/:id
func GetBookmarks(ctx context.Context, id string) (*GetResponse, error) {
	rows, err := sqldb.Query(ctx, `
		SELECT id, url, owner, note, created_at
		FROM bookmark
		WHERE owner = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []*Bookmark
	for rows.Next() {
		var b Bookmark
		if err := rows.Scan(&b.ID, &b.URL, &b.OWNER, &b.NOTE, &b.CREATED_AT); err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, &b)
	}

	return &GetResponse{Bookmarks: bookmarks}, nil
}
