package bot

import (
	"fmt"
	"github.com/checkTG/db"
	tg "github.com/mymmrac/telego"
	"log"
	"time"
)

func ProcessStateMachine(bot *tg.Bot, user *db.User) {
	switch user.Status {
	case "start":
		SendTextWithKeyboard(bot, user,
			"Привет! Помогу вам подзаработать очень простым способом.", ChooseKeyboard)

		user.Status = StateChoosePic
		if err := database.UpdateUser(*user); err != nil {
			return
		}

	case "choose":
		// записываем последнию фотографию, если она поменяется, то аннулируем таймер на выплату
		repeatIndex := user.SelectPic

		SendPhotoKeyboard(bot, user)

		// проверку на изменение фотографии
		if repeatIndex != user.SelectPic {
			user.CountCheck = 0
		}

		user.Status = StateCheckPic
		if err := database.UpdateUser(*user); err != nil {
			return
		}

	case "check":
		go func() {
			ticker := time.NewTicker(3 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					downloadPhoto(bot, user)
					compare := ComparePhotos(fmt.Sprintf("src/%d.jpg", user.SelectPic),
						fmt.Sprintf("src/photos/%d.jpg", user.TgID))

					if compare {
						user.CountCheck++
						if user.CountCheck >= 24 {
							// TODO сделать свою систему выплат
							SendText(bot, user, fmt.Sprintf("Ура друг! Тебе назначена выплата в %s рублей!",
								Price[user.SelectPic]))
							user.CountCheck = 0
						}
					} else {
						// если были успешные проверки и пользователь резко сменил аватарку, то говорим ему об этом
						if user.CountCheck != 0 {
							SendText(bot, user, "Фотография не совпадает. Выплата аннулируется.")
						}

						user.CountCheck = 0
					}

					err := database.UpdateUser(*user)
					if err != nil {
						log.Printf("error updating user: %v", err)
						return
					}
				}
			}
		}()
	}
}
