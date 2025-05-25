// stats_repository.go
package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"statistics-service/internal/domain/entities"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseStatsRepository struct {
	DB *sql.DB
}

func NewClickHouseStatsRepository(dsn string) (*ClickHouseStatsRepository, error) {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{dsn},
		Auth: clickhouse.Auth{
			Database: "stats",
			Username: "admin",
			Password: "admin",
		},
	})

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	return &ClickHouseStatsRepository{DB: conn}, nil
}

func (r *ClickHouseStatsRepository) GetPostStats(ctx context.Context, postID string) (*entities.PostStats, error) {
	fmt.Println("test1")
	query := `
		SELECT 
			countIf(event_type = 'view') as views,
			countIf(event_type = 'like') as likes,
			countIf(event_type = 'comment') as comments
		FROM post_events
		WHERE post_id = ?
	`

	var stats entities.PostStats
	stats.PostID = postID
	err := r.DB.QueryRowContext(ctx, query, postID).Scan(
		&stats.Views,
		&stats.Likes,
		&stats.Comments,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get post stats: %w", err)
	}
	return &stats, nil
}

func (r *ClickHouseStatsRepository) GetPostViewsDynamics(ctx context.Context, postID string, from, to time.Time) ([]entities.DailyStat, error) {
	return r.getPostDynamics(ctx, postID, "view", from, to)
}

func (r *ClickHouseStatsRepository) GetPostLikesDynamics(ctx context.Context, postID string, from, to time.Time) ([]entities.DailyStat, error) {
	return r.getPostDynamics(ctx, postID, "like", from, to)
}

func (r *ClickHouseStatsRepository) GetPostCommentsDynamics(ctx context.Context, postID string, from, to time.Time) ([]entities.DailyStat, error) {
	return r.getPostDynamics(ctx, postID, "comment", from, to)
}

func (r *ClickHouseStatsRepository) getPostDynamics(ctx context.Context, postID string, eventType string, from, to time.Time) ([]entities.DailyStat, error) {
	query := `
		SELECT 
			toDate(event_time) as date,
			count() as count
		FROM post_events
		WHERE post_id = ? AND event_type = ? AND event_time >= ? AND event_time <= ?
		GROUP BY date
		ORDER BY date
	`

	rows, err := r.DB.QueryContext(ctx, query, postID, eventType, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to query dynamics: %w", err)
	}
	defer rows.Close()

	var stats []entities.DailyStat
	for rows.Next() {
		var stat entities.DailyStat
		var date time.Time

		if err := rows.Scan(&date, &stat.Count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		stat.Date = date
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return stats, nil
}

func (r *ClickHouseStatsRepository) GetTopPosts(ctx context.Context, statType entities.StatType, limit int) ([]entities.TopPost, error) {
	var eventType string
	switch statType {
	case entities.StatTypeViews:
		eventType = "view"
	case entities.StatTypeLikes:
		eventType = "like"
	case entities.StatTypeComments:
		eventType = "comment"
	default:
		return nil, fmt.Errorf("invalid stat type: %s", statType)
	}

	query := `
		SELECT 
			post_id,
			count() as count
		FROM post_events
		WHERE event_type = ?
		GROUP BY post_id
		ORDER BY count DESC
		LIMIT ?
	`

	rows, err := r.DB.QueryContext(ctx, query, eventType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top posts: %w", err)
	}
	defer rows.Close()

	var topPosts []entities.TopPost
	for rows.Next() {
		var post entities.TopPost
		if err := rows.Scan(&post.PostID, &post.Count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		topPosts = append(topPosts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return topPosts, nil
}

func (r *ClickHouseStatsRepository) GetTopUsers(ctx context.Context, statType entities.StatType, limit int) ([]entities.TopUser, error) {
	var eventType string
	switch statType {
	case entities.StatTypeViews:
		eventType = "view"
	case entities.StatTypeLikes:
		eventType = "like"
	case entities.StatTypeComments:
		eventType = "comment"
	default:
		return nil, fmt.Errorf("invalid stat type: %s", statType)
	}

	query := `
		SELECT 
			user_id,
			count() as count
		FROM post_events
		WHERE event_type = ?
		GROUP BY user_id
		ORDER BY count DESC
		LIMIT ?
	`

	rows, err := r.DB.QueryContext(ctx, query, eventType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top users: %w", err)
	}
	defer rows.Close()

	var topUsers []entities.TopUser
	for rows.Next() {
		var user entities.TopUser
		if err := rows.Scan(&user.UserID, &user.Count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		topUsers = append(topUsers, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return topUsers, nil
}

func (r *ClickHouseStatsRepository) RecordView(ctx context.Context, postID, userID string) error {
	return r.recordEvent(ctx, postID, userID, "view")
}

func (r *ClickHouseStatsRepository) RecordLike(ctx context.Context, postID, userID string) error {
	return r.recordEvent(ctx, postID, userID, "like")
}

func (r *ClickHouseStatsRepository) RecordComment(ctx context.Context, postID, userID string) error {
	return r.recordEvent(ctx, postID, userID, "comment")
}

func (r *ClickHouseStatsRepository) recordEvent(ctx context.Context, postID, userID string, eventType string) error {
	query := `
		INSERT INTO post_events (event_time, post_id, user_id, event_type)
		VALUES (?, ?, ?, ?)
	`
	fmt.Println("test")
	_, err := r.DB.ExecContext(ctx, query, time.Now(), postID, userID, eventType)
	if err != nil {
		return fmt.Errorf("failed to record event: %w", err)
	}

	return nil
}
