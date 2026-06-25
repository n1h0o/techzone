package service

import (
	"context"
	"log"
)

type NotificationJob struct {
	OrderID int64
	UserID  int64
}

type NotificationWorkerPool struct {
	jobs chan NotificationJob

	notificationRepo *NotificationService
}

func NewNotificationWorkerPool(
	workers int,
	service *NotificationService,
) *NotificationWorkerPool {

	pool := &NotificationWorkerPool{
		jobs:             make(chan NotificationJob, 100),
		notificationRepo: service,
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

		err := p.notificationRepo.CreateNotification(
			context.Background(),
			job.UserID,
			job.OrderID,
			"Order created successfully",
		)

		if err != nil {
			log.Printf(
				"worker=%d notification error=%v",
				id,
				err,
			)
		}

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
