package modules

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
)

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
