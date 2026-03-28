package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type threadSession interface {
	ThreadStartComplex(
		channelID string, data *discordgo.ThreadStart, options ...discordgo.RequestOption,
	) (*discordgo.Channel, error)
	ChannelMessageSend(
		channelID string, content string, options ...discordgo.RequestOption,
	) (*discordgo.Message, error)
}

// ThreadCreator implements port.ThreadCreator using the Discord API.
type ThreadCreator struct {
	s threadSession
}

func NewThreadCreator(s *discordgo.Session) *ThreadCreator {
	return &ThreadCreator{s: s}
}

func (tc *ThreadCreator) CreateThread(_ context.Context, channelID string, name string, message string) error {
	thread, err := tc.s.ThreadStartComplex(channelID, &discordgo.ThreadStart{
		Name: name,
		Type: discordgo.ChannelTypeGuildPublicThread,
	})
	if err != nil {
		return fmt.Errorf("error creating thread: %w", err)
	}

	if message != "" {
		_, err = tc.s.ChannelMessageSend(thread.ID, message)
		if err != nil {
			return fmt.Errorf("error sending thread message: %w", err)
		}
	}

	return nil
}
