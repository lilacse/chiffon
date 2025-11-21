package handler

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/google/uuid"
	"github.com/lilacse/chiffon/handler/messagehandlers/gcm"
	"github.com/lilacse/chiffon/logger"
	"github.com/lilacse/chiffon/store"
)

type onMessageCreateHandler struct {
	store *store.Store
}

type messageHandler func(ctx context.Context, e *gateway.MessageCreateEvent) bool

func (h *onMessageCreateHandler) Handle(e *gateway.MessageCreateEvent) {
	handleMessage(e, h)
}

func handleMessage(e *gateway.MessageCreateEvent, h *onMessageCreateHandler) {
	traceId := uuid.NewString()
	ctx := context.WithValue(h.store.Bot.Context(), logger.TraceId, traceId)

	handlers := []messageHandler{
		gcm.NewMaiInfoHandler(h.store).Handle,
	}

	defer func() {
		r := recover()
		if r != nil {
			logger.Error(ctx, fmt.Sprintf("error handling command: %s\nstack trace: %s", r, debug.Stack()))
			sendHandleError(ctx, r, h.store.Bot.State(), e.ID, e.ChannelID)
		}
	}()

	for _, handler := range handlers {
		handled := handler(ctx, e)
		if handled {
			return
		}
	}
}
