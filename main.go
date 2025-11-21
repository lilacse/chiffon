package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/lilacse/chiffon/handler"
	"github.com/lilacse/chiffon/logger"
	"github.com/lilacse/chiffon/store"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	store := store.GetStore()
	store.Bot.SetContext(ctx)
	defer stop()

	logger.Info(ctx, "starting up...")

	token := os.Getenv("CHIFFON_TOKEN")
	if token == "" {
		logger.Fatal(ctx, "environment variable CHIFFON_TOKEN is not set")
	}

	s := state.New("Bot " + token)
	store.Bot.SetState(s)

	s.AddIntents(gateway.IntentGuildMessages)

	hfactory := handler.NewFactory(store)
	s.AddHandler(hfactory.NewOnMessageCreateHandler().Handle)

	u, err := s.Me()
	if err != nil {
		logger.Fatal(ctx, "failed to get bot user with error "+err.Error())
	}

	logger.Info(ctx, fmt.Sprintf("bot user is: %s#%s (%v)", u.Username, u.Discriminator, u.ID))
	store.Bot.SetBotId(u.ID)

	logger.Info(ctx, "starting connection to Discord. bot should be ready!")
	err = s.Connect(ctx)
	if err != nil {
		logger.Fatal(ctx, "connection to Discord is broken with error "+err.Error())
	}

	logger.Info(ctx, "received stopping signal, bot exiting")
}
