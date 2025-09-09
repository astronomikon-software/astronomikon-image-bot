package main

import (
	"context"
	"ortemios/imgbot/handlers"
	"ortemios/imgbot/handlers/post_images"
	"ortemios/imgbot/util"
	"os"
	"os/signal"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const SetGroupCommand = "/setgroup"

func main() {
	botToken := util.MustGetEnv("BOT_TOKEN")

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
		text := strings.TrimSpace(update.Message.Text)
		if text == SetGroupCommand {
			handlers.SetGroup(ctx, b, update)
		} else if len(text) > 0 {
			post_images.PostImages(ctx, b, update)
		}
	}()
}
