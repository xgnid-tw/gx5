package discord

import (
	"context"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/xgnid-tw/gx5/domain"
	"github.com/xgnid-tw/gx5/usecase"
)

// NewOrderCommand returns the Discord slash command definition for /neworder.
func NewOrderCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "neworder",
		Description: "Create a new group purchase order",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "ordertitle",
				Description: "Name of the order",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "deadline",
				Description: "Order deadline (YYYY-MM-DD)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "shopurl",
				Description: "Shop URL",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "tags",
				Description: "Tag",
				Required:    true,
				Choices:     tagChoices(),
			},
		},
	}
}

// HandleNewOrder returns a Discord interaction handler for the /neworder command.
func HandleNewOrder(uc *usecase.CreateOrder) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		opts := i.ApplicationCommandData().Options
		optMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(opts))

		for _, opt := range opts {
			optMap[opt.Name] = opt
		}

		order := domain.Order{}

		if v, ok := optMap["ordertitle"]; ok {
			order.ThreadName = v.StringValue()
		}

		if v, ok := optMap["deadline"]; ok {
			order.Deadline = v.StringValue()
		}

		if v, ok := optMap["shopurl"]; ok {
			order.ShopURL = v.StringValue()
		}

		if v, ok := optMap["tags"]; ok {
			order.Tag = domain.Tag(v.StringValue())
		}

		var callerID string
		if i.Member != nil && i.Member.User != nil {
			callerID = i.Member.User.ID
		} else if i.User != nil {
			callerID = i.User.ID
		}

		err := uc.Execute(context.Background(), callerID, i.ChannelID, order)
		if err != nil {
			respondToInteraction(s, i, "Error: "+err.Error())
			return
		}

		respondToInteraction(s, i, "Order created: "+order.ThreadName)
	}
}

func tagChoices() []*discordgo.ApplicationCommandOptionChoice {
	choices := make([]*discordgo.ApplicationCommandOptionChoice, len(domain.ValidTags))
	for i, tag := range domain.ValidTags {
		choices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  string(tag),
			Value: string(tag),
		}
	}

	return choices
}

func respondToInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	if err != nil {
		log.Printf("error responding to interaction: %s", err)
	}
}
