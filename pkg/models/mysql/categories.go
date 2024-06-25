package mysql

import (
	"database/sql"

	"thredly.com/thredly/pkg/models"
)

type CategoriesMoodel struct {
	DB *sql.DB
}

func (m *CategoriesMoodel) Latest() ([]*models.Categories, error) {
	// Пишем SQL запрос, который мы хотим выполнить.
	stmt := `SELECT * from category`

	// Используем метод Query() для выполнения нашего SQL запроса.
	// В ответ мы получим sql.Rows, который содержит результат нашего запроса.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var categories []*models.Categories

	for rows.Next() {
		s := &models.Categories{}
		err = rows.Scan(&s.ID, &s.Name)
		if err != nil {
			return nil, err
		}
		// Добавляем структуру в срез.
		categories = append(categories, s)
	}

	// Когда цикл rows.Next() завершается, вызываем метод rows.Err(), чтобы узнать
	// если в ходе работы у нас не возникла какая либо ошибка.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Если все в порядке, возвращаем срез с данными.
	return categories, nil
}
