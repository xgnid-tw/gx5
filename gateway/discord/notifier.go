package discord

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/xgnid-tw/gx5/domain"
)

type discordSession interface {
	UserChannelCreate(recipientID string, options ...discordgo.RequestOption) (*discordgo.Channel, error)
	ChannelMessageSend(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
}

// Notifier implements port.Notifier using Discord DMs.
type Notifier struct {
	s            discordSession
	logChannelID string
	debug        bool
}

func NewNotifier(s *discordgo.Session, logChannelID string, debug bool) *Notifier {
	return &Notifier{s: s, logChannelID: logChannelID, debug: debug}
}

func (n *Notifier) Notify(_ context.Context, user domain.User) error {
	message := fmt.Sprintf(
		"[欠費提醒] https://www.notion.so/%s (如果有漏登聯絡一下XG) ",
		user.NotionID,
	)

	return n.sendDM(user.DiscordID, message)
}

func (n *Notifier) sendDM(discordID string, message string) error {
	channel, err := n.s.UserChannelCreate(discordID)
	if err != nil {
		return fmt.Errorf("error creating channel: %w", err)
	}

	_, err = n.s.ChannelMessageSend(n.logChannelID, message)
	if err != nil {
		return fmt.Errorf("error sending to log channel: %w", err)
	}

	if n.debug {
		log.Print("debug mode on")
		return nil
	}

	_, err = n.s.ChannelMessageSend(channel.ID, message)
	if err != nil {
		return fmt.Errorf("error sending dm: %w", err)
	}

	return nil
}
