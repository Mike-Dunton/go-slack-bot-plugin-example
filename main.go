package main

import (
	"encoding/json"
	"fmt"
	"os"

	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/maximRnback/gin-oidc"

	"github.com/go-chat-bot/bot/slack"
	_ "github.com/mike-dunton/go-slack-bot-plugin-example/example"
)

type PocUser struct {
	Groups  []string
	Id      string
	IsAuthd bool
}

func main() {
	go slack.Run(os.Getenv("dunton_dev_SLACK_TOKEN"))

	r := gin.Default()
	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	issuerUrl, _ := url.Parse(os.Getenv("dunton_dev_OKTA_URL") + "/oauth2/default")
	clientUrl, _ := url.Parse("http://localhost:5050/")
	logoutUrl, _ := url.Parse(os.Getenv("dunton_dev_OKTA_URL"))
	initParams := gin_oidc.InitParams{
		Router:       r,
		ClientId:     os.Getenv("dunton_dev_OKTA_CLIENT_ID"),
		ClientSecret: os.Getenv("dunton_dev_OKTA_CLIENT_SECRET"),
		Issuer:       *issuerUrl,
		ClientUrl:    *clientUrl,
		Scopes:       []string{"openid", "profile"},
		ErrorHandler: func(c *gin.Context) {
			c.Redirect(503, "http://localhost:5050/error")
		},
		PostLogoutUrl: *logoutUrl,
	}
	r.Use(gin_oidc.Init(initParams))
	r.GET("/link", func(c *gin.Context) {
		keyString := c.Query("key")
		redisClient := redis.NewClient(&redis.Options{
			Addr:     "localhost:32768",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		serverSession := sessions.Default(c)
		userId, err := redisClient.Get(keyString).Result()
		var claims map[string]interface{}
		sessionClaims := serverSession.Get("oidcClaims")
		json.Unmarshal([]byte(sessionClaims.(string)), &claims)
		user := &PocUser{
			IsAuthd: true,
			Id:      userId,
			Groups:  convertGroups(claims["groups"]),
		}
		serialized, err := json.Marshal(user)
		err = redisClient.Set(userId, serialized, 0).Err()
		if err != nil {
			c.JSON(503, gin.H{
				"message": err,
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "Thank you for linking your slack user",
		})

	})
	r.Run(":5050") // listen and serve on 0.0.0.0:8080
}

func convertGroups(groups interface{}) []string {
	groupsSlice := groups.([]interface{})
	s := make([]string, len(groupsSlice))
	for i, v := range groupsSlice {
		s[i] = fmt.Sprint(v)
	}
	return s
}
