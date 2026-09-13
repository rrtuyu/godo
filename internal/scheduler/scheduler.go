package scheduler

import (
	"context"
	"errors"
	"fmt"
	"time"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Task struct {
	ID			string
	ExecFunc	func() error
	Timeout		time.Duration
}

type Scheduler struct {
	tasks 	chan Task
	workers int
	ctx		context.Context
	stopChan chan struct{}
}

func NewScheduler(ctx context.Context, workers int) *Scheduler{
	return &Scheduler{
		tasks: 		make(chan Task, 10),
		workers: 	workers,
		ctx:		ctx,
	}
}

func (s *Scheduler) Run() {
	for i := 0; i < s.workers; i++{
		go func(workerID int) {
			for {
				select {
				case task := <-s.tasks:
					err := task.ExecFunc()
					if err != nil {
						message := fmt.Sprintf("Task failed: %w", err)
						wailsruntime.LogError(s.ctx, message)
					}
				case <-s.stopChan:
					return
				}
			}
		}(i)
	}
}

func (s *Scheduler) AddTask(task Task) error {
	select {
	case <-s.stopChan:
		return errors.New("Scheduler is stopped, cannot add task")
	default:
		s.tasks <- task
		return nil
	}
}

func (s *Scheduler) AddIntervalTask(task Task, interval time.Duration) {
	go func(){
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
		select {
			case <-ticker.C:
				err := s.AddTask(task)
				if err != nil {
					message := fmt.Sprintf("Failed to edit interval task: %w", err)
					wailsruntime.LogError(s.ctx, message)
				}
			case <-s.stopChan:
				return
			case <-s.ctx.Done():
				return
			}
		}
	}()
}