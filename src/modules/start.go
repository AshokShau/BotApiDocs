package modules

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"time"
	_ "time/tzdata"
)

var StartTime = time.Now()

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	text := fmt.Sprintf("👋 Hello! I'm your handy Telegram Bot API assistant, built with GoTgBot.\n\n💡 Usage: <code>@%s your_query</code> - Quickly search for any method or type in the Telegram Bot API documentation.", b.User.Username)

	_, err := msg.Reply(b, text, &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	if err != nil {
		log.Printf("[start] Error sending message: %v", err)
		return err
	}

	return ext.EndGroups
}

func ping(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	startTime := time.Now()

	rest, err := msg.Reply(b, "<code>Pinging</code>", &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	if err != nil {
		return fmt.Errorf("ping: %v", err)
	}

	// Calculate latency
	elapsedTime := time.Since(startTime)

	// Calculate uptime
	uptime := time.Since(StartTime)
	formattedUptime := getFormattedDuration(uptime)

	location, _ := time.LoadLocation("Asia/Kolkata")
	responseText := fmt.Sprintf("Pinged in %vms (Latency: %.2fs) at %s\n\nUptime: %s", elapsedTime.Milliseconds(), elapsedTime.Seconds(), time.Now().In(location).Format(time.RFC1123), formattedUptime)

	_, _, err = rest.EditText(b, responseText, nil)
	if err != nil {
		return fmt.Errorf("ping: %v", err)
	}

	return ext.EndGroups
}
