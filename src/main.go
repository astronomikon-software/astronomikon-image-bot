package main

import (
	"context"
	"ortemios/imgbot/handlers"
	"ortemios/imgbot/handlers/post_images"
	"ortemios/imgbot/messages"
	"ortemios/imgbot/util"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const SetGroupCommand = "/setgroup"

var botToken string
var allowedUsers []string

func init() {
	botToken = util.MustGetEnv("BOT_TOKEN")
	allowedUsers = strings.Split(util.MustGetEnv("ALLOWED_USERS"), ",")
}

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	go func() {
		if update == nil || update.Message == nil {
			return
		}
		if !isUserAllowed(update.Message.From) {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.From.ID,
				Text:   messages.AccessDenied,
			})
			return
		}
		text := strings.TrimSpace(update.Message.Text)
		if text == SetGroupCommand {
			handlers.SetGroup(ctx, b, update)
		} else if len(text) > 0 {
			post_images.PostImages(ctx, b, update)
		}
	}()
}

func isUserAllowed(user *models.User) bool {
	for _, u := range allowedUsers {
		if u == strconv.Itoa(int(user.ID)) || u == user.Username {
			return true
		}
	}
	return false
}
