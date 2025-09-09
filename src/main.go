package main

import (
	"context"
	"log"
	"ortemios/imgbot/handlers"
	"ortemios/imgbot/handlers/post_images"
	"ortemios/imgbot/messages"
	"ortemios/imgbot/types"
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
		bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			if b != nil {
				err := handler(ctx, b, buildUpdate(update))
				if err != nil {
					log.Println(err)
				}
			}
		}),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *types.Update) error {
	if b == nil || update == nil {
		return nil
	}
	if !isUserAllowed(update.From) {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.From.ID,
			Text:   messages.AccessDenied,
		})
		return err
	}
	text := strings.TrimSpace(update.Text)
	if strings.HasPrefix(text, SetGroupCommand) {
		return handlers.SetGroup(ctx, b, update)
	} else if len(text) > 0 {
		return post_images.PostImages(ctx, b, update)
	}
	return nil
}

func buildUpdate(update *models.Update) *types.Update {
	if update.Message != nil {
		return &types.Update{
			ChatID: strconv.Itoa(int(update.Message.Chat.ID)),
			From:   buildUser(update.Message.From),
			Text:   update.Message.Text,
		}
	}
	return nil
}

func buildUser(user *models.User) *types.User {
	return &types.User{ID: user.ID, Username: user.Username}
}

func isUserAllowed(user *types.User) bool {
	for _, u := range allowedUsers {
		if u == strconv.Itoa(int(user.ID)) || u == user.Username {
			return true
		}
	}
	return false
}
