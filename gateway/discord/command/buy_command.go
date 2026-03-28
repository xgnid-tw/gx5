package command

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/xgnid-tw/gx5/port"
)

const (
	buyCommandName     = "buy"
	buyModalPrefix     = "buy_modal"
	amountInputID      = "jpy_amount"
	itemNameInputID    = "item_name"
	modalCustomIDParts = 3
)

// RegisterBuyCommand registers the /buy message command and its modal handler.
func RegisterBuyCommand(ch *Handler, uc port.BuyRecordRegisterer) {
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
		respondError(s, i, "無法取得目標訊息")
		return
	}

	targetDiscordID := targetMsg.Author.ID

	// Get thread title from the channel (must be in a thread)
	channel, err := s.Channel(i.ChannelID)
	if err != nil {
		respondError(s, i, "無法取得頻道資訊")
		return
	}

	if !channel.IsThread() {
		respondError(s, i, "此指令只能在討論串中使用")
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
			Title:    "確認購買",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    amountInputID,
							Label:       "日幣",
							Style:       discordgo.TextInputShort,
							Placeholder: "例: 3000",
							Required:    true,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: itemNameInputID,
							Label:    "品項",
							Style:    discordgo.TextInputShort,
							Required: false,
							Value:    threadTitle,
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
	s *discordgo.Session, i *discordgo.InteractionCreate, uc port.BuyRecordRegisterer,
) {
	data := i.ModalSubmitData()

	// Parse customID: buy_modal:<targetDiscordID>:<threadTitle>
	parts := strings.SplitN(data.CustomID, ":", modalCustomIDParts)
	if len(parts) != modalCustomIDParts {
		respondError(s, i, "無效的表單資料")
		return
	}

	targetDiscordID := parts[1]
	threadTitle := parts[2]

	// Extract JPY amount and item name from modal inputs
	var jpyStr string

	var itemName string

	for _, row := range data.Components {
		if ar, ok := row.(*discordgo.ActionsRow); ok {
			for _, comp := range ar.Components {
				if ti, ok := comp.(*discordgo.TextInput); ok {
					switch ti.CustomID {
					case amountInputID:
						jpyStr = ti.Value
					case itemNameInputID:
						itemName = ti.Value
					}
				}
			}
		}
	}

	if itemName == "" {
		itemName = threadTitle
	}

	jpyAmount, err := strconv.ParseFloat(jpyStr, 64)
	if err != nil || jpyAmount <= 0 {
		respondError(s, i, "無效的日幣金額")
		return
	}

	ctx := context.Background()

	err = uc.Execute(ctx, targetDiscordID, jpyAmount, itemName)
	if err != nil {
		log.Printf("register buy record failed: %s", err)
		respondError(s, i, "登記失敗")

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
