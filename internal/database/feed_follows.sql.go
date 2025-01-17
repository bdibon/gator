// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollows = `-- name: CreateFeedFollows :one
WITH inserted_feed_follows AS (
    INSERT INTO feed_follows (
        id,
        created_at,
        updated_at,
        user_id,
        feed_id
    ) VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING id, created_at, updated_at, user_id, feed_id
)
SELECT
    inserted_feed_follows.id, inserted_feed_follows.created_at, inserted_feed_follows.updated_at, inserted_feed_follows.user_id, inserted_feed_follows.feed_id,
    users.name AS username,
    feeds.name AS feedname
FROM inserted_feed_follows
INNER JOIN users ON users.id = inserted_feed_follows.user_id
INNER JOIN feeds ON feeds.id = inserted_feed_follows.feed_id
`

type CreateFeedFollowsParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

type CreateFeedFollowsRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	Username  string
	Feedname  string
}

func (q *Queries) CreateFeedFollows(ctx context.Context, arg CreateFeedFollowsParams) (CreateFeedFollowsRow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollows,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i CreateFeedFollowsRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.Username,
		&i.Feedname,
	)
	return i, err
}

const deleteFeedFollowsForUser = `-- name: DeleteFeedFollowsForUser :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2
`

type DeleteFeedFollowsForUserParams struct {
	UserID uuid.UUID
	FeedID uuid.UUID
}

func (q *Queries) DeleteFeedFollowsForUser(ctx context.Context, arg DeleteFeedFollowsForUserParams) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollowsForUser, arg.UserID, arg.FeedID)
	return err
}

const getFeedFollowsForUser = `-- name: GetFeedFollowsForUser :many
SELECT users.name AS username, feeds.name AS feedname FROM feed_follows
INNER JOIN users ON users.id = feed_follows.user_id
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE users.id = $1
`

type GetFeedFollowsForUserRow struct {
	Username string
	Feedname string
}

func (q *Queries) GetFeedFollowsForUser(ctx context.Context, id uuid.UUID) ([]GetFeedFollowsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsForUser, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowsForUserRow
	for rows.Next() {
		var i GetFeedFollowsForUserRow
		if err := rows.Scan(&i.Username, &i.Feedname); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
