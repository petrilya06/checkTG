package bot

import (
	"fmt"
	"github.com/checkTG/db"
	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"io"
	"net/http"
	"os"
)

func SendText(bot *tg.Bot, user *db.User, message string) {
	msg, _ := bot.SendMessage(tu.Message(
		tu.ID(user.TgID),
		message,
	))

	user.LastMessageID = msg.MessageID
	err := database.UpdateUser(*user)
	if err != nil {
		return
	}
}

func SendPhoto(bot *tg.Bot, user *db.User, keyboard tg.InlineKeyboardMarkup) {
	msg, _ := bot.SendPhoto(tu.Photo(
		tu.ID(user.TgID),
		tu.File(MustOpen(fmt.Sprintf("src/%d.jpg", Index))),
	).WithReplyMarkup(&keyboard))

	user.LastPhotoID = msg.MessageID
	err := database.UpdateUser(*user)
	if err != nil {
		return
	}
}

func SendPhotoKeyboard(bot *tg.Bot, user *db.User) {
	SendText(bot, user, fmt.Sprintf("Выбери подходяющую для вас аватарку\nЗа данную аватарку будет выплачиться %s рублей", Price[Index]))
	SendPhoto(bot, user, InlineKeyboard)
}

func EditPhotoKeyboard(bot *tg.Bot, user *db.User) {
	EditText(bot, user, fmt.Sprintf("За данную аватарку будет выплачиться %s рублей", Price[Index]))
	EditMedia(bot, user)
	EditReplyMarkup(bot, user, InlineKeyboard)
}

func EditText(bot *tg.Bot, user *db.User, msg string) {
	_, _ = bot.EditMessageText(&tg.EditMessageTextParams{
		ChatID:    tu.ID(user.TgID),
		MessageID: user.LastMessageID,
		Text:      msg,
	})
}

func EditMedia(bot *tg.Bot, user *db.User) {
	_, _ = bot.EditMessageMedia(&tg.EditMessageMediaParams{
		ChatID:    tu.ID(user.TgID),
		MessageID: user.LastPhotoID,
		Media:     tu.MediaPhoto(tu.File(MustOpen(fmt.Sprintf("src/%d.jpg", Index)))),
	})
}

func EditReplyMarkup(bot *tg.Bot, user *db.User, keyboard tg.InlineKeyboardMarkup) {
	_, _ = bot.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParams{
		ChatID:      tu.ID(user.TgID),
		MessageID:   user.LastPhotoID,
		ReplyMarkup: &keyboard,
	})
}

func DeleteMessage(bot *tg.Bot, user *db.User, message []int) {
	_ = bot.DeleteMessages(&tg.DeleteMessagesParams{
		ChatID:     tu.ID(user.TgID),
		MessageIDs: message,
	})
}

func downloadPhoto(bot *tg.Bot, user *db.User) {
	profilePhotos, err := bot.GetUserProfilePhotos(&tg.GetUserProfilePhotosParams{
		UserID: user.TgID,
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		fmt.Println("error in get user photo:", err)
		return
	}

	if len(profilePhotos.Photos) == 0 {
		SendText(bot, user, "У вас нет фотографий!")
		return
	}

	photo := profilePhotos.Photos[0][2]
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

	out, err := os.Create(fmt.Sprintf("src/photos/%d.jpg", user.TgID))
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
