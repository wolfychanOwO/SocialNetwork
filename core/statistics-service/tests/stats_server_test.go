package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"statistics-service/internal/domain/entities"
	"statistics-service/internal/domain/usecases/mocks"
	server "statistics-service/internal/interfaces/grpc"
	"statistics-service/proto"
)

func TestStatsServer_GetPostStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockStatsUseCase(ctrl)
	server := server.NewStatsServer(mockUC)

	tests := []struct {
		name    string
		setup   func()
		req     *proto.GetPostStatsRequest
		want    *proto.PostStatsResponse
		wantErr bool
	}{
		{
			name: "success",
			setup: func() {
				mockUC.EXPECT().GetPostStats(gomock.Any(), "post123").Return(
					&entities.PostStats{Views: 100, Likes: 50, Comments: 10}, nil)
			},
			req:     &proto.GetPostStatsRequest{PostId: "post123"},
			want:    &proto.PostStatsResponse{Views: 100, Likes: 50, Comments: 10},
			wantErr: false,
		},
		{
			name: "not found",
			setup: func() {
				mockUC.EXPECT().GetPostStats(gomock.Any(), "post456").Return(
					nil, errors.New("not found"))
			},
			req:     &proto.GetPostStatsRequest{PostId: "post456"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			resp, err := server.GetPostStats(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, resp)
			}
		})
	}
}

func TestStatsServer_GetPostViewsDynamics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockStatsUseCase(ctrl)
	server := server.NewStatsServer(mockUC)

	from := "2023-01-01"
	to := "2023-01-07"
	fromTime, _ := time.Parse("2006-01-02", from)
	toTime, _ := time.Parse("2006-01-02", to)

	tests := []struct {
		name    string
		setup   func()
		req     *proto.GetDynamicsRequest
		want    *proto.DynamicsResponse
		wantErr bool
	}{
		{
			name: "success",
			setup: func() {
				mockUC.EXPECT().GetPostViewsDynamics(gomock.Any(), "post123", fromTime, toTime).Return(
					[]entities.DailyStat{
						{Date: fromTime, Count: 10},
						{Date: fromTime.Add(24 * time.Hour), Count: 20},
					}, nil)
			},
			req: &proto.GetDynamicsRequest{
				PostId:   "post123",
				FromDate: from,
				ToDate:   to,
			},
			want: &proto.DynamicsResponse{
				Stats: []*proto.DailyStat{
					{Date: from, Count: 10},
					{Date: "2023-01-02", Count: 20},
				},
			},
			wantErr: false,
		},
		{
			name:  "invalid date format",
			setup: func() {},
			req: &proto.GetDynamicsRequest{
				PostId:   "post123",
				FromDate: "invalid-date",
				ToDate:   to,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			resp, err := server.GetPostViewsDynamics(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, resp)
			}
		})
	}
}

func TestStatsServer_GetTopPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockStatsUseCase(ctrl)
	server := server.NewStatsServer(mockUC)

	tests := []struct {
		name    string
		setup   func()
		req     *proto.GetTopRequest
		want    *proto.TopPostsResponse
		wantErr bool
	}{
		{
			name: "success - views",
			setup: func() {
				mockUC.EXPECT().GetTopPosts(gomock.Any(), entities.StatTypeViews, 5).Return(
					[]entities.TopPost{
						{PostID: "post1", Count: 100},
						{PostID: "post2", Count: 90},
					}, nil)
			},
			req: &proto.GetTopRequest{
				StatType: "views",
				Limit:    5,
			},
			want: &proto.TopPostsResponse{
				Posts: []*proto.TopPost{
					{PostId: "post1", Count: 100},
					{PostId: "post2", Count: 90},
				},
			},
			wantErr: false,
		},
		{
			name:  "invalid stat type",
			setup: func() {},
			req: &proto.GetTopRequest{
				StatType: "invalid",
				Limit:    5,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			resp, err := server.GetTopPosts(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, resp)
			}
		})
	}
}

func TestStatsServer_GetTopUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockStatsUseCase(ctrl)
	server := server.NewStatsServer(mockUC)

	tests := []struct {
		name    string
		setup   func()
		req     *proto.GetTopRequest
		want    *proto.TopUsersResponse
		wantErr bool
	}{
		{
			name: "success - comments",
			setup: func() {
				mockUC.EXPECT().GetTopUsers(gomock.Any(), entities.StatTypeComments, 3).Return(
					[]entities.TopUser{
						{UserID: "user1", Count: 50},
						{UserID: "user2", Count: 40},
					}, nil)
			},
			req: &proto.GetTopRequest{
				StatType: "comments",
				Limit:    3,
			},
			want: &proto.TopUsersResponse{
				Users: []*proto.TopUser{
					{UserId: "user1", Count: 50},
					{UserId: "user2", Count: 40},
				},
			},
			wantErr: false,
		},
		{
			name: "default limit",
			setup: func() {
				mockUC.EXPECT().GetTopUsers(gomock.Any(), entities.StatTypeLikes, 10).Return(
					[]entities.TopUser{}, nil)
			},
			req: &proto.GetTopRequest{
				StatType: "likes",
				Limit:    0,
			},
			want:    &proto.TopUsersResponse{Users: []*proto.TopUser{}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			resp, err := server.GetTopUsers(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, resp)
			}
		})
	}
}
