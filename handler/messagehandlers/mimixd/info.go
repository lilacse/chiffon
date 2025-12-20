package mimixd

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/lilacse/chiffon/embedbuilder"
	"github.com/lilacse/chiffon/resp"
	"github.com/lilacse/chiffon/store"
)

type infoHandler struct {
	store *store.Store
}

type levelCc struct {
	level string
	cc    float64
}

func NewInfoHandler(store *store.Store) *infoHandler {
	return &infoHandler{
		store: store,
	}
}

func (h *infoHandler) Handle(ctx context.Context, e *gateway.MessageCreateEvent) bool {
	st := h.store.Bot.State()

	mentionedBot := false
	for _, user := range e.Mentions {
		if user.ID == h.store.Bot.BotId() {
			mentionedBot = true
			break
		}
	}

	if !mentionedBot || e.Message.ReferencedMessage == nil {
		return false
	}

	refMsg := e.Message.ReferencedMessage
	if len(refMsg.Embeds) == 0 {
		return false
	}

	embed := refMsg.Embeds[0]
	descLines := strings.Split(embed.Description, "\n")
	dxCcs := make([]levelCc, 0)
	stCcs := make([]levelCc, 0)

	for _, line := range descLines {
		if strings.HasPrefix(line, "**DX:** ") {
			dxCcs = parseLevelCc(line)
		}

		if strings.HasPrefix(line, "**STD:** ") {
			stCcs = parseLevelCc(line)
		}
	}

	if len(dxCcs) == 0 && len(stCcs) == 0 {
		return false
	}

	descBuilder := strings.Builder{}
	descBuilder.WriteString("Here are the ratings for each rank border for the song!\n")
	descBuilder.WriteString("Legend: **SSS+** / SSS / SS+ / SS / S+ / S\n\n")

	if len(dxCcs) > 0 {
		descBuilder.WriteString("**DX Charts**\n")
		for _, lc := range dxCcs {
			rts := calculateRt(lc.cc)
			writeLevelInfo(&descBuilder, lc, rts)
		}
		descBuilder.WriteString("\n")
	}

	if len(stCcs) > 0 {
		descBuilder.WriteString("**ST Charts**\n")
		for _, lc := range stCcs {
			rts := calculateRt(lc.cc)
			writeLevelInfo(&descBuilder, lc, rts)
		}
		descBuilder.WriteString("\n")
	}

	respEmbed := discord.Embed{
		Title:       embed.Title,
		Description: descBuilder.String(),
	}

	if embed.Thumbnail != nil && embed.Thumbnail.URL != "" {
		respEmbed.Thumbnail = &discord.EmbedThumbnail{
			URL: embed.Thumbnail.URL,
		}
	}

	resp.SendReplyToMessage(st, embedbuilder.Info(respEmbed), refMsg)

	return true
}

func parseLevelCc(line string) []levelCc {
	var ccs []levelCc
	parts := strings.Split(line, " / ")
	ccRegex := regexp.MustCompile(`\([0-9]{1,2}\.[0-9]\)`)
	lvls := []string{"B", "A", "E", "M", "R"}

	for i, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		if trimmedPart == "" {
			continue
		}

		ccMatch := ccRegex.FindString(trimmedPart)
		if ccMatch == "" {
			continue
		}

		cc, _ := strconv.ParseFloat(strings.TrimPrefix(strings.TrimSuffix(ccMatch, ")"), "("), 64)
		ccs = append(ccs, levelCc{level: lvls[i], cc: cc})
	}

	return ccs
}

func calculateRt(cc float64) []int {
	borders := []float64{1.005, 1.00, 0.995, 0.99, 0.98, 0.97}
	multipliers := []float64{22.4, 21.6, 21.1, 20.8, 20.3, 20}

	results := make([]int, len(borders))
	for i, border := range borders {
		results[i] = int(math.Floor(cc * border * multipliers[i]))
	}
	return results
}

func writeLevelInfo(sb *strings.Builder, lc levelCc, rts []int) {
	levelEmoteMap := map[string]string{
		"B": ":green_square:",
		"A": ":yellow_square:",
		"E": ":red_square:",
		"M": ":purple_square:",
		"R": ":white_large_square:",
	}

	fmt.Fprintf(sb, "%s (%.1f): ", levelEmoteMap[lc.level], lc.cc)

	for i, rt := range rts {
		if i > 0 {
			sb.WriteString(" / ")
		}
		if i == 0 {
			fmt.Fprintf(sb, "**%d**", rt)
		} else {
			fmt.Fprintf(sb, "%d", rt)
		}
	}
	sb.WriteString("\n")
}
