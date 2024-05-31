package main

func main() {
	channelName := GetChannelName()
	databaseName := GetDatabaseName()
	userName := GetBotUsername()
	oauthToken := GetOAuthToken()

	db := ConnectToDatabase(databaseName)
	CreateInitialDatabaseData(db)

	client := NewTwitchChatClient(userName, oauthToken, channelName)

	chatCommandChannel := make(chan ChatCommand, 100)
	outgoingMessagesChannel := make(chan string, 100)

	go client.StartListening(chatCommandChannel)
	go client.StartSaying(outgoingMessagesChannel)
	go ExecuteCommands(db, chatCommandChannel, outgoingMessagesChannel) //TODO: Preload commands into a syncmap and pre-assemble template variables

	client.JoinChannel()
}
