package handler

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/lilacse/chiffon/embedbuilder"
	"github.com/lilacse/chiffon/store"
)

type factory struct {
	store *store.Store
}

func NewFactory(store *store.Store) *factory {
	return &factory{
		store: store,
	}
}

func (f *factory) NewOnMessageCreateHandler() *onMessageCreateHandler {
	return &onMessageCreateHandler{
		store: f.store,
	}
}

func sendHandleError(ctx context.Context, r any, st *state.State, messageId discord.MessageID, channelId discord.ChannelID) {
	d := api.SendMessageData{
		Embeds: []discord.Embed{
			embedbuilder.Error(ctx, fmt.Sprintf("%s", r)),
		},
		Reference: &discord.MessageReference{
			MessageID: messageId,
		},
		AllowedMentions: &api.AllowedMentions{
			RepliedUser: option.False,
		},
	}

	st.SendMessageComplex(channelId, d)
}
