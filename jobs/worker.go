// Package jobs provides background job processing with worker pools,
// retry strategies, and cron scheduling.
package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// Job represents a unit of work.
type Job interface {
	Name() string
	Execute(ctx context.Context) error
}

// JobFunc wraps a function as a Job.
type JobFunc struct {
	JobName string
	Fn      func(ctx context.Context) error
}

func (j *JobFunc) Name() string                      { return j.JobName }
func (j *JobFunc) Execute(ctx context.Context) error { return j.Fn(ctx) }

// RetryStrategy defines how failed jobs are retried.
type RetryStrategy struct {
	MaxRetries int
	Delay      time.Duration
	Backoff    float64 // multiplier per retry (e.g., 2.0 for exponential)
}

// DefaultRetryStrategy returns a sensible default.
func DefaultRetryStrategy() RetryStrategy {
	return RetryStrategy{
		MaxRetries: 3,
		Delay:      time.Second,
		Backoff:    2.0,
	}
}

// QueuedJob wraps a job with metadata.
type QueuedJob struct {
	Job       Job
	Retries   int
	Strategy  RetryStrategy
	CreatedAt time.Time
}

// WorkerPool manages background workers.
type WorkerPool struct {
	mu        sync.Mutex
	queue     chan *QueuedJob
	workers   int
	logger    *slog.Logger
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	scheduler *cron.Cron
	running   bool
}

// NewWorkerPool creates a new worker pool.
func NewWorkerPool(workers int, queueSize int, logger *slog.Logger) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		queue:     make(chan *QueuedJob, queueSize),
		workers:   workers,
		logger:    logger,
		ctx:       ctx,
		cancel:    cancel,
		scheduler: cron.New(cron.WithSeconds()),
	}
}

// Start begins processing jobs.
func (wp *WorkerPool) Start() {
	wp.mu.Lock()
	if wp.running {
		wp.mu.Unlock()
		return
	}
	wp.running = true
	wp.mu.Unlock()

	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	wp.scheduler.Start()
	wp.logger.Info("worker pool started", "workers", wp.workers)
}

// Stop gracefully shuts down the worker pool.
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	if !wp.running {
		wp.mu.Unlock()
		return
	}
	wp.running = false
	wp.mu.Unlock()

	wp.cancel()
	wp.scheduler.Stop()
	close(wp.queue)
	wp.wg.Wait()
	wp.logger.Info("worker pool stopped")
}

// Enqueue adds a job to the queue.
func (wp *WorkerPool) Enqueue(job Job, strategy ...RetryStrategy) error {
	s := DefaultRetryStrategy()
	if len(strategy) > 0 {
		s = strategy[0]
	}

	qj := &QueuedJob{
		Job:       job,
		Strategy:  s,
		CreatedAt: time.Now(),
	}

	select {
	case wp.queue <- qj:
		wp.logger.Debug("job enqueued", "job", job.Name())
		return nil
	default:
		return fmt.Errorf("jobs: queue is full, cannot enqueue %s", job.Name())
	}
}

// EnqueueFunc enqueues a function as a job.
func (wp *WorkerPool) EnqueueFunc(name string, fn func(ctx context.Context) error) error {
	return wp.Enqueue(&JobFunc{JobName: name, Fn: fn})
}

// Schedule adds a cron job.
func (wp *WorkerPool) Schedule(spec string, job Job) error {
	_, err := wp.scheduler.AddFunc(spec, func() {
		if err := wp.Enqueue(job); err != nil {
			wp.logger.Error("failed to enqueue scheduled job",
				"job", job.Name(),
				"error", err.Error(),
			)
		}
	})
	if err != nil {
		return fmt.Errorf("jobs: invalid cron spec %q: %w", spec, err)
	}
	wp.logger.Info("job scheduled", "job", job.Name(), "spec", spec)
	return nil
}

// ScheduleFunc schedules a function as a cron job.
func (wp *WorkerPool) ScheduleFunc(spec, name string, fn func(ctx context.Context) error) error {
	return wp.Schedule(spec, &JobFunc{JobName: name, Fn: fn})
}

// worker processes jobs from the queue.
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for qj := range wp.queue {
		select {
		case <-wp.ctx.Done():
			return
		default:
			wp.processJob(id, qj)
		}
	}
}

func (wp *WorkerPool) processJob(workerID int, qj *QueuedJob) {
	logger := wp.logger.With("worker", workerID, "job", qj.Job.Name())

	for attempt := 0; attempt <= qj.Strategy.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(float64(qj.Strategy.Delay) * math.Pow(qj.Strategy.Backoff, float64(attempt-1)))
			logger.Info("retrying job", "attempt", attempt, "delay", delay)

			select {
			case <-wp.ctx.Done():
				return
			case <-time.After(delay):
			}
		}

		err := qj.Job.Execute(wp.ctx)
		if err == nil {
			logger.Debug("job completed")
			return
		}

		logger.Error("job failed",
			"attempt", attempt+1,
			"max_retries", qj.Strategy.MaxRetries,
			"error", err.Error(),
		)
	}

	logger.Error("job exhausted all retries", "retries", qj.Strategy.MaxRetries)
}

// OnInit implements core.OnInit.
func (wp *WorkerPool) OnInit() error {
	wp.Start()
	return nil
}

// OnShutdown implements core.OnShutdown.
func (wp *WorkerPool) OnShutdown() error {
	wp.Stop()
	return nil
}
