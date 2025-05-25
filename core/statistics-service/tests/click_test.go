package tests

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"statistics-service/internal/domain/entities"
	clickhouse "statistics-service/internal/infrastructure/db"

	test_clickhouse "github.com/testcontainers/testcontainers-go/modules/clickhouse"
)

var repo *clickhouse.ClickHouseStatsRepository

func TestMain(m *testing.M) {
	ctx := context.Background()

	user := "admin"
	password := "admin"
	dbname := "stats"

	clickHouseContainer, err := test_clickhouse.Run(ctx,
		"clickhouse/clickhouse-server:23.3.8.21-alpine",
		test_clickhouse.WithUsername(user),
		test_clickhouse.WithPassword(password),
		test_clickhouse.WithDatabase(dbname),
		test_clickhouse.WithInitScripts(filepath.Join("clickhouse", "init.sql")),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	dsn, err := clickHouseContainer.ConnectionHost(ctx)
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	repo, err = clickhouse.NewClickHouseStatsRepository(dsn)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	code := m.Run()

	if err := clickHouseContainer.Terminate(ctx); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}

	os.Exit(code)
}

func TestGetPostStats(t *testing.T) {
	ctx := context.Background()
	postID := "test-post-1"

	err := repo.RecordView(ctx, postID, "user1")
	if err != nil {
		t.Fatalf("Failed to record view: %v", err)
	}

	err = repo.RecordLike(ctx, postID, "user2")
	if err != nil {
		t.Fatalf("Failed to record like: %v", err)
	}

	err = repo.RecordComment(ctx, postID, "user3")
	if err != nil {
		t.Fatalf("Failed to record comment: %v", err)
	}

	stats, err := repo.GetPostStats(ctx, postID)
	if err != nil {
		t.Fatalf("GetPostStats failed: %v", err)
	}

	if stats.PostID != postID {
		t.Errorf("Expected post ID %s, got %s", postID, stats.PostID)
	}

	if stats.Views != 1 {
		t.Errorf("Expected 1 view, got %d", stats.Views)
	}

	if stats.Likes != 1 {
		t.Errorf("Expected 1 like, got %d", stats.Likes)
	}

	if stats.Comments != 1 {
		t.Errorf("Expected 1 comment, got %d", stats.Comments)
	}
}

func TestGetPostViewsDynamics(t *testing.T) {
	ctx := context.Background()
	postID := "test-post-dyn"

	testDates := []time.Time{
		time.Now().AddDate(0, 0, -2),
		time.Now().AddDate(0, 0, -1),
		time.Now(),
	}

	for _, date := range testDates {
		_, err := repo.DB.ExecContext(ctx, `
			INSERT INTO post_events (event_time, post_id, user_id, event_type)
			VALUES (?, ?, ?, 'view')
		`, date, postID, "test-user")
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}

	from := time.Now().AddDate(0, 0, -3)
	to := time.Now()
	dynamics, err := repo.GetPostViewsDynamics(ctx, postID, from, to)
	if err != nil {
		t.Fatalf("GetPostViewsDynamics failed: %v", err)
	}

	if len(dynamics) != 3 {
		t.Fatalf("Expected 3 days of dynamics, got %d", len(dynamics))
	}

	for _, stat := range dynamics {
		if stat.Count != 1 {
			t.Errorf("Expected count 1 for date %v, got %d", stat.Date, stat.Count)
		}
	}
}

func TestGetTopPosts(t *testing.T) {
	ctx := context.Background()

	posts := []struct {
		id    string
		views int
	}{
		{"post-top-1", 5},
		{"post-top-2", 3},
		{"post-top-3", 1},
	}

	for _, p := range posts {
		for i := 0; i < p.views; i++ {
			err := repo.RecordView(ctx, p.id, "user")
			if err != nil {
				t.Fatalf("Failed to record view: %v", err)
			}
		}
	}

	topPosts, err := repo.GetTopPosts(ctx, entities.StatTypeViews, 2)
	if err != nil {
		t.Fatalf("GetTopPosts failed: %v", err)
	}

	if len(topPosts) != 2 {
		t.Fatalf("Expected 2 top posts, got %d", len(topPosts))
	}

	if topPosts[0].PostID != "post-top-1" || topPosts[0].Count != 5 {
		t.Errorf("Expected top post to be post-top-1 with 5 views")
	}

	if topPosts[1].PostID != "post-top-2" || topPosts[1].Count != 3 {
		t.Errorf("Expected second post to be post-top-2 with 3 views")
	}
}

func TestRecordEvent(t *testing.T) {
	ctx := context.Background()
	postID := "test-record-event"
	userID := "test-user"

	tests := []struct {
		name      string
		eventType string
		method    func(context.Context, string, string) error
	}{
		{"RecordView", "view", repo.RecordView},
		{"RecordLike", "like", repo.RecordLike},
		{"RecordComment", "comment", repo.RecordComment},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.method(ctx, postID, userID)
			if err != nil {
				t.Fatalf("%s failed: %v", tt.name, err)
			}

			var count int
			err = repo.DB.QueryRowContext(ctx, `
				SELECT count()
				FROM post_events
				WHERE post_id = ? AND user_id = ? AND event_type = ?
			`, postID, userID, tt.eventType).Scan(&count)
			if err != nil {
				t.Fatalf("Failed to verify event: %v", err)
			}

			if count != 1 {
				t.Errorf("Expected 1 event, got %d", count)
			}
		})
	}
}

func TestGetTopUsers(t *testing.T) {
	ctx := context.Background()
	postID := "test-post-top-users"

	users := []struct {
		id    string
		views int
	}{
		{"user1", 5},
		{"user2", 3},
		{"user3", 1},
	}

	for _, user := range users {
		for i := 0; i < user.views; i++ {
			err := repo.RecordView(ctx, postID, user.id)
			if err != nil {
				t.Fatalf("Failed to record view for %s: %v", user.id, err)
			}
		}
	}

	topUsers, err := repo.GetTopUsers(ctx, entities.StatTypeViews, 2)
	if err != nil {
		t.Fatalf("GetTopUsers failed: %v", err)
	}

	if len(topUsers) != 2 {
		t.Fatalf("Expected 2 top users, got %d", len(topUsers))
	}

	if topUsers[0].UserID != "user1" || topUsers[0].Count != 5 {
		t.Errorf("Expected top user to be user1 with 5 views, got %s with %d views",
			topUsers[0].UserID, topUsers[0].Count)
	}

	if topUsers[1].UserID != "user2" || topUsers[1].Count != 3 {
		t.Errorf("Expected second user to be user2 with 3 views, got %s with %d views",
			topUsers[1].UserID, topUsers[1].Count)
	}
}
