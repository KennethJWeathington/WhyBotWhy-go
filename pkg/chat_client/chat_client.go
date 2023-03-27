package chat_client

import (
	"github.com/gempir/go-twitch-irc/v4"
)

type ChatClient struct {
	twitchClient *twitch.Client
	channelName  string
}

func (client *ChatClient) JoinChannel() {
	client.twitchClient.Join(client.channelName)
	if err := client.twitchClient.Connect(); err != nil {
		panic(err)
	}
}

func (client *ChatClient) StartListening(inputChannel chan<- struct {
	Message     string
	IsModerator bool
}) {
	parseIncomingMessage := client.parseIncomingMessageCallback(inputChannel)
	client.twitchClient.OnPrivateMessage(parseIncomingMessage)
}

func (client *ChatClient) StartChatting(outputChannel <-chan string) {
	go client.sayIncomingMessages(outputChannel)
}

func (client *ChatClient) parseIncomingMessageCallback(inputChannel chan<- struct {
	Message     string
	IsModerator bool
}) func(message twitch.PrivateMessage) {
	return func(message twitch.PrivateMessage) {
		messageText := message.Message
		isAdmin := false

		inputChannel <- struct {
			Message     string
			IsModerator bool
		}{messageText, isAdmin}
	}
}

func (client *ChatClient) sayIncomingMessages(channel <-chan string) {
	for message := range channel {
		client.twitchClient.Say(client.channelName, message)
	}
}

func NewChatClient(channel string, username string, oauth string) *ChatClient {
	return &ChatClient{
		twitchClient: twitch.NewClient(username, oauth),
		channelName:  channel,
	}
}
