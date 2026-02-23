package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"

	"github.com/xgnid-tw/gx5/config"
	discordgw "github.com/xgnid-tw/gx5/gateway/discord"
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

	ctx := context.Background()

	// Initialize external service clients
	dc, err := discordgo.New(cfg.DiscordToken)
	if err != nil {
		log.Fatalf("can not create discord session: %s", err)
	}

	dc.Identify.Intents = discordgo.IntentsAll

	// Load Asia/Tokyo timezone for scheduler and use case day guard
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("invalid location: %s", err)
	}

	// Wire dependencies: gateway adapters -> use case
	repo := notiongw.NewRepository(cfg.NotionToken, cfg.NotionUserDBID, cfg.NotionOthersDBID)
	notifier := discordgw.NewNotifier(dc, cfg.DiscordLogChannelID, cfg.Debug)
	uc := usecase.NewNotifyUnpaid(repo, notifier, cfg.NotionOthersDBID, loc)

	// In debug mode, fake the clock and run the job every minute
	crontab := cfg.WorkerCrontab

	if cfg.Debug {
		clk := clock.NewMock()
		clk.Set(time.Date(time.Now().Year(), time.Now().Month(), 1, 9, 0, 0, 0, loc))
		uc.Clock = clk
		crontab = "*/1 * * * *"
	}

	s, err := gocron.NewScheduler(gocron.WithLocation(loc))
	if err != nil {
		log.Fatalf("can not create scheduler: %s", err)
	}

	// Register the unpaid notification job on the configured cron schedule
	_, err = s.NewJob(gocron.CronJob(crontab, false), gocron.NewTask(func() {
		log.Print("run job")

		err := uc.Execute(ctx)
		if err != nil {
			log.Printf("worker: %s", err)
		}
	}))
	if err != nil {
		log.Fatalf("can not start scheduler: %s", err)
	}

	// Open Discord connection and start the scheduler
	err = dc.Open()
	if err != nil {
		log.Fatalf("error opening connection: %s", err)
	}
	defer dc.Close()

	log.Print("Bot is now running. Press CTRL-C to exit.")

	s.Start()

	// Block until termination signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	_ = s.Shutdown()
}
