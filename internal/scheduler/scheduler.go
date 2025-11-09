package scheduler

import (
	"context"
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
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()
	s.wg.Add(1)
	go s.run()

}

func (s *Scheduler) run() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.task(context.Background())

	for {
		select {
		case <-ticker.C:
			s.task(context.Background())
		case <-s.stop:
			return
		}
	}
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		return
	}
	s.running = false
	close(s.stop)
	s.wg.Wait()
}

func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}
