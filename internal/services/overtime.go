package services

import (
	"fmt"
	"time"
)

// обновление баллов юзера на основе переработок из таблицы overtimes
func (s *Service) UpdatePointsFromOvertimes(userID string) error {
	rows, err := s.DB.Query(`
		SELECT hours, recorded_at 
		FROM overtimes 
		WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("ошибка выборки переработок: %w", err)
	}
	defer rows.Close()

	totalHours := 0
	for rows.Next() {
		var hours int
		var recordedAt time.Time
		if err := rows.Scan(&hours, &recordedAt); err != nil {
			return fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		totalHours += hours
	}

	points := totalHours

	_, err = s.DB.Exec(`
		INSERT INTO points (user_id, points) 
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET points = EXCLUDED.points
	`, userID, points)
	if err != nil {
		return fmt.Errorf("ошибка обновления баллов: %w", err)
	}

	return nil
}
