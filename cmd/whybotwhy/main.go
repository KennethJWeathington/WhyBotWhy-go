package main

import (
	"github.com/jake-weath/whybotwhy_go/pkg/chat_client"
	"github.com/jake-weath/whybotwhy_go/pkg/command"
	"github.com/jake-weath/whybotwhy_go/pkg/database_client"

	"github.com/joho/godotenv"
)

func main() {
	env, err := godotenv.Read()
	if err != nil {
		panic(err)
	}

	channelName, userName, oauthToken, databaseName := env["CHANNEL_NAME"], env["BOT_USERNAME"], env["OAUTH_TOKEN"], env["DATABASE_NAME"]

	db := database_client.SetUpDatabase(databaseName)
	database_client.CreateInitialDatabaseData(db)
	commands := database_client.GetAllCommands(db)

	client := chat_client.NewChatClient(channelName, userName, oauthToken)

	incomingMessagesChannel := make(chan struct {
		UserName    string
		Message     string
		IsModerator bool
	}, 100)
	commandsChannel := make(chan command.CommandExecutionDetails, 100)
	outgoingMessagesChannel := make(chan string, 100)

	go client.StartListening(incomingMessagesChannel)
	go client.StartSaying(outgoingMessagesChannel)
	go command.ParseIncomingMessagesToCommands(incomingMessagesChannel, commandsChannel)
	go command.ExecuteCommands(commands, commandsChannel, outgoingMessagesChannel)

	//go testChatHandler(incomingMessagesChannel, outgoingMessagesChannel)

	client.JoinChannel()
}

// func testChatHandler(inputChannel <-chan struct {
// 	UserName    string
// 	Message     string
// 	IsModerator bool
// },
// 	outputChannel chan<- string) {
// 	for messageDetails := range inputChannel {
// 		handleChatMessage(messageDetails.Message, outputChannel)
// 	}
// }

// func handleChatMessage(message string, outputChannel chan<- string) {
// 	testTemplate(outputChannel)

// 	if command, err := command.ParseCommand(message); err == nil {
// 		commandName := "The command name was: " + command.Name
// 		commandArgs := "The command arguments were: " + strings.Join(command.Arguments, " ")

// 		outputChannel <- commandName
// 		outputChannel <- commandArgs
// 	}
// }

// func testTemplate(outputChannel chan<- string) {
// 	testMap := map[string]string{"name": "Jake"}
// 	builder := &strings.Builder{}

// 	tmpl, err := template.New("test").Parse("{{.name}} is the test name")
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = tmpl.Execute(builder, testMap)
// 	if err != nil {
// 		panic(err)
// 	}

// 	templateMessage := builder.String()

// 	outputChannel <- templateMessage
// }
