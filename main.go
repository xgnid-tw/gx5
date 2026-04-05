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
	// Load environment variables from .env file
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

	// Load Asia/Tokyo timezone for scheduler
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

	// Register Discord application commands
	cmdHandler := discordcmd.NewHandler(dc, cfg.DiscordAppID)

	// Scheduler for one-shot delayed jobs
	s, err := gocron.NewScheduler(gocron.WithLocation(loc))
	if err != nil {
		log.Fatalf("can not create scheduler: %s", err)
	}

	discordcmd.RegisterNewOrderCommand(cmdHandler, createOrderUC)
	discordcmd.RegisterBuyCommand(cmdHandler, buyUC)
	discordcmd.RegisterDebtReminderCommand(cmdHandler, notifyUnpaidUC, s)

	// Open Discord connection and start the scheduler
	err = dc.Open()
	if err != nil {
		log.Fatalf("error opening connection: %s", err)
	}

	// Sync application commands after connection is open
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
