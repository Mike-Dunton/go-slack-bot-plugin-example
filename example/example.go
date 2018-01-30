package example

import (
	"encoding/json"
	"fmt"

	"github.com/agext/uuid"
	"github.com/go-chat-bot/bot"
	"github.com/go-redis/redis"
)

type PocUser struct {
	Groups  []string
	Id      string
	IsAuthd bool
}

func getPocUser(redisClient *redis.Client, Id string) *PocUser {
	user, err := redisClient.Get(Id).Result()
	if err != nil {
		return &PocUser{
			Groups:  []string{},
			Id:      Id,
			IsAuthd: false,
		}
	}
	var deserializedUser PocUser
	err = json.Unmarshal([]byte(user), &deserializedUser)
	return &deserializedUser
}

func deauth(command *bot.Cmd) (msg string, err error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:32768",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	user := getPocUser(redisClient, command.User.ID)
	if user.IsAuthd {
		redisClient.Del(command.User.ID)
		msg = fmt.Sprintf("You are no longer authorized")
		return
	}
	msg = fmt.Sprintf("You are not currently authorized. Please Authorize yourself in order to use this command")
	return
}

func auth(command *bot.Cmd) (msg string, err error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:32768",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	user := getPocUser(redisClient, command.User.ID)
	if user.IsAuthd {
		msg = fmt.Sprintf("You are already authd, Your Id is %s, and your groups are %s", user.Id, user.Groups)
		return
	}
	key := uuid.New().String()
	redisClient.Set(key, command.User.ID, 0)
	msg = fmt.Sprintf("Visit http://localhost:5050/link?key=%s", key)
	return
}

func init() {
	bot.RegisterCommand(
		"auth",
		"Authorizes slack user with okta backend",
		"",
		auth)
	bot.RegisterCommand(
		"deauth",
		"Deauthorizes slack user",
		"",
		deauth)
}
