package mysql

import (
	"database/sql"

	"thredly.com/thredly/pkg/models"
)

type EventsMoodel struct {
	DB *sql.DB
}

// Insert - Метод для создания новой заметки в базе данных.
func (m *EventsMoodel) Insert(UserId int, ParentTredId, Content, Categoryid string) (int, error) {
	stmt := `
	INSERT INTO t_events (user_id, publication_date, views_count, content, photo) VALUES
	(?, NOW(), 100, ?,'photo1.jpg');
	`

	result, err := m.DB.Exec(stmt, UserId, Content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *EventsMoodel) Latest() ([]*models.Events, error) {

	stmt := `SELECT * from t_events`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var events []*models.Events

	for rows.Next() {
		s := &models.Events{}
		err = rows.Scan(&s.ID, &s.UserId, &s.PublicationDate, &s.ViewsCount, &s.Content, &s.CategoryId, &s.Photo)
		if err != nil {
			return nil, err
		}
		events = append(events, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
