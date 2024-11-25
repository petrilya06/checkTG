package bot

import (
	"fmt"
	"io"
	"net/http"
	"os"

	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func MustOpen(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return file
}

func SendText(bot *tg.Bot, chatID int64, message string) {
	msg, _ := bot.SendMessage(tu.Message(
		tu.ID(chatID),
		message,
	))

	users[chatID].LastMessageID = msg.MessageID
}

func SendPhoto(bot *tg.Bot, chatID int64, keyboard tg.InlineKeyboardMarkup) {
	msg, _ := bot.SendPhoto(tu.Photo(
		tu.ID(chatID),
		tu.File(MustOpen(fmt.Sprintf("src/%d.jpg", Index))),
	).WithReplyMarkup(&keyboard))

	users[chatID].LastPhotoID = msg.MessageID
}

func SendPhotoKeyboard(bot *tg.Bot, chatID int64) {
	SendText(bot, chatID, fmt.Sprintf("Выбери подходяющую для вас аватарку\nЗа данную аватарку будет выплачиться %s рублей", Price[Index]))
	SendPhoto(bot, chatID, InlineKeyboard)
}

func EditPhotoKeyboard(bot *tg.Bot, chatID int64) {
	EditText(bot, chatID, fmt.Sprintf("За данную аватарку будет выплачиться %s рублей", Price[Index]))
	EditMedia(bot, chatID)
	EditReplyMarkup(bot, chatID, InlineKeyboard)
}

func EditText(bot *tg.Bot, chatID int64, msg string) {
	_, _ = bot.EditMessageText(&tg.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: users[chatID].LastMessageID,
		Text:      msg,
	})
}

func EditMedia(bot *tg.Bot, chatID int64) {
	_, _ = bot.EditMessageMedia(&tg.EditMessageMediaParams{
		ChatID:    tu.ID(chatID),
		MessageID: users[chatID].LastPhotoID,
		Media:     tu.MediaPhoto(tu.File(MustOpen(fmt.Sprintf("src/%d.jpg", Index)))),
	})
}

func EditReplyMarkup(bot *tg.Bot, chatID int64, keyboard tg.InlineKeyboardMarkup) {
	_, _ = bot.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParams{
		ChatID:      tu.ID(chatID),
		MessageID:   users[chatID].LastPhotoID,
		ReplyMarkup: &keyboard,
	})
}

func DeleteMessage(bot *tg.Bot, chatID int64, message []int) {
	_ = bot.DeleteMessages(&tg.DeleteMessagesParams{
		ChatID:     tu.ID(chatID),
		MessageIDs: message,
	})
}

func downloadPhoto(bot *tg.Bot, userID int64) {
	profilePhotos, err := bot.GetUserProfilePhotos(&tg.GetUserProfilePhotosParams{
		UserID: userID,
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		fmt.Println("error in get user photo:", err)
		return
	}

	if len(profilePhotos.Photos) == 0 {
		fmt.Println("user doesn`t have a photo.")
		return
	}

	photo := profilePhotos.Photos[0][0]
	file, err := bot.GetFile(&tg.GetFileParams{FileID: photo.FileID})
	if err != nil {
		fmt.Println("error in get file:", err)
		return
	}

	response, err := http.Get("https://api.telegram.org/file/bot" + os.Getenv("TOKEN") + "/" + file.FilePath)
	if err != nil {
		fmt.Println("error in download:", err)
		return
	}
	defer response.Body.Close()

	out, err := os.Create(fmt.Sprintf("src/photos/%d.jpg", users[userID].ID))
	if err != nil {
		fmt.Println("error in make a file:", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		fmt.Println("error in write to file:", err)
		return
	}
}
