package resp

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func SendReplyToMessage(st *state.State, em discord.Embed, e *discord.Message) {
	SendReplyWithComponents(st, em, []discord.ContainerComponent{}, e.ChannelID, e.ID)
}

func SendReplyWithComponents(st *state.State, em discord.Embed, cc []discord.ContainerComponent, channelId discord.ChannelID, replyId discord.MessageID) {
	d := api.SendMessageData{
		Embeds: []discord.Embed{
			em,
		},
		Components: cc,
		Reference: &discord.MessageReference{
			MessageID: replyId,
		},
		AllowedMentions: &api.AllowedMentions{
			RepliedUser: option.False,
		},
	}

	st.SendMessageComplex(channelId, d)
}
