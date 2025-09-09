package handlers

import (
	"context"
	"errors"
	"ortemios/imgbot/managedgroup"
	"ortemios/imgbot/messages"
	"ortemios/imgbot/types"
	"strings"

	"github.com/go-telegram/bot"
)

var ErrInvalidSyntax = errors.New("set group: invalid command syntax")

func SetGroup(ctx context.Context, b *bot.Bot, update *types.Update) error {
	tokens := strings.Split(update.Text, " ")
	if len(tokens) != 2 {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.ChatID,
			Text:   messages.InvalidSyntax,
		})
		return ErrInvalidSyntax
	}
	err := managedgroup.Set(tokens[1])
	if err != nil {
		return err
	}
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.ChatID,
		Text:   messages.GroupSet,
	})
	return err
}
