package worker

import (
	"context"
	"log"
	"techzone/internal/metrics"
)

type NotificationJob struct {
	OrderID int64
	UserID  int64

	Message string
}

type NotificationWorkerPool struct {
	jobs chan NotificationJob

	notificationCreator NotificationCreator
}

type NotificationCreator interface {
	CreateNotification(
		ctx context.Context,
		userID int64,
		orderID int64,
		message string,
	) error
}

func NewNotificationWorkerPool(
	workers int,
	creator NotificationCreator,
) *NotificationWorkerPool {

	pool := &NotificationWorkerPool{
		jobs:                make(chan NotificationJob, 100),
		notificationCreator: creator,
	}

	for i := 0; i < workers; i++ {
		go pool.worker(i)
	}

	return pool
}

func (p *NotificationWorkerPool) worker(
	id int,
) {
	for job := range p.jobs {

		err := p.notificationCreator.CreateNotification(
			context.Background(),
			job.UserID,
			job.OrderID,
			job.Message,
		)

		if err != nil {
			log.Printf(
				"worker=%d notification error=%v",
				id,
				err,
			)
			continue
		}

		metrics.NotificationsCreatedTotal.Inc()
		log.Printf(
			"[worker=%d] sending notification order=%d user=%d",
			id,
			job.OrderID,
			job.UserID,
		)
	}
}

func (p *NotificationWorkerPool) Submit(
	job NotificationJob,
) {
	p.jobs <- job
}
