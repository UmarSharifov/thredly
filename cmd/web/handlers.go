package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"thredly.com/thredly/pkg/models"
)

// isAuthenticated проверяет, авторизован ли пользователь
func (app *application) isAuthenticated(r *http.Request) bool {
	session, err := app.sessionStore.Get(r, "session-name")
	if err != nil {
		app.errorLog.Printf("Ошибка получения сессии: %v\n", err)
		return false
	}

	_, ok := session.Values["userID"]
	return ok
}

// Функция для получения ID текущего пользователя из сессии
func (app *application) getCurrentUserID(r *http.Request) (int, error) {
	session, err := app.sessionStore.Get(r, "session-name")
	if err != nil {
		return 0, err
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		return 0, errors.New("не удалось получить ID пользователя из сессии")
	}

	return userID, nil
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if !app.isAuthenticated(r) {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	userID, err := app.getCurrentUserID(r)
	if err != nil {
		app.errorLog.Printf("Ошибка получения ID пользователя: %v\n", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	s, err := app.users.Get(userID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := &templateData{User: s}

	files := []string{
		"..\\..\\ui\\html\\home.page.tmpl",
		"..\\..\\ui\\html\\base.layout.tmpl",
		"..\\..\\ui\\html\\footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		// Поскольку обработчик home теперь является методом структуры application
		// он может получить доступ к логгерам из структуры.
		// Используем их вместо стандартного логгера от Go.
		app.errorLog.Println(err.Error())
		http.Error(w, "Внутренняя ошибка сервера", 500)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		// Обновляем код для использования логгера-ошибок
		// из структуры application.
		app.errorLog.Println(err.Error())
		http.Error(w, "Внутренняя ошибка сервера", 500)
	}
}

func (app *application) profile(w http.ResponseWriter, r *http.Request) {
	if !app.isAuthenticated(r) {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	userID, err := app.getCurrentUserID(r)
	if err != nil {
		app.errorLog.Printf("Ошибка получения ID пользователя: %v\n", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	userData, err := app.users.Get(userID)
	if err != nil {
		app.errorLog.Printf("Ошибка получения данных пользователя: %v\n", err)
		http.Error(w, "Ошибка получения данных пользователя", http.StatusInternalServerError)
		return
	}
	tredsData, err := app.treds.GetLatest()
	if err != nil {
		app.errorLog.Printf("Ошибка получения последних тредов: %v\n", err)
		http.Error(w, "Ошибка получения последних тредов", http.StatusInternalServerError)
		return
	}
	tredsChildsData, err := app.treds.GetChildTreds()

	data := &templateData{
		User:       userData,
		Treds:      tredsData,
		TredsChild: tredsChildsData,
	}

	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Ваш код для отображения профиля пользователя
	files := []string{
		"..\\..\\ui\\html\\profile.page.tmpl",
		"..\\..\\ui\\html\\base.layout.tmpl",
		"..\\..\\ui\\html\\footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}

func (app *application) auth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		userID, err := app.users.Authenticate(email, password)
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				app.errorLog.Printf("Неверный email или пароль: %v\n", err)
				http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
				return
			} else {
				app.errorLog.Printf("Ошибка аутентификации: %v\n", err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				return
			}
		}

		session, err := app.sessionStore.Get(r, "session-name")
		if err != nil {
			app.errorLog.Printf("Ошибка получения сессии: %v\n", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}

		session.Values["userID"] = userID
		err = session.Save(r, w)
		if err != nil {
			app.errorLog.Printf("Ошибка сохранения сессии: %v\n", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	files := []string{
		"..\\..\\ui\\html\\auth.page.tmpl",
		"..\\..\\ui\\html\\base.auth.layout.tmpl",
		"..\\..\\ui\\html\\footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Printf("Ошибка парсинга шаблонов: %v\n", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Printf("Ошибка выполнения шаблона: %v\n", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := app.sessionStore.Get(r, "session-name")
	delete(session.Values, "userID")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Страница регистрация пользователя
func (app *application) registration(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Получаем логин и пароль из формы POST
		email := r.FormValue("email")
		password1 := r.FormValue("password1")
		password2 := r.FormValue("password2")

		lastname := "Пользователь1"
		firstname := "Пользователь1"
		photo := "img.jpg"
		phoneNumber := "+123451234"

		if password1 == password2 {
			_, err := app.users.Insert(lastname, firstname, photo, email, phoneNumber, email, password1)
			if err != nil {
				log.Println("Ошибка в insert-e:", err)
				return
			}
			// w.Write([]byte("Создание новой заметки..."))
			s, err := app.users.GetUser(email, password1)
			if err != nil {
				if errors.Is(err, models.ErrNoRecord) {
					log.Println("ID не найден:", err)
					http.NotFound(w, r)
				} else {
					log.Println("Иная ошибка", err)
					http.NotFound(w, r)
				}
				return
			}
			http.Redirect(w, r, fmt.Sprintf("/profile/detail?id=%d", s.ID), http.StatusSeeOther)
		}
	}

	s, err := app.snippets.Latest()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := &templateData{Snippets: s}

	files := []string{
		"..\\..\\ui\\html\\registration.page.tmpl",
		"..\\..\\ui\\html\\base.auth.layout.tmpl",
		"..\\..\\ui\\html\\footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		// Поскольку обработчик home теперь является методом структуры application
		// он может получить доступ к логгерам из структуры.
		// Используем их вместо стандартного логгера от Go.
		app.errorLog.Println(err.Error())
		http.Error(w, "Внутренняя ошибка сервера", 500)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		// Обновляем код для использования логгера-ошибок
		// из структуры application.
		app.errorLog.Println(err.Error())
		http.Error(w, "Внутренняя ошибка сервера", 500)
	}
}

// Подробная информация по пользователю
func (app *application) profileDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Получаем логин и пароль из формы POST
		id := r.FormValue("id")
		lastName := r.FormValue("last_name")
		firstName := r.FormValue("first_name")
		email := r.FormValue("email")
		phoneNumber := r.FormValue("phone_number")

		_, err := app.users.Update(id, lastName, firstName, "", email, phoneNumber)
		if err != nil {
			log.Println("Ошибка в update-e:", err)
			return
		}
		// http.Redirect(w, r, "/profile", http.StatusSeeOther)
		http.Redirect(w, r, fmt.Sprintf("/profile?id=%s", id), http.StatusSeeOther)

	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		log.Println("ID не подходит или не найден:", err)
		http.NotFound(w, r)
		return
	}

	s, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			log.Println("ID не найден:", err)
			http.NotFound(w, r)
		} else {
			log.Println("Иная ошибка", err)
			http.NotFound(w, r)
		}
		return
	}
	// Создаем экземпляр структуры templateData, содержащей данные заметки.
	data := &templateData{User: s}

	files := []string{
		"..\\..\\ui\\html\\detail.profile.page.tmpl",
		"..\\..\\ui\\html\\base.layout.tmpl",
		"..\\..\\ui\\html\\footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println("Шаблон не загрузился", err)
		http.NotFound(w, r)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println("Execute не выполнился", err)
		http.NotFound(w, r)
	}

}

// Обработчик для создания новой заметки.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(405)
		fmt.Fprintln(w, "Метод запрещен!")
		return
	}

	title := "История про улитку"
	content := "Улитка выползла из раковины,\nвытянула рожки,\nи опять подобрала их."
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		return
	}
	// w.Write([]byte("Создание новой заметки..."))
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

// Обработчик для создании треда
func (app *application) createThred(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		log.Println("ID не подходит или не найден:", err)
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodPost {
		// Получаем логин и пароль из формы POST
		id := r.FormValue("id")
		content := r.FormValue("tredArea")

		_, err := app.treds.Insert(id, content)
		if err != nil {
			log.Println("Ошибка в insert-e:", err)
			return
		}
		// http.Redirect(w, r, "/profile", http.StatusSeeOther)
		http.Redirect(w, r, fmt.Sprintf("/profile?id=%s", id), http.StatusSeeOther)

	}

	s, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			log.Println("ID не найден:", err)
			http.NotFound(w, r)
		} else {
			log.Println("Иная ошибка", err)
			http.NotFound(w, r)
		}
		return
	}

	// Создаем экземпляр структуры templateData, содержащей данные заметки.
	data := &templateData{User: s}

	files := []string{
		"..\\..\\ui\\html\\create.page.tmpl",
		"..\\..\\ui\\html\\base.layout.tmpl",
		"..\\..\\ui\\html\\footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		// Поскольку обработчик home теперь является методом структуры application
		// он может получить доступ к логгерам из структуры.
		// Используем их вместо стандартного логгера от Go.
		app.errorLog.Println(err.Error())
		http.Error(w, "Внутренняя ошибка сервера", 500)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		// Обновляем код для использования логгера-ошибок
		// из структуры application.
		app.errorLog.Println(err.Error())
		http.Error(w, "Внутренняя ошибка сервера", 500)
	}
}
