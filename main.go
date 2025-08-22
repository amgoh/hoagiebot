package main

import (
	"log"
	"os"

	bot "github.com/amgoh/hoagiebot/bot"
	dotenv "github.com/joho/godotenv"
)

func main() {
	dotenv.Load()
	bot.BotToken = os.Getenv("BOT_TOKEN")
	if(len(bot.BotToken) == 0) {
		log.Fatal("env variable BOT_TOKEN not set")
	}
	
	bot.Run()
}
