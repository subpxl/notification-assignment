package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"insider-assignment/internal/models"
	"time"
)

type MessageRepoPsql struct {
	db *sql.DB
}

func New(host, port, user, password, dbname string) (*MessageRepoPsql, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	r := &MessageRepoPsql{db: db}
	r.initSchema()
	return r, nil

}

func (r *MessageRepoPsql) initSchema() error {
	r.db.Exec(`
	create table if not exists messages (
		id serial primary key,
		recipient varchar(255) not null,
		content text not null,
		message_id varchar(255),
		status varchar(50) not null,
		created_at timestamp not null default now(),
		sent_at timestamp
	);
	`)
	return nil
}

func (r *MessageRepoPsql) Close() error {
	return r.db.Close()
}

func (r *MessageRepoPsql) GetUnset(ctx context.Context, limit int) ([]models.Message, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, recipient, content, status FROM messages where status='pending' order by created_at  limit $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.Recipient, &msg.Content, &msg.Status); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *MessageRepoPsql) MarkSent(ctx context.Context, id int64, msgID string, sentAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`update messages set status='sent', message_id=$1, sent_at=$2 where id=$3`,
		msgID, sentAt, id)
	return err
}

func (r *MessageRepoPsql) GetSent(ctx context.Context, limiy, offdet int) ([]models.Message, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, recipient, content, message_id, sent_at from messages where status='sent' order by sent_at DESC limit $1 offset $2`, limiy, offdet)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.Recipient, &message.Content, &message.MessageID, &message.SentAt); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
