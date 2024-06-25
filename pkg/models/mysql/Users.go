package mysql

import (
	"database/sql"
	"errors"
	"strconv"

	"thredly.com/thredly/pkg/models"
)

// UserModel - Определяем тип который обертывает пул подключения sql.DB
type UserModel struct {
	DB *sql.DB
}

// Insert - Метод для создания новой заметки в базе дынных.
func (m *UserModel) Insert(LastName, FirstName, Photo, Email, PhoneNumber, UserLogin, UserPwd string) (int, error) {
	stmt := `INSERT INTO userAccount (last_name, first_name, photo, date_of_birth, email, phone_number, userLogin, userPassword) VALUES
	(?, ?, ?, '1990-05-15', ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, LastName, FirstName, Photo, Email, PhoneNumber, UserLogin, UserPwd)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Update - Метод для обновления данных пользователя в базе данных.
func (m *UserModel) Update(idU, LastName, FirstName, Photo, Email, PhoneNumber string) (int, error) {
	// SQL запрос для обновления данных пользователя
	stmt := `UPDATE useraccount SET last_name = ?, first_name = ?, 
	photo = ?, email = ?, 
	phone_number = ? WHERE id = ?`

	// Выполняем запрос
	_, err := m.DB.Exec(stmt, LastName, FirstName, Photo, Email, PhoneNumber, idU)
	if err != nil {
		return 0, err
	}

	// Конвертируем idU в int
	id, err := strconv.Atoi(idU)
	if err != nil {
		return 0, err
	}

	// Возвращаем ID пользователя
	return id, nil
}

func (m *UserModel) GetSubscribleCount(userId int) (int, error) {
	// SQL запрос для обновления данных пользователя
	stmt := `select count(*) from userSubscription where subscriber_id = ?`

	row := m.DB.QueryRow(stmt, userId)

	// Переменная для сохранения результата
	var s int

	// Считываем результат запроса
	err := row.Scan(&s)
	if err != nil {
		return 0, err
	}

	// Возвращаем результат
	return s, nil
}

func (m *UserModel) GetSubscribleToCount(userId int) (int, error) {
	// SQL запрос для обновления данных пользователя
	stmt := `select count(*) from userSubscription where subscribed_to_id = ?`

	row := m.DB.QueryRow(stmt, userId)

	// Переменная для сохранения результата
	var s int

	// Считываем результат запроса
	err := row.Scan(&s)
	if err != nil {
		return 0, err
	}

	// Возвращаем результат
	return s, nil
}

// Get - Метод для возвращения данных пользователя по его идентификатору ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	// SQL запрос для получения данных одной записи.
	stmt := `SELECT id, first_name, last_name, email, date_of_birth, phone_number, photo FROM useraccount WHERE id = ?`

	// Используем метод QueryRow() для выполнения SQL запроса,
	// передавая id в качестве значения для плейсхолдера.
	row := m.DB.QueryRow(stmt, id)

	// Инициализируем указатель на новую структуру User.
	s := &models.User{}

	// Используем row.Scan() для копирования значений из каждого поля от sql.Row в
	// соответствующее поле в структуре User. Обратите внимание, что аргументы
	// для row.Scan - это указатели на место, куда требуется скопировать данные
	// и количество аргументов должно быть точно таким же, как количество
	// столбцов в таблице базы данных.
	err := row.Scan(&s.ID, &s.FirstName, &s.LastName, &s.Email, &s.DateOfBirthDay, &s.PhoneNumber, &s.Photo)
	if err != nil {
		// Проверяем, если запрос был выполнен с ошибкой. Если ошибка обнаружена, то
		// возвращаем нашу ошибку из модели models.ErrNoRecord.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	// Если все хорошо, возвращается объект User.
	return s, nil
}

// Get - Метод для возвращения данных пользователя по его идентификатору ID.
func (m *UserModel) GetSubscrible(userId int) ([]*models.User, error) {
	// SQL запрос для получения данных одной записи.
	stmt := `SELECT id, first_name, last_name, email, date_of_birth, phone_number, photo 
			FROM useraccount AS u
			WHERE EXISTS (
				SELECT 1
				FROM userSubscription AS s
				WHERE s.subscriber_id = ? AND s.subscribed_to_id = u.id
			) and id != ?
LIMIT 3
	`

	rows, err := m.DB.Query(stmt, userId, userId)
	if err != nil {
		return nil, err
	}

	var users []*models.User

	for rows.Next() {
		s := &models.User{}
		err := rows.Scan(&s.ID, &s.FirstName, &s.LastName, &s.Email, &s.DateOfBirthDay, &s.PhoneNumber, &s.Photo)
		if err != nil {
			// Проверяем, если запрос был выполнен с ошибкой. Если ошибка обнаружена, то
			// возвращаем нашу ошибку из модели models.ErrNoRecord.
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, err
			}
		}
		users = append(users, s)

	}

	// Если все хорошо, возвращается объект User.
	return users, nil
}

// Get - Метод для возвращения данных пользователя по его идентификатору ID.
func (m *UserModel) GetSubscribleAll(userId int) ([]*models.User, error) {
	// SQL запрос для получения данных одной записи.
	stmt := `SELECT id, first_name, last_name, email, date_of_birth, phone_number, photo 
			FROM useraccount AS u
			WHERE NOT EXISTS (
				SELECT 1
				FROM userSubscription AS s
				WHERE s.subscriber_id = ? AND s.subscribed_to_id = u.id
			) and id != ?;
	`

	rows, err := m.DB.Query(stmt, userId, userId)
	if err != nil {
		return nil, err
	}

	var users []*models.User

	for rows.Next() {
		s := &models.User{}
		err := rows.Scan(&s.ID, &s.FirstName, &s.LastName, &s.Email, &s.DateOfBirthDay, &s.PhoneNumber, &s.Photo)
		if err != nil {
			// Проверяем, если запрос был выполнен с ошибкой. Если ошибка обнаружена, то
			// возвращаем нашу ошибку из модели models.ErrNoRecord.
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, err
			}
		}
		users = append(users, s)

	}

	// Если все хорошо, возвращается объект User.
	return users, nil
}

// Get - Метод для возвращения данных заметки по её логину и паролю
func (m *UserModel) GetUser(login string, password string) (*models.User, error) {
	// SQL запрос для получения данных одной записи.
	stmt := `SELECT ID, first_name, last_name, email, date_of_birth, phone_number FROM useraccount WHERE userLogin = ? and userPassword = ? `
	// (&s.ID, &s.FirstName, &s.LastName, &s.Email, &s.DateOfBirthDay, &s.PhoneNumber)
	// Используем метод QueryRow() для выполнения SQL запроса,
	// передавая ненадежную переменную id в качестве значения для плейсхолдера
	// Возвращается указатель на объект sql.Row, который содержит данные записи.
	row := m.DB.QueryRow(stmt, login, password)

	// Инициализируем указатель на новую структуру Snippet.
	s := &models.User{}

	// Используйте row.Scan(), чтобы скопировать значения из каждого поля от sql.Row в
	// соответствующее поле в структуре User. Обратите внимание, что аргументы
	// для row.Scan - это указатели на место, куда требуется скопировать данные
	// и количество аргументов должно быть точно таким же, как количество
	// столбцов в таблице базы данных.

	err := row.Scan(&s.ID, &s.FirstName, &s.LastName, &s.Email, &s.DateOfBirthDay, &s.PhoneNumber)
	if err != nil {
		// Специально для этого случая, мы проверим при помощи функции errors.Is()
		// если запрос был выполнен с ошибкой. Если ошибка обнаружена, то
		// возвращаем нашу ошибку из модели models.ErrNoRecord.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	// Если все хорошо, возвращается объект User.
	return s, nil
}

// Authenticate - Метод для проверки учетных данных пользователя.
func (m *UserModel) Authenticate(UserLogin, UserPwd string) (int, error) {
	var id int
	var storedPassword string
	stmt := `SELECT id, userPassword FROM userAccount WHERE userLogin = ?`
	row := m.DB.QueryRow(stmt, UserLogin)
	err := row.Scan(&id, &storedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	if UserPwd != storedPassword {
		return 0, models.ErrInvalidCredentials
	}

	return id, nil
}
