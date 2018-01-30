package example

import (
	"fmt"

	"github.com/go-chat-bot/bot"
)

func auth(command *bot.Cmd) (msg string, err error) {
	msg = fmt.Sprintf("Whats up? %s", command.User.RealName)
	return
}

func init() {
	bot.RegisterCommand(
		"auth",
		"Authorizes slack user with okta backend",
		"",
		auth)
}
