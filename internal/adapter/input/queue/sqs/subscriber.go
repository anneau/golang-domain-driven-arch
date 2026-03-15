package sqs

import (
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/hkobori/golang-domain-driven-arch/internal/app/port"
)

type SubscriberConfig struct {
	QueueURL          string
	EndpointURL       string
	MaxMessages       int32
	WaitTimeSeconds   int32
	VisibilityTimeout int32
}

type Subscriber struct {
	client   *sqs.Client
	config   SubscriberConfig
	messages chan port.Message
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

func NewSubscriber(ctx context.Context, cfg SubscriberConfig) (*Subscriber, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion("us-east-1"),
	}
	if cfg.EndpointURL != "" {
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("test", "test", ""),
		))
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}

	if cfg.MaxMessages <= 0 {
		cfg.MaxMessages = 10
	}
	if cfg.WaitTimeSeconds <= 0 {
		cfg.WaitTimeSeconds = 20
	}
	if cfg.VisibilityTimeout <= 0 {
		cfg.VisibilityTimeout = 30
	}

	sqsOpts := []func(*sqs.Options){}
	if cfg.EndpointURL != "" {
		sqsOpts = append(sqsOpts, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(cfg.EndpointURL)
		})
	}

	return &Subscriber{
		client:   sqs.NewFromConfig(awsCfg, sqsOpts...),
		config:   cfg,
		messages: make(chan port.Message, cfg.MaxMessages),
		stopCh:   make(chan struct{}),
	}, nil
}

func (s *Subscriber) Subscribe(ctx context.Context) (<-chan port.Message, error) {
	s.wg.Add(1)
	go s.poll(ctx)
	return s.messages, nil
}

func (s *Subscriber) poll(ctx context.Context) {
	defer s.wg.Done()
	defer close(s.messages)

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		default:
			s.receiveMessages(ctx)
		}
	}
}

func (s *Subscriber) receiveMessages(ctx context.Context) {
	output, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &s.config.QueueURL,
		MaxNumberOfMessages: s.config.MaxMessages,
		WaitTimeSeconds:     s.config.WaitTimeSeconds,
		VisibilityTimeout:   s.config.VisibilityTimeout,
	})
	if err != nil {
		log.Printf("failed to receive messages: %v", err)
		return
	}

	for _, msg := range output.Messages {
		select {
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		case s.messages <- &sqsMessage{
			client:        s.client,
			queueURL:      s.config.QueueURL,
			msg:           msg,
			receiptHandle: *msg.ReceiptHandle,
		}:
		}
	}
}

func (s *Subscriber) Close() error {
	close(s.stopCh)
	s.wg.Wait()
	return nil
}

type sqsMessage struct {
	client        *sqs.Client
	queueURL      string
	msg           types.Message
	receiptHandle string
}

func (m *sqsMessage) Body() []byte {
	if m.msg.Body == nil {
		return nil
	}
	return []byte(*m.msg.Body)
}

func (m *sqsMessage) Ack(ctx context.Context) error {
	_, err := m.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &m.queueURL,
		ReceiptHandle: &m.receiptHandle,
	})
	return err
}

func (m *sqsMessage) Nack(ctx context.Context) error {
	zero := int32(0)
	_, err := m.client.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          &m.queueURL,
		ReceiptHandle:     &m.receiptHandle,
		VisibilityTimeout: zero,
	})
	return err
}

func (m *sqsMessage) MessageID() string {
	if m.msg.MessageId == nil {
		return ""
	}
	return *m.msg.MessageId
}
