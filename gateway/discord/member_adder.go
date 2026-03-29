package discord

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/bwmarrin/discordgo"
)

type memberAdderSession interface {
	GuildMembers(
		guildID string, after string, limit int, options ...discordgo.RequestOption,
	) ([]*discordgo.Member, error)
	ThreadMemberAdd(threadID, memberID string, options ...discordgo.RequestOption) error
	ChannelMessages(
		channelID string, limit int, beforeID, afterID, aroundID string,
		options ...discordgo.RequestOption,
	) ([]*discordgo.Message, error)
	ChannelMessageDelete(
		channelID, messageID string, options ...discordgo.RequestOption,
	) error
}

// MemberAdder implements port.MemberAdder using the Discord API.
type MemberAdder struct {
	s       memberAdderSession
	guildID string
}

func NewMemberAdder(s *discordgo.Session, guildID string) *MemberAdder {
	return &MemberAdder{s: s, guildID: guildID}
}

func (ma *MemberAdder) AddRoleMembersToThread(
	_ context.Context, threadID string, roleID string,
) error {
	members, err := ma.fetchMembersWithRole(roleID)
	if err != nil {
		return fmt.Errorf("fetch members with role %s: %w", roleID, err)
	}

	for _, m := range members {
		err = ma.s.ThreadMemberAdd(threadID, m.User.ID)
		if err != nil {
			log.Printf("failed to add member %s to thread %s: %s", m.User.ID, threadID, err)
		}
	}

	ma.deleteSystemMessages(threadID)

	return nil
}

const recentMessageLimit = 50

func (ma *MemberAdder) deleteSystemMessages(threadID string) {
	messages, err := ma.s.ChannelMessages(threadID, recentMessageLimit, "", "", "")
	if err != nil {
		log.Printf("failed to fetch messages for cleanup in thread %s: %s", threadID, err)
		return
	}

	for _, msg := range messages {
		if msg.Type != discordgo.MessageTypeDefault {
			delErr := ma.s.ChannelMessageDelete(threadID, msg.ID)
			if delErr != nil {
				log.Printf("failed to delete system message %s (type %d): %s", msg.ID, msg.Type, delErr)
			}
		}
	}
}

func (ma *MemberAdder) fetchMembersWithRole(roleID string) ([]*discordgo.Member, error) {
	var result []*discordgo.Member

	after := ""

	const pageSize = 100

	for {
		members, err := ma.s.GuildMembers(ma.guildID, after, pageSize)
		if err != nil {
			return nil, fmt.Errorf("guild members: %w", err)
		}

		for _, m := range members {
			if slices.Contains(m.Roles, roleID) {
				result = append(result, m)
			}
		}

		if len(members) < pageSize {
			break
		}

		after = members[len(members)-1].User.ID
	}

	return result, nil
}
