package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	Version   int64     `json:"version"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

// For the user feed
type PostWithMetaData struct {
	Post
	CommentsCount int64 `json:"comments_count"`
}
type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts(content,title,user_id,tags)
	VALUES ($1, $2, $3, $4) RETURNING	id, created_at,updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
func (s *PostStore) GetById(ctx context.Context, id int64) (*Post, error) {
	query := `
	SELECT id,user_id,title,content,created_at,updated_at,tags FROM posts
	WHERE id = $1  
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	post := Post{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}
func (s *PostStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil

}
func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
	UPDATE posts
	SET title = $1, content = $2, version = version + 1
	WHERE id = $3 and version = $4
	RETURNING version 
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}
func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetaData, error) {
	query := `
	SELECT
    p.id,
    p.title,
    p.tags,
    p.content,
    p.created_at,
    p.version,
    p.user_id,
		u.username,
    COUNT(c.id) as comments_count
FROM
    posts p
    LEFT JOIN comments c ON c.post_id = p.id  
		LEFT JOIN users u on p.user_id = u.id
    JOIN followers f ON f.follower_id = p.user_id  OR p.user_id = $1
		WHERE f.user_id = $1 AND
		(p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
		(p.tags @> $5 OR $5 = '{}') AND
		(CASE WHEN $6::text = '' THEN true ELSE p.created_at >= $6::timestamp END) AND
    (CASE WHEN $7::text = '' THEN true ELSE p.created_at <= $7::timestamp END)
		GROUP BY p.id,u.username
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $2 OFFSET $3
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags), fq.Since, fq.Until)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var feed []PostWithMetaData
	for rows.Next() {
		var post PostWithMetaData
		err := rows.Scan(
			&post.ID,
			&post.Title,
			pq.Array(&post.Tags),
			&post.Content,
			&post.CreatedAt,
			&post.Version,
			&post.UserID,
			&post.User.Username,
			&post.CommentsCount,
		)
		if err != nil {
			return nil, err
		}
		feed = append(feed, post)

	}

	return feed, nil
}
