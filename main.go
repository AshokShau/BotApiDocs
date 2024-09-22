package main

import (
	"fmt"
	"github.com/AshokShau/BotApiDocs/Telegram/config"
	"github.com/AshokShau/BotApiDocs/Telegram/modules"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"time"
)

const secretToken = "thisIsASecretToken"

var allowedUpdates = []string{"message", "inline_query"}

// initBot initializes the bot and updater with the provided token from the config.
// Returns the bot, updater, and an error if any.
func initBot() (*gotgbot.Bot, *ext.Updater, error) {
	if config.Token == "" {
		return nil, nil, fmt.Errorf("no token provided Add `TOKEN` to .env file")
	}

	bot, err := gotgbot.NewBot(config.Token, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create new bot: %w", err)
	}

	updater := ext.NewUpdater(modules.Dispatcher, nil)
	return bot, updater, nil
}

// configureWebhook sets up the webhook for the bot using the provided URL and token from the config.
// Returns an error if the webhook setup fails.
func configureWebhook(bot *gotgbot.Bot, updater *ext.Updater) error {
	if config.WebhookUrl == "" {
		return fmt.Errorf("WEBHOOK_URL is not provided")
	}

	_, err := bot.SetWebhook(config.WebhookUrl+config.Token, &gotgbot.SetWebhookOpts{
		MaxConnections:     40,
		DropPendingUpdates: true,
		AllowedUpdates:     allowedUpdates,
		SecretToken:        secretToken,
	})
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	return updater.StartWebhook(bot, config.Token, ext.WebhookOpts{
		ListenAddr:  "0.0.0.0:" + config.Port,
		SecretToken: secretToken,
	})
}

// startPolling starts polling for updates from the bot.
// Returns an error if polling fails.
func startPolling(bot *gotgbot.Bot, updater *ext.Updater) error {
	return updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates:    true,
		EnableWebhookDeletion: true,
		GetUpdatesOpts:        &gotgbot.GetUpdatesOpts{AllowedUpdates: allowedUpdates},
	})
}

// The main is the entry point of the application.
// It initializes the bot, configures the webhook or starts polling based on the configuration,
// and handles the bot's lifecycle.
func main() {
	bot, updater, err := initBot()
	if err != nil {
		log.Fatalf("Initialization error: %s", err)
	}

	mode := "Webhook"
	if err = configureWebhook(bot, updater); err != nil {
		log.Printf("Webhook configuration failed: %s", err)
		mode = "Polling"
		if err = startPolling(bot, updater); err != nil {
			log.Fatalf("Polling start failed: %s", err)
		}
	}
	// Start API cache updater with a 1-hour interval
	go modules.StartAPICacheUpdater(1 * time.Hour)
	log.Printf("Bot has been started as %s[%s] using %s", bot.FirstName, bot.Username, mode)
	updater.Idle()

	log.Printf("Bot has been stopped")
}
