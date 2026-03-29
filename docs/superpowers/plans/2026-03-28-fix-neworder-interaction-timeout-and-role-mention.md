# Fix /neworder Interaction Timeout & Role Mention

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix two bugs in the /neworder command: (1) Discord interaction times out because the bot responds after slow work completes, (2) tag mentions render as plain text instead of actual Discord role mentions.

**Architecture:** Add deferred interaction response (`DeferredChannelMessageWithSource`) to acknowledge within 3s, then do async work and follow up with `InteractionResponseEdit`. For role mentions, add a tag-to-role-ID mapping via environment variable, and format mentions as `<@&ROLE_ID>`.

**Tech Stack:** Go, discordgo, testify

---

## File Structure

| File | Action | Responsibility |
|---|---|---|
| `gateway/discord/command/respond.go` | Modify | Add `respondDeferred` and `respondDeferredEdit` helpers |
| `gateway/discord/command/new_order_handler.go` | Modify | Use deferred response pattern in `handleNewOrder` |
| `config/config.go` | Modify | Add `TagRoleMap` field (parsed from env) |
| `domain/order.go` | No change | Tag type stays the same |
| `usecase/create_order.go` | Modify | Accept and use role ID map for building thread message |
| `usecase/create_order_test.go` | Modify | Update tests for new role mention format |
| `main.go` | Modify | Pass `TagRoleMap` to `CreateOrder` use case |

---

## Task 1: Add Deferred Response Helpers

**Files:**
- Modify: `gateway/discord/command/respond.go`

- [ ] **Step 1: Add `respondDeferred` function**

```go
func respondDeferred(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("error deferring interaction response: %s", err)
	}
}
```

- [ ] **Step 2: Add `editDeferredResponse` function**

```go
func editDeferredResponse(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
	if err != nil {
		log.Printf("error editing deferred response: %s", err)
	}
}
```

**Note:** These are thin wrappers around discordgo API calls — not unit-tested, consistent with the existing `respondError`/`respondSuccess` pattern. Verified at integration level.

- [ ] **Step 3: Commit**

```bash
git add gateway/discord/command/respond.go
git commit -m "feat: add deferred interaction response helpers"
```

---

## Task 2: Use Deferred Response in /neworder Handler

**Behavioral change:** With deferred responses, error messages (e.g., "建立訂單失敗") will be visible to all channel members. Previously, `respondError` used ephemeral flags (visible only to the invoking user). This is a known tradeoff — ephemeral status cannot be changed after deferring. Acceptable for this bot's use case (admin-only command in a private guild).

**Files:**
- Modify: `gateway/discord/command/new_order_handler.go`

- [ ] **Step 1: Change `handleNewOrder` to defer first, then execute, then edit**

Replace the current flow:
```
execute -> respondError/respondSuccess
```
With:
```
respondDeferred -> execute -> editDeferredResponse (success or error)
```

Updated `handleNewOrder`:

```go
func handleNewOrder(
	s *discordgo.Session, i *discordgo.InteractionCreate, uc port.OrderCreator,
) {
	respondDeferred(s, i)

	opts := i.ApplicationCommandData().Options
	optMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(opts))

	for _, opt := range opts {
		optMap[opt.Name] = opt
	}

	order := domain.Order{}

	if v, ok := optMap["ordertitle"]; ok {
		order.ThreadName = v.StringValue()
	}

	if v, ok := optMap["deadline"]; ok {
		order.Deadline = v.StringValue()
	}

	if v, ok := optMap["shopurl"]; ok {
		order.ShopURL = v.StringValue()
	}

	if v, ok := optMap["tags"]; ok {
		order.Tag = domain.Tag(v.StringValue())
	}

	err := uc.Execute(context.Background(), i.ChannelID, order)
	if err != nil {
		log.Printf("create order failed: %s", err)
		editDeferredResponse(s, i, "建立訂單失敗")

		return
	}

	editDeferredResponse(s, i, "訂單已建立: "+order.ThreadName)
}
```

- [ ] **Step 2: Run existing tests to ensure no breakage**

Run: `go test ./...`
Expected: All existing tests PASS (handler tests are at integration level, mock-based usecase tests unaffected)

- [ ] **Step 3: Commit**

```bash
git add gateway/discord/command/new_order_handler.go
git commit -m "fix: defer interaction response in /neworder to avoid 3s timeout"
```

---

## Task 3: Add Tag-to-Role-ID Mapping via Config

**Files:**
- Modify: `config/config.go`

- [ ] **Step 1: Add `TagRoleMap` field to Config**

Add a new field and parser. The env var format is comma-separated `tag=roleID` pairs:

```
TAG_ROLE_MAP=315pro=123456789,学マス=987654321,283pro=111222333,346pro=444555666,765pro=777888999
```

Add `"strings"` to the import block (existing imports: `"fmt"`, `"os"`, `"strconv"`).

Add to `Config` struct:
```go
TagRoleMap map[string]string
```

Add parsing function:
```go
func parseTagRoleMap(raw string) map[string]string {
	m := make(map[string]string)
	if raw == "" {
		return m
	}

	for _, pair := range strings.Split(raw, ",") {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return m
}
```

Call in `Load()`:
```go
cfg.TagRoleMap = parseTagRoleMap(os.Getenv("TAG_ROLE_MAP"))
```

- [ ] **Step 2: Write unit test for `parseTagRoleMap`**

Create `config/config_test.go`:

```go
package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTagRoleMap(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want map[string]string
	}{
		{"empty", "", map[string]string{}},
		{"single", "283pro=111", map[string]string{"283pro": "111"}},
		{"multiple", "283pro=111,315pro=222", map[string]string{"283pro": "111", "315pro": "222"}},
		{"with spaces", " 283pro = 111 , 315pro = 222 ", map[string]string{"283pro": "111", "315pro": "222"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTagRoleMap(tt.raw)
			require.Equal(t, tt.want, got)
		})
	}
}
```

- [ ] **Step 3: Run test**

Run: `go test ./config/ -v`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add config/config.go config/config_test.go
git commit -m "feat: add TAG_ROLE_MAP env var for tag-to-Discord-role-ID mapping"
```

---

## Task 4: Update `buildThreadMessage` to Use Role Mentions

**Files:**
- Modify: `usecase/create_order.go`
- Modify: `usecase/create_order_test.go`

- [ ] **Step 1: Update test expectations first (TDD)**

Update `TestCreateOrder_Success_AllFields` — expected message should use `<@&ROLE_ID>` format:

```go
func TestCreateOrder_Success_AllFields(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)

	order := domain.Order{
		ThreadName: "test order",
		Deadline:   "2026-04-01",
		ShopURL:    "https://shop.example.com",
		Tag:        domain.Tag315Pro,
	}

	tagRoleMap := map[string]string{"315pro": "123456"}
	expectedMessage := "https://shop.example.com\n<@&123456>\n截止時間: 2026-04-01"

	tc.On("CreateThread", mock.Anything, "ch-1", "test order", expectedMessage).
		Return(nil)
	repo.On("CreateOrder", mock.Anything, order).
		Return(nil)

	uc := usecase.NewCreateOrder(repo, tc, tagRoleMap)

	err := uc.Execute(context.Background(), "ch-1", order)

	require.NoError(t, err)
}
```

Update `TestCreateOrder_Success_TagOnly` — tag with role ID:

```go
func TestCreateOrder_Success_TagOnly(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)

	order := domain.Order{
		ThreadName: "tag order",
		Tag:        domain.TagGakumas,
	}

	tagRoleMap := map[string]string{"学マス": "789012"}
	expectedMessage := "<@&789012>"

	tc.On("CreateThread", mock.Anything, "ch-1", "tag order", expectedMessage).
		Return(nil)
	repo.On("CreateOrder", mock.Anything, order).
		Return(nil)

	uc := usecase.NewCreateOrder(repo, tc, tagRoleMap)

	err := uc.Execute(context.Background(), "ch-1", order)

	require.NoError(t, err)
}
```

Add `TestCreateOrder_Success_TagWithoutRoleID` — tag not in map falls back to plain text:

```go
func TestCreateOrder_Success_TagWithoutRoleID(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)

	order := domain.Order{
		ThreadName: "fallback order",
		Tag:        domain.Tag283Pro,
	}

	tagRoleMap := map[string]string{} // empty map, no role ID for 283pro
	expectedMessage := "@283pro"

	tc.On("CreateThread", mock.Anything, "ch-1", "fallback order", expectedMessage).
		Return(nil)
	repo.On("CreateOrder", mock.Anything, order).
		Return(nil)

	uc := usecase.NewCreateOrder(repo, tc, tagRoleMap)

	err := uc.Execute(context.Background(), "ch-1", order)

	require.NoError(t, err)
}
```

Update all other tests that call `NewCreateOrder(repo, tc)` to `NewCreateOrder(repo, tc, nil)`. These are:
- `TestCreateOrder_MissingOrderTitle`
- `TestCreateOrder_ThreadCreationError`
- `TestCreateOrder_NotionError`
- `TestCreateOrder_Success_OnlyTitle`
- `TestCreateOrder_Success_PartialFields`
- `TestCreateOrder_Success_ShopURLOnly`

Search-and-replace: `NewCreateOrder(repo, tc)` → `NewCreateOrder(repo, tc, nil)` for these 6 tests.

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./usecase/ -v`
Expected: FAIL — `NewCreateOrder` doesn't accept `tagRoleMap` yet

- [ ] **Step 3: Update `CreateOrder` to accept and use `tagRoleMap`**

In `usecase/create_order.go`:

```go
type CreateOrder struct {
	repo          port.OrderRepository
	threadCreator port.ThreadCreator
	tagRoleMap    map[string]string
}

func NewCreateOrder(
	repo port.OrderRepository, threadCreator port.ThreadCreator, tagRoleMap map[string]string,
) *CreateOrder {
	return &CreateOrder{
		repo:          repo,
		threadCreator: threadCreator,
		tagRoleMap:    tagRoleMap,
	}
}
```

Update `Execute` to pass `tagRoleMap`:
```go
message := buildThreadMessage(order, uc.tagRoleMap)
```

Update `buildThreadMessage`:
```go
func buildThreadMessage(order domain.Order, tagRoleMap map[string]string) string {
	var lines []string

	if order.ShopURL != "" {
		lines = append(lines, order.ShopURL)
	}

	if order.Tag != "" {
		if roleID, ok := tagRoleMap[string(order.Tag)]; ok {
			lines = append(lines, fmt.Sprintf("<@&%s>", roleID))
		} else {
			lines = append(lines, fmt.Sprintf("@%s", order.Tag))
		}
	}

	if order.Deadline != "" {
		lines = append(lines, fmt.Sprintf("截止時間: %s", order.Deadline))
	}

	return strings.Join(lines, "\n")
}
```

- [ ] **Step 4: Run tests**

Run: `go test ./usecase/ -v`
Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add usecase/create_order.go usecase/create_order_test.go
git commit -m "fix: use Discord role mention format for tag in thread message"
```

---

## Task 5: Wire TagRoleMap in main.go

**Files:**
- Modify: `main.go`

- [ ] **Step 1: Pass `cfg.TagRoleMap` to `NewCreateOrder`**

Change the `usecase.NewCreateOrder(orderRepo, threadCreator)` call to include `cfg.TagRoleMap`:
```go
createOrderUC := usecase.NewCreateOrder(orderRepo, threadCreator, cfg.TagRoleMap)
```

- [ ] **Step 2: Build and verify**

Run: `go build ./...`
Expected: Compiles with no errors

- [ ] **Step 3: Run all tests**

Run: `go test ./...`
Expected: All PASS

- [ ] **Step 4: Commit**

```bash
git add main.go
git commit -m "feat: wire TAG_ROLE_MAP config to create order use case"
```

---

## Task 6: Update .env and CLAUDE.md

**Files:**
- Modify: `.env` (add `TAG_ROLE_MAP`)
- Modify: `CLAUDE.md` (add `TAG_ROLE_MAP` to env var table)

- [ ] **Step 1: Add to `.env`**

```
TAG_ROLE_MAP=315pro=ROLE_ID_HERE,学マス=ROLE_ID_HERE,283pro=ROLE_ID_HERE,346pro=ROLE_ID_HERE,765pro=ROLE_ID_HERE
```

(User fills in actual role IDs from Discord server settings → Roles)

- [ ] **Step 2: Add to CLAUDE.md env var table**

Add row:
```
| `TAG_ROLE_MAP` | Comma-separated tag=roleID pairs for Discord role mentions |
```

- [ ] **Step 3: Commit**

```bash
git add CLAUDE.md
git commit -m "docs: add TAG_ROLE_MAP to environment variable documentation"
```

---

## Summary of Changes

| Bug | Root Cause | Fix |
|---|---|---|
| "アプリケーションが応答しませんでした" | `InteractionRespond` called after slow Notion/thread work exceeds 3s limit | Defer response immediately, edit after work completes |
| "@283pro" plain text | `fmt.Sprintf("@%s", tag)` outputs literal text, not Discord mention | Use `<@&ROLE_ID>` format with role ID lookup from `TAG_ROLE_MAP` env |
