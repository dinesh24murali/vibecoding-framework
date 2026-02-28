package notifications

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dinesh/vibecoding-framework/backend/internal/queue"
)

type Worker struct {
	queue  *queue.Asynq
	client SMSClient
	pace   time.Duration
}

func NewWorker(jobQueue *queue.Asynq, client SMSClient) *Worker {
	return &Worker{
		queue:  jobQueue,
		client: client,
		pace:   2 * time.Second,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.pace)
	defer ticker.Stop()

	for {
		if err := w.processOnce(ctx); err != nil {
			log.Printf("sms worker cycle error: %v", err)
		}

		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func (w *Worker) processOnce(ctx context.Context) error {
	now := time.Now().UTC()
	jobs, err := w.queue.ListPendingSMS(ctx, 25)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		if !IsRetryDue(job.Status, job.AttemptCount, job.UpdatedAt, now) {
			continue
		}

		sendErr := w.client.Send(ctx, job.PhoneNumber, job.Template)
		if sendErr != nil {
			if err := w.queue.MarkFailed(ctx, job.ID, sendErr.Error()); err != nil {
				log.Printf("sms mark failed error id=%d: %v", job.ID, err)
			}
			continue
		}

		if err := w.queue.MarkSent(ctx, job.ID); err != nil {
			return fmt.Errorf("mark sms sent: %w", err)
		}
	}

	return nil
}
