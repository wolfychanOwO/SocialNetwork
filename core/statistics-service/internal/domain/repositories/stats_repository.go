package repositories

import (
	"context"
	"statistics-service/internal/domain/entities"
	"time"
)

type StatsRepository interface {
	GetPostStats(ctx context.Context, postID string) (*entities.PostStats, error)
	GetPostViewsDynamics(ctx context.Context, postID string, from, to time.Time) ([]entities.DailyStat, error)
	GetPostLikesDynamics(ctx context.Context, postID string, from, to time.Time) ([]entities.DailyStat, error)
	GetPostCommentsDynamics(ctx context.Context, postID string, from, to time.Time) ([]entities.DailyStat, error)
	GetTopPosts(ctx context.Context, statType entities.StatType, limit int) ([]entities.TopPost, error)
	GetTopUsers(ctx context.Context, statType entities.StatType, limit int) ([]entities.TopUser, error)

	RecordView(ctx context.Context, postID string, userID string) error
	RecordLike(ctx context.Context, postID string, userID string) error
	RecordComment(ctx context.Context, postID string, userID string) error
}
