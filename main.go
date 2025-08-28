package main

import (
	"log"
	"os"

	twitch "github.com/amgoh/hoagiebot/bot/twitch"
	bot "github.com/amgoh/hoagiebot/bot"
	dotenv "github.com/joho/godotenv"
)

func main() {
	dotenv.Load()
	bot.BotToken = os.Getenv("BOT_TOKEN")
	twitch.WebhookSecret = os.Getenv("TWITCH_WEBHOOK_SECRET")
	twitch.TwitchClientID = os.Getenv("TWITCH_CLIENT_ID")
	twitch.TwitchClientSecret = os.Getenv("TWITCH_CLIENT_SECRET")

	if(len(bot.BotToken) == 0) {
		log.Fatal("env variable BOT_TOKEN not set")
	}
	
	bot.Run()
}
