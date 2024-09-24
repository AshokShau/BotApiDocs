package modules

import (
	"encoding/json"
	"fmt"
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

const maxMessageLength = 4096 // Telegram's maximum message length

// apiCache is a global cache for storing API methods and types.
var apiCache struct {
	sync.RWMutex
	Methods map[string]Method
	Types   map[string]Type
}

// Method represents an API method with its details.
type Method struct {
	Name        string   `json:"name"`
	Description []string `json:"description"`
	Href        string   `json:"href"`
	Returns     []string `json:"returns"`
	Fields      []Field  `json:"fields,omitempty"`
}

// Type represents an API type with its details.
type Type struct {
	Name        string   `json:"name"`
	Description []string `json:"description"`
	Href        string   `json:"href"`
	Fields      []Field  `json:"fields,omitempty"`
}

// Field represents a field in an API method or type.
type Field struct {
	Name        string   `json:"name"`
	Types       []string `json:"types"`
	Required    bool     `json:"required"`
	Description string   `json:"description"`
}

// fetchAPI fetches the API documentation from a remote source and updates the apiCache.
func fetchAPI() error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://github.com/PaulSonOfLars/telegram-bot-api-spec/raw/main/api.json")
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
	apiCache.Methods = apiDocs.Methods
	apiCache.Types = apiDocs.Types
	apiCache.Unlock()

	return nil
}

// StartAPICacheUpdater starts a goroutine that periodically updates the API cache.
func StartAPICacheUpdater(interval time.Duration) {
	go func() {
		for {
			if err := fetchAPI(); err != nil {
				log.Println("Error updating API documentation:", err)
			}
			time.Sleep(interval)
		}
	}()
}

// inlineQueryHandler handles inline queries from the bot.
func inlineQueryHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	query := strings.TrimSpace(ctx.InlineQuery.Query)
	parts := strings.Fields(query)

	if len(parts) < 1 {
		return sendEmptyQueryResponse(bot, ctx)
	}

	kueri := strings.Join(parts, " ")
	if strings.ToLower(parts[0]) == "botapi" {
		kueri = strings.Join(parts[1:], " ")
	}

	apiCache.RLock()
	methods := apiCache.Methods
	types := apiCache.Types
	apiCache.RUnlock()

	results := searchAPI(kueri, methods, types)

	if len(results) == 0 {
		return sendNoResultsResponse(bot, ctx, kueri)
	}

	if len(results) > 50 {
		results = results[:50]
	}

	_, err := ctx.InlineQuery.Answer(bot, results, &gotgbot.AnswerInlineQueryOpts{IsPersonal: true})
	return err
}

// sendEmptyQueryResponse sends a response for an empty inline query.
func sendEmptyQueryResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.InlineQuery.Answer(bot, []gotgbot.InlineQueryResult{}, &gotgbot.AnswerInlineQueryOpts{
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

	for name, method := range methods {
		if strings.Contains(strings.ToLower(name), strings.ToLower(query)) {
			msg := buildMethodMessage(method)
			results = append(results, createInlineResult(name, method.Href, msg, method.Href))
		}
	}

	for name, typ := range types {
		if strings.Contains(strings.ToLower(name), strings.ToLower(query)) {
			msg := buildTypeMessage(typ)
			results = append(results, createInlineResult(name, typ.Href, msg, typ.Href))
		}
	}

	return results
}

// sendNoResultsResponse sends a response when no results are found for the query.
func sendNoResultsResponse(bot *gotgbot.Bot, ctx *ext.Context, query string) error {
	_, err := ctx.InlineQuery.Answer(bot, []gotgbot.InlineQueryResult{noResultsArticle(query)}, &gotgbot.AnswerInlineQueryOpts{
		IsPersonal: true,
		CacheTime:  5,
	})
	return err
}

// buildMethodMessage builds a message string for a given API method.
func buildMethodMessage(method Method) string {
	var msgBuilder strings.Builder
	msgBuilder.WriteString(fmt.Sprintf("<b>%s</b>\n", method.Name))
	msgBuilder.WriteString(fmt.Sprintf("Description: %s\n\n", sanitizeHTML(strings.Join(method.Description, ", "))))
	msgBuilder.WriteString("<b>Returns:</b> " + strings.Join(method.Returns, ", ") + "\n")

	if len(method.Fields) > 0 {
		msgBuilder.WriteString("<b>Fields:</b>\n")
		for _, field := range method.Fields {
			msgBuilder.WriteString(fmt.Sprintf("<code>%s</code> (<b>%s</b>) - Required: <code>%t</code>\n", field.Name, strings.Join(field.Types, ", "), field.Required))
			msgBuilder.WriteString(sanitizeHTML(field.Description) + "\n\n")
		}
	}

	message := msgBuilder.String()
	if len(message) > maxMessageLength {
		return fmt.Sprintf("See full documentation: %s", method.Href)
	}
	return message
}

// buildTypeMessage builds a message string for a given API type.
func buildTypeMessage(typ Type) string {
	var msgBuilder strings.Builder
	msgBuilder.WriteString(fmt.Sprintf("<b>%s</b>\n", typ.Name))
	msgBuilder.WriteString(fmt.Sprintf("Description: %s\n\n", sanitizeHTML(strings.Join(typ.Description, ", "))))

	if len(typ.Fields) > 0 {
		msgBuilder.WriteString("<b>Fields:</b>\n")
		for _, field := range typ.Fields {
			msgBuilder.WriteString(fmt.Sprintf("<code>%s</code> (<b>%s</b>) - Required: <code>%t</code>\n", field.Name, strings.Join(field.Types, ", "), field.Required))
			msgBuilder.WriteString(sanitizeHTML(field.Description) + "\n\n")
		}
	}

	message := msgBuilder.String()
	if len(message) > maxMessageLength {
		return fmt.Sprintf("See full documentation: %s", typ.Href)
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
			MessageText: message,
			ParseMode:   gotgbot.ParseModeHTML,
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
	ok := "botapi"
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
				{{Text: "Search Again", SwitchInlineQueryCurrentChat: &ok}},
			},
		},
	}
}

// sanitizeHTML removes unsupported HTML tags from the message
func sanitizeHTML(input string) string {
	// This regex matches any HTML tags that are not supported
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(input, "")
}
