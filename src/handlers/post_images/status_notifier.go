package post_images

import (
	"context"
	"fmt"
	"ortemios/imgbot/messages"
	"ortemios/imgbot/types"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type StatusNotifier struct {
	msg    *models.Message
	bot    *bot.Bot
	userID types.UserID
}

func NewStatusNotifier(bot *bot.Bot, userID types.UserID) *StatusNotifier {
	return &StatusNotifier{
		msg:    nil,
		bot:    bot,
		userID: userID,
	}
}

func (n *StatusNotifier) Notify(ctx context.Context, index int, total int) {
	if n.bot == nil {
		return
	}
	if n.msg != nil {
		_, _ = n.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    n.userID,
			MessageID: n.msg.ID,
		})
	}
	n.msg, _ = n.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: n.userID,
		Text:   fmt.Sprintf(messages.LoadingImage, index+1, total),
	})
}

func (n *StatusNotifier) Finish(ctx context.Context) {
	if n.bot == nil {
		return
	}
	_, _ = n.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
		MessageID: n.msg.ID,
		ChatID:    n.userID,
	})
	_, _ = n.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: n.userID,
		Text:   messages.LoadingDone,
	})
}
