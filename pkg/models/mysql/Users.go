package mysql

import (
	"database/sql"
	"errors"

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

// Insert - Метод для создания новой заметки в базе дынных.
func (m *UserModel) Update(idU, LastName, FirstName, Photo, Email, PhoneNumber string) (int, error) {
	// Ниже будет SQL запрос, который мы хотим выполнить. Мы разделили его на две строки
	// для удобства чтения (поэтому он окружен обратными кавычками
	// вместо обычных двойных кавычек).
	// stmt := `INSERT INTO snippets (title, content, created, expires)
	// VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	stmt := `update useraccount set last_name = ?, first_name = ?, 
	photo = ?, email = ?, 
	phone_number = ? where id = ?`

	// Используем метод Exec() из встроенного пула подключений для выполнения
	// запроса. Первый параметр это сам SQL запрос, за которым следует
	// заголовок заметки, содержимое и срока жизни заметки. Этот
	// метод возвращает объект sql.Result, который содержит некоторые основные
	// данные о том, что произошло после выполнении запроса.
	result, err := m.DB.Exec(stmt, LastName, FirstName, Photo, Email, PhoneNumber, idU)
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

// Get - Метод для возвращения данных заметки по её идентификатору ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	// SQL запрос для получения данных одной записи.
	stmt := `SELECT ID, first_name, last_name, email, date_of_birth, phone_number FROM useraccount WHERE id = ?`
	// (&s.ID, &s.FirstName, &s.LastName, &s.Email, &s.DateOfBirthDay, &s.PhoneNumber)
	// Используем метод QueryRow() для выполнения SQL запроса,
	// передавая ненадежную переменную id в качестве значения для плейсхолдера
	// Возвращается указатель на объект sql.Row, который содержит данные записи.
	row := m.DB.QueryRow(stmt, id)

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
