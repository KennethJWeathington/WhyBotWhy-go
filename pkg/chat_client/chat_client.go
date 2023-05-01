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

func (client *ChatClient) StartListening(incomingMessagesChannel chan<- struct {
	UserName    string
	Message     string
	IsModerator bool
}) {
	parseIncomingMessageCallback := client.parseIncomingMessage(incomingMessagesChannel)
	client.twitchClient.OnPrivateMessage(parseIncomingMessageCallback)
}

func (client *ChatClient) parseIncomingMessage(incomingMessagesChannel chan<- struct {
	UserName    string
	Message     string
	IsModerator bool
}) func(message twitch.PrivateMessage) {
	return func(message twitch.PrivateMessage) {
		userName := message.User.DisplayName
		messageText := message.Message
		isAdmin := isModerator(message.User.Badges)

		incomingMessagesChannel <- struct {
			UserName    string
			Message     string
			IsModerator bool
		}{userName, messageText, isAdmin}
	}
}

func (client *ChatClient) StartSaying(outgoingMessagesChannel <-chan string) {
	for message := range outgoingMessagesChannel {
		client.twitchClient.Say(client.channelName, message)
	}
}

func NewChatClient(channel string, username string, oauth string) *ChatClient {
	return &ChatClient{
		twitchClient: twitch.NewClient(username, oauth),
		channelName:  channel,
	}
}

func isModerator(badges map[string]int) bool {
	return badges["broadcaster"] > 0 || badges["moderator"] > 0
}
