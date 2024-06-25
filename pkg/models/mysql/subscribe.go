package mysql

import (
	"database/sql"
)

// TredModel - Определяем тип который обертывает пул подключения sql.DB
type SubscribeModel struct {
	DB *sql.DB
}

// Insert - Метод для создания новой жалобы в базе данных.
func (m *SubscribeModel) Insert(subscriber_id int, subscribed_to_id string) (int, error) {
	// SQL запрос для вставки новой жалобы
	stmt := `INSERT INTO userSubscription (subscriber_id, subscribed_to_id, subscription_date)
				SELECT ?, ?, NOW()
				FROM DUAL
				WHERE NOT EXISTS (
					SELECT 1
					FROM userSubscription
					WHERE subscriber_id = ? AND subscribed_to_id = ?
				) LIMIT 1;`

	// Используем метод Exec() для выполнения запроса
	result, err := m.DB.Exec(stmt, subscriber_id, subscribed_to_id, subscriber_id, subscribed_to_id)
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
