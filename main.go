package main

import (
	"context"
	"embed"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

//go:embed frontend/dist
var server embed.FS

//go:embed frontend/dist/index.html
var index []byte

type ClientCreds struct {
	Web struct {
		ClientID          string   `json:"client_id"`
		ClientSecret      string   `json:"client_secret"`
		RedirectURIs      []string `json:"redirect_uris"`
		JavascriptOrigins []string `json:"javascript_origins"`
	} `json:"web"`
	SessionSecret string `json:"session_secret"`
}

func LoadClientCreds(filename string) (*ClientCreds, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var creds ClientCreds
	if err := json.Unmarshal(content, &creds); err != nil {
		return nil, err
	}
	return &creds, nil
}

type providerContextKey struct {
	key string
}

var ProviderParamKey = &providerContextKey{"provider"}

var isProd = os.Getenv("PROD") == "1"

func setupOauthRoutes(r *gin.Engine, clientCreds *ClientCreds) {
	goth.UseProviders(
		google.New(
			clientCreds.Web.ClientID,
			clientCreds.Web.ClientSecret,
			clientCreds.Web.RedirectURIs[0]),
	)

	gothic.GetProviderName = func(req *http.Request) (string, error) {
		provider := req.Context().Value(ProviderParamKey).(string)
		return provider, nil
	}
	secret, err := hex.DecodeString(clientCreds.SessionSecret)
	if err != nil {
		panic(err)
	}
	store := sessions.NewCookieStore(secret)
	store.MaxAge(86400 * 30)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd // TODO: change to true in production

	gothic.Store = store

	r.GET("/api/auth/:provider/callback", func(c *gin.Context) {
		provider := c.Param("provider")
		ctx := context.WithValue(c.Request.Context(), ProviderParamKey, provider)
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request.WithContext(ctx))
		if err != nil {
			log.Printf("Error completing user auth: %v", err)
			c.String(http.StatusUnauthorized, "Unauthorized")
			return
		}
		c.JSON(http.StatusOK, user)
	})

	r.GET("/api/auth/:provider", func(c *gin.Context) {
		provider := c.Param("provider")
		ctx := context.WithValue(c.Request.Context(), ProviderParamKey, provider)
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request.WithContext(ctx))
		if err != nil {
			gothic.BeginAuthHandler(c.Writer, c.Request.WithContext(ctx))
			return
		}
		c.JSON(http.StatusOK, user)
	})
}

func main() {
	clientCreds, err := LoadClientCreds("client_secret.json")
	if err != nil {
		log.Fatalf("Failed to load client creds: %v", err)
	}

	r := gin.Default()
	setupOauthRoutes(r, clientCreds)

	r.GET("/api/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.Use(static.Serve("/", static.EmbedFolder(server, "frontend/dist")))
	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/api") {
			c.Status(http.StatusOK)
			c.Writer.Write(index)
		}
	})

	if err := r.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
