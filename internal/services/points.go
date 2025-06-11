package services

import (
	"fmt"
)

// добавить переработку юзеру
func (s *Service) AddOvertime(userID string, hours int) string {
	_, err := s.DB.Exec(`
		INSERT INTO overtimes (user_id, hours, recorded_at)
		VALUES ($1, $2, NOW())`, userID, hours)
	if err != nil {
		return "Ошибка при добавлении переработки"
	}
	return fmt.Sprintf("Добавлено %d часов переработки.", hours)
}

// получение баллов юзера
func (s *Service) GetPoints(userID string) string {
	var points int
	err := s.DB.QueryRow(`
		SELECT points 
		FROM points 
		WHERE user_id = $1`, userID).Scan(&points)
	if err != nil {
		return "У вас пока нет баллов."
	}
	return fmt.Sprintf("У вас %d баллов.", points)
}

// обмен баллов на отгулы
func (s *Service) SpendPoints(userID string, reward string) string {
	var points int
	err := s.DB.QueryRow(`
		SELECT points 
		FROM points 
		WHERE user_id = $1`, userID).Scan(&points)
	if err != nil {
		return "У вас пока нет баллов."
	}
	if points < 8 {
		return "Недостаточно баллов для обмена (нужно минимум 8)."
	}

	_, err = s.DB.Exec(`
		UPDATE points 
		SET points = points - 8 
		WHERE user_id = $1`, userID)
	if err != nil {
		return "Ошибка при списании баллов."
	}

	return fmt.Sprintf("Баллы списаны! Вы получили: %s", reward)
}
