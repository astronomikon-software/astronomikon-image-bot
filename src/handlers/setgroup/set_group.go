package setgroup

import (
	"context"
	"errors"
	"log"
	"ortemios/imgbot/managedgroup"
	"ortemios/imgbot/messages"
	"ortemios/imgbot/types"
	"strings"

	"github.com/go-telegram/bot"
)

var logPrefix = "SetGroup"

var ErrInvalidSyntax = errors.New("set group: invalid command syntax")

func SetGroup(ctx context.Context, b *bot.Bot, update *types.Update) error {
	logCalled(update)
	tokens := strings.Split(update.Text, " ")
	if len(tokens) != 2 {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.ChatID,
			Text:   messages.InvalidSyntax,
		})
		logInvalidSyntax(update)
		return ErrInvalidSyntax
	}
	var groupID = tokens[1]
	log.Printf("%s '%v': called\n", logPrefix, groupID)
	err := managedgroup.Set(groupID)
	if err != nil {
		return err
	}
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.ChatID,
		Text:   messages.GroupSet,
	})
	logDone(update, groupID)
	return err
}

func logCalled(update *types.Update) {
	log.Printf("[%v %s] called\n", update.ID, logPrefix)
}

func logInvalidSyntax(update *types.Update) {
	log.Printf("[%v %s] invalid syntax\n", update.ID, logPrefix)
}

func logDone(update *types.Update, id types.GroupID) {
	log.Printf("[%v %s] group set to %v\n", update.ID, logPrefix, id)
}
