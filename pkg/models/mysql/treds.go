package mysql

import (
	"database/sql"

	"thredly.com/thredly/pkg/models"
)

type TredModel struct {
	DB *sql.DB
}

// Insert - Метод для создания новой заметки в базе данных.
func (m *TredModel) Insert(UserId int, ParentTredId, Content, Categoryid string) (int, error) {
	stmt := `
	INSERT INTO tred (user_id, publication_date, views_count, content, photo, parent_tred_id, category_id) 
VALUES (
    CAST(? AS UNSIGNED), 
    NOW(), 
    FLOOR(100 + RAND() * (1000 - 100 + 1)), 
    ?, 
    'img.jpg', 
    ?, 
    IFNULL(NULLIF(?, ''), 0)
);
	`

	result, err := m.DB.Exec(stmt, UserId, Content, ParentTredId, Categoryid)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *TredModel) GetLatest() ([]*models.Tred, error) {
	stmt := `SELECT t.id, u.first_name, u.last_name, u.photo, publication_date, views_count, content, t.photo, parent_tred_id, COALESCE(c.name, 'Все категории')
		 	FROM tred t
			JOIN useraccount u on u.id = t.user_id
			LEFT JOIN category c ON c.id = t.category_id
			WHERE parent_tred_id IS NULL or parent_tred_id = 0 
			ORDER BY publication_date;`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var treds []*models.Tred
	for rows.Next() {
		s := &models.Tred{}
		err = rows.Scan(&s.ID, &s.UFirstName, &s.ULastName, &s.UPhoto, &s.PublicationDate, &s.ViewsCount, &s.Content, &s.Photo, &s.ParentId, &s.Category)
		if err != nil {
			return nil, err
		}
		// Получаем теги для текущего треда
		s.Tags, err = m.getTagsForThread(s.ID)
		if err != nil {
			return nil, err
		}
		treds = append(treds, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return treds, nil
}

func (m *TredModel) GetChildTreds(cat string) ([]*models.Tred, error) {
	var rows *sql.Rows
	var err error

	if cat == "all" {
		stmt := `SELECT t.id, u.first_name, u.last_name, u.photo, publication_date, views_count, content, t.photo, parent_tred_id, COALESCE(c.name, 'Все категории') 
				FROM tred t
				JOIN useraccount u on u.id = t.user_id
				LEFT JOIN category c ON c.id = t.category_id
				WHERE parent_tred_id IS NULL OR parent_tred_id = 0
				ORDER BY publication_date DESC;`
		rows, err = m.DB.Query(stmt)
	} else {
		stmt := `SELECT t.id, u.first_name, u.last_name, u.photo, publication_date, views_count, content, u.photo, parent_tred_id, COALESCE(c.name, 'Все категории')
                 FROM tred t
				 JOIN useraccount u on u.id = t.user_id
				 LEFT JOIN category c ON c.id = t.category_id
                 WHERE (parent_tred_id IS NULL OR parent_tred_id = 0) 
                 AND c.Name = ?
                 ORDER BY publication_date DESC;`
		rows, err = m.DB.Query(stmt, cat)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mainTreds []*models.Tred
	for rows.Next() {
		s := &models.Tred{}
		err = rows.Scan(&s.ID, &s.UFirstName, &s.ULastName, &s.UPhoto, &s.PublicationDate, &s.ViewsCount, &s.Content, &s.Photo, &s.ParentId, &s.Category)
		if err != nil {
			return nil, err
		}
		// Получаем теги для текущего треда
		s.Tags, err = m.getTagsForThread(s.ID)
		if err != nil {
			return nil, err
		}
		mainTreds = append(mainTreds, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Получаем дочерние треды рекурсивно
	for _, mainTred := range mainTreds {
		err = m.populateChildTreds(mainTred)
		if err != nil {
			return nil, err
		}
	}

	return mainTreds, nil
}

func (m *TredModel) GetRecommendedTreds(userID int, cat string) ([]*models.Tred, error) {
	var rows *sql.Rows
	var err error

	if cat == "all" {
		stmt := `SELECT t.id, u.first_name, u.last_name, u.photo, publication_date, views_count, content, t.photo, parent_tred_id, COALESCE(c.name, 'Все категории') 
				FROM tred t
				JOIN useraccount u on u.id = t.user_id
				LEFT JOIN category c ON c.id = t.category_id
				WHERE parent_tred_id IS NULL OR parent_tred_id = 0
                AND user_id in (select distinct subscribed_to_id from userSubscription where subscriber_id = ?)
				ORDER BY publication_date DESC;`
		rows, err = m.DB.Query(stmt, userID)
	} else {
		stmt := `SELECT t.id, u.first_name, u.last_name, u.photo, publication_date, views_count, content, u.photo, parent_tred_id, COALESCE(c.name, 'Все категории')
                 FROM tred t
				 JOIN useraccount u on u.id = t.user_id
				 LEFT JOIN category c ON c.id = t.category_id
                 WHERE (parent_tred_id IS NULL OR parent_tred_id = 0) 
                 AND c.Name = ?
				 AND user_id in (select distinct subscribed_to_id from userSubscription where subscriber_id = ?)
                 ORDER BY publication_date DESC;`
		rows, err = m.DB.Query(stmt, cat, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mainTreds []*models.Tred
	for rows.Next() {
		s := &models.Tred{}
		err = rows.Scan(&s.ID, &s.UFirstName, &s.ULastName, &s.UPhoto, &s.PublicationDate, &s.ViewsCount, &s.Content, &s.Photo, &s.ParentId, &s.Category)
		if err != nil {
			return nil, err
		}
		// Получаем теги для текущего треда
		s.Tags, err = m.getTagsForThread(s.ID)
		if err != nil {
			return nil, err
		}
		mainTreds = append(mainTreds, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Получаем дочерние треды рекурсивно
	for _, mainTred := range mainTreds {
		err = m.populateChildTreds(mainTred)
		if err != nil {
			return nil, err
		}
	}

	return mainTreds, nil
}

func (m *TredModel) populateChildTreds(parent *models.Tred) error {
	stmt := `SELECT t.id, u.first_name, u.last_name, u.photo, publication_date, views_count, content, u.photo, parent_tred_id, COALESCE(c.name, 'Все категории') 
			    FROM tred t
				JOIN useraccount u on u.id = t.user_id
				LEFT JOIN category c ON c.id = t.category_id
				WHERE parent_tred_id = ? 
				order by publication_date desc;`
	rows, err := m.DB.Query(stmt, parent.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var childTreds []*models.Tred
	for rows.Next() {
		s := &models.Tred{}
		err = rows.Scan(&s.ID, &s.UFirstName, &s.ULastName, &s.UPhoto, &s.PublicationDate, &s.ViewsCount, &s.Content, &s.Photo, &s.ParentId, &s.Category)
		if err != nil {
			return err
		}
		// Получаем теги для текущего треда
		s.Tags, err = m.getTagsForThread(s.ID)
		if err != nil {
			return err
		}
		childTreds = append(childTreds, s)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	parent.ChildTreds = childTreds

	// Рекурсивно получаем дочерние треды для каждого дочернего треда
	for _, childTred := range childTreds {
		err = m.populateChildTreds(childTred)
		if err != nil {
			return err
		}
	}

	return nil
}

// getTagsForThread - новый метод для получения тегов для конкретного треда
func (m *TredModel) getTagsForThread(threadID int) ([]string, error) {
	stmt := `SELECT tag FROM thread_tags WHERE thread_id = ?`
	rows, err := m.DB.Query(stmt, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

// Вотвращает количество публикаций созданные пользователем
func (m *TredModel) GetCount(userId int) (int, error) {
	stmt := `select count(*) from tred where user_id = ?;`
	row := m.DB.QueryRow(stmt, userId)

	var s int
	err := row.Scan(&s)
	if err != nil {
		return 0, err
	}

	return s, nil
}
