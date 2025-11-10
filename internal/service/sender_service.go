package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"insider-assignment/internal/cache"
	"insider-assignment/internal/models"
	"insider-assignment/internal/repository"
	"insider-assignment/internal/scheduler"
	"sync"

	"net/http"
	"time"
)

type SenderService struct {
	mu sync.Mutex

	Repo      *repository.MessageRepoPsql
	cache     *cache.Cache
	Scheduler *scheduler.Scheduler
	webHook   string
	authKey   string
	bacthSize int
	client    *http.Client
}

func NewSenderService(repo *repository.MessageRepoPsql, cache *cache.Cache, webhook string, authKey string, interval time.Duration, batchSize int) *SenderService {

	svc := &SenderService{Repo: repo,
		cache:     cache,
		webHook:   webhook,
		authKey:   authKey,
		bacthSize: batchSize,
		client:    &http.Client{Timeout: 10 * time.Second},
	}
	svc.Scheduler = scheduler.NewScheduler(interval, svc.process)
	return svc
}

func (s *SenderService) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Scheduler.IsRunning() {
		return
	}
	s.Scheduler.Start()
}

func (s *SenderService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.Scheduler.IsRunning() {
		return
	}
	s.Scheduler.Stop()
}

func (s *SenderService) IsRunning() bool {
	return s.Scheduler.IsRunning()
}

func (s *SenderService) process(ctx context.Context) error {
	msgs, err := s.Repo.GetUnset(ctx, s.bacthSize)
	if err != nil {
		return err
	}
	for _, msg := range msgs {
		if len(msg.Content) > 500 {
			continue
		}

		if err := s.send(ctx, &msg); err != nil {
			fmt.Printf("failed to send msg %d: %v\n", msg.ID, err)
		}
	}
	return nil
}

func (s *SenderService) send(ctx context.Context, msg *models.Message) error {
	req := models.MessageRequest{
		To:      msg.Recipient,
		Content: msg.Content,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.webHook, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-ins-auth-key", s.authKey)
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var webhookResp models.WebhookResponse
	json.NewDecoder(resp.Body).Decode(&webhookResp)

	sentAt := time.Now()
	s.Repo.MarkSent(ctx, msg.ID, webhookResp.MessageID, sentAt)
	if s.cache != nil {
		s.cache.Set(ctx, msg.ID, webhookResp.MessageID, sentAt)
	}

	fmt.Printf("sent msg %d, messageId: %s\n", msg.ID, webhookResp.MessageID)
	return nil
}

func (s *SenderService) GetSent(ctx context.Context, limit, offset int) ([]models.Message, error) {
	return s.Repo.GetSent(ctx, limit, offset)
}
