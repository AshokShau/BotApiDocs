package modules

import (
	"github.com/AshokShau/BotApiDocs/Telegram/config"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/inlinequery"
)

var Dispatcher = newDispatcher()

// newDispatcher creates a new dispatcher and loads modules.
func newDispatcher() *ext.Dispatcher {
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			_, _ = b.SendMessage(config.OwnerId, "An error occurred: "+err.Error(), nil)
			return ext.DispatcherActionNoop
		},
	})

	loadModules(dispatcher)
	return dispatcher
}

func loadModules(d *ext.Dispatcher) {
	d.AddHandler(handlers.NewCommand("start", start))
	d.AddHandler(handlers.NewInlineQuery(inlinequery.All, inlineQueryHandler))
}
