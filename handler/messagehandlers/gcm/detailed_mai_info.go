package gcm

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/lilacse/chiffon/embedbuilder"
	"github.com/lilacse/chiffon/resp"
	"github.com/lilacse/chiffon/store"
)

type detailedMaiInfoHandler struct {
	store *store.Store
}

func NewDetailedMaiInfoHandler(store *store.Store) *detailedMaiInfoHandler {
	return &detailedMaiInfoHandler{
		store: store,
	}
}

func (h *detailedMaiInfoHandler) Handle(ctx context.Context, e *gateway.MessageCreateEvent) bool {
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
	dxLevels := make([]levelDetails, 0)
	stLevels := make([]levelDetails, 0)
	isReadDx := false
	isReadSt := false
	currLevel := ""

	for _, line := range descLines {
		if strings.TrimSpace(line) == "**DX Chart Info:**" {
			isReadDx = true
			continue
		}
		if strings.TrimSpace(line) == "**ST Chart Info:**" {
			isReadSt = true
			continue
		}
		if !isReadDx && !isReadSt {
			continue
		}

		if strings.TrimSpace(line) == "" {
			isReadDx = false
			continue
		}
		if currLevel == "" {
			level := parseDetailedLevelCc(line)
			if isReadDx {
				dxLevels = append(dxLevels, level)
			} else if isReadSt {
				stLevels = append(stLevels, level)
			}
			currLevel = level.level
		} else {
			if isReadDx {
				parseDetailedLevelNoteCounts(line, &dxLevels[len(dxLevels)-1])
			} else if isReadSt {
				parseDetailedLevelNoteCounts(line, &stLevels[len(stLevels)-1])
			}
			currLevel = ""
		}
	}

	if len(dxLevels) == 0 && len(stLevels) == 0 {
		return false
	}

	descBuilder := strings.Builder{}
	descBuilder.WriteString("Here are the lost percentages and ratings for each rank border for the song!\n")
	descBuilder.WriteString("Rating Legend: **SSS+** / SSS / SS+ / SS / S+ / S\n")
	descBuilder.WriteString("Lost Percentage Legend: Great / Good / Miss\n")
	descBuilder.WriteString("Break Lost Percentage Legend: HPerf-LPerf / HGreat-MGreat-LGreat / Good / Miss\n\n")

	if len(dxLevels) > 0 {
		descBuilder.WriteString("**DX Charts**\n")
		for _, ld := range dxLevels {
			calculateLostRate(&ld)
			rts := calculateRt(ld.cc)
			writeLevelInfo(&descBuilder, ld, rts)
			writeLostPercentageInfo(&descBuilder, ld)
		}
		descBuilder.WriteString("\n")
	}
	if len(stLevels) > 0 {
		descBuilder.WriteString("**ST Charts**\n")
		for _, ld := range stLevels {
			calculateLostRate(&ld)
			rts := calculateRt(ld.cc)
			writeLevelInfo(&descBuilder, ld, rts)
			writeLostPercentageInfo(&descBuilder, ld)
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

func writeLostPercentageInfo(sb *strings.Builder, ld levelDetails) {
	if ld.tapCount > 0 {
		fmt.Fprintf(sb, "**Tap:** %.4f%% / %.4f%% / %.4f%%\n", ld.tapGreatLostRate, ld.tapGoodLostRate, ld.tapMissLostRate)
	}
	if ld.holdCount > 0 {
		fmt.Fprintf(sb, "**Hold:** %.4f%% / %.4f%% / %.4f%%\n", ld.holdGreatLostRate, ld.holdGoodLostRate, ld.holdMissLostRate)
	}
	if ld.slideCount > 0 {
		fmt.Fprintf(sb, "**Slide:** %.4f%% / %.4f%% / %.4f%%\n", ld.slideGreatLostRate, ld.slideGoodLostRate, ld.slideMissLostRate)
	}
	if ld.touchCount > 0 {
		fmt.Fprintf(sb, "**Touch:** %.4f%% / %.4f%% / %.4f%%\n", ld.touchGreatLostRate, ld.touchGoodLostRate, ld.touchMissLostRate)
	}
	if ld.breakCount > 0 {
		fmt.Fprintf(sb, "**Break:** %.4f-%.4f%% / %.4f-%.4f-%.4f%% / %.4f%% / %.4f%%\n",
			ld.breakHighPerfectBonusLostRate,
			ld.breakLowPerfectBonusLostRate,
			ld.breakHighGreatLostRate+ld.breakGreatBonusLostRate,
			ld.breakMidGreatLostRate+ld.breakGreatBonusLostRate,
			ld.breakLowGreatLostRate+ld.breakGreatBonusLostRate,
			ld.breakGoodLostRate+ld.breakGoodBonusLostRate,
			ld.breakMissLostRate+ld.breakMissBonusLostRate,
		)
	}
}

func parseDetailedLevelCc(line string) levelDetails {
	res := levelDetails{}
	splitted := strings.Fields(line)
	if len(splitted) < 3 {
		return res
	}

	levelEmote := strings.TrimSpace(splitted[0])
	level := ""
	switch levelEmote {
	case ":green_square:":
		level = "B"
	case ":yellow_square:":
		level = "A"
	case ":red_square:":
		level = "E"
	case ":purple_square:":
		level = "M"
	case ":white_large_square:":
		level = "R"
	}

	if level == "" {
		return res
	}

	ccStr := strings.TrimSpace(splitted[2])
	if len(ccStr) < 5 || ccStr[0] != '(' || ccStr[len(ccStr)-1] != ')' {
		return res
	}
	ccStr = ccStr[1 : len(ccStr)-1]
	cc, err := strconv.ParseFloat(ccStr, 64)
	if err != nil {
		return res
	}

	return levelDetails{
		level: level,
		cc:    cc,
	}
}

func parseDetailedLevelNoteCounts(line string, ld *levelDetails) {
	parts := strings.Split(line, " / ")
	if len(parts) != 6 {
		return
	}

	tapCount, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return
	}

	holdCount, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return
	}

	slideCount, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return
	}

	touchCount, err := strconv.ParseInt(parts[4], 10, 64)
	if err != nil {
		return
	}

	breakCount, err := strconv.ParseInt(parts[5], 10, 64)
	if err != nil {
		return
	}

	ld.tapCount = tapCount
	ld.holdCount = holdCount
	ld.slideCount = slideCount
	ld.touchCount = touchCount
	ld.breakCount = breakCount
}

func calculateLostRate(ld *levelDetails) {
	tapWeight := int64(1)
	holdWeight := int64(2)
	slideWeight := int64(3)
	touchWeight := int64(1)
	breakWeight := int64(5)

	totalWeight := ld.tapCount*tapWeight +
		ld.holdCount*holdWeight +
		ld.slideCount*slideWeight +
		ld.touchCount*touchWeight +
		ld.breakCount*breakWeight

	greatPenaltyRatio := 0.2
	goodPenaltyRatio := 0.5

	breakHighGreatPenaltyRatio := 0.2
	breakMidGreatPenaltyRatio := 0.4
	breakLowGreatPenaltyRatio := 0.5
	breakGoodPenaltyRatio := 0.6

	breakBonusHighPerfectPenaltyRatio := 0.25
	breakBonusLowPerfectPenaltyRatio := 0.5
	breakBonusGreatPenaltyRatio := 0.6
	breakBonusGoodPenaltyRatio := 0.7

	ld.tapMissLostRate = float64(tapWeight) / float64(totalWeight) * 100
	ld.tapGreatLostRate = ld.tapMissLostRate * greatPenaltyRatio
	ld.tapGoodLostRate = ld.tapMissLostRate * goodPenaltyRatio

	ld.holdMissLostRate = float64(holdWeight) / float64(totalWeight) * 100
	ld.holdGreatLostRate = ld.holdMissLostRate * greatPenaltyRatio
	ld.holdGoodLostRate = ld.holdMissLostRate * goodPenaltyRatio

	ld.slideMissLostRate = float64(slideWeight) / float64(totalWeight) * 100
	ld.slideGreatLostRate = ld.slideMissLostRate * greatPenaltyRatio
	ld.slideGoodLostRate = ld.slideMissLostRate * goodPenaltyRatio

	ld.touchMissLostRate = float64(touchWeight) / float64(totalWeight) * 100
	ld.touchGreatLostRate = ld.touchMissLostRate * greatPenaltyRatio
	ld.touchGoodLostRate = ld.touchMissLostRate * goodPenaltyRatio

	ld.breakMissLostRate = float64(breakWeight) / float64(totalWeight) * 100
	ld.breakHighGreatLostRate = ld.breakMissLostRate * breakHighGreatPenaltyRatio
	ld.breakMidGreatLostRate = ld.breakMissLostRate * breakMidGreatPenaltyRatio
	ld.breakLowGreatLostRate = ld.breakMissLostRate * breakLowGreatPenaltyRatio
	ld.breakGoodLostRate = ld.breakMissLostRate * breakGoodPenaltyRatio

	ld.breakMissBonusLostRate = 1.0 / float64(ld.breakCount)
	ld.breakHighPerfectBonusLostRate = ld.breakMissBonusLostRate * breakBonusHighPerfectPenaltyRatio
	ld.breakLowPerfectBonusLostRate = ld.breakMissBonusLostRate * breakBonusLowPerfectPenaltyRatio
	ld.breakGreatBonusLostRate = ld.breakMissBonusLostRate * breakBonusGreatPenaltyRatio
	ld.breakGoodBonusLostRate = ld.breakMissBonusLostRate * breakBonusGoodPenaltyRatio
}
