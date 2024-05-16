package mysql

import (
	"database/sql"

	"thredly.com/thredly/pkg/models"
)

// TredModel - Определяем тип который обертывает пул подключения sql.DB
type TredModel struct {
	DB *sql.DB
}

// Insert - Метод для создания новой заметки в базе дынных.
func (m *TredModel) Insert(UserId, Content string) (int, error) {
	// Ниже будет SQL запрос, который мы хотим выполнить. Мы разделили его на две строки
	// для удобства чтения (поэтому он окружен обратными кавычками
	// вместо обычных двойных кавычек).
	// stmt := `INSERT INTO snippets (title, content, created, expires)
	// VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	stmt := `INSERT INTO tred (user_id, publication_date, views_count, content, photo) 
	VALUES (CAST(? AS UNSIGNED), NOW(), FLOOR(100 + RAND() * (1000 - 100 + 1)), ?, 'img.jpg');`

	// Используем метод Exec() из встроенного пула подключений для выполнения
	// запроса. Первый параметр это сам SQL запрос, за которым следует
	// заголовок заметки, содержимое и срока жизни заметки. Этот
	// метод возвращает объект sql.Result, который содержит некоторые основные
	// данные о том, что произошло после выполнении запроса.
	result, err := m.DB.Exec(stmt, UserId, Content)
	if err != nil {
		return 0, err
	}

	// Используем метод LastInsertId(), чтобы получить последний ID
	// созданной записи из таблицу snippets.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Возвращаемый ID имеет тип int64, поэтому мы конвертируем его в тип int
	// перед возвратом из метода.
	return int(id), nil
}
func (m *TredModel) GetLatest() ([]*models.Tred, error) {
	stmt := `SELECT * FROM tred ORDER BY id DESC LIMIT 10;`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	// Создаем пустой массив, который будет содержать все заметки
	var treds []*models.Tred
	// Цикл, который будет проходить по всем строкам в наборе результатов
	for rows.Next() {
		s := &models.Tred{}
		err = rows.Scan(&s.ID, &s.UserId, &s.PublicationDate, &s.ViewsCount, &s.Content, &s.Photo)
		if err != nil {
			return nil, err
		}
		treds = append(treds, s)
	}
	// Когда цикл rows.Next() завершается, вызываем метод rows.Err(), чтобы узнать
	// если в ходе работы у нас не возникла какая либо ошибка.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Если все в порядке, возвращаем срез с данными.
	return treds, nil

}
