package scheduler

import (
	"context"
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	interval time.Duration
	task     func(context.Context) error
	running  bool
	mu       sync.RWMutex
	stop     chan struct{}
	wg       sync.WaitGroup
}

func NewScheduler(interval time.Duration, task func(context.Context) error) *Scheduler {
	return &Scheduler{
		interval: interval,
		task:     task,
		stop:     make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		log.Println("[scheduler] already running")

		return
	}
	s.stop = make(chan struct{})
	s.running = true

	s.wg.Add(1)
	go s.run()

}

func (s *Scheduler) run() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	log.Println("[scheduler] started")
	ctx := context.Background()
	s.task(ctx)

	for {
		select {
		case <-ticker.C:
			s.task(ctx)
		case <-s.stop:
			log.Println("[scheduler] stopped")
			return
		}
	}
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		log.Println("[scheduler] stop called but not running")
		return
	}

	s.running = false
	close(s.stop)
	s.mu.Unlock()

	s.wg.Wait()
	log.Println("[scheduler] gracefully stopped")

}
func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}
