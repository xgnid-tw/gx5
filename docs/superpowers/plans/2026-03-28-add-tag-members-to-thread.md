# Add Tag Members to Thread on Order Creation

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** When `/neworder` creates a thread with a tag, automatically add all Discord guild members who have that tag's role to the thread.

**Architecture:** Add a `MemberAdder` port that the `CreateOrder` use case calls after thread creation. The gateway implementation uses `GuildMembers` to find members with the role, then `ThreadMemberAdd` to add each one. Failures to add individual members are logged but do not block the order creation.

**Tech Stack:** Go, discordgo v0.29.0, testify

---

## File Structure

| File | Action | Responsibility |
|---|---|---|
| `port/member_adder.go` | Create | `MemberAdder` interface |
| `gateway/discord/member_adder.go` | Create | Discord implementation: fetch guild members by role, add to thread |
| `usecase/create_order.go` | Modify | Call `MemberAdder` after thread creation |
| `usecase/create_order_test.go` | Modify | Add tests for member addition flow |
| `mocks/MemberAdder.go` | Create | Mock for testing |
| `config/config.go` | Modify | Add `DiscordGuildID` field |
| `main.go` | Modify | Wire `MemberAdder` and guild ID |
| `.circleci/config.yml` | Modify | Add `DISCORD_GUILD_ID` env var |
| `CLAUDE.md` | Modify | Document `DISCORD_GUILD_ID` env var |

---

## Task 1: Add `DiscordGuildID` to Config

**Files:**
- Modify: `config/config.go`

- [ ] **Step 1: Add field to Config struct**

```go
DiscordGuildID string
```

Add to the `Load()` assignments:
```go
DiscordGuildID: os.Getenv("DISCORD_GUILD_ID"),
```

Add validation:
```go
if cfg.DiscordGuildID == "" {
    return Config{}, fmt.Errorf("DISCORD_GUILD_ID is required")
}
```

- [ ] **Step 2: Verify build**

Run: `go build ./config/...`

- [ ] **Step 3: Commit**

```bash
git add config/config.go
git commit -m "feat: add DISCORD_GUILD_ID to config"
```

---

## Task 2: Create `MemberAdder` Port

**Files:**
- Create: `port/member_adder.go`

- [ ] **Step 1: Define the interface**

```go
package port

import "context"

// MemberAdder adds guild members with a specific role to a thread.
type MemberAdder interface {
    AddRoleMembersToThread(ctx context.Context, threadID string, roleID string) error
}
```

- [ ] **Step 2: Commit**

```bash
git add port/member_adder.go
git commit -m "feat: add MemberAdder port interface"
```

---

## Task 3: Update `ThreadCreator` to Return Thread ID

The `CreateThread` port currently returns `error`. We need the thread ID to add members later.

**Files:**
- Modify: `port/thread_creator.go`
- Modify: `gateway/discord/thread_creator.go`
- Regenerate: `mocks/ThreadCreator.go`

- [ ] **Step 1: Update port interface**

`port/thread_creator.go`:
```go
type ThreadCreator interface {
    CreateThread(ctx context.Context, channelID string, name string, message string) (string, error)
}
```

- [ ] **Step 2: Update gateway implementation**

`gateway/discord/thread_creator.go` — update return type and all return paths:

```go
func (tc *ThreadCreator) CreateThread(_ context.Context, channelID string, name string, message string) (string, error) {
    thread, err := tc.s.ThreadStartComplex(channelID, &discordgo.ThreadStart{
        Name: name,
        Type: discordgo.ChannelTypeGuildPublicThread,
    })
    if err != nil {
        return "", fmt.Errorf("error creating thread: %w", err)
    }

    if message != "" {
        _, err = tc.s.ChannelMessageSendComplex(thread.ID, &discordgo.MessageSend{
            Content: message,
            AllowedMentions: &discordgo.MessageAllowedMentions{
                Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeRoles},
            },
        })
        if err != nil {
            return "", fmt.Errorf("error sending thread message: %w", err)
        }
    }

    return thread.ID, nil
}
```

- [ ] **Step 3: Regenerate mock**

Run: `go run github.com/vektra/mockery/v2@latest --name=ThreadCreator --dir=./port --output=./mocks --outpkg=mocks`

- [ ] **Step 4: Update `usecase/create_order.go`** to capture thread ID

Change:
```go
err := uc.threadCreator.CreateThread(ctx, channelID, order.ThreadName, message)
```
To:
```go
_, err := uc.threadCreator.CreateThread(ctx, channelID, order.ThreadName, message)
```

(Using `_` for now; Task 5 will use the thread ID.)

- [ ] **Step 5: Update all test expectations in `usecase/create_order_test.go`**

All `CreateThread` mock `.Return(...)` calls must change:

| Test function | Old return | New return |
|---|---|---|
| `TestCreateOrder_ThreadCreationError` | `.Return(errors.New("discord error"))` | `.Return("", errors.New("discord error"))` |
| `TestCreateOrder_NotionError` | `.Return(nil)` | `.Return("thread-id", nil)` |
| `TestCreateOrder_Success_AllFields` | `.Return(nil)` | `.Return("thread-id", nil)` |
| `TestCreateOrder_Success_OnlyTitle` | `.Return(nil)` | `.Return("thread-id", nil)` |
| `TestCreateOrder_Success_PartialFields` | `.Return(nil)` | `.Return("thread-id", nil)` |
| `TestCreateOrder_Success_ShopURLOnly` | `.Return(nil)` | `.Return("thread-id", nil)` |
| `TestCreateOrder_Success_TagOnly` | `.Return(nil)` | `.Return("thread-id", nil)` |
| `TestCreateOrder_Success_TagWithoutRoleID` | `.Return(nil)` | `.Return("thread-id", nil)` |

(`TestCreateOrder_MissingOrderTitle` does not call `CreateThread`, so no change.)

- [ ] **Step 6: Verify build and tests pass**

Run: `go build ./... && go test ./usecase/ -v`
Expected: All PASS

- [ ] **Step 7: Commit**

```bash
git add port/thread_creator.go gateway/discord/thread_creator.go mocks/ThreadCreator.go usecase/create_order.go usecase/create_order_test.go
git commit -m "refactor: return thread ID from CreateThread"
```

---

## Task 4: Generate MemberAdder Mock

**Files:**
- Create: `mocks/MemberAdder.go`

- [ ] **Step 1: Generate mock using mockery**

Run: `go run github.com/vektra/mockery/v2@latest --name=MemberAdder --dir=./port --output=./mocks --outpkg=mocks`

This keeps consistency with other mockery-generated mocks in the project.

- [ ] **Step 2: Verify build**

Run: `go build ./mocks/...`

- [ ] **Step 3: Commit**

```bash
git add mocks/MemberAdder.go
git commit -m "feat: add MemberAdder mock"
```

---

## Task 5: Update `CreateOrder` Use Case to Add Members (TDD)

**Files:**
- Modify: `usecase/create_order.go`
- Modify: `usecase/create_order_test.go`

- [ ] **Step 1: Update tests first**

**A. Update `NewCreateOrder` signature in ALL test functions.**

Change `NewCreateOrder(repo, tc, tagRoleMap)` to `NewCreateOrder(repo, tc, nil, tagRoleMap)` in these 6 tests (pass `nil` for `memberAdder`):
- `TestCreateOrder_MissingOrderTitle`
- `TestCreateOrder_ThreadCreationError`
- `TestCreateOrder_NotionError`
- `TestCreateOrder_Success_OnlyTitle`
- `TestCreateOrder_Success_PartialFields`
- `TestCreateOrder_Success_ShopURLOnly`

Change `NewCreateOrder(repo, tc, tagRoleMap)` to `NewCreateOrder(repo, tc, nil, tagRoleMap)` in this test (tag has no role mapping, so no member add expected):
- `TestCreateOrder_Success_TagWithoutRoleID`

**B. Add `MemberAdder` mock to tests with tags that have roleID mapping:**

For `TestCreateOrder_Success_AllFields`:
```go
ma := mocks.NewMemberAdder(t)
ma.On("AddRoleMembersToThread", mock.Anything, "thread-id", "123456").Return(nil)
uc := usecase.NewCreateOrder(repo, tc, ma, tagRoleMap)
```

For `TestCreateOrder_Success_TagOnly`:
```go
ma := mocks.NewMemberAdder(t)
ma.On("AddRoleMembersToThread", mock.Anything, "thread-id", "789012").Return(nil)
uc := usecase.NewCreateOrder(repo, tc, ma, tagRoleMap)
```

**C. Add new test for member addition failure (non-fatal):**

```go
func TestCreateOrder_MemberAddFailure_NonFatal(t *testing.T) {
    repo := mocks.NewOrderRepository(t)
    tc := mocks.NewThreadCreator(t)
    ma := mocks.NewMemberAdder(t)

    order := domain.Order{
        ThreadName: "test order",
        Tag:        domain.Tag315Pro,
    }

    tagRoleMap := map[string]string{"315pro": "123456"}

    tc.On("CreateThread", mock.Anything, "ch-1", "test order", "<@&123456>").
        Return("thread-id", nil)
    ma.On("AddRoleMembersToThread", mock.Anything, "thread-id", "123456").
        Return(errors.New("discord api error"))
    repo.On("CreateOrder", mock.Anything, order).
        Return(nil)

    uc := usecase.NewCreateOrder(repo, tc, ma, tagRoleMap)
    err := uc.Execute(context.Background(), "ch-1", order)

    require.NoError(t, err) // member add failure is non-fatal
}
```

- [ ] **Step 2: Verify tests fail (compilation error expected)**

Run: `go test ./usecase/ -v`
Expected: FAIL — `NewCreateOrder` signature mismatch (3 args vs 4)

- [ ] **Step 3: Update implementation**

In `usecase/create_order.go`:

```go
type CreateOrder struct {
    repo          port.OrderRepository
    threadCreator port.ThreadCreator
    memberAdder   port.MemberAdder
    tagRoleMap    map[string]string
}

func NewCreateOrder(
    repo port.OrderRepository, threadCreator port.ThreadCreator,
    memberAdder port.MemberAdder, tagRoleMap map[string]string,
) *CreateOrder {
    return &CreateOrder{
        repo:          repo,
        threadCreator: threadCreator,
        memberAdder:   memberAdder,
        tagRoleMap:    tagRoleMap,
    }
}

func (uc *CreateOrder) Execute(
    ctx context.Context, channelID string, order domain.Order,
) error {
    if order.ThreadName == "" {
        return fmt.Errorf("orderTitle is required")
    }

    message := buildThreadMessage(order, uc.tagRoleMap)

    threadID, err := uc.threadCreator.CreateThread(ctx, channelID, order.ThreadName, message)
    if err != nil {
        return fmt.Errorf("create thread: %w", err)
    }

    if order.Tag != "" {
        if roleID, ok := uc.tagRoleMap[string(order.Tag)]; ok && uc.memberAdder != nil {
            if addErr := uc.memberAdder.AddRoleMembersToThread(ctx, threadID, roleID); addErr != nil {
                log.Printf("add role members to thread: %s", addErr)
            }
        }
    }

    err = uc.repo.CreateOrder(ctx, order)
    if err != nil {
        return fmt.Errorf("create order record: %w", err)
    }

    return nil
}
```

Key design decisions:
- `memberAdder` can be `nil` (skips member addition)
- Member addition failure is **logged but non-fatal** (consistent with project pattern)
- Only attempts if tag has a role ID mapping

- [ ] **Step 4: Run tests**

Run: `go test ./usecase/ -v`
Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add usecase/create_order.go usecase/create_order_test.go
git commit -m "feat: add tag members to thread after order creation"
```

---

## Task 6: Implement Discord `MemberAdder` Gateway

**Files:**
- Create: `gateway/discord/member_adder.go`
- Create: `gateway/discord/member_adder_test.go`

- [ ] **Step 1: Write the implementation**

```go
package discord

import (
    "context"
    "fmt"
    "log"

    "github.com/bwmarrin/discordgo"
)

type memberAdderSession interface {
    GuildMembers(guildID string, after string, limit int, options ...discordgo.RequestOption) ([]*discordgo.Member, error)
    ThreadMemberAdd(threadID, memberID string, options ...discordgo.RequestOption) error
}

type MemberAdder struct {
    s       memberAdderSession
    guildID string
}

func NewMemberAdder(s *discordgo.Session, guildID string) *MemberAdder {
    return &MemberAdder{s: s, guildID: guildID}
}

func (ma *MemberAdder) AddRoleMembersToThread(_ context.Context, threadID string, roleID string) error {
    members, err := ma.fetchMembersWithRole(roleID)
    if err != nil {
        return fmt.Errorf("fetch members with role %s: %w", roleID, err)
    }

    for _, m := range members {
        if err := ma.s.ThreadMemberAdd(threadID, m.User.ID); err != nil {
            log.Printf("failed to add member %s to thread %s: %s", m.User.ID, threadID, err)
        }
    }

    return nil
}

func (ma *MemberAdder) fetchMembersWithRole(roleID string) ([]*discordgo.Member, error) {
    var result []*discordgo.Member
    after := ""
    const pageSize = 100

    for {
        members, err := ma.s.GuildMembers(ma.guildID, after, pageSize)
        if err != nil {
            return nil, err
        }

        for _, m := range members {
            for _, r := range m.Roles {
                if r == roleID {
                    result = append(result, m)
                    break
                }
            }
        }

        if len(members) < pageSize {
            break
        }

        after = members[len(members)-1].User.ID
    }

    return result, nil
}
```

**Note:** `GuildMembers` fetches ALL guild members and filters client-side. Discord API does not support server-side role filtering. This is fine for a small private guild. Requires **Server Members Intent** (privileged) — already enabled via `IntentsAll` in `main.go`.

- [ ] **Step 2: Write unit tests**

Create `gateway/discord/member_adder_test.go` with a mock `memberAdderSession`:

```go
package discord

import (
    "context"
    "errors"
    "testing"

    "github.com/bwmarrin/discordgo"
    "github.com/stretchr/testify/require"
)

type mockMemberAdderSession struct {
    guildMembersFn    func(guildID, after string, limit int, options ...discordgo.RequestOption) ([]*discordgo.Member, error)
    threadMemberAddFn func(threadID, memberID string, options ...discordgo.RequestOption) error
    addedMembers      []string
}

func (m *mockMemberAdderSession) GuildMembers(guildID, after string, limit int, options ...discordgo.RequestOption) ([]*discordgo.Member, error) {
    return m.guildMembersFn(guildID, after, limit, options...)
}

func (m *mockMemberAdderSession) ThreadMemberAdd(threadID, memberID string, options ...discordgo.RequestOption) error {
    m.addedMembers = append(m.addedMembers, memberID)
    return m.threadMemberAddFn(threadID, memberID, options...)
}
```

Tests to write:
- **Success**: 2 members with role, 1 without → only 2 added to thread
- **GuildMembers error**: returns error
- **ThreadMemberAdd partial failure**: one member fails to add, others succeed, no error returned (logged only)
- **No members with role**: no ThreadMemberAdd calls

- [ ] **Step 3: Run tests**

Run: `go test ./gateway/discord/ -v -run TestMemberAdder`
Expected: All PASS

- [ ] **Step 4: Commit**

```bash
git add gateway/discord/member_adder.go gateway/discord/member_adder_test.go
git commit -m "feat: implement Discord MemberAdder gateway with tests"
```

---

## Task 7: Wire in main.go and Update CI/Docs

**Files:**
- Modify: `main.go`
- Modify: `.circleci/config.yml`
- Modify: `CLAUDE.md`

- [ ] **Step 1: Wire `MemberAdder` in main.go**

After `threadCreator := discordgw.NewThreadCreator(dc)`:
```go
memberAdder := discordgw.NewMemberAdder(dc, cfg.DiscordGuildID)
```

Update `NewCreateOrder` call:
```go
createOrderUC := usecase.NewCreateOrder(orderRepo, threadCreator, memberAdder, cfg.TagRoleMap)
```

- [ ] **Step 2: Add `DISCORD_GUILD_ID` to CI config**

Add to `.circleci/config.yml` env creation:
```
echo "DISCORD_GUILD_ID=${DISCORD_GUILD_ID}" >> .env
```

- [ ] **Step 3: Add to CLAUDE.md env var table**

```
| `DISCORD_GUILD_ID` | Discord guild (server) ID |
```

- [ ] **Step 4: Build and test**

Run: `go build ./... && go test ./...`
Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add main.go .circleci/config.yml CLAUDE.md
git commit -m "feat: wire MemberAdder and add DISCORD_GUILD_ID config"
```

---

## Summary

| Step | What happens |
|---|---|
| User runs `/neworder` with tag "283pro" | Bot defers response |
| Bot creates thread | Thread ID returned |
| Bot looks up role ID for "283pro" | From `TAG_ROLE_MAP` |
| Bot fetches guild members with that role | `GuildMembers` + filter by role |
| Bot adds each member to thread | `ThreadMemberAdd` per member, failures logged |
| Bot creates Notion order record | Existing behavior |
| Bot edits deferred response | "訂單已建立: ..." |

**Prerequisites:**
- `DISCORD_GUILD_ID` env var set in CircleCI and `.env`
- **Server Members Intent** enabled in Discord Developer Portal (already on via `IntentsAll`)
