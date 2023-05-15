package env_reader

import (
	"github.com/joho/godotenv"
)

const (
	ChannelNameEnvVar  = "CHANNEL_NAME"
	BotUsernameEnvVar  = "BOT_USERNAME"
	OAuthTokenEnvVar   = "OAUTH_TOKEN"
	DatabaseNameEnvVar = "DATABASE_NAME"
)

var env map[string]string
var err error

func init() {
	env, err = godotenv.Read()
	if err != nil {
		panic(err)
	}
}

func GetChannelName() string {
	return env[ChannelNameEnvVar]
}

func GetBotUsername() string {
	return env[BotUsernameEnvVar]
}

func GetOAuthToken() string {
	return env[OAuthTokenEnvVar]
}

func GetDatabaseName() string {
	return env[DatabaseNameEnvVar]
}
