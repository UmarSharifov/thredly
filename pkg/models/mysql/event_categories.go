package mysql

import (
	"database/sql"

	"thredly.com/thredly/pkg/models"
)

type EventsCategoriesMoodel struct {
	DB *sql.DB
}

func (m *EventsCategoriesMoodel) Latest() ([]*models.EventCategories, error) {

	stmt := `SELECT * from event_category`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var eventCats []*models.EventCategories

	for rows.Next() {
		s := &models.EventCategories{}
		err = rows.Scan(&s.ID, &s.Name)
		if err != nil {
			return nil, err
		}
		eventCats = append(eventCats, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return eventCats, nil
}
