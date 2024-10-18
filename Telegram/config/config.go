package config

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

var (
	Token      string
	OwnerId    int64
	WebhookUrl string
	Port       string
	Vercel     string
)

func init() {
	Token = os.Getenv("TOKEN")
	OwnerId = toInt64(os.Getenv("OWNER_ID"))
	Vercel = os.Getenv("Vercel")
	WebhookUrl = os.Getenv("WEBHOOK_URL")
	Port = os.Getenv("PORT")
}

func toInt64(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}
