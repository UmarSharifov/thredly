package mysql

import (
	"database/sql"
	"strings"
)

// TredModel - Определяем тип который обертывает пул подключения sql.DB
type TagsModel struct {
	DB *sql.DB
}

// Insert - Метод для создания новых тегов в базе данных.
func (m *TagsModel) Insert(threadID int, tags string) error {
	// Разделяем строку тегов на отдельные теги
	tagList := strings.Split(tags, ",")

	for _, tag := range tagList {
		// Обрезаем пробелы и проверяем, что тег не пустой
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}

		// Проверяем, существует ли уже тег для данного thread_id
		var exists bool
		err := m.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM thread_tags WHERE thread_id = ? AND tag = ?)", threadID, tag).Scan(&exists)
		if err != nil {
			return err
		}

		// Если тег не существует, добавляем его
		if !exists {
			stmt := `INSERT INTO thread_tags (thread_id, tag) VALUES (?, ?)`
			_, err := m.DB.Exec(stmt, threadID, tag)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
