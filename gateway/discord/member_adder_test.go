package discord

import (
	"context"
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
)

type mockMemberAdderSession struct {
	guildMembersFn func(
		guildID, after string, limit int, options ...discordgo.RequestOption,
	) ([]*discordgo.Member, error)
	threadMemberAddFn func(threadID, memberID string, options ...discordgo.RequestOption) error
	addedMembers      []string
}

func (m *mockMemberAdderSession) GuildMembers(
	guildID, after string, limit int, options ...discordgo.RequestOption,
) ([]*discordgo.Member, error) {
	return m.guildMembersFn(guildID, after, limit, options...)
}

func (m *mockMemberAdderSession) ThreadMemberAdd(
	threadID, memberID string, options ...discordgo.RequestOption,
) error {
	m.addedMembers = append(m.addedMembers, memberID)

	return m.threadMemberAddFn(threadID, memberID, options...)
}

func makeMember(id string, roles []string) *discordgo.Member {
	return &discordgo.Member{
		User:  &discordgo.User{ID: id},
		Roles: roles,
	}
}

func TestMemberAdder_Success(t *testing.T) {
	targetRole := "role-abc"
	members := []*discordgo.Member{
		makeMember("user-1", []string{targetRole}),
		makeMember("user-2", []string{"other-role"}),
		makeMember("user-3", []string{targetRole, "other-role"}),
	}

	mock := &mockMemberAdderSession{
		guildMembersFn: func(_, _ string, _ int, _ ...discordgo.RequestOption) ([]*discordgo.Member, error) {
			return members, nil
		},
		threadMemberAddFn: func(_, _ string, _ ...discordgo.RequestOption) error {
			return nil
		},
	}

	ma := &MemberAdder{s: mock, guildID: "guild-1"}

	err := ma.AddRoleMembersToThread(context.Background(), "thread-1", targetRole)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mock.addedMembers) != 2 {
		t.Fatalf("expected 2 members added, got %d", len(mock.addedMembers))
	}

	added := map[string]bool{}
	for _, id := range mock.addedMembers {
		added[id] = true
	}

	if !added["user-1"] {
		t.Error("expected user-1 to be added")
	}

	if !added["user-3"] {
		t.Error("expected user-3 to be added")
	}
}

func TestMemberAdder_GuildMembersError(t *testing.T) {
	fetchErr := errors.New("discord API error")

	mock := &mockMemberAdderSession{
		guildMembersFn: func(_, _ string, _ int, _ ...discordgo.RequestOption) ([]*discordgo.Member, error) {
			return nil, fetchErr
		},
		threadMemberAddFn: func(_, _ string, _ ...discordgo.RequestOption) error {
			return nil
		},
	}

	ma := &MemberAdder{s: mock, guildID: "guild-1"}

	err := ma.AddRoleMembersToThread(context.Background(), "thread-1", "role-abc")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, fetchErr) {
		t.Errorf("expected wrapped fetchErr, got: %v", err)
	}
}

func TestMemberAdder_ThreadMemberAddPartialFailure(t *testing.T) {
	targetRole := "role-abc"
	members := []*discordgo.Member{
		makeMember("user-1", []string{targetRole}),
		makeMember("user-2", []string{targetRole}),
	}

	addErr := errors.New("add failed")

	mock := &mockMemberAdderSession{
		guildMembersFn: func(_, _ string, _ int, _ ...discordgo.RequestOption) ([]*discordgo.Member, error) {
			return members, nil
		},
		threadMemberAddFn: func(_, memberID string, _ ...discordgo.RequestOption) error {
			if memberID == "user-1" {
				return addErr
			}

			return nil
		},
	}

	ma := &MemberAdder{s: mock, guildID: "guild-1"}

	err := ma.AddRoleMembersToThread(context.Background(), "thread-1", targetRole)
	if err != nil {
		t.Fatalf("expected no error (failures are logged only), got: %v", err)
	}

	if len(mock.addedMembers) != 2 {
		t.Fatalf("expected both members to be attempted, got %d", len(mock.addedMembers))
	}
}

func TestMemberAdder_NoMembersWithRole(t *testing.T) {
	members := []*discordgo.Member{
		makeMember("user-1", []string{"other-role"}),
		makeMember("user-2", []string{"another-role"}),
	}

	mock := &mockMemberAdderSession{
		guildMembersFn: func(_, _ string, _ int, _ ...discordgo.RequestOption) ([]*discordgo.Member, error) {
			return members, nil
		},
		threadMemberAddFn: func(_, _ string, _ ...discordgo.RequestOption) error {
			return nil
		},
	}

	ma := &MemberAdder{s: mock, guildID: "guild-1"}

	err := ma.AddRoleMembersToThread(context.Background(), "thread-1", "role-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mock.addedMembers) != 0 {
		t.Fatalf("expected no members added, got %d", len(mock.addedMembers))
	}
}
