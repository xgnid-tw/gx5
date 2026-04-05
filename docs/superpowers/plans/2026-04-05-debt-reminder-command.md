# Debt Reminder Slash Command Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the cron-based unpaid notification (UC-001) with an on-demand `/debt-reminder` Discord slash command (UC-004) that runs immediately and schedules a one-shot delayed production run.

**Architecture:** The existing `NotifyUnpaid` use case is simplified (remove day guard, clock dependency) and its `Execute` method gains a `debug` parameter. A new port interface `DebtReminder` wraps the use case for the command handler. A new slash command handler schedules the delayed run via gocron. The cron job and `WORKER_CORNTAB` config are removed.

**Tech Stack:** Go, discordgo, gocron/v2, notionapi, testify

---

## File Map

| Action | File | Responsibility |
|--------|------|---------------|
| Modify | `port/notifier.go` | Add `debug` param to `Notifier.Notify` |
| Modify | `gateway/discord/notifier.go` | Accept `debug` per-call instead of struct field |
| Modify | `gateway/discord/mock_test.go` | Update `newTestNotifier` (no debug field) |
| Modify | `gateway/discord/notifier_test.go` | Pass `debug` to `Notify` calls |
| Modify | `mocks/Notifier.go` | Regenerate (or hand-update) for new signature |
| Modify | `usecase/notify_unpaid.go` | Remove day guard, clock; add `debug` param to `Execute` |
| Modify | `usecase/notify_unpaid_test.go` | Remove day-guard tests, clock setup; pass `debug` param |
| Create | `port/debt_reminder.go` | `DebtReminder` interface with `Execute(ctx, debug)` |
| Create | `gateway/discord/command/debt_reminder_handler.go` | `/debt-reminder` slash command handler |
| Modify | `config/config.go` | Remove `WorkerCrontab` field and validation |
| Modify | `config/config_test.go` | Remove `WORKER_CORNTAB` test cases |
| Modify | `main.go` | Remove cron job, wire new command, pass scheduler to handler |

---

### Task 1: Update `port/notifier.go` — add `debug` parameter

**Files:**
- Modify: `port/notifier.go`

- [ ] **Step 1: Update the Notifier interface**

Change `Notify` to accept a `debug bool` parameter:

```go
package port

import (
	"context"

	"github.com/xgnid-tw/gx5/domain"
)

type Notifier interface {
	Notify(ctx context.Context, user domain.User, debug bool) error
}
```

- [ ] **Step 2: Commit**

```bash
git add port/notifier.go
git commit -m "refactor: add debug parameter to Notifier.Notify interface"
```

---

### Task 2: Update `gateway/discord/notifier.go` — per-call debug

**Files:**
- Modify: `gateway/discord/notifier.go`
- Modify: `gateway/discord/mock_test.go`
- Modify: `gateway/discord/notifier_test.go`

- [ ] **Step 1: Update Notifier struct and methods**

Remove the `debug` field from the struct. Remove it from `NewNotifier`. Accept `debug` as a parameter in `Notify` and `sendDM`:

```go
type Notifier struct {
	s            discordSession
	logChannelID string
}

func NewNotifier(s *discordgo.Session, logChannelID string) *Notifier {
	return &Notifier{s: s, logChannelID: logChannelID}
}

func (n *Notifier) Notify(_ context.Context, user domain.User, debug bool) error {
	message := fmt.Sprintf(
		"[欠費提醒] https://www.notion.so/%s (如果有漏登聯絡一下XG) ",
		user.NotionID,
	)

	return n.sendDM(user.DiscordID, message, debug)
}

func (n *Notifier) sendDM(discordID string, message string, debug bool) error {
	channel, err := n.s.UserChannelCreate(discordID)
	if err != nil {
		return fmt.Errorf("error creating channel: %w", err)
	}

	_, err = n.s.ChannelMessageSend(n.logChannelID, message)
	if err != nil {
		return fmt.Errorf("error sending to log channel: %w", err)
	}

	if debug {
		log.Print("debug mode on")
		return nil
	}

	_, err = n.s.ChannelMessageSend(channel.ID, message)
	if err != nil {
		return fmt.Errorf("error sending dm: %w", err)
	}

	return nil
}
```

- [ ] **Step 2: Update mock_test.go — remove debug from newTestNotifier**

```go
func newTestNotifier(s discordSession, logChannelID string) *Notifier {
	return &Notifier{s: s, logChannelID: logChannelID}
}
```

- [ ] **Step 3: Update notifier_test.go — pass debug to Notify**

Replace all `n.Notify(context.Background(), testUser)` calls:
- `TestNotify_DebugMode_SkipsDM`: use `newTestNotifier(m, "log-chan")`, call `n.Notify(context.Background(), testUser, true)`
- `TestNotify_NormalMode_SendsDM`: use `newTestNotifier(m, "log-chan")`, call `n.Notify(context.Background(), testUser, false)`
- `TestNotify_UserChannelCreateFails`: use `newTestNotifier(m, "log-chan")`, call `n.Notify(context.Background(), testUser, false)`
- `TestNotify_LogChannelSendFails`: use `newTestNotifier(m, "log-chan")`, call `n.Notify(context.Background(), testUser, false)`
- `TestNotify_DMSendFails`: use `newTestNotifier(m, "log-chan")`, call `n.Notify(context.Background(), testUser, false)`

- [ ] **Step 4: Run tests**

Run: `go test ./gateway/discord/...`
Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add gateway/discord/notifier.go gateway/discord/mock_test.go gateway/discord/notifier_test.go
git commit -m "refactor: make Notifier.debug a per-call parameter"
```

---

### Task 3: Update `usecase/notify_unpaid.go` — remove day guard and clock, add debug param

**Files:**
- Modify: `usecase/notify_unpaid.go`
- Modify: `usecase/notify_unpaid_test.go`

- [ ] **Step 1: Simplify NotifyUnpaid**

Remove `location`, `Clock` fields. Remove day guard. Add `debug` param to `Execute`, pass it to `notifier.Notify`:

```go
package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/xgnid-tw/gx5/domain"
	"github.com/xgnid-tw/gx5/port"
)

const (
	twdNotificationThreshold = 2000
	jpyNotificationThreshold = 8000
)

var notificationAmountLimit = map[domain.Currency]float64{
	domain.CurrencyTWD: twdNotificationThreshold,
	domain.CurrencyJPY: jpyNotificationThreshold,
}

type NotifyUnpaid struct {
	repo       port.UserRepository
	notifier   port.Notifier
	othersDBID string
}

func NewNotifyUnpaid(
	repo port.UserRepository, notifier port.Notifier,
	othersDBID string,
) *NotifyUnpaid {
	return &NotifyUnpaid{
		repo: repo, notifier: notifier,
		othersDBID: othersDBID,
	}
}

func (uc *NotifyUnpaid) Execute(ctx context.Context, debug bool) error {
	users, err := uc.repo.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("get users: %w", err)
	}

	for _, u := range users {
		shouldNotify, err := uc.shouldNotifyUser(ctx, u)
		if err != nil {
			return err
		}

		if shouldNotify {
			err = uc.notifier.Notify(ctx, *u, debug)
			if err != nil {
				log.Printf("notify %s: %s", u.Name, err)
			}
		}
	}

	return nil
}
```

(`shouldNotifyUser` stays unchanged.)

- [ ] **Step 2: Update tests**

Remove `TestExecute_NotFirstOfMonth_SkipsAll` (day guard no longer exists).

Remove `mockClock` helper, `testLocation` var, and all `clock` imports.

Update `NewNotifyUnpaid` calls — remove the `testLocation` parameter:
```go
uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID)
```

Remove all `uc.Clock = mockClock(...)` lines.

Update `uc.Execute(context.Background())` → `uc.Execute(context.Background(), false)` in all tests.

Update `notifier.On("Notify", ...)` — add `false` as the third argument after `*user`:
```go
notifier.On("Notify", mock.Anything, *user, false).Return(nil)
```

- [ ] **Step 3: Regenerate mocks**

Run: `go generate ./mocks/... || mockery`

If mockery is not available, hand-update `mocks/Notifier.go`:
- Change `Notify` signature to `Notify(ctx context.Context, user domain.User, debug bool) error`
- Update `_m.Called(ctx, user, debug)`
- Update the type assertion to `func(context.Context, domain.User, bool) error`

- [ ] **Step 4: Run tests**

Run: `go test ./usecase/... ./gateway/discord/...`
Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add usecase/notify_unpaid.go usecase/notify_unpaid_test.go mocks/Notifier.go
git commit -m "refactor: remove day guard and clock from NotifyUnpaid, add debug param"
```

---

### Task 4: Create `port/debt_reminder.go` — DebtReminder interface

**Files:**
- Create: `port/debt_reminder.go`

- [ ] **Step 1: Create the interface**

```go
package port

import "context"

type DebtReminder interface {
	Execute(ctx context.Context, debug bool) error
}
```

- [ ] **Step 2: Verify NotifyUnpaid satisfies the interface**

Run: `go build ./...`
Expected: compiles without error (NotifyUnpaid already has `Execute(ctx, debug)`)

- [ ] **Step 3: Commit**

```bash
git add port/debt_reminder.go
git commit -m "feat: add DebtReminder port interface"
```

---

### Task 5: Create `/debt-reminder` slash command handler

**Files:**
- Create: `gateway/discord/command/debt_reminder_handler.go`

- [ ] **Step 1: Create the handler**

```go
package command

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"

	"github.com/xgnid-tw/gx5/port"
)

const (
	debtReminderCommandName = "debt-reminder"
	defaultDays             = 15
)

// RegisterDebtReminderCommand registers the /debt-reminder slash command and its handler.
func RegisterDebtReminderCommand(ch *Handler, uc port.DebtReminder, scheduler gocron.Scheduler) {
	adminPerm := int64(discordgo.PermissionAdministrator)

	cmd := &discordgo.ApplicationCommand{
		Name:                     debtReminderCommandName,
		Description:              "立即執行欠費提醒，並排程於指定天數後再次執行",
		DefaultMemberPermissions: &adminPerm,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "days",
				Description: "幾天後再次執行（預設 15 天）",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "debug",
				Description: "除錯模式（僅傳送至 log 頻道，不發送 DM）",
				Required:    false,
			},
		},
	}

	ch.RegisterCommand(cmd, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handleDebtReminder(s, i, uc, scheduler)
	})
}

func handleDebtReminder(
	s *discordgo.Session, i *discordgo.InteractionCreate,
	uc port.DebtReminder, scheduler gocron.Scheduler,
) {
	respondDeferred(s, i)

	opts := i.ApplicationCommandData().Options
	optMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(opts))
	for _, opt := range opts {
		optMap[opt.Name] = opt
	}

	days := int64(defaultDays)
	if v, ok := optMap["days"]; ok {
		days = v.IntValue()
	}

	debug := false
	if v, ok := optMap["debug"]; ok {
		debug = v.BoolValue()
	}

	// Immediate run
	err := uc.Execute(context.Background(), debug)
	if err != nil {
		log.Printf("debt-reminder immediate run failed: %s", err)
	}

	// Schedule delayed production run
	runAt := time.Now().Add(time.Duration(days) * 24 * time.Hour)

	_, err = scheduler.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(runAt)),
		gocron.NewTask(func() {
			log.Print("debt-reminder scheduled run")
			if err := uc.Execute(context.Background(), false); err != nil {
				log.Printf("debt-reminder scheduled run failed: %s", err)
			}
		}),
	)
	if err != nil {
		log.Printf("debt-reminder schedule failed: %s", err)
		editDeferredResponse(s, i, fmt.Sprintf(
			"提醒已執行（模式: %s），但排程失敗: %s",
			modeLabel(debug), err,
		))

		return
	}

	editDeferredResponse(s, i, fmt.Sprintf(
		"提醒已執行（模式: %s）。下次執行: %s（正式模式）",
		modeLabel(debug),
		runAt.Format("2006-01-02 15:04"),
	))
}

func modeLabel(debug bool) string {
	if debug {
		return "除錯"
	}

	return "正式"
}
```

- [ ] **Step 2: Verify compilation**

Run: `go build ./...`
Expected: compiles without error

- [ ] **Step 3: Commit**

```bash
git add gateway/discord/command/debt_reminder_handler.go
git commit -m "feat: add /debt-reminder slash command handler"
```

---

### Task 6: Remove `WORKER_CORNTAB` from config

**Files:**
- Modify: `config/config.go`
- Modify: `config/config_test.go`

- [ ] **Step 1: Remove WorkerCrontab from Config struct and Load**

In `config/config.go`:
- Remove `WorkerCrontab string` from the `Config` struct
- Remove `WorkerCrontab: os.Getenv("WORKER_CORNTAB"),` from `Load()`
- Remove the validation block:
  ```go
  if cfg.WorkerCrontab == "" {
      return Config{}, fmt.Errorf("WORKER_CORNTAB is required")
  }
  ```

- [ ] **Step 2: Update config_test.go**

Remove any test cases that reference `WORKER_CORNTAB`. Check the file first — if tests set this env var, remove those lines.

- [ ] **Step 3: Run tests**

Run: `go test ./config/...`
Expected: all PASS

- [ ] **Step 4: Commit**

```bash
git add config/config.go config/config_test.go
git commit -m "refactor: remove WORKER_CORNTAB from config"
```

---

### Task 7: Update `main.go` — remove cron job, wire new command

**Files:**
- Modify: `main.go`

- [ ] **Step 1: Rewrite main.go**

Key changes:
- Remove `clock` import
- Remove `WORKER_CORNTAB` / crontab / debug-clock block (lines 74-81)
- Remove the cron job registration (lines 89-99)
- Keep scheduler creation (needed for one-shot jobs)
- Update `NewNotifier` call — remove `cfg.Debug` parameter
- Update `NewNotifyUnpaid` call — remove `loc` parameter
- Register the new `/debt-reminder` command, passing the scheduler
- Keep scheduler `Start()` and `Shutdown()` (scheduler still manages delayed jobs)

Updated main.go:

```go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"

	"github.com/xgnid-tw/gx5/config"
	discordgw "github.com/xgnid-tw/gx5/gateway/discord"
	discordcmd "github.com/xgnid-tw/gx5/gateway/discord/command"
	notiongw "github.com/xgnid-tw/gx5/gateway/notion"
	"github.com/xgnid-tw/gx5/usecase"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("can not fetch env")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("invalid config: %s", err)
	}

	// Initialize external service clients
	dc, err := discordgo.New(cfg.DiscordToken)
	if err != nil {
		log.Fatalf("can not create discord session: %s", err)
	}

	dc.Identify.Intents = discordgo.IntentsAll

	notionClient := notionapi.NewClient(notionapi.Token(cfg.NotionToken))

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("invalid location: %s", err)
	}

	// Wire dependencies: gateway adapters -> use cases
	repo := notiongw.NewRepository(notionClient.Database, cfg.NotionUserDBID, cfg.NotionOthersDBID)
	notifier := discordgw.NewNotifier(dc, cfg.DiscordLogChannelID)
	notifyUnpaidUC := usecase.NewNotifyUnpaid(repo, notifier, cfg.NotionOthersDBID)

	orderRepo := notiongw.NewOrderRepository(notionClient.Page, cfg.NotionOrderDBID)
	threadCreator := discordgw.NewThreadCreator(dc)
	memberAdder := discordgw.NewMemberAdder(dc, cfg.DiscordGuildID)
	createOrderUC := usecase.NewCreateOrder(orderRepo, threadCreator, memberAdder, cfg.TagRoleMap)

	txRepo := notiongw.NewTransactionRepository(notionClient.Page)
	buyUC := usecase.NewRegisterBuyRecord(repo, txRepo, cfg.ExchangeRateJPYTWD)

	// Scheduler for one-shot delayed jobs
	s, err := gocron.NewScheduler(gocron.WithLocation(loc))
	if err != nil {
		log.Fatalf("can not create scheduler: %s", err)
	}

	// Register Discord application commands
	cmdHandler := discordcmd.NewHandler(dc, cfg.DiscordAppID)

	discordcmd.RegisterNewOrderCommand(cmdHandler, createOrderUC)
	discordcmd.RegisterBuyCommand(cmdHandler, buyUC)
	discordcmd.RegisterDebtReminderCommand(cmdHandler, notifyUnpaidUC, s)

	// Open Discord connection and start the scheduler
	err = dc.Open()
	if err != nil {
		log.Fatalf("error opening connection: %s", err)
	}

	err = cmdHandler.SyncCommands()
	if err != nil {
		_ = dc.Close()

		log.Fatalf("error syncing commands: %s", err)
	}

	defer cmdHandler.UnregisterAll()
	defer dc.Close()

	log.Print("Bot is now running. Press CTRL-C to exit.")

	s.Start()

	// Block until termination signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	_ = s.Shutdown()
}
```

- [ ] **Step 2: Verify compilation**

Run: `go build ./...`
Expected: compiles without error

- [ ] **Step 3: Run all tests**

Run: `go test ./...`
Expected: all PASS

- [ ] **Step 4: Commit**

```bash
git add main.go
git commit -m "feat: wire /debt-reminder command, remove cron-based scheduler"
```

---

### Task 8: Clean up — remove unused dependencies

**Files:**
- Modify: `go.mod` / `go.sum`

- [ ] **Step 1: Tidy modules**

Run: `go mod tidy`

This should remove `github.com/benbjohnson/clock` if no longer imported anywhere.

- [ ] **Step 2: Verify**

Run: `go build ./... && go test ./...`
Expected: all PASS

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: go mod tidy, remove unused clock dependency"
```
