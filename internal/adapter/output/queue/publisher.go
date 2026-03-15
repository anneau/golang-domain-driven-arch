package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	messagingevent "github.com/hkobori/golang-domain-driven-arch/internal/adapter/messaging/event"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/event"
)

type Publisher struct {
	client   *sqs.Client
	queueURL string
}

func NewPublisher(ctx context.Context, queueURL string, endpointURL string) (*Publisher, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion("us-east-1"),
	}
	if endpointURL != "" {
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("test", "test", ""),
		))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}

	sqsOpts := []func(*sqs.Options){}
	if endpointURL != "" {
		sqsOpts = append(sqsOpts, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(endpointURL)
		})
	}

	return &Publisher{
		client:   sqs.NewFromConfig(cfg, sqsOpts...),
		queueURL: queueURL,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, evt event.Event) error {
	var payload any
	switch e := evt.(type) {
	case *event.UserCreatedEvent:
		payload = messagingevent.UserCreatedEventDTO{
			EventType:  e.EventType(),
			UserID:     e.UserID(),
			Name:       e.Name(),
			Email:      e.Email(),
			OccurredAt: e.OccurredAt(),
		}
	default:
		return fmt.Errorf("unknown event type: %T", evt)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	bodyStr := string(body)
	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &p.queueURL,
		MessageBody: &bodyStr,
	})
	return err
}
