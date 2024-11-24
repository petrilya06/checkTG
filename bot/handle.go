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
		SendPhotoKeyboard(bot, userID, Index)
	case "choose":
		SendPhotoKeyboard(bot, userID, Index)
		users[userID].State = StateCheckPic
		fmt.Println(users[userID].State)
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

		EditText(bot, chatID, fmt.Sprintf("За данную аватарку будет выплачиться %s рублей", Price[Index]))
		EditMedia(bot, chatID)
		EditReplyMarkup(bot, chatID, InlineKeyboard)

	case "1", "2":
		EditText(bot, chatID, "Хотите оставить аватарку?")
		EditReplyMarkup(bot, chatID, InlineKeyboardConfirm)
	}
}
