package discord

import "github.com/bwmarrin/discordgo"

type mockDiscordSession struct {
	userChannelCreateFn  func(recipientID string, options ...discordgo.RequestOption) (*discordgo.Channel, error)
	channelMessageSendFn func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
	sentMessages         []struct{ channelID, content string }
}

func (m *mockDiscordSession) UserChannelCreate(recipientID string, options ...discordgo.RequestOption) (*discordgo.Channel, error) {
	return m.userChannelCreateFn(recipientID, options...)
}

func (m *mockDiscordSession) ChannelMessageSend(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	m.sentMessages = append(m.sentMessages, struct{ channelID, content string }{channelID, content})
	return m.channelMessageSendFn(channelID, content, options...)
}

func newTestNotifier(s discordSession, logChannelID string, debug bool) *Notifier {
	return &Notifier{s: s, logChannelID: logChannelID, debug: debug}
}
