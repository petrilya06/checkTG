package bot

import (
	"fmt"
	"log"
	"os"

	"github.com/checkTG/db" // Убедитесь, что путь к вашему модулю правильный
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func Run() {
	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Создание нового бота
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	// Создание соединения с базой данных
	connStr := "user=postgres password=qwerty123456 host=localhost sslmode=disable"
	database, err := db.NewDatabase(connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	if err := database.CreateTable(); err != nil {
		log.Fatalf("Ошибка при создании таблицы: %v", err)
	}

	defer database.Close()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Добавление пользователя в базу данных
		err := database.AddUser("eorkgoergerg", int(update.Message.From.ID), true, true)
		if err != nil {
			log.Printf("error in add: %v", err)
		}

		allUsers, err := database.GetAllUsers()
		if err != nil {
			log.Printf("error in get all users: %v", err)
		}

		fmt.Println("\n\n\n\n", allUsers)
	}
}
