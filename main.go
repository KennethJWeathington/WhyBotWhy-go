package main

func main() {
	channelName := GetChannelName()
	databaseName := GetDatabaseName()
	userName := GetBotUsername()
	oauthToken := GetOAuthToken()

	InitDatabase(databaseName)
	db := GetConnection()
	CreateInitialDatabaseData(db)

	client := NewTwitchChatClient(userName, oauthToken, channelName)

	chatCommandChannel := make(chan ChatCommand)
	outgoingMessagesChannel := make(chan string)

	go client.StartListening(chatCommandChannel)
	go client.StartSaying(outgoingMessagesChannel)
	go ExecuteCommands(db, chatCommandChannel, outgoingMessagesChannel) //TODO: Preload commands into a syncmap and pre-assemble template variables

	client.JoinChannel()
}
