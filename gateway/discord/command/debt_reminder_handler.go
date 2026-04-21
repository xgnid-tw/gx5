package command

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/xgnid-tw/gx5/port"
)

const (
	debtReminderCommandName = "debt-reminder"
	debtReminderOptionDays  = "days"
	debtReminderOptionDebug = "debug"
	defaultDays             = 15
	minDays                 = 1
	timestampLayout         = "2006-01-02 15:04"
)

// debtReminderScheduler is the subset of *DebtReminderScheduler used by the handler.
// Declared as an interface so tests can substitute a fake.
type debtReminderScheduler interface {
	ScheduleProductionRun(runAt time.Time, task func()) (ScheduleResult, error)
}

// RegisterDebtReminderCommand registers the /debt-reminder slash command and its handler.
func RegisterDebtReminderCommand(
	ch *Handler, uc port.DebtReminder, sched debtReminderScheduler,
) {
	adminPerm := int64(discordgo.PermissionAdministrator)

	cmd := &discordgo.ApplicationCommand{
		Name:                     debtReminderCommandName,
		Description:              "立即執行欠費提醒，並排程於指定天數後再次執行",
		DefaultMemberPermissions: &adminPerm,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        debtReminderOptionDays,
				Description: "幾天後再次執行（預設 15 天）",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        debtReminderOptionDebug,
				Description: "除錯模式（僅傳送至 log 頻道，不發送 DM，不排程下次執行）",
				Required:    false,
			},
		},
	}

	ch.RegisterCommand(cmd, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handleDebtReminder(s, i, uc, sched)
	})
}

func handleDebtReminder(
	s *discordgo.Session, i *discordgo.InteractionCreate,
	uc port.DebtReminder, sched debtReminderScheduler,
) {
	respondDeferred(s, i)

	opts := i.ApplicationCommandData().Options

	optMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(opts))
	for _, opt := range opts {
		optMap[opt.Name] = opt
	}

	days := int64(defaultDays)
	if v, ok := optMap[debtReminderOptionDays]; ok {
		days = v.IntValue()
	}

	if days < minDays {
		editDeferredResponse(s, i, fmt.Sprintf("天數必須至少為 %d", minDays))
		return
	}

	debug := false
	if v, ok := optMap[debtReminderOptionDebug]; ok {
		debug = v.BoolValue()
	}

	if debug {
		handleDebtReminderDebug(s, i, uc)
		return
	}

	handleDebtReminderProduction(s, i, uc, sched, days)
}

func handleDebtReminderDebug(
	s *discordgo.Session, i *discordgo.InteractionCreate,
	uc port.DebtReminder,
) {
	err := uc.Execute(context.Background(), true)
	if err != nil {
		log.Printf("debt-reminder immediate run failed: %s", err)
		editDeferredResponse(s, i, fmt.Sprintf("提醒執行失敗: %s", err))

		return
	}

	editDeferredResponse(s, i, "提醒已執行（模式: 除錯）。未排程下次執行。")
}

func handleDebtReminderProduction(
	s *discordgo.Session, i *discordgo.InteractionCreate,
	uc port.DebtReminder, sched debtReminderScheduler, days int64,
) {
	err := uc.Execute(context.Background(), false)
	if err != nil {
		log.Printf("debt-reminder immediate run failed: %s", err)
		editDeferredResponse(s, i, fmt.Sprintf("提醒執行失敗，未排程下次執行: %s", err))

		return
	}

	runAt := time.Now().Add(time.Duration(days) * 24 * time.Hour)

	task := func() {
		log.Print("debt-reminder scheduled run")

		execErr := uc.Execute(context.Background(), false)
		if execErr != nil {
			log.Printf("debt-reminder scheduled run failed: %s", execErr)
		}
	}

	result, err := sched.ScheduleProductionRun(runAt, task)
	if err != nil {
		log.Printf("debt-reminder schedule failed: %s", err)
		editDeferredResponse(s, i, fmt.Sprintf(
			"提醒已執行（模式: 正式），但排程失敗: %s", err,
		))

		return
	}

	editDeferredResponse(s, i, buildProductionResponse(runAt, result))
}

func buildProductionResponse(runAt time.Time, result ScheduleResult) string {
	msg := fmt.Sprintf(
		"提醒已執行（模式: 正式）。下次執行: %s",
		runAt.Format(timestampLayout),
	)

	if !result.ReplacedAt.IsZero() {
		msg += fmt.Sprintf("（已取代先前排程: %s）", result.ReplacedAt.Format(timestampLayout))
	}

	if result.RemoveWarn != nil {
		msg += fmt.Sprintf("（注意：先前排程移除失敗: %s）", result.RemoveWarn)
	}

	return msg
}
