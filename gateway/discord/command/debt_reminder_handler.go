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

			err := uc.Execute(context.Background(), false)
			if err != nil {
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
