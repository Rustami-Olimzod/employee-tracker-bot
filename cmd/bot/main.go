package main

import (
	"database/sql"
	"employee-tracker-bot/internal/services"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"employee-tracker-bot/internal/handlers"
	"employee-tracker-bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	_ = godotenv.Load()

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is not set")
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	repo := models.NewPgRepository(db)
	lateHandler := handlers.NewLateHandler(repo)
	pointService := services.NewService(db)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// 1. Обработка callback-запросов
		if update.CallbackQuery != nil {
			sess, hasSession := handlers.UserSessions[update.CallbackQuery.From.ID]
			if hasSession && sess.State != "" && len(sess.State) >= 4 && sess.State[:4] == "edit" {
				lateHandler.HandleEditCallback(bot, update.CallbackQuery)
			} else {
				lateHandler.HandleCallback(bot, update.CallbackQuery)
			}
			continue // <--- Очень важно!
		}
		// 2. Игнорируем всё, что не является сообщением
		if update.Message == nil {
			continue
		}

		// 3. Обработка редактирования, если это flow редактирования
		sess, hasSession := handlers.UserSessions[update.Message.From.ID]
		if hasSession && sess.State != "" && (sess.State == "edit_select" || sess.State == "edit_reason" || sess.State == "edit_time" || sess.State == "edit_manual_time" || sess.State == "edit_confirm") {
			lateHandler.HandleEditFlow(bot, update.Message)
			continue
		}

		// 4. Основные команды
		switch update.Message.Text {
		case "/start":
			sticker := tgbotapi.NewSticker(update.Message.Chat.ID, tgbotapi.FileID("CAACAgIAAxkBAAEGfKlkSJBOzP5xvU9vK9y1r4kAAQwAAuQAA8oAAWZc8pQAAQXbgSME"))
			bot.Send(sticker)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать! Выберите действие:")
			msg.ReplyMarkup = handlers.GetMainKeyboard()
			bot.Send(msg)
		case "Опоздание":
			lateHandler.HandleLateStart(bot, update.Message)
		case "Баллы":
			userID := fmt.Sprintf("%d", update.Message.From.ID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, pointService.GetPoints(userID))
			msg.ReplyMarkup = services.GetPointsKeyboard() // подключаем клавиатуру
			bot.Send(msg)
		case "Изменить":
			lateHandler.HandleEdit(bot, update.Message)
		default:
			lateHandler.HandleLateFlow(bot, update.Message)
		}
	}
}
