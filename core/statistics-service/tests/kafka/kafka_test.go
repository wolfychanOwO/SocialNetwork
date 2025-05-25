package tests

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"

	clickhouse "statistics-service/internal/infrastructure/db"
	kafkaConsumer "statistics-service/internal/interfaces/kafka"
)

func TestConsumer_ProcessesLikeEvent(t *testing.T) {
	ctx := context.Background()

	repo, err := clickhouse.NewClickHouseStatsRepository("localhost:9000")
	require.NoError(t, err)

	postID := uuid.New().String()
	groupID := "test-group-" + uuid.New().String()

	consumer := kafkaConsumer.NewConsumer(
		[]string{"localhost:9092"},
		"post_events",
		groupID,
		repo,
	)

	go func() {
		err := consumer.Consume(ctx)
		if err != nil {
			t.Errorf("Consumer error: %v", err)
		}
	}()

	time.Sleep(300 * time.Millisecond)

	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "post_events",
	})
	defer producer.Close()

	event := kafkaConsumer.Event{
		Type:   "like",
		PostID: postID,
		UserID: uuid.New().String(),
	}

	value, err := json.Marshal(event)
	require.NoError(t, err)

	err = producer.WriteMessages(ctx, kafka.Message{Value: value})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		stats, err := repo.GetPostStats(ctx, postID)

		return err == nil && stats.Likes == 1
	}, 10*time.Second, 200*time.Millisecond)

}

func TestConsumer_ProcessesViewEvent(t *testing.T) {
	ctx := context.Background()

	repo, err := clickhouse.NewClickHouseStatsRepository("localhost:9000")
	require.NoError(t, err)

	postID := uuid.New().String()
	groupID := "test-group-" + uuid.New().String()

	consumer := kafkaConsumer.NewConsumer(
		[]string{"localhost:9092"},
		"post_events",
		groupID,
		repo,
	)

	go func() {
		err := consumer.Consume(ctx)
		if err != nil {
			t.Errorf("Consumer error: %v", err)
		}
	}()

	time.Sleep(300 * time.Millisecond)

	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "post_events",
	})
	defer producer.Close()

	event := kafkaConsumer.Event{
		Type:   "view",
		PostID: postID,
		UserID: uuid.New().String(),
	}

	value, err := json.Marshal(event)
	require.NoError(t, err)

	err = producer.WriteMessages(ctx, kafka.Message{Value: value})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		stats, err := repo.GetPostStats(ctx, postID)
		return err == nil && stats.Views == 1
	}, 10*time.Second, 200*time.Millisecond)
}
