package main

import (
    "os"

    "github.com/gin-gonic/gin"
    "github.com/maximRnback/gin-oidc"
    "net/url"
    "github.com/gin-contrib/sessions"


    "github.com/go-chat-bot/bot/slack"
    _ "github.com/mike-dunton/go-slack-bot-plugin-example/example"
)

func main() {
   go slack.Run(os.Getenv("dunton_dev_SLACK_TOKEN"))

    r := gin.Default()
    store := sessions.NewCookieStore([]byte("secret"))
	  r.Use(sessions.Sessions("mysession", store))
    issuerUrl, _ := url.Parse(os.Getenv("dunton_dev_OKTA_URL")+"/oauth2/default")
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
  	r.GET("/ping", func(c *gin.Context) {
      serverSession := sessions.Default(c)
  		c.JSON(200, gin.H{
  			"message": serverSession.Get("oidcClaims"),
  		})
  	})
  	r.Run(":5050") // listen and serve on 0.0.0.0:8080
}
