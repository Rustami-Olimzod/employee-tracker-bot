package models

import (
	"database/sql"
	"employee-tracker-bot/internal/entity"
	"fmt"
)

type PgRepository struct {
	db *sql.DB
}

func NewPgRepository(db *sql.DB) *PgRepository {
	return &PgRepository{db: db}
}

func (r *PgRepository) AddLateEvent(event entity.Late) error {
	_, err := r.db.Exec(
		"INSERT INTO late_events (user_id, reason, time) VALUES ($1, $2, $3)",
		event.UserID, event.Reason, event.Time,
	)
	return err
}

func (r *PgRepository) DeleteLateEvent(id int64, userID int64) error {
	res, err := r.db.Exec(
		"DELETE FROM late_events WHERE id=$1 AND user_id=$2",
		id, userID,
	)
	count, _ := res.RowsAffected()
	if count == 0 {
		return fmt.Errorf("not found")
	}
	return err
}

func (r *PgRepository) ListLateEvents(userID int64) ([]entity.Late, error) {
	rows, err := r.db.Query(
		"SELECT id, reason, time FROM late_events WHERE user_id=$1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []entity.Late
	for rows.Next() {
		var e entity.Late
		var id int64
		if err := rows.Scan(&id, &e.Reason, &e.Time); err != nil {
			return nil, err
		}
		e.UserID = userID
		e.ID = id
		res = append(res, e)
	}
	return res, nil
}

func (r *PgRepository) ClearUser(userID int64) {
	// Не требуется для PostgreSQL, оставьте пустым
}

func (r *PgRepository) GetLate(userID int64) entity.Late {
	row := r.db.QueryRow(
		"SELECT id, reason, time FROM late_events WHERE user_id=$1 ORDER BY created_at DESC LIMIT 1",
		userID,
	)
	var e entity.Late
	var id int64
	if err := row.Scan(&id, &e.Reason, &e.Time); err == nil {
		e.UserID = userID
		e.ID = id
		return e
	}
	return entity.Late{}
}

// добавьте этот метод!

func (r *PgRepository) UpdateLateEvent(event entity.Late) error {
	_, err := r.db.Exec(
		"UPDATE late_events SET reason = $1, time = $2 WHERE id = $3 AND user_id = $4",
		event.Reason, event.Time, event.ID, event.UserID,
	)
	return err
}
