package main

import (
	"github.com/gempir/go-twitch-irc/v4"
)

func main() {
	channelName := GetChannelName()
	databaseName := GetDatabaseName()
	userName := GetBotUsername()
	oauthToken := GetOAuthToken()

	db := ConnectToDatabase(databaseName)
	CreateInitialDatabaseData(db)

	twitchClient := twitch.NewClient(userName, oauthToken)

	client := NewChatClient(channelName, twitchClient)

	incomingMessagesChannel := make(chan struct {
		UserName    string
		Message     string
		IsModerator bool
	}, 100)
	commandExecutionMetadataChannel := make(chan CommandExecutionMetadata, 100)
	outgoingMessagesChannel := make(chan string, 100)

	go client.StartListening(incomingMessagesChannel)
	go client.StartSaying(outgoingMessagesChannel)
	go ParseIncomingMessagesToCommands(incomingMessagesChannel, commandExecutionMetadataChannel)
	go ExecuteCommands(db, commandExecutionMetadataChannel, outgoingMessagesChannel) //TODO: Preload commands into a syncmap and pre-assemble template variables

	client.JoinChannel()
}
