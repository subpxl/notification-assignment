package models

import "time"

type MessageStatus string

const (
	MessageStatusPending MessageStatus = "pending"
	MessageStatusSent    MessageStatus = "sent"
	MessageStatusFailed  MessageStatus = "failed"
)

type Message struct {
	ID        int64         `json:"id"`
	Recipient string        `json:"recipient"`
	Content   string        `json:"content"`
	MessageID string        `json:"message_id,omitempty"`
	Status    MessageStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	SentAt    *time.Time    `json:"sent_at,omitempty"`
}

type MessageRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type WebhookResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}
