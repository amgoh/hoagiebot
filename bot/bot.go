package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var BotToken string

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}

func Run() {
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	discord.AddHandler(newMessage) // command handler
	discord.AddHandler(userJoin) // welcome message handler
	discord.AddHandler(verifyMember) // verify channel

	discord.Open()
	defer discord.Close()

	fmt.Println("Bot active!")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func userJoin(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID {
		return
	}
	if message.Type != discordgo.MessageTypeGuildMemberJoin {
		return
	}

	guildPrev, err := discord.GuildPreview(message.GuildID)
	if err != nil {
		return
	}

	guild, err := discord.Guild(message.GuildID)
	if err != nil {
		return
	}


	embeds := []*discordgo.MessageEmbed{ 
		{
			Title: "Welcome to HOAGIE", 
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: message.Author.AvatarURL(""),
				ProxyURL: message.Author.AvatarURL(""),
				Width: 300,
				Height: 300,
			},
			Description: "> Read the [Rules](https://discord.com/channels/"+message.GuildID+"/"+message.ChannelID+")\n> Chat & Boost\n> You are User #"+strconv.Itoa(guildPrev.ApproximateMemberCount),
		},
	}

	welcome_msg := discordgo.MessageSend {
		Content: message.Author.Mention(),
		Embeds: embeds,
	}
	
	discord.ChannelMessageDelete(message.ChannelID, message.Message.ID)
	discord.ChannelMessageSendComplex(guild.SystemChannelID, &welcome_msg)
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
 	if message.Author.ID == discord.State.User.ID {
  	return
 	}

	tokens := strings.Split(message.Content, " ")

	req_user := message.Author.Mention()

	if tokens[0] == "!help" {
	  discord.ChannelMessageSend(message.ChannelID, "Hello World😃")
	}
	if tokens[0] == "!youtube" {
		discord.ChannelMessageSend(message.ChannelID, req_user + " CHECK OUT THE CHANNEL AND SUBSCRIBE!\nhttps://youtube.com/@amoghiehoagie")
	}
}

func verifyMember(discord *discordgo.Session, event *discordgo.MessageReactionAdd) {
	
}
