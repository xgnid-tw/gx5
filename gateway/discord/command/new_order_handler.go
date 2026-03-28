package command

import (
	"context"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/xgnid-tw/gx5/domain"
	"github.com/xgnid-tw/gx5/port"
)

const (
	newOrderCommandName = "neworder"
)

// RegisterNewOrderCommand registers the /neworder slash command and its handler.
func RegisterNewOrderCommand(ch *Handler, uc port.OrderCreator) error {
	adminPerm := int64(discordgo.PermissionAdministrator)

	cmd := &discordgo.ApplicationCommand{
		Name:                     newOrderCommandName,
		Description:              "Create a new group purchase order",
		DefaultMemberPermissions: &adminPerm,
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

	return ch.RegisterCommand(cmd, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handleNewOrder(s, i, uc)
	})
}

func handleNewOrder(
	s *discordgo.Session, i *discordgo.InteractionCreate, uc port.OrderCreator,
) {
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

	err := uc.Execute(context.Background(), i.ChannelID, order)
	if err != nil {
		log.Printf("create order failed: %s", err)
		respondError(s, i, "failed to create order")

		return
	}

	respondSuccess(s, i, "Order created: "+order.ThreadName)
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
