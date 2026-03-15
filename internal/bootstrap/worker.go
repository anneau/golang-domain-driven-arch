package bootstrap

import (
	"context"
	"fmt"
	"log"

	identityusecase "github.com/hkobori/golang-domain-driven-arch/internal/app/usecase/identity"
	adapterqueue "github.com/hkobori/golang-domain-driven-arch/internal/adapter/input/queue"
	queuehandler "github.com/hkobori/golang-domain-driven-arch/internal/adapter/input/queue/handler"
	sqssubscriber "github.com/hkobori/golang-domain-driven-arch/internal/adapter/input/queue/sqs"
	"github.com/hkobori/golang-domain-driven-arch/internal/adapter/output/auth"
)

type WorkerApp struct {
	worker *adapterqueue.Worker
	queue  interface{ Close() error }
}

func NewWorkerApp() (*WorkerApp, error) {
	cfg := loadConfig()

	sqsSub, err := sqssubscriber.NewSubscriber(context.Background(), sqssubscriber.SubscriberConfig{
		QueueURL:          cfg.SQS.QueueURL,
		EndpointURL:       cfg.SQS.EndpointURL,
		Region:            cfg.SQS.Region,
		MaxMessages:       int32(cfg.SQS.MaxMessages),
		WaitTimeSeconds:   int32(cfg.SQS.WaitTimeSeconds),
		VisibilityTimeout: int32(cfg.SQS.VisibilityTimeout),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create SQS subscriber: %w", err)
	}

	log.Printf("using SQS queue: %s (endpoint: %s)", cfg.SQS.QueueURL, cfg.SQS.EndpointURL)

	auth0Client := auth.NewClient(auth.Config{
		Domain:       cfg.Auth0.Domain,
		ClientID:     cfg.Auth0.ClientID,
		ClientSecret: cfg.Auth0.ClientSecret,
	})

	registerIdentity := identityusecase.NewRegisterIdentityUseCase(auth0Client)
	eventHandler := queuehandler.NewEventHandler(registerIdentity)
	worker := adapterqueue.NewWorker(sqsSub, eventHandler)

	return &WorkerApp{
		worker: worker,
		queue:  sqsSub,
	}, nil
}

func (a *WorkerApp) Start(ctx context.Context) error {
	return a.worker.Start(ctx)
}

func (a *WorkerApp) Stop() {
	a.worker.Stop()
}

func (a *WorkerApp) Close() {
	if a.queue != nil {
		_ = a.queue.Close()
	}
}
