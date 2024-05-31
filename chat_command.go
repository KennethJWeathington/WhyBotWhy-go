package main

type ChatCommand struct {
	UserName    string
	IsModerator bool
	CommandName string
	Arguments   []string
}
