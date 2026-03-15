package queue

import (
	"context"
	"log"

	"github.com/hkobori/golang-domain-driven-arch/internal/app/port"
)

// MessageHandler はキューメッセージを処理するハンドラーのインターフェース。
// Handle がエラーを返した場合、Worker は次のハンドラーを試みる。
// すべてのハンドラーが失敗した場合はメッセージを Nack する。
type MessageHandler interface {
	Handle(ctx context.Context, msg port.Message) error
}

type Worker struct {
	subscriber port.EventSubscriber
	handlers   []MessageHandler
	stopCh     chan struct{}
}

func NewWorker(subscriber port.EventSubscriber, handlers ...MessageHandler) *Worker {
	return &Worker{
		subscriber: subscriber,
		handlers:   handlers,
		stopCh:     make(chan struct{}),
	}
}

func (w *Worker) Start(ctx context.Context) error {
	messages, err := w.subscriber.Subscribe(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-w.stopCh:
			return nil
		case msg, ok := <-messages:
			if !ok {
				return nil
			}
			w.dispatch(ctx, msg)
		}
	}
}

func (w *Worker) dispatch(ctx context.Context, msg port.Message) {
	for _, h := range w.handlers {
		if err := h.Handle(ctx, msg); err != nil {
			log.Printf("handler failed for message %s: %v", msg.MessageID(), err)
			continue
		}
		if err := msg.Ack(ctx); err != nil {
			log.Printf("failed to ack message %s: %v", msg.MessageID(), err)
		}
		return
	}
	if err := msg.Nack(ctx); err != nil {
		log.Printf("failed to nack message %s: %v", msg.MessageID(), err)
	}
}

func (w *Worker) Stop() {
	close(w.stopCh)
}
