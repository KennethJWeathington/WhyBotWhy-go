package main

import (
	"github.com/gempir/go-twitch-irc/v4"
)

type ChatClient struct {
	chatConnection ChatConnection
	channelName    string
}

type ChatConnection interface {
	Join(channels ...string)
	Connect() error
	Say(channel, message string)
	OnPrivateMessage(handler func(message twitch.PrivateMessage))
}

func (client *ChatClient) JoinChannel() {
	client.chatConnection.Join(client.channelName)
	if err := client.chatConnection.Connect(); err != nil {
		panic(err)
	}
}

func (client *ChatClient) StartListening(incomingMessagesChannel chan<- struct {
	UserName    string
	Message     string
	IsModerator bool
}) {
	parseIncomingMessageCallback := client.parseIncomingMessage(incomingMessagesChannel)
	client.chatConnection.OnPrivateMessage(parseIncomingMessageCallback)
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
		client.chatConnection.Say(client.channelName, message)
	}
}

func NewChatClient(channel string, chatConnection ChatConnection) *ChatClient {
	return &ChatClient{
		chatConnection: chatConnection,
		channelName:    channel,
	}
}

func isModerator(badges map[string]int) bool {
	return badges["broadcaster"] > 0 || badges["moderator"] > 0
}
