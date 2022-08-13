package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Dinmukhamet/gostandup/commands"
	"github.com/Dinmukhamet/gostandup/constants"
	"github.com/Dinmukhamet/gostandup/tasks"
)

type BotCommand func(tg.BotAPI, tg.Message) (tg.Message, error)

func startScheduler(bot tg.BotAPI) {
	s := gocron.NewScheduler(constants.LOCATION)
	deadline, err := time.Parse(constants.TIME_FORMAT, constants.DEADLINE_TIME)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return
	}
	before := deadline.Add(time.Duration(-15) * time.Minute)

	if _, err := s.Every(1).Day().At(before.Format(constants.TIME_CRON_FORMAT)).Do(func() {
		tasks.DailyReminderTask(bot)
	}); err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}

	if _, err := s.Every(1).Day().At(deadline.Format(constants.TIME_CRON_FORMAT)).Do(func() {
		tasks.ShowMissedTask(bot)
	}); err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}

	s.StartAsync()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	os.Setenv("TZ", constants.DEFAULT_TIMEZONE)
	l, _ := time.LoadLocation(constants.DEFAULT_TIMEZONE)
	time.Local = l
	constants.LOCATION = l

	bot, err := tg.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Error statring the bot: %s", err.Error())
	}

	startScheduler(*bot)

	mongoURI := fmt.Sprintf("mongodb://%s:27017", os.Getenv("MONGO_HOST"))
	err = mgm.SetDefaultConfig(nil, "standups", options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error connecting to mongo: %s", err.Error())
	}

	bot.Debug = true

	u := tg.NewUpdate(0)
	u.Timeout = 60

	registeredCommands := map[string]BotCommand{
		"register_me":      commands.RegisterUserCommand,
		"start":            commands.StartCommand,
		"pause_activity":   commands.PauseActivityCommand,
		"unpause_activity": commands.UnpauseActivityCommand,
		"test":             commands.TestCommand,
	}

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			if strings.HasPrefix(update.Message.Text, "@mega_standup_bot standup") {
				commands.ReceiveStandupCommand(*bot, *update.Message)
			} else if update.Message.IsCommand() {
				command := registeredCommands[update.Message.Command()]
				if _, err := command(*bot, *update.Message); err != nil {
					log.Printf("ERROR: %s", err.Error())
				}
			}
		}
	}
}
