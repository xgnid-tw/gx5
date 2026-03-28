package command

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// Handler manages Discord application command registration and dispatch.
type Handler struct {
	session    *discordgo.Session
	appID      string
	commands   []*discordgo.ApplicationCommand
	handlers   map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	registered []*discordgo.ApplicationCommand
}

func NewHandler(session *discordgo.Session, appID string) *Handler {
	h := &Handler{
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
			customID := i.ModalSubmitData().CustomID
			prefix := customID[:modalPrefixLen(customID)]

			if handler, ok := h.handlers["modal:"+prefix]; ok {
				handler(s, i)
			}
		case discordgo.InteractionPing,
			discordgo.InteractionMessageComponent,
			discordgo.InteractionApplicationCommandAutocomplete:
			// not handled
		}
	})

	return h
}

// RegisterCommand registers an application command with its handler.
func (h *Handler) RegisterCommand(
	cmd *discordgo.ApplicationCommand,
	handler func(s *discordgo.Session, i *discordgo.InteractionCreate),
) {
	h.commands = append(h.commands, cmd)
	h.handlers[cmd.Name] = handler
}

// RegisterModalHandler registers a handler for modal submissions with a given prefix.
func (h *Handler) RegisterModalHandler(
	prefix string,
	handler func(s *discordgo.Session, i *discordgo.InteractionCreate),
) {
	h.handlers["modal:"+prefix] = handler
}

// SyncCommands creates all registered commands with the Discord API.
func (h *Handler) SyncCommands() error {
	for _, cmd := range h.commands {
		registered, err := h.session.ApplicationCommandCreate(h.appID, "", cmd)
		if err != nil {
			return fmt.Errorf("create command %q: %w", cmd.Name, err)
		}

		h.registered = append(h.registered, registered)
	}

	return nil
}

// UnregisterAll removes all registered commands from the Discord API.
func (h *Handler) UnregisterAll() {
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
