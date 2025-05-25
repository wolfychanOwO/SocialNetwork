package entities

import "time"

type StatType string

const (
	StatTypeViews    StatType = "views"
	StatTypeLikes    StatType = "likes"
	StatTypeComments StatType = "comments"
)

type PostStats struct {
	PostID   string
	Views    int64
	Likes    int64
	Comments int64
}

type DailyStat struct {
	Date  time.Time
	Count int64
}

type TopPost struct {
	PostID string
	Count  int64
}

type TopUser struct {
	UserID string
	Count  int64
}
