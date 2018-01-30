package example

import (
	"testing"

	"github.com/go-chat-bot/bot"
)

func TestAuth(t *testing.T) {
	bot := &bot.Cmd{
		Command: "auth",
		User: &bot.User{
			Nick:     "nick",
			RealName: "Real Name",
		},
	}
	want := "Whats up? Real Name"
	got, error := hello(bot)

	if got != want {
		t.Errorf("Expected '%v' got '%v'", want, got)
	}

	if error != nil {
		t.Errorf("Expected '%v' got '%v'", nil, error)
	}
}
