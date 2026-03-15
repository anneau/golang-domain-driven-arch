package handler

import (
	"context"
	"encoding/json"
	"log"

	identityusecase "github.com/hkobori/golang-domain-driven-arch/internal/app/usecase/identity"
	messagingevent "github.com/hkobori/golang-domain-driven-arch/internal/adapter/messaging/event"
	"github.com/hkobori/golang-domain-driven-arch/internal/app/port"
)

type EventHandler struct {
	registerIdentity identityusecase.RegisterIdentityUseCase
}

func NewEventHandler(registerIdentity identityusecase.RegisterIdentityUseCase) *EventHandler {
	return &EventHandler{
		registerIdentity: registerIdentity,
	}
}

func (h *EventHandler) Handle(ctx context.Context, msg port.Message) error {
	var dto messagingevent.UserCreatedEventDTO
	if err := json.Unmarshal(msg.Body(), &dto); err != nil {
		log.Printf("failed to unmarshal event: %v", err)
		return err
	}

	if err := h.registerIdentity.Execute(ctx, identityusecase.RegisterIdentityInput{
		UserID: dto.UserID,
		Email:  dto.Email,
	}); err != nil {
		log.Printf("failed to register identity for user %s: %v", dto.UserID, err)
		return err
	}

	log.Printf("successfully registered identity for user %s", dto.UserID)
	return nil
}
