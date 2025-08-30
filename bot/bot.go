package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	twitch "github.com/amgoh/hoagiebot/bot/twitch"
	"github.com/bwmarrin/discordgo"
)

// User Environment Variables
var BotToken string

// Guild-Specific Variables
var commandPrefix string = "!" // default prefix is "!" -> !command

func Run() {
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal("Error message")
	}

//	BotSession := discord

	// TO-DO: receive twitch event notification and send ping in discord
	go twitch.SubscribeAndListen()
	discord.Identify.Intents = discordgo.IntentsAll
	
	discord.AddHandler(newMessage) // command handler
	discord.AddHandler(memberJoin) // welcome message handler
	discord.AddHandler(verifyMember) // verify channel

	discord.Open()
	defer discord.Close()

	fmt.Println("Bot active!")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func memberJoin(discord *discordgo.Session, user *discordgo.GuildMemberAdd) {
	guildPrev, err := discord.GuildPreview(user.GuildID)
	if err != nil {
		return
	}

	guild, err := discord.Guild(user.GuildID)
	if err != nil {
		return
	}

	embeds := []*discordgo.MessageEmbed{ 
		{
			Title: "Welcome to "+guild.Name, 
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: user.AvatarURL(""),
				ProxyURL: user.AvatarURL(""),
				Width: 300,
				Height: 300,
			},
			Description: 
			"> Read the [Rules](https://discord.com/channels/"+user.GuildID+"/"+guild.SystemChannelID+")"+
			"\n> Chat & Boost"+
			"\n> You are User #"+strconv.Itoa(guildPrev.ApproximateMemberCount),
		},
	}

	welcome_msg := discordgo.MessageSend {
		Content: user.Mention(),
		Embeds: embeds,
	}

//	discord.GuildMemberRoleAdd(user.GuildID, user.User.ID, )
	discord.ChannelMessageSendComplex(guild.SystemChannelID, &welcome_msg)	
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
 	if message.Author.ID == discord.State.User.ID {
  	return
 	}

	tokens := strings.Split(message.Content, " ")

	req_user := message.Author.Mention()



	switch (tokens[0]) {
	case commandPrefix+"help":
	  discord.ChannelMessageSend(message.ChannelID, "Hello World😃")
	
	case commandPrefix+"youtube":
		discord.ChannelMessageSend(message.ChannelID, req_user + " CHECK OUT THE CHANNEL AND SUBSCRIBE!\nhttps://youtube.com/") // ----- INCLUDE YOUTUBE USER
	case commandPrefix+"setPrefix":
		if len(tokens) != 2 {
			discord.ChannelMessageSend(message.ChannelID, req_user + ": Incorrect Usage. Try !setPrefix [symbol]")
			return
		}
		
		discord.ChannelMessageSend(message.ChannelID, req_user + ": Prefix set to " + tokens[1])
		commandPrefix = tokens[1]
	}
}

func verifyMember(discord *discordgo.Session, event *discordgo.MessageReactionAdd) {
	// TO-DO
	// Give member roles when reacting to the Rules Message in the selected Rules channel

}
