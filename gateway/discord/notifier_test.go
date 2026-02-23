package discord

import (
	"context"
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/require"

	"github.com/xgnid-tw/gx5/domain"
)

var testUser = domain.User{DiscordID: "111", Name: "Alice", NotionID: "abc"}

func TestNotify_DebugMode_SkipsDM(t *testing.T) {
	m := &mockDiscordSession{
		userChannelCreateFn: func(string, ...discordgo.RequestOption) (*discordgo.Channel, error) {
			return &discordgo.Channel{ID: "dm-chan"}, nil
		},
		channelMessageSendFn: func(string, string, ...discordgo.RequestOption) (*discordgo.Message, error) {
			return &discordgo.Message{}, nil
		},
	}

	n := newTestNotifier(m, "log-chan", true)
	err := n.Notify(context.Background(), testUser)

	require.NoError(t, err)
	require.Len(t, m.sentMessages, 1)
	require.Equal(t, "log-chan", m.sentMessages[0].channelID)
}

func TestNotify_NormalMode_SendsDM(t *testing.T) {
	m := &mockDiscordSession{
		userChannelCreateFn: func(string, ...discordgo.RequestOption) (*discordgo.Channel, error) {
			return &discordgo.Channel{ID: "dm-chan"}, nil
		},
		channelMessageSendFn: func(string, string, ...discordgo.RequestOption) (*discordgo.Message, error) {
			return &discordgo.Message{}, nil
		},
	}

	n := newTestNotifier(m, "log-chan", false)
	err := n.Notify(context.Background(), testUser)

	require.NoError(t, err)
	require.Len(t, m.sentMessages, 2)
	require.Equal(t, "log-chan", m.sentMessages[0].channelID)
	require.Equal(t, "dm-chan", m.sentMessages[1].channelID)
	require.Contains(t, m.sentMessages[1].content, testUser.NotionID)
}

func TestNotify_UserChannelCreateFails(t *testing.T) {
	m := &mockDiscordSession{
		userChannelCreateFn: func(string, ...discordgo.RequestOption) (*discordgo.Channel, error) {
			return nil, errors.New("api down")
		},
		channelMessageSendFn: func(string, string, ...discordgo.RequestOption) (*discordgo.Message, error) {
			return &discordgo.Message{}, nil
		},
	}

	n := newTestNotifier(m, "log-chan", false)
	err := n.Notify(context.Background(), testUser)

	require.Error(t, err)
	require.ErrorContains(t, err, "error creating channel")
}

func TestNotify_LogChannelSendFails(t *testing.T) {
	m := &mockDiscordSession{
		userChannelCreateFn: func(string, ...discordgo.RequestOption) (*discordgo.Channel, error) {
			return &discordgo.Channel{ID: "dm-chan"}, nil
		},
		channelMessageSendFn: func(string, string, ...discordgo.RequestOption) (*discordgo.Message, error) {
			return nil, errors.New("send failed")
		},
	}

	n := newTestNotifier(m, "log-chan", false)
	err := n.Notify(context.Background(), testUser)

	require.Error(t, err)
	require.ErrorContains(t, err, "error sending to log channel")
}

func TestNotify_DMSendFails(t *testing.T) {
	callCount := 0
	m := &mockDiscordSession{
		userChannelCreateFn: func(string, ...discordgo.RequestOption) (*discordgo.Channel, error) {
			return &discordgo.Channel{ID: "dm-chan"}, nil
		},
		channelMessageSendFn: func(string, string, ...discordgo.RequestOption) (*discordgo.Message, error) {
			callCount++
			if callCount == 1 {
				return &discordgo.Message{}, nil
			}
			return nil, errors.New("dm failed")
		},
	}

	n := newTestNotifier(m, "log-chan", false)
	err := n.Notify(context.Background(), testUser)

	require.Error(t, err)
	require.ErrorContains(t, err, "error sending dm")
}
