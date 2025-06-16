package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"

	"github.com/xgnid-tw/gx5/discord"
	"github.com/xgnid-tw/gx5/model"
	"github.com/xgnid-tw/gx5/notion"
)

const ChanBuffer = 20

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("can not fetch env")
	}

	ctx := context.Background()
	notionKey := os.Getenv("NOTION_TOKEN")
	notionUserDBID := os.Getenv("NOTION_USER_DB_ID")

	discordToken := os.Getenv("DISCORD_TOKEN")
	discordLogChannelID := os.Getenv("DISCORD_GUILD_LOG_CHANNEL_ID")

	dc, err := discordgo.New(discordToken)
	if err != nil {
		log.Fatalf("can not create discord session, %s", err)
	}
	defer dc.Close()

	nToDch := make(chan model.User, ChanBuffer)

	// run worker
	runWorker(ctx, notionKey, nToDch, notionUserDBID)

	// run discord bot
	runDiscordBot(ctx, dc, nToDch, discordLogChannelID)
}

func runWorker(ctx context.Context, nsKey string, nToDch chan model.User, userDBID string) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("invalid location :%s", err)
	}

	s, err := gocron.NewScheduler(
		gocron.WithLocation(loc),
	)
	if err != nil {
		log.Fatalf("can not create scheduler, %s", err)
	}

	ns, err := notion.NewNotion(nsKey, nToDch, userDBID)
	if err != nil {
		log.Fatalf("can not create notion service, %s", err)
	}

	_, err = s.NewJob(gocron.DailyJob(
		1, gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0)),
	), gocron.NewTask(
		func() {
			log.Print("run job ")

			err = ns.SendNotPaidInformation(ctx)
			if err != nil {
				log.Fatalf("worker: %s", err)
			}
		},
	))
	if err != nil {
		log.Fatalf("can not start scheduler, %s", err)
	}

	s.Start()
}

func runDiscordBot(ctx context.Context,
	dc *discordgo.Session, nToDch chan model.User, logChannelID string,
) {
	des := discord.NewDiscordEventService(dc, nToDch, logChannelID)

	dc.Identify.Intents = discordgo.IntentsAll

	err := dc.Open()
	if err != nil {
		log.Fatalf("error opening connection: %s", err)
	}

	log.Print("Bot is now running. Press CTRL-C to exit.")

	go des.GetChanMsgAndDM(ctx)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
