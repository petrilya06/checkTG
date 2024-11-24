package bot

import (
	"fmt"
	"log"

	"github.com/checkTG/db"
	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleMessage(bot *tg.Bot, update *tg.Update) {
	userID := int64(update.Message.Chat.ID)

	if update.Message.Text == "/start" {
		if _, exists := users[userID]; !exists {
			users[userID] = &UserState{ID: userID, State: StateStart}

			newUser := db.User{
				TgID:      userID,
				SelectPic: 0,
				HasText:   false,
				HasPic:    false,
			}

			if err := database.AddUser(newUser); err != nil {
				log.Printf("error in add user: %v", err)
			}
		}
	}

	switch users[userID].State {
	case "start":
		_, _ = bot.SendMessage(tu.Message(
			tu.ID(userID),
			"Привет! Помогу вам подзаработать очень простым способом.",
		))

		users[userID].State = StateChoosePic
		SendPhotoKeyboard(bot, userID)
	case "choose":
		SendPhotoKeyboard(bot, userID)
		users[userID].State = StateCheckPic
	case "check":
		users[userID].State = StateWait
	case "wait":
	}
}

func HandleCallbackQuery(bot *tg.Bot, update *tg.Update) {
	chatID := update.CallbackQuery.From.ID

	switch update.CallbackQuery.Data {
	case "next":
		Index += 1
		if Index == len(Photos) {
			Index = 0 // Сброс индекса, если достигнут конец списка
		}

		EditPhotoKeyboard(bot, chatID)

	case "150", "300":
		EditText(bot, chatID, "Хотите оставить аватарку?")
		EditReplyMarkup(bot, chatID, InlineKeyboardConfirm)

	case "no":
		EditPhotoKeyboard(bot, chatID)

	case "yes":
		if err := database.UpdateUser(db.User{
			TgID:      chatID,
			SelectPic: Index,
			HasPic:    false,
			HasText:   true,
		}); err != nil {
			log.Printf("error in update user pic: %v", err)
		}

		DeleteMessage(bot, chatID, []int{users[chatID].LastPhotoID, users[chatID].LastMessageID})
		SendText(bot, chatID, fmt.Sprintf("Спасибо за подтверждение! Вам назначена выплата %s рублей", Price[Index]))
		users[chatID].State = StateWait
	}
}
