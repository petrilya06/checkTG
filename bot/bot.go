package bot

import (
	"fmt"
	"log"
	"os"

	"github.com/checkTG/db"
	"github.com/joho/godotenv"
	tg "github.com/mymmrac/telego"
)

func Run() {
	var err error
	database, err = db.NewDatabase()
	if err != nil {
		log.Fatalf("error in database connection: %v", err)
	}

	if err := database.CreateTable(); err != nil {
		log.Fatalf("error in create table: %v", err)
	}

	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tg.NewBot(os.Getenv("TOKEN"), tg.WithDefaultDebugLogger())
	if err != nil {
		panic(err)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)

	defer bot.StopLongPolling()
	defer database.Close()

	for update := range updates {
		if update.Message != nil {
			userMessageID := update.Message.From.ID
			userMessage, _ := database.GetUserByID(userMessageID)

			fmt.Println(userMessage)
			HandleMessage(&update, userMessage, userMessageID)
			ProcessStateMachine(bot, userMessage)
		}

		if update.CallbackQuery != nil {
			userCallbackID := update.CallbackQuery.From.ID
			userCallback, _ := database.GetUserByID(userCallbackID)

			fmt.Println(userCallback)
			HandleCallbackQuery(bot, &update, userCallback)
			ProcessStateMachine(bot, userCallback)
		}
	}
}
