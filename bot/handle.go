package bot

import (
	"fmt"
	"github.com/checkTG/db"
	tg "github.com/mymmrac/telego"
	"log"
)

func HandleMessage(update *tg.Update, user *db.User, userID int64) {
	if update.Message.Text == "/start" {
		newUser := db.User{
			TgID:          userID,
			SelectPic:     0,
			HasText:       false,
			HasPic:        false,
			Status:        StateStart,
			CountCheck:    0,
			LastMessageID: 0,
			LastPhotoID:   0,
		}

		if err := database.AddUser(newUser); err != nil {
			log.Printf("error in add user: %v", err)
		}

		user.Status = StateStart
		if err := database.UpdateUser(*user); err != nil {
			log.Printf("error updating user status: %v", err)
		}
	}

	if update.Message.Text == "Выбрать аватарку" {
		user.Status = StateChoosePic
		if err := database.UpdateUser(*user); err != nil {
			log.Printf("error updating user status: %v", err)
		}

		fmt.Println(user)
	}
}

func HandleCallbackQuery(bot *tg.Bot, update *tg.Update, user *db.User) {
	switch update.CallbackQuery.Data {
	case "next":
		Index += 1
		if Index == len(Photos) {
			Index = 0 // Сброс индекса, если достигнут конец списка
		}

		EditPhotoKeyboard(bot, user)

	case "150", "300":
		EditText(bot, user, "Хотите оставить аватарку?")
		EditReplyMarkup(bot, user, InlineKeyboardConfirm)

	case "no":
		EditPhotoKeyboard(bot, user)

	case "yes":
		user.SelectPic = Index
		if err := database.UpdateUser(*user); err != nil {
			return
		}

		DeleteMessages(bot, user, []int{user.LastPhotoID, user.LastMessageID})
		SendTextWithKeyboard(bot, user, fmt.Sprintf("Спасибо за подтверждение! "+
			"Вам назначена выплата %s рублей", Price[Index]), ChooseKeyboard)
	}
}
