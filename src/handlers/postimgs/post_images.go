package postimgs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"ortemios/imgbot/handlers/postimgs/imgproc"
	"ortemios/imgbot/managedgroup"
	"ortemios/imgbot/messages"
	"ortemios/imgbot/types"
	"ortemios/imgbot/util"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const maxUrlsPerMessage = 10
const maxImageDimensionsSum = 10000
const maxImageSizeBytes = 10485760

var logPrefix = "PostImages"

var ErrMaxUrlsExceeded = errors.New("post images: max urls per messages exceeded")

var isBusy = false

func PostImages(ctx context.Context, b *bot.Bot, update *types.Update) error {
	logCalled(update)
	if isBusy {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.From.ID,
			Text:   messages.ServiceIsBusy,
		})
		logRejectedBusy(update)
		return err
	}
	isBusy = true
	defer func() { isBusy = false }()

	groupID, err := getGroupIdOrNotify(ctx, b, update)
	if err != nil {
		return err
	}

	urls, err := extractUrlsOrNotify(ctx, b, update)
	if err != nil {
		return err
	}

	notifier := NewStatusNotifier(b, update.From.ID)
	err = func() error {
		images := make([]*Image, 0, len(urls))
		for res := range loadImages(ctx, urls) {
			if res.err != nil {
				logImageLoadFailed(update, res.url, err)
				return res.err
			} else {
				notifier.NotifyCount(ctx, len(images), len(urls))
				images = append(images, res.image)
				logImageLoaded(update, res.image.Url)
			}
		}

		resizedImages, err := resizeImages(images)
		if err != nil {
			logError(update, err)
			return err
		}

		logUploadingSD(update)
		notifier.NotifyUploadingSD(ctx)

		_, err = b.SendMediaGroup(ctx, &bot.SendMediaGroupParams{
			ChatID: groupID,
			Media: util.Map(resizedImages, func(image *Image) models.InputMedia {
				return &models.InputMediaPhoto{
					Media:           fmt.Sprintf("attach://%s", image.Filename),
					MediaAttachment: bytes.NewReader(image.Data),
				}
			}),
		})
		if err != nil {
			logUploadingSDFailed(update)
			return err
		}

		logUploadingHD(update)
		notifier.NotifyUploadingHD(ctx)
		_, err = b.SendMediaGroup(ctx, &bot.SendMediaGroupParams{
			ChatID: groupID,
			Media: util.Map(images, func(image *Image) models.InputMedia {
				return &models.InputMediaDocument{
					Media:           fmt.Sprintf("attach://%s", image.Filename),
					MediaAttachment: bytes.NewReader(image.Data),
				}
			}),
		})
		if err != nil {
			logUploadingHDFailed(update)
			return err
		}

		return nil
	}()

	logFinished(update, err == nil)
	notifier.Finish(ctx, err == nil)
	return err
}

type loadImageResult struct {
	url   string
	image *Image
	err   error
}

func loadImages(ctx context.Context, urls []string) chan loadImageResult {
	out := make(chan loadImageResult)
	go func() {
		defer close(out)
		for _, url := range urls {
			image, err := LoadImage(ctx, url)
			if err != nil {
				out <- loadImageResult{err: err, url: url}
			} else {
				out <- loadImageResult{image: image, url: url}
			}
		}
	}()
	return out
}

func getGroupIdOrNotify(ctx context.Context, b *bot.Bot, update *types.Update) (types.GroupID, error) {
	groupID, err := managedgroup.Get()
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.From.ID,
			Text:   messages.GroupNotSet,
		})
		logGroupNotSet(update)
	}
	return groupID, err
}

func extractUrlsOrNotify(ctx context.Context, b *bot.Bot, update *types.Update) ([]string, error) {
	urls, err := ExtractUrls(update.Text)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.From.ID,
			Text:   messages.InvalidUrl,
		})
		logInvalidUrl(update)
		return nil, err
	}
	if len(urls) > maxUrlsPerMessage {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.From.ID,
			Text:   messages.MaxUrlsExceeded,
		})
		logMaxUrlsExceeded(update)
		return nil, ErrMaxUrlsExceeded
	}
	return urls, nil
}

func resizeImages(images []*Image) ([]*Image, error) {
	out := make([]*Image, 0)
	for index, image := range images {
		resizedImageData, err := imgproc.ResizeImage(
			image.Data,
			maxImageDimensionsSum,
			maxImageSizeBytes,
		)
		if err != nil {
			return nil, fmt.Errorf("error resizing image [index=%v]: %w", index, err)
		}
		out = append(out, &Image{
			Url:      image.Url,
			Filename: image.Filename,
			Data:     resizedImageData,
		})
	}
	return out, nil
}

func logCalled(update *types.Update) {
	log.Printf("[%v %s] called\n", update.ID, logPrefix)
}

func logRejectedBusy(update *types.Update) {
	log.Printf("[%v %s] service is busy, rejected\n", update.ID, logPrefix)
}

func logMaxUrlsExceeded(update *types.Update) {
	log.Printf("[%v %s] max urls exceeded\n", update.ID, logPrefix)
}

func logGroupNotSet(update *types.Update) {
	log.Printf("[%v %s] group not set\n", update.ID, logPrefix)
}

func logInvalidUrl(update *types.Update) {
	log.Printf("[%v %s] invalid url\n", update.ID, logPrefix)
}

func logImageLoaded(update *types.Update, url string) {
	log.Printf("[%v %s] %q image loaded\n", update.ID, logPrefix, url)
}

func logImageLoadFailed(update *types.Update, url string, err error) {
	log.Printf("[%v %s] %q image loading failed (%v)\n", update.ID, logPrefix, url, err)
}

func logUploadingSD(update *types.Update) {
	log.Printf("[%v %s] uploading SD\n", update.ID, logPrefix)
}

func logUploadingSDFailed(update *types.Update) {
	log.Printf("[%v %s] uploading SD failed\n", update.ID, logPrefix)
}

func logUploadingHD(update *types.Update) {
	log.Printf("[%v %s] uploading HD\n", update.ID, logPrefix)
}

func logUploadingHDFailed(update *types.Update) {
	log.Printf("[%v %s] uploading HD failed\n", update.ID, logPrefix)
}

func logError(update *types.Update, err error) {
	log.Printf("[%v %s] %v\n", update.ID, logPrefix, err)
}

func logFinished(update *types.Update, success bool) {
	log.Printf("[%v %s] finished, success=%v\n", update.ID, logPrefix, success)
}
