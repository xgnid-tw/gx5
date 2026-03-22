package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// CommandHandler manages Discord application command registration and dispatch.
type CommandHandler struct {
	session    *discordgo.Session
	appID      string
	commands   []*discordgo.ApplicationCommand
	handlers   map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	registered []*discordgo.ApplicationCommand
}

func NewCommandHandler(session *discordgo.Session, appID string) *CommandHandler {
	h := &CommandHandler{
		session:  session,
		appID:    appID,
		handlers: make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)),
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if handler, ok := h.handlers[i.ApplicationCommandData().Name]; ok {
				handler(s, i)
			}
		case discordgo.InteractionModalSubmit:
			if handler, ok := h.handlers["modal:"+i.ModalSubmitData().CustomID[:modalPrefixLen(i.ModalSubmitData().CustomID)]]; ok {
				handler(s, i)
			}
		}
	})

	return h
}

// RegisterCommand registers an application command with its handler.
func (h *CommandHandler) RegisterCommand(
	cmd *discordgo.ApplicationCommand,
	handler func(s *discordgo.Session, i *discordgo.InteractionCreate),
) {
	h.commands = append(h.commands, cmd)
	h.handlers[cmd.Name] = handler
}

// RegisterModalHandler registers a handler for modal submissions with a given prefix.
func (h *CommandHandler) RegisterModalHandler(
	prefix string,
	handler func(s *discordgo.Session, i *discordgo.InteractionCreate),
) {
	h.handlers["modal:"+prefix] = handler
}

// SyncCommands creates all registered commands with the Discord API.
func (h *CommandHandler) SyncCommands() error {
	for _, cmd := range h.commands {
		registered, err := h.session.ApplicationCommandCreate(h.appID, "", cmd)
		if err != nil {
			return err
		}

		h.registered = append(h.registered, registered)
	}

	return nil
}

// UnregisterAll removes all registered commands from the Discord API.
func (h *CommandHandler) UnregisterAll() {
	for _, cmd := range h.registered {
		err := h.session.ApplicationCommandDelete(h.appID, "", cmd.ID)
		if err != nil {
			log.Printf("error removing command %s: %s", cmd.Name, err)
		}
	}
}

// modalPrefixLen returns the length up to the first ':' separator in a modal custom ID,
// or the full length if no separator is found.
func modalPrefixLen(customID string) int {
	for i, c := range customID {
		if c == ':' {
			return i
		}
	}

	return len(customID)
}
