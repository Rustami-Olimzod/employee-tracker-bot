package models

import "employee-tracker-bot/internal/entity"

// LateEventRepository интерфейс для работы с событиями
type LateEventRepository interface {
	AddLateEvent(event entity.Late) error
	ListLateEvents(userID int64) ([]entity.Late, error)
	DeleteLateEvent(id int64, userID int64) error
	UpdateLateEvent(event entity.Late) error // исправить тип!
}
