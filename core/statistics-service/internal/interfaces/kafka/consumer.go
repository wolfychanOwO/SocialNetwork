package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"statistics-service/internal/domain/repositories"

	"github.com/segmentio/kafka-go"
)

type Event struct {
	Type   string `json:"event"`
	PostID string `json:"post_id"`
	UserID string `json:"user_id"`
	Time   string `json:"time"`
}

type Consumer struct {
	reader *kafka.Reader
	repo   repositories.StatsRepository
}

func NewConsumer(brokers []string, topic string, groupID string, repo repositories.StatsRepository) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
		repo: repo,
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			continue
		}
		fmt.Println("test_consumer")
		fmt.Println(event)
		fmt.Println("type:", event.Type)
		switch event.Type {
		case "view":
			c.repo.RecordView(ctx, event.PostID, event.UserID)
		case "like":
			fmt.Println("record_start")
			c.repo.RecordLike(ctx, event.PostID, event.UserID)
			fmt.Println("record_done")
		case "comment":
			c.repo.RecordComment(ctx, event.PostID, event.UserID)
		}
	}
}
