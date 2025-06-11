package handlers

import (
	"employee-tracker-bot/internal/entity"
	"employee-tracker-bot/internal/models"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserSession struct {
	Reason    string
	Time      string
	State     string
	EditEvent *entity.Late
}

// Экспортируемая переменная!
var UserSessions = map[int64]*UserSession{}

type LateHandler struct {
	repo models.LateEventRepository
}

func NewLateHandler(repo models.LateEventRepository) *LateHandler {
	return &LateHandler{repo: repo}
}

func (h *LateHandler) HandleLateStart(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	UserSessions[m.From.ID] = &UserSession{State: "reason"}
	msg := tgbotapi.NewMessage(m.Chat.ID, "Выберите причину опоздания:")
	msg.ReplyMarkup = GetLateReasonsKeyboard()
	bot.Send(msg)
}

func (h *LateHandler) HandleLateFlow(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	sess, ok := UserSessions[m.From.ID]
	if !ok {
		h.HandleLateStart(bot, m)
		return
	}

	switch sess.State {
	case "reason":
		if m.Text == "Назад" {
			msg := tgbotapi.NewMessage(m.Chat.ID, "Главное меню.")
			msg.ReplyMarkup = GetMainKeyboard()
			bot.Send(msg)
			delete(UserSessions, m.From.ID)
			return
		}
		sess.Reason = m.Text
		sess.State = "time"
		msg := tgbotapi.NewMessage(m.Chat.ID, "На сколько вы опоздали?")
		msg.ReplyMarkup = GetTimeKeyboard()
		bot.Send(msg)
	case "time":
		if m.Text == "Назад" {
			sess.State = "reason"
			msg := tgbotapi.NewMessage(m.Chat.ID, "Выберите причину опоздания:")
			msg.ReplyMarkup = GetLateReasonsKeyboard()
			bot.Send(msg)
			return
		}
		if m.Text == "Указать вручную" {
			sess.State = "manual_time"
			msg := tgbotapi.NewMessage(m.Chat.ID, "Пожалуйста, введите время опоздания вручную:")
			msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{}
			bot.Send(msg)
			return
		}
		sess.Time = m.Text
		sess.State = "confirm"
		text := fmt.Sprintf("Причина: %s\nВремя: %s\nПодтвердить отправку?", sess.Reason, sess.Time)
		msg := tgbotapi.NewMessage(m.Chat.ID, text)
		msg.ReplyMarkup = GetConfirmationInlineKeyboard()
		bot.Send(msg)
	case "manual_time":
		// Пользователь вручную вводит время
		sess.Time = m.Text
		sess.State = "confirm"
		text := fmt.Sprintf("Причина: %s\nВремя: %s\nПодтвердить отправку?", sess.Reason, sess.Time)
		msg := tgbotapi.NewMessage(m.Chat.ID, text)
		msg.ReplyMarkup = GetConfirmationInlineKeyboard()
		bot.Send(msg)
	default:
		msg := tgbotapi.NewMessage(m.Chat.ID, "Выберите действие:")
		msg.ReplyMarkup = GetMainKeyboard()
		bot.Send(msg)
	}
}

func (h *LateHandler) HandleCallback(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
	sess, ok := UserSessions[cq.From.ID]
	if !ok || sess.State != "confirm" {
		return
	}

	switch cq.Data {
	case "confirm":
		late := entity.Late{
			UserID: cq.From.ID,
			Reason: sess.Reason,
			Time:   sess.Time,
		}
		err := h.repo.AddLateEvent(late)
		text := "Событие сохранено!"
		if err != nil {
			text = "Ошибка при сохранении: " + err.Error()
		}
		msg := tgbotapi.NewMessage(cq.Message.Chat.ID, text)
		msg.ReplyMarkup = GetMainKeyboard()
		bot.Send(msg)
		delete(UserSessions, cq.From.ID)
	case "cancel":
		msg := tgbotapi.NewMessage(cq.Message.Chat.ID, "Отправка отменена.")
		msg.ReplyMarkup = GetMainKeyboard()
		bot.Send(msg)
		delete(UserSessions, cq.From.ID)
	}
}

func (h *LateHandler) HandleEdit(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	// Получаем последнее событие пользователя
	events, err := h.repo.ListLateEvents(m.From.ID)
	if err != nil || len(events) == 0 {
		msg := tgbotapi.NewMessage(m.Chat.ID, "Нет событий для изменения.")
		msg.ReplyMarkup = GetMainKeyboard()
		bot.Send(msg)
		return
	}
	latest := events[0]
	UserSessions[m.From.ID] = &UserSession{
		State:     "edit_select",
		Reason:    latest.Reason,
		Time:      latest.Time,
		EditEvent: &latest,
	}
	text := fmt.Sprintf("Последнее событие:\nПричина: %s\nВремя: %s\nЧто изменить?", latest.Reason, latest.Time)
	msg := tgbotapi.NewMessage(m.Chat.ID, text)
	msg.ReplyMarkup = GetEditKeyboard()
	bot.Send(msg)
}

func (h *LateHandler) HandleEditFlow(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	sess, ok := UserSessions[m.From.ID]
	if !ok || sess.State == "" {
		msg := tgbotapi.NewMessage(m.Chat.ID, "Нет события для редактирования.")
		msg.ReplyMarkup = GetMainKeyboard()
		bot.Send(msg)
		return
	}

	switch sess.State {
	case "edit_select":
		switch m.Text {
		case "Причину":
			sess.State = "edit_reason"
			msg := tgbotapi.NewMessage(m.Chat.ID, "Выберите новую причину:")
			msg.ReplyMarkup = GetLateReasonsKeyboard()
			bot.Send(msg)
			return
		case "Время":
			sess.State = "edit_time"
			msg := tgbotapi.NewMessage(m.Chat.ID, "Выберите новое время:")
			msg.ReplyMarkup = GetTimeKeyboard()
			bot.Send(msg)
			return
		case "Назад":
			msg := tgbotapi.NewMessage(m.Chat.ID, "Главное меню.")
			msg.ReplyMarkup = GetMainKeyboard()
			bot.Send(msg)
			delete(UserSessions, m.From.ID)
			return
		}
	case "edit_reason":
		if m.Text == "Назад" {
			sess.State = "edit_select"
			text := fmt.Sprintf("Причина: %s\nВремя: %s\nЧто изменить?", sess.Reason, sess.Time)
			msg := tgbotapi.NewMessage(m.Chat.ID, text)
			msg.ReplyMarkup = GetEditKeyboard()
			bot.Send(msg)
			return
		}
		sess.Reason = m.Text
		sess.State = "edit_confirm"
	case "edit_time":
		if m.Text == "Назад" {
			sess.State = "edit_select"
			text := fmt.Sprintf("Причина: %s\nВремя: %s\nЧто изменить?", sess.Reason, sess.Time)
			msg := tgbotapi.NewMessage(m.Chat.ID, text)
			msg.ReplyMarkup = GetEditKeyboard()
			bot.Send(msg)
			return
		}
		if m.Text == "Указать вручную" {
			sess.State = "edit_manual_time"
			msg := tgbotapi.NewMessage(m.Chat.ID, "Пожалуйста, введите новое время вручную:")
			msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{}
			bot.Send(msg)
			return
		}
		sess.Time = m.Text
		sess.State = "edit_confirm"
	case "edit_manual_time":
		sess.Time = m.Text
		sess.State = "edit_confirm"
	}

	if sess.State == "edit_confirm" {
		text := fmt.Sprintf("Изменить событие на:\nПричина: %s\nВремя: %s\nПодтвердить изменения?", sess.Reason, sess.Time)
		msg := tgbotapi.NewMessage(m.Chat.ID, text)
		msg.ReplyMarkup = GetConfirmationInlineKeyboard()
		bot.Send(msg)
	}
}

func (h *LateHandler) HandleEditCallback(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
	sess, ok := UserSessions[cq.From.ID]
	if !ok || sess.State != "edit_confirm" {
		msg := tgbotapi.NewMessage(cq.Message.Chat.ID, "Нет события для изменения.")
		msg.ReplyMarkup = GetMainKeyboard()
		bot.Send(msg)
		return
	}

	switch cq.Data {
	case "confirm":
		// Обновляем событие
		event := sess.EditEvent
		event.Reason = sess.Reason
		event.Time = sess.Time
		err := h.repo.UpdateLateEvent(*event)
		text := "Событие изменено!"
		if err != nil {
			text = "Ошибка при изменении: " + err.Error()
		}
		msg := tgbotapi.NewMessage(cq.Message.Chat.ID, text)
		msg.ReplyMarkup = GetMainKeyboard()
		bot.Send(msg)
		delete(UserSessions, cq.From.ID)
	case "cancel":
		msg := tgbotapi.NewMessage(cq.Message.Chat.ID, "Редактирование отменено.")
		msg.ReplyMarkup = GetMainKeyboard()
		bot.Send(msg)
		delete(UserSessions, cq.From.ID)
	}
}
