package models

import "time"

// LateRequest представляет заявку на опоздание
type LateRequest struct {
	ID       int64         `json:"id"`
	UserID   int64         `json:"user_id"`
	Reason   string        `json:"reason"`
	Duration time.Duration `json:"duration"`
	Date     time.Time     `json:"date"`
	Status   string        `json:"status"` // pending, approved, rejected
}

// UserState хранит текущее состояние пользователя в процессе создания заявки
type UserState struct {
	UserID      int64       `json:"user_id"`
	CurrentStep string      `json:"current_step"` // late_reason, late_time, late_confirmation
	RequestType string      `json:"request_type"` // late, absence, etc.
	RequestData interface{} `json:"request_data"` // Может быть LateRequest, AbsenceRequest и т.д.
}

// Repository интерфейс для работы с хранилищем
type LateRequestRepository interface {
	SaveLateRequest(request LateRequest) error
	GetLateRequests(userID int64) ([]LateRequest, error)
	DeleteLateRequest(requestID int64) error

	SaveUserState(state UserState) error
	GetUserState(userID int64) (UserState, error)
	ClearUserState(userID int64) error
}
