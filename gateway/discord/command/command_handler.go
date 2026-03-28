package command

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// Handler manages Discord slash command registration and dispatch.
type Handler struct {
	s       *discordgo.Session
	appID   string
	cmds    []*discordgo.ApplicationCommand
	routes  map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	cleanup func()
}

func NewHandler(s *discordgo.Session, appID string) *Handler {
	h := &Handler{
		s:      s,
		appID:  appID,
		routes: make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)),
	}

	h.cleanup = s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		name := i.ApplicationCommandData().Name
		if handler, ok := h.routes[name]; ok {
			handler(s, i)
		}
	})

	return h
}

func (h *Handler) RegisterCommand(
	cmd *discordgo.ApplicationCommand,
	handler func(s *discordgo.Session, i *discordgo.InteractionCreate),
) error {
	created, err := h.s.ApplicationCommandCreate(h.appID, "", cmd)
	if err != nil {
		return fmt.Errorf("error registering command %s: %w", cmd.Name, err)
	}

	h.cmds = append(h.cmds, created)
	h.routes[cmd.Name] = handler

	return nil
}

func (h *Handler) UnregisterAll() error {
	if h.cleanup != nil {
		h.cleanup()
	}

	for _, cmd := range h.cmds {
		err := h.s.ApplicationCommandDelete(h.appID, "", cmd.ID)
		if err != nil {
			log.Printf("error deleting command %s: %s", cmd.Name, err)
		}
	}

	return nil
}
