package postimgs

import (
	"context"
	"fmt"
	"log"
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

func (n *StatusNotifier) NotifyCount(ctx context.Context, index int, total int) {
	n.updateMessage(ctx, fmt.Sprintf(messages.LoadingImage, index+1, total))
}

func (n *StatusNotifier) NotifyUploadingSD(ctx context.Context) {
	n.updateMessage(ctx, messages.UploadingImagesSD)
}

func (n *StatusNotifier) NotifyUploadingHD(ctx context.Context) {
	n.updateMessage(ctx, messages.UploadingImagesHD)
}

func (n *StatusNotifier) updateMessage(ctx context.Context, text string) {
	if n.bot == nil {
		return
	}

	var err error

	if n.msg != nil {
		_, err = n.bot.EditMessageText(ctx, &bot.EditMessageTextParams{
			MessageID: n.msg.ID,
			ChatID:    n.userID,
			Text:      text,
		})
	} else {
		n.msg, err = n.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: n.userID,
			Text:   text,
		})
	}
	if err != nil {
		log.Println(err)
	}
}
