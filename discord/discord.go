package discord

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/xgnid-tw/gx5/model"
)

type Discord interface {
	GetChanMsgAndDM(ctx context.Context)
	// MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate)
}

type discord struct {
	s  *discordgo.Session
	ch chan model.User
}

func NewDiscordEventService(
	s *discordgo.Session,
	ch chan model.User,
) Discord {
	return &discord{
		s:  s,
		ch: ch,
	}
}

func (de *discord) GetChanMsgAndDM(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case result := <-de.ch:
		message := fmt.Sprintf(
			"[欠費提醒] https://www.notion.so/%s (如果有漏登聯絡一下XG) ",
			result.NotionID,
		)

		err := de.sendDM(result.DiscordID, message)
		if err != nil {
			log.Fatalf("send dm: %s", err)
		}
	}
}

func (de *discord) sendDM(discordID string, message string) error {
	channel, err := de.s.UserChannelCreate(discordID)
	if err != nil {
		return fmt.Errorf("error creating channel %w", err)
	}

	_, err = de.s.ChannelMessageSend(channel.ID, message)
	if err != nil {
		return fmt.Errorf("error when sending dm: %w", err)
	}

	return nil
}

/*func (de *discordEvent) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content != "ping" {
		return
	}

	de.sendDM(m.Author.ID, "pong")
}*/
