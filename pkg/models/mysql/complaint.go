package mysql

import (
	"database/sql"

	"thredly.com/thredly/pkg/models"
)

type ComplaintModel struct {
	DB *sql.DB
}

// Insert - Метод для создания новой жалобы в базе данных.
func (m *ComplaintModel) Insert(tredID string, userID int) (int, error) {
	// SQL запрос для вставки новой жалобы
	stmt := `INSERT INTO complaint (tred_id, user_id, complaint_date) 
		SELECT ?, ?, NOW() FROM DUAL 
		WHERE NOT EXISTS (
			SELECT id FROM complaint WHERE tred_id = ? AND user_id = ?
		);`

	// Используем метод Exec() для выполнения запроса
	result, err := m.DB.Exec(stmt, tredID, userID, tredID, userID)
	if err != nil {
		return 0, err
	}

	// Получаем ID последней вставленной записи
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Конвертируем ID в тип int перед возвратом
	return int(id), nil
}

// GetAll возвращает все жалобы с информацией о пользователе и треде.
func (m *ComplaintModel) GetAll() ([]*models.ComplaintWithDetails, error) {
	// SQL запрос для выборки всех жалоб с информацией о пользователе и треде
	stmt := `
		SELECT c.id, c.tred_id, c.user_id, c.complaint_date, u.first_name, u.last_name, t.content
		FROM complaint c
		JOIN userAccount u ON c.user_id = u.id
		JOIN tred t ON c.tred_id = t.id
		ORDER BY c.complaint_date DESC;
	`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var complaints []*models.ComplaintWithDetails
	for rows.Next() {
		c := &models.ComplaintWithDetails{}
		err := rows.Scan(&c.ID, &c.TredID, &c.UserID, &c.ComplaintDate, &c.UserFirstName, &c.UserLastName, &c.TredContent)
		if err != nil {
			return nil, err
		}
		complaints = append(complaints, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return complaints, nil
}
