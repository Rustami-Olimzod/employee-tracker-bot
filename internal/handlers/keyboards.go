package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetMainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Опоздание"),
			tgbotapi.NewKeyboardButton("Изменить"),
		),
	)
}
func GetPointsKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Баллы"),
			tgbotapi.NewKeyboardButton("Добавить переработку"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Обновить баллы"),
			tgbotapi.NewKeyboardButton("Получить награду"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Главное меню"),
		),
	)
}

func GetLateReasonsKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Пробки"),
			tgbotapi.NewKeyboardButton("Общественный транспорт"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Проблемы со здоровьем"),
			tgbotapi.NewKeyboardButton("Семейные обстоятельства"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Другая причина"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
}

func GetTimeKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("5 минут"),
			tgbotapi.NewKeyboardButton("10 минут"),
			tgbotapi.NewKeyboardButton("15 минут"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("30 минут"),
			tgbotapi.NewKeyboardButton("1 час"),
			tgbotapi.NewKeyboardButton("2 часа"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Указать вручную"),
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
}

func GetEditKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Причину"),
			tgbotapi.NewKeyboardButton("Время"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
}

func GetConfirmationInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Отправить", "confirm"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", "cancel"),
		),
	)
}
