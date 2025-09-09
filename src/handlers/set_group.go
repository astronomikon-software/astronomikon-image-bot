package handlers

import (
	"context"
	"ortemios/imgbot/managedgroup"
	"ortemios/imgbot/messages"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func SetGroup(ctx context.Context, b *bot.Bot, update *models.Update) {
	groupID := update.Message.Chat.ID
	managedgroup.Set(groupID)
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: groupID,
		Text:   messages.GroupSet,
	})
}
