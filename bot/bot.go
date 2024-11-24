package bot

import (
	"log"
	"os"

	"github.com/checkTG/db"
	"github.com/joho/godotenv"
	tg "github.com/mymmrac/telego"
)

const (
	StateStart     = "start"
	StateChoosePic = "choose"
	StateCheckPic  = "check"
	StateWait      = "wait"
)

type UserState struct {
	ID            int64
	State         string
	LastPhotoID   int
	LastMessageID int
}

var users = make(map[int64]*UserState)
var database *db.Database

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
			HandleMessage(bot, &update)
		}

		if update.CallbackQuery != nil {
			HandleCallbackQuery(bot, &update)
		}
	}
}
