package utils

import (
	"log"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/zenv"
	"github.com/joho/godotenv"
)

type TEnv struct {
	DiscordClientID     string `zog:"DISCORD_CLIENT_ID"`
	DiscordClientSecret string `zog:"DISCORD_CLIENT_SECRET"`
	RedirectURI         string `zog:"REDIRECT_URI"`
	DiscordURL          string `zog:"DISCORD_URL"`
	DBUser              string `zog:"DB_USER"`
	DBPass              string `zog:"DB_PASS"`
	DashboardURL        string `zog:"DASHBOARD_URL"`
}

var Env TEnv

func LoadEnv() {
	godotenv.Load()

	schema := z.Struct(z.Shape{
		"DiscordClientID":     z.String().Required(),
		"DiscordClientSecret": z.String().Required(),
		"RedirectURI":         z.String().URL().Required(),
		"DiscordURL":          z.String().URL().Default("https://discord.com"),
		"DBUser":              z.String().Required(),
		"DBPass":              z.String().Required(),
		"DashboardURL":        z.String().URL().Required(),
	})

	err := schema.Parse(zenv.NewDataProvider(), &Env)
	if err != nil {
		log.Fatal(err)
	}
}
