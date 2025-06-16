// server.go
package grpc

import (
	"context"
	"fmt"
	"statistics-service/internal/domain/entities"
	"statistics-service/internal/domain/usecases"
	"statistics-service/proto"
	"time"
)

type StatsServer struct {
	proto.UnimplementedStatsServiceServer
	uc usecases.StatsUseCase
}

func NewStatsServer(uc usecases.StatsUseCase) *StatsServer {
	return &StatsServer{uc: uc}
}

func (s *StatsServer) GetPostStats(ctx context.Context, req *proto.GetPostStatsRequest) (*proto.PostStatsResponse, error) {
	stats, err := s.uc.GetPostStats(ctx, req.PostId)
	if err != nil {
		return nil, err
	}

	return &proto.PostStatsResponse{
		Views:    stats.Views,
		Likes:    stats.Likes,
		Comments: stats.Comments,
	}, nil
}

func (s *StatsServer) GetPostViewsDynamics(ctx context.Context, req *proto.GetDynamicsRequest) (*proto.DynamicsResponse, error) {
	from, to, err := parseDateRange(req.FromDate, req.ToDate)
	if err != nil {
		return nil, err
	}

	dynamics, err := s.uc.GetPostViewsDynamics(ctx, req.PostId, from, to)
	if err != nil {
		return nil, err
	}

	return convertDailyStatsToProto(dynamics), nil
}

func (s *StatsServer) GetPostLikesDynamics(ctx context.Context, req *proto.GetDynamicsRequest) (*proto.DynamicsResponse, error) {
	from, to, err := parseDateRange(req.FromDate, req.ToDate)
	if err != nil {
		return nil, err
	}

	dynamics, err := s.uc.GetPostLikesDynamics(ctx, req.PostId, from, to)
	if err != nil {
		return nil, err
	}

	return convertDailyStatsToProto(dynamics), nil
}

func (s *StatsServer) GetPostCommentsDynamics(ctx context.Context, req *proto.GetDynamicsRequest) (*proto.DynamicsResponse, error) {
	from, to, err := parseDateRange(req.FromDate, req.ToDate)
	if err != nil {
		return nil, err
	}

	dynamics, err := s.uc.GetPostCommentsDynamics(ctx, req.PostId, from, to)
	if err != nil {
		return nil, err
	}

	return convertDailyStatsToProto(dynamics), nil
}

func (s *StatsServer) GetTopPosts(ctx context.Context, req *proto.GetTopRequest) (*proto.TopPostsResponse, error) {
	statType, err := parseStatType(req.StatType)
	if err != nil {
		return nil, err
	}

	limit := 10
	if req.Limit > 0 {
		limit = int(req.Limit)
	}

	topPosts, err := s.uc.GetTopPosts(ctx, statType, limit)
	if err != nil {
		return nil, err
	}

	response := &proto.TopPostsResponse{
		Posts: make([]*proto.TopPost, len(topPosts)),
	}

	for i, post := range topPosts {
		response.Posts[i] = &proto.TopPost{
			PostId: post.PostID,
			Count:  post.Count,
		}
	}

	return response, nil
}

func (s *StatsServer) GetTopUsers(ctx context.Context, req *proto.GetTopRequest) (*proto.TopUsersResponse, error) {
	statType, err := parseStatType(req.StatType)
	if err != nil {
		return nil, err
	}

	limit := 10
	if req.Limit > 0 {
		limit = int(req.Limit)
	}

	topUsers, err := s.uc.GetTopUsers(ctx, statType, limit)
	if err != nil {
		return nil, err
	}

	response := &proto.TopUsersResponse{
		Users: make([]*proto.TopUser, len(topUsers)),
	}

	for i, user := range topUsers {
		response.Users[i] = &proto.TopUser{
			UserId: user.UserID,
			Count:  user.Count,
		}
	}

	return response, nil
}

func parseDateRange(fromStr, toStr string) (time.Time, time.Time, error) {
	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return from, to, nil
}

func convertDailyStatsToProto(stats []entities.DailyStat) *proto.DynamicsResponse {
	response := &proto.DynamicsResponse{
		Stats: make([]*proto.DailyStat, len(stats)),
	}

	for i, stat := range stats {
		response.Stats[i] = &proto.DailyStat{
			Date:  stat.Date.Format("2006-01-02"),
			Count: stat.Count,
		}
	}

	return response
}

func parseStatType(statType string) (entities.StatType, error) {
	switch statType {
	case "views":
		return entities.StatTypeViews, nil
	case "likes":
		return entities.StatTypeLikes, nil
	case "comments":
		return entities.StatTypeComments, nil
	default:
		return "", fmt.Errorf("invalid stat type: %s", statType)
	}
}
