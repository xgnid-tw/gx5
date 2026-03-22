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
	"github.com/jomei/notionapi"

	"github.com/xgnid-tw/gx5/config"
	"github.com/xgnid-tw/gx5/domain"
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

	notionClient := notionapi.NewClient(notionapi.Token(cfg.NotionToken))

	// Load Asia/Tokyo timezone for scheduler and use case day guard
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("invalid location: %s", err)
	}

	// Wire dependencies: gateway adapters -> use cases
	repo := notiongw.NewRepository(cfg.NotionToken, cfg.NotionUserDBID, cfg.NotionOthersDBID)
	notifier := discordgw.NewNotifier(dc, cfg.DiscordLogChannelID, cfg.Debug)
	notifyUnpaidUC := usecase.NewNotifyUnpaid(repo, notifier, cfg.NotionOthersDBID, loc)

	orderRepo := notiongw.NewOrderRepository(notionClient.Page, cfg.NotionOrderDBID)
	threadCreator := discordgw.NewThreadCreator(dc)
	createOrderUC := usecase.NewCreateOrder(orderRepo, threadCreator, cfg.DiscordOwnerID)

	// In debug mode, fake the clock and run the job every minute
	crontab := cfg.WorkerCrontab

	if cfg.Debug {
		clk := clock.NewMock()
		clk.Set(time.Date(time.Now().Year(), time.Now().Month(), 1, 9, 0, 0, 0, loc))
		notifyUnpaidUC.Clock = clk
		crontab = "*/1 * * * *"
	}

	s, err := gocron.NewScheduler(gocron.WithLocation(loc))
	if err != nil {
		log.Fatalf("can not create scheduler: %s", err)
	}

	// Register the unpaid notification job on the configured cron schedule
	_, err = s.NewJob(gocron.CronJob(crontab, false), gocron.NewTask(func() {
		log.Print("run job")

		err := notifyUnpaidUC.Execute(ctx)
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

	// Register slash commands
	cmdHandler := discordgw.NewCommandHandler(dc, cfg.DiscordAppID)
	defer cmdHandler.UnregisterAll()

	err = cmdHandler.RegisterCommand(newOrderCommand(), func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handleNewOrder(s, i, createOrderUC)
	})
	if err != nil {
		log.Fatalf("error registering newOrder command: %s", err)
	}

	log.Print("Bot is now running. Press CTRL-C to exit.")

	s.Start()

	// Block until termination signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	_ = s.Shutdown()
}

func newOrderCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "neworder",
		Description: "Create a new group purchase order",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "ordertitle",
				Description: "Name of the order",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "deadline",
				Description: "Order deadline (YYYY-MM-DD)",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "shopurl",
				Description: "Shop URL",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "tags",
				Description: "Tag (315pro, 学マス, 283pro, 346pro, 765pro)",
				Required:    false,
			},
		},
	}
}

func handleNewOrder(s *discordgo.Session, i *discordgo.InteractionCreate, uc *usecase.CreateOrder) {
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

	callerID := i.Member.User.ID

	err := uc.Execute(context.Background(), callerID, i.ChannelID, order)
	if err != nil {
		respondToInteraction(s, i, "Error: "+err.Error())
		return
	}

	respondToInteraction(s, i, "Order created: "+order.ThreadName)
}

func respondToInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	if err != nil {
		log.Printf("error responding to interaction: %s", err)
	}
}
