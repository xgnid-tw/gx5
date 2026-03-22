package discord

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/xgnid-tw/gx5/usecase"
)

const (
	buyCommandName  = "buy"
	buyModalPrefix  = "buy_modal"
	amountInputID   = "jpy_amount"
)

// RegisterBuyCommand registers the /buy message command and its modal handler.
func RegisterBuyCommand(ch *CommandHandler, uc *usecase.RegisterBuyRecord) {
	cmd := &discordgo.ApplicationCommand{
		Name: buyCommandName,
		Type: discordgo.MessageApplicationCommand,
	}

	ch.RegisterCommand(cmd, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handleBuyCommand(s, i)
	})

	ch.RegisterModalHandler(buyModalPrefix, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handleBuyModal(s, i, uc)
	})
}

func handleBuyCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	// Get the target message from the resolved data
	targetMsg, ok := data.Resolved.Messages[data.TargetID]
	if !ok {
		respondError(s, i, "could not resolve target message")
		return
	}

	targetDiscordID := targetMsg.Author.ID

	// Get thread title from the channel (must be in a thread)
	channel, err := s.Channel(i.ChannelID)
	if err != nil {
		respondError(s, i, "could not get channel info")
		return
	}

	if !channel.IsThread() {
		respondError(s, i, "this command must be used in a thread")
		return
	}

	threadTitle := channel.Name

	// Encode targetDiscordID and threadTitle into modal CustomID
	// Format: buy_modal:<targetDiscordID>:<threadTitle>
	customID := fmt.Sprintf("%s:%s:%s", buyModalPrefix, targetDiscordID, threadTitle)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: customID,
			Title:    "Register Buy Record",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    amountInputID,
							Label:       "JPY Amount",
							Style:       discordgo.TextInputShort,
							Placeholder: "e.g. 3000",
							Required:    true,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("error responding with modal: %s", err)
	}
}

func handleBuyModal(
	s *discordgo.Session, i *discordgo.InteractionCreate, uc *usecase.RegisterBuyRecord,
) {
	data := i.ModalSubmitData()

	// Parse customID: buy_modal:<targetDiscordID>:<threadTitle>
	parts := strings.SplitN(data.CustomID, ":", 3)
	if len(parts) != 3 {
		respondError(s, i, "invalid modal data")
		return
	}

	targetDiscordID := parts[1]
	threadTitle := parts[2]

	// Extract JPY amount from modal input
	var jpyStr string

	for _, row := range data.Components {
		if ar, ok := row.(*discordgo.ActionsRow); ok {
			for _, comp := range ar.Components {
				if ti, ok := comp.(*discordgo.TextInput); ok && ti.CustomID == amountInputID {
					jpyStr = ti.Value
				}
			}
		}
	}

	jpyAmount, err := strconv.ParseFloat(jpyStr, 64)
	if err != nil || jpyAmount <= 0 {
		respondError(s, i, "invalid JPY amount")
		return
	}

	ctx := context.Background()

	err = uc.Execute(ctx, targetDiscordID, jpyAmount, threadTitle)
	if err != nil {
		log.Printf("register buy record failed: %s", err)
		respondError(s, i, "failed to register buy record")

		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "登記完畢",
		},
	})
	if err != nil {
		log.Printf("error responding to modal: %s", err)
	}
}

func respondError(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("error responding with error message: %s", err)
	}
}
