package repository

import (
	"context"
	"insider-assignment/internal/models"
	"time"
)

type MessageRepository interface {
	GetUnsent(ctx context.Context, limit int) ([]models.Message, error)
	MarkSent(ctx context.Context, id int64, sentAt time.Time, messageID string) error
	GetSent(ctx context.Context, limit int, offset int) ([]models.Message, error)
}
