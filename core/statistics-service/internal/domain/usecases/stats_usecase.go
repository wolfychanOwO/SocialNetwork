// stats_usecase.go
package usecases

import (
	"context"
	"statistics-service/internal/domain/entities"
	"statistics-service/internal/domain/repositories"
	"time"
)

type StatsUseCase interface {
	GetPostStats(ctx context.Context, postID string) (*entities.PostStats, error)
	GetPostViewsDynamics(ctx context.Context, postID string, from, to time.Time) ([]entities.DailyStat, error)
	GetPostLikesDynamics(ctx context.Context, postID string, from, to time.Time) ([]entities.DailyStat, error)
	GetPostCommentsDynamics(ctx context.Context, postID string, from, to time.Time) ([]entities.DailyStat, error)
	GetTopPosts(ctx context.Context, statType entities.StatType, limit int) ([]entities.TopPost, error)
	GetTopUsers(ctx context.Context, statType entities.StatType, limit int) ([]entities.TopUser, error)
}

type statsUseCase struct {
	statsRepo repositories.StatsRepository
}

func NewStatsUseCase(statsRepo repositories.StatsRepository) StatsUseCase {
	return &statsUseCase{statsRepo: statsRepo}
}

func (uc *statsUseCase) GetPostStats(ctx context.Context, postID string) (*entities.PostStats, error) {
	return uc.statsRepo.GetPostStats(ctx, postID)
}

func (uc *statsUseCase) GetPostViewsDynamics(ctx context.Context, postID string, from, to time.Time) ([]entities.DailyStat, error) {
	return uc.statsRepo.GetPostViewsDynamics(ctx, postID, from, to)
}

func (uc *statsUseCase) GetPostLikesDynamics(ctx context.Context, postID string, from, to time.Time) ([]entities.DailyStat, error) {
	return uc.statsRepo.GetPostLikesDynamics(ctx, postID, from, to)
}

func (uc *statsUseCase) GetPostCommentsDynamics(ctx context.Context, postID string, from, to time.Time) ([]entities.DailyStat, error) {
	return uc.statsRepo.GetPostCommentsDynamics(ctx, postID, from, to)
}

func (uc *statsUseCase) GetTopPosts(ctx context.Context, statType entities.StatType, limit int) ([]entities.TopPost, error) {
	return uc.statsRepo.GetTopPosts(ctx, statType, limit)
}

func (uc *statsUseCase) GetTopUsers(ctx context.Context, statType entities.StatType, limit int) ([]entities.TopUser, error) {
	return uc.statsRepo.GetTopUsers(ctx, statType, limit)
}
