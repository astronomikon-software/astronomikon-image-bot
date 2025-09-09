package post_images

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"ortemios/imgbot/imgboard"
	"ortemios/imgbot/managedgroup"
	"ortemios/imgbot/messages"
	"ortemios/imgbot/types"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const maxUrlsPerMessage = 10

var ErrMaxUrlsExceeded = errors.New("post images: max urls per messages exceeded")

var isBusy = false

func PostImages(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.Chat.ID

	if isBusy {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   messages.ServiceIsBusy,
		})
		return
	}
	isBusy = true
	defer func() { isBusy = false }()

	groupID, err := getGroupIdOrNotify(ctx, b, userID)
	if err != nil {
		return
	}

	urls, err := extractUrlsOrNotify(ctx, b, userID, update.Message.Text)
	if err != nil {
		return
	}

	notifier := NewStatusNotifier(b, userID)
	defer notifier.Finish(ctx)

	for index, url := range urls {
		notifier.Notify(ctx, index, len(urls))
		err := postImage(ctx, b, url, groupID)
		if err != nil {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: userID,
				Text:   messages.LoadingImageFailed,
			})
			break
		}
	}
}

func postImage(
	ctx context.Context,
	b *bot.Bot,
	url string,
	groupID types.ChatID,
) error {
	sender := func(ctx context.Context, image imgboard.Image) error {
		imageData, err := readImageData(image)
		if err != nil {
			return err
		}

		_, err = b.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID: groupID,
			Photo: &models.InputFileUpload{
				Filename: image.Filename,
				Data:     bytes.NewReader(imageData),
			},
		})
		if err != nil {
			return err
		}

		_, err = b.SendDocument(ctx, &bot.SendDocumentParams{
			ChatID: groupID,
			Document: &models.InputFileUpload{
				Filename: image.Filename,
				Data:     bytes.NewReader(imageData),
			},
		})
		if err != nil {
			return err
		}

		return nil
	}
	return imgboard.PostImage(ctx, url, imgboard.FetchHtml, sender)
}

func getGroupIdOrNotify(ctx context.Context, b *bot.Bot, userID types.UserID) (types.ChatID, error) {
	groupID, err := managedgroup.Get()
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   messages.GroupNotSet,
		})
	}
	return groupID, err
}

func extractUrlsOrNotify(ctx context.Context, b *bot.Bot, userID types.UserID, message string) ([]string, error) {
	urls, err := imgboard.ExtractUrls(message)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   messages.InvalidUrl,
		})
		return nil, err
	}
	if len(urls) > maxUrlsPerMessage {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   messages.MaxUrlsExceeded,
		})
		return nil, ErrMaxUrlsExceeded
	}
	return urls, nil
}

func readImageData(image imgboard.Image) ([]byte, error) {
	resp, err := http.Get(image.Url)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}
