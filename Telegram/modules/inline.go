package modules

import (
	"encoding/json"
	"fmt"
	"github.com/AshokShau/BotApiDocs/Telegram/config"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

const (
	maxMessageLength = 4096 // Telegram's maximum message length
	apiURL           = "https://github.com/PaulSonOfLars/telegram-bot-api-spec/raw/main/api.json"
)

// apiCache is a global cache for storing API methods and types.
var apiCache struct {
	sync.RWMutex
	Methods map[string]Method
	Types   map[string]Type
}

type Method struct {
	Name        string   `json:"name"`
	Description []string `json:"description"`
	Href        string   `json:"href"`
	Returns     []string `json:"returns"`
	Fields      []Field  `json:"fields,omitempty"`
}

type Type struct {
	Name        string   `json:"name"`
	Description []string `json:"description"`
	Href        string   `json:"href"`
	Fields      []Field  `json:"fields,omitempty"`
}

type Field struct {
	Name        string   `json:"name"`
	Types       []string `json:"types"`
	Required    bool     `json:"required"`
	Description string   `json:"description"`
}

// isVercel checks if the application is running on Vercel.
func isVercel() bool {
	return config.VERCEL == "1"
}

// fetchAPI fetches the API documentation from a remote source and updates the apiCache.
func fetchAPI() error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch API: %w", err)
	}
	defer resp.Body.Close()

	var apiDocs struct {
		Methods map[string]Method `json:"methods"`
		Types   map[string]Type   `json:"types"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiDocs); err != nil {
		return fmt.Errorf("failed to decode API response: %w", err)
	}

	apiCache.Lock()
	defer apiCache.Unlock()
	apiCache.Methods = apiDocs.Methods
	apiCache.Types = apiDocs.Types

	return nil
}

// StartAPICacheUpdater starts a goroutine that periodically updates the API cache.
func StartAPICacheUpdater(interval time.Duration) {
	go func() {
		for {
			if !isVercel() { // Only fetch if not on Vercel
				if err := fetchAPI(); err != nil {
					log.Println("Error updating API documentation:", err)
				}
			}
			time.Sleep(interval)
		}
	}()
}

// getAPICache returns a snapshot of the current API cache.
func getAPICache() (map[string]Method, map[string]Type, error) {
	if isVercel() {
		// Fetch directly if on Vercel
		if err := fetchAPI(); err != nil {
			return nil, nil, err
		}
	}

	apiCache.RLock()
	defer apiCache.RUnlock()
	return apiCache.Methods, apiCache.Types, nil
}

// inlineQueryHandler handles inline queries from the bot.
func inlineQueryHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	query := strings.TrimSpace(ctx.InlineQuery.Query)
	parts := strings.Fields(query)

	if len(parts) == 0 {
		return sendEmptyQueryResponse(bot, ctx)
	}

	kueri := strings.Join(parts[1:], " ")
	if strings.EqualFold(parts[0], "botapi") {
		query = kueri
	}

	methods, types, err := getAPICache()
	if err != nil {
		return fmt.Errorf("failed to get API cache: %w", err)
	}
	results := searchAPI(query, methods, types)

	if len(results) == 0 {
		return sendNoResultsResponse(bot, ctx, query)
	}

	if len(results) > 50 {
		results = results[:50]
	}

	_, err = ctx.InlineQuery.Answer(bot, results, &gotgbot.AnswerInlineQueryOpts{IsPersonal: true})
	return err
}

// sendEmptyQueryResponse sends a response for an empty inline query.
func sendEmptyQueryResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.InlineQuery.Answer(bot, nil, &gotgbot.AnswerInlineQueryOpts{
		IsPersonal: true,
		CacheTime:  5,
		Button: &gotgbot.InlineQueryResultsButton{
			Text:           "Type 'your_query' to search!",
			StartParameter: "start",
		},
	})
	return err
}

// searchAPI searches the API methods and types for the given query.
func searchAPI(query string, methods map[string]Method, types map[string]Type) []gotgbot.InlineQueryResult {
	var results []gotgbot.InlineQueryResult
	lowerQuery := strings.ToLower(query)

	search := func(name, href, msg string) {
		results = append(results, createInlineResult(name, href, msg, href))
	}

	for name, method := range methods {
		if strings.Contains(strings.ToLower(name), lowerQuery) {
			search(name, method.Href, buildMethodMessage(method))
		}
	}

	for name, typ := range types {
		if strings.Contains(strings.ToLower(name), lowerQuery) {
			search(name, typ.Href, buildTypeMessage(typ))
		}
	}

	return results
}

// sendNoResultsResponse sends a response when no results are found for the query.
func sendNoResultsResponse(bot *gotgbot.Bot, ctx *ext.Context, query string) error {
	_, err := ctx.InlineQuery.Answer(bot, []gotgbot.InlineQueryResult{noResultsArticle(query)}, &gotgbot.AnswerInlineQueryOpts{
		IsPersonal: true,
		CacheTime:  500,
	})
	return err
}

// buildMethodMessage builds a message string for a given API method.
func buildMethodMessage(method Method) string {
	return buildMessage(method.Name, method.Description, method.Returns, method.Fields, method.Href)
}

// buildTypeMessage builds a message string for a given API type.
func buildTypeMessage(typ Type) string {
	return buildMessage(typ.Name, typ.Description, nil, typ.Fields, typ.Href)
}

func buildMessage(name string, description []string, returns []string, fields []Field, href string) string {
	var msgBuilder strings.Builder
	msgBuilder.WriteString(fmt.Sprintf("<b>%s</b>\n", name))
	msgBuilder.WriteString(fmt.Sprintf("Description: %s\n\n", sanitizeHTML(strings.Join(description, ", "))))
	if returns != nil {
		msgBuilder.WriteString("<b>Returns:</b> " + strings.Join(returns, ", ") + "\n")
	}

	if len(fields) > 0 {
		msgBuilder.WriteString("<b>Fields:</b>\n")
		for _, field := range fields {
			msgBuilder.WriteString(fmt.Sprintf("<code>%s</code> (<b>%s</b>) - Required: <code>%t</code>\n", field.Name, strings.Join(field.Types, ", "), field.Required))
			msgBuilder.WriteString(sanitizeHTML(field.Description) + "\n\n")
		}
	}

	message := msgBuilder.String()
	if len(message) > maxMessageLength {
		return fmt.Sprintf("See full documentation: %s", href)
	}

	return message
}

// createInlineResult creates an inline query result for a given API method or type.
func createInlineResult(title, url, message, methodUrl string) gotgbot.InlineQueryResult {
	return gotgbot.InlineQueryResultArticle{
		Id:      strconv.Itoa(rand.Intn(100000)),
		Title:   title,
		Url:     url,
		HideUrl: true,
		InputMessageContent: gotgbot.InputTextMessageContent{
			MessageText:        message,
			ParseMode:          gotgbot.ParseModeHTML,
			LinkPreviewOptions: &gotgbot.LinkPreviewOptions{PreferSmallMedia: true},
		},
		Description: "View more details",
		ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{Text: "Open Docs", Url: methodUrl}},
				{{Text: "Search Again", SwitchInlineQueryCurrentChat: &title}},
			},
		},
	}
}

// noResultsArticle creates an inline query result indicating no results were found.
func noResultsArticle(query string) gotgbot.InlineQueryResult {
	return gotgbot.InlineQueryResultArticle{
		Id:    strconv.Itoa(rand.Intn(100000)),
		Title: "No Results Found!",
		InputMessageContent: gotgbot.InputTextMessageContent{
			MessageText: fmt.Sprintf("<i>👋 Sorry, I couldn't find any results for '%s'. Try searching with a different keyword!</i>", query),
			ParseMode:   gotgbot.ParseModeHTML,
		},
		Description: "No results found for your query.",
		ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{Text: "Search Again", SwitchInlineQueryCurrentChat: &query}},
			},
		},
	}
}

// sanitizeHTML removes unsupported HTML tags from the message.
func sanitizeHTML(input string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(input, "")
}
