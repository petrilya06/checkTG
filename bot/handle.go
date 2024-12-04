package bot

import (
	"fmt"
	"github.com/checkTG/db"
	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"log"
)

func HandleMessage(bot *tg.Bot, update *tg.Update, user *db.User, userID int64) {
	if update.Message.Text == "/start" {
		if _, exists := users[userID]; !exists {
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
		}
	}

	fmt.Println(user)
	switch user.Status {
	case "start":
		_, _ = bot.SendMessage(tu.Message(
			tu.ID(userID),
			"Привет! Помогу вам подзаработать очень простым способом.",
		).WithReplyMarkup(ChooseKeyboard))

		user.Status = StateChoosePic
		err := database.UpdateUser(*user)
		if err != nil {
			return
		}

	case "choose":
		SendPhotoKeyboard(bot, user)
		downloadPhoto(bot, user)
		fmt.Println(user.LastMessageID, user.LastPhotoID)

		user.Status = StateCheckPic
		err := database.UpdateUser(*user)
		if err != nil {
			return
		}

	case "check":
		fmt.Println(user)
		downloadPhoto(bot, user)
		ComparePhotos(fmt.Sprintf("src/%d.jpg", Index),
			fmt.Sprintf("src/photos/%d.jpg", user.TgID))

		user.Status = StateWait
		err := database.UpdateUser(*user)
		if err != nil {
			return
		}

	case "wait":
	}
}

func HandleCallbackQuery(bot *tg.Bot, update *tg.Update, user *db.User, userID int64) {
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
		if err := database.UpdateUser(db.User{
			TgID:      userID,
			SelectPic: Index,
		}); err != nil {
			log.Printf("error in update user pic: %v", err)
		}

		DeleteMessage(bot, user, []int{user.LastPhotoID, user.LastMessageID})
		SendText(bot, user, fmt.Sprintf("Спасибо за подтверждение! Вам назначена выплата %s рублей", Price[Index]))
	}
}
