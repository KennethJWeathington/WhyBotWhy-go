package main

import (
	"github.com/jake-weath/whybotwhy_go/pkg/chat_client"
	"github.com/jake-weath/whybotwhy_go/pkg/command"
	"github.com/jake-weath/whybotwhy_go/pkg/database_client"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/joho/godotenv"
)

const (
	ChannelNameEnvVar  = "CHANNEL_NAME"
	BotUsernameEnvVar  = "BOT_USERNAME"
	OAuthTokenEnvVar   = "OAUTH_TOKEN"
	DatabaseNameEnvVar = "DATABASE_NAME"
)

func main() {
	env, err := godotenv.Read()
	if err != nil {
		panic(err)
	}

	channelName, userName, oauthToken, databaseName := env[ChannelNameEnvVar], env[BotUsernameEnvVar], env[OAuthTokenEnvVar], env[DatabaseNameEnvVar]

	db := database_client.ConnectToDatabase(databaseName)
	database_client.CreateInitialDatabaseData(db)

	twitchClient := twitch.NewClient(userName, oauthToken)

	client := chat_client.NewChatClient(channelName, twitchClient)

	incomingMessagesChannel := make(chan struct {
		UserName    string
		Message     string
		IsModerator bool
	}, 100)
	commandExecutionMetadataChannel := make(chan command.CommandExecutionMetadata, 100)
	outgoingMessagesChannel := make(chan string, 100)

	go client.StartListening(incomingMessagesChannel)
	go client.StartSaying(outgoingMessagesChannel)
	go command.ParseIncomingMessagesToCommands(incomingMessagesChannel, commandExecutionMetadataChannel)
	go command.ExecuteCommands(db, commandExecutionMetadataChannel, outgoingMessagesChannel) //TODO: Preload commands into a syncmap and pre-assemble template variables

	client.JoinChannel()
}
