package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	openai "github.com/sashabaranov/go-openai"
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

	catStr := r.URL.Query().Get("cat")
	if catStr == "" {
		catStr = "all"
	}

	s, err := app.users.Get(userID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	tredsData, err := app.treds.GetChildTreds(catStr)
	if err != nil {
		app.errorLog.Printf("Ошибка получения тредов: %v\n", err)
		http.Error(w, "Ошибка получения тредов", http.StatusInternalServerError)
		return
	}
	categories, err := app.categories.Latest()
	if err != nil {
		log.Println("Ошибка при получении категорий:", err)
		http.NotFound(w, r)
		return
	}

	menu := `
		<div class="flex space-x-4">
		<!-- Current: "bg-gray-900 text-white", Default: "text-gray-300 hover:bg-gray-700 hover:text-white" -->
		<a href="http://127.0.0.1:4000/" class="text-gray-950 hover:text-gray-950 px-3 py-2 text-xl font-medium">Главная</a>
		<a href="http://127.0.0.1:4000/recs" class="text-gray-300 hover:text-gray-950 rounded-md px-3 py-2 text-xl font-medium">Рекомендовано</a>
		<a href="http://127.0.0.1:4000/events" class="text-gray-300 hover:text-gray-950 rounded-md px-3 py-2 text-xl font-medium">События</a>
		</div>
	`

	data := &templateData{
		User:            s,
		Treds:           tredsData,
		Categories:      categories,
		CurrentCategory: catStr,
		CurrentMenu:     template.HTML(menu),
		ActiveValue:     "300",
	}

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

func (app *application) recs(w http.ResponseWriter, r *http.Request) {
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

	s, err := app.users.Get(userID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	tredsData, err := app.treds.GetRecommendedTreds(userID, "all")
	if err != nil {
		app.errorLog.Printf("Ошибка получения тредов: %v\n", err)
		http.Error(w, "Ошибка получения тредов", http.StatusInternalServerError)
		return
	}
	categories, err := app.categories.Latest()
	if err != nil {
		log.Println("Ошибка при получении категорий:", err)
		http.NotFound(w, r)
		return
	}

	menu := `
	<div class="flex space-x-4">
	<!-- Current: "bg-gray-900 text-white", Default: "text-gray-300 hover:bg-gray-700 hover:text-white" -->
	<a href="http://127.0.0.1:4000/" class="text-gray-300 hover:text-gray-950 px-3 py-2 text-xl font-medium">Главная</a>
	<a href="http://127.0.0.1:4000/recs" class="text-gray-950 hover:text-gray-950 rounded-md px-3 py-2 text-xl font-medium">Рекомендовано</a>
	<a href="http://127.0.0.1:4000/events" class="text-gray-300 hover:text-gray-950 rounded-md px-3 py-2 text-xl font-medium">События</a>
	</div>
`
	data := &templateData{
		User:        s,
		Treds:       tredsData,
		Categories:  categories,
		CurrentMenu: template.HTML(menu),
		ActiveValue: "300",
	}

	files := []string{
		"..\\..\\ui\\html\\recs.page.tmpl",
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

func (app *application) events(w http.ResponseWriter, r *http.Request) {
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

	s, err := app.users.Get(userID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	eventsData, err := app.eventsObject.Latest()
	if err != nil {
		app.errorLog.Printf("Ошибка получения событий: %v\n", err)
		http.Error(w, "Ошибка получения событий", http.StatusInternalServerError)
		return
	}
	categories, err := app.eventsCategoy.Latest()
	if err != nil {
		log.Println("Ошибка при получении категорий:", err)
		http.NotFound(w, r)
		return
	}

	menu := `
	<div class="flex space-x-4">
	<a href="http://127.0.0.1:4000/" class="text-gray-300 hover:text-gray-950 px-3 py-2 text-xl font-medium">Главная</a>
	<a href="http://127.0.0.1:4000/recs" class="text-gray-300 hover:text-gray-950 rounded-md px-3 py-2 text-xl font-medium">Рекомендовано</a>
	<a href="http://127.0.0.1:4000/events" class="text-gray-950 hover:text-gray-950 rounded-md px-3 py-2 text-xl font-medium">События</a>
	</div>
`
	data := &templateData{
		User:             s,
		Events:           eventsData,
		EventsCategories: categories,
		CurrentMenu:      template.HTML(menu),
		ActiveValue:      "300",
	}

	files := []string{
		"..\\..\\ui\\html\\events.page.tmpl",
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

	tredsData, err := app.treds.GetChildTreds("all")
	if err != nil {
		app.errorLog.Printf("Ошибка получения тредов: %v\n", err)
		http.Error(w, "Ошибка получения тредов", http.StatusInternalServerError)
		return
	}

	usersSubs, err := app.users.GetSubscrible(userID)
	if err != nil {
		app.errorLog.Printf("Ошибка получения пользователей: %v\n", err)
		http.Error(w, "Ошибка получения пользователей", http.StatusInternalServerError)
		return
	}

	// подписки
	userSubsTo, err := app.users.GetSubscribleCount(userID)
	if err != nil {
		app.errorLog.Printf("Ошибка получения подписок: %v\n", err)
		http.Error(w, "Ошибка получения подписок", http.StatusInternalServerError)
		return
	}

	// подписчики
	userSubsFrom, err := app.users.GetSubscribleToCount(userID)
	if err != nil {
		app.errorLog.Printf("Ошибка получения подписчиков: %v\n", err)
		http.Error(w, "Ошибка получения подписчиков", http.StatusInternalServerError)
		return
	}

	// Количество публикаций пользователя
	userTredsCount, err := app.treds.GetCount(userID)
	if err != nil {
		app.errorLog.Printf("Ошибка получения Количества тредов: %v\n", err)
		http.Error(w, "Ошибка получения количество тредов", http.StatusInternalServerError)
		return
	}
	menu := `
	<div class="flex space-x-4">
	<!-- Current: "bg-gray-900 text-white", Default: "text-gray-300 hover:bg-gray-700 hover:text-white" -->
	<a href="http://127.0.0.1:4000/" class="text-gray-300 hover:text-gray-950 px-3 py-2 text-xl font-medium">Главная</a>
	<a href="http://127.0.0.1:4000/recs" class="text-gray-300 hover:text-gray-950 rounded-md px-3 py-2 text-xl font-medium">Рекомендовано</a>
	<a href="http://127.0.0.1:4000/events" class="text-gray-300 hover:text-gray-950 rounded-md px-3 py-2 text-xl font-medium">События</a>
	</div>
`

	data := &templateData{
		User:           userData,
		Treds:          tredsData,
		CurrentMenu:    template.HTML(menu),
		ActiveValue:    "950",
		Users:          usersSubs,
		UserSubsTo:     userSubsTo,
		UserSubsFrom:   userSubsFrom,
		UserTredsCount: userTredsCount,
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

		http.Redirect(w, r, "/", http.StatusSeeOther)
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

	data := &templateData{}

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
		// Получаем данные из формы POST
		id := r.FormValue("id")
		lastName := r.FormValue("last_name")
		firstName := r.FormValue("first_name")
		email := r.FormValue("email")
		phoneNumber := r.FormValue("phone_number")

		// Обработка загруженного фото профиля
		var photoFilename string
		file, handler, err := r.FormFile("profile_picture")
		if err != nil && err != http.ErrMissingFile {
			log.Println("Ошибка при загрузке файла:", err)
			http.Error(w, "Ошибка при загрузке файла", http.StatusInternalServerError)
			return
		} else if err == nil {
			defer file.Close()

			// Убедимся, что директория uploads существует
			uploadDir := "..\\..\\ui\\static\\img"
			if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
				err = os.Mkdir(uploadDir, 0755)
				if err != nil {
					log.Println("Ошибка при создании директории:", err)
					http.Error(w, "Ошибка при создании директории", http.StatusInternalServerError)
					return
				}
			}

			// Создаем файл на сервере
			photoFilename = handler.Filename
			f, err := os.OpenFile(fmt.Sprintf("%s/%s", uploadDir, photoFilename), os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				log.Println("Ошибка при создании файла:", err)
				http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
				return
			}
			defer f.Close()

			// Копируем содержимое загруженного файла в новый файл
			_, err = io.Copy(f, file)
			if err != nil {
				log.Println("Ошибка при копировании файла:", err)
				http.Error(w, "Ошибка при копировании файла", http.StatusInternalServerError)
				return
			}
		}

		_, err = app.users.Update(id, lastName, firstName, photoFilename, email, phoneNumber)
		if err != nil {
			log.Println("Ошибка в update-е:", err)
			http.Error(w, "Ошибка при обновлении данных", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/profile?id=%s", id), http.StatusSeeOther)
		return
	}

	userID, err := app.getCurrentUserID(r)
	if err != nil {
		app.errorLog.Printf("Ошибка получения ID пользователя: %v\n", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	s, err := app.users.Get(userID)
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
	// Создаем экземпляр структуры templateData, содержащей данные пользователя.
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

// Обработчик для создании треда
func (app *application) createThred(w http.ResponseWriter, r *http.Request) {

	id, err := app.getCurrentUserID(r)
	if err != nil {
		app.errorLog.Printf("Ошибка получения ID пользователя: %v\n", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	parentThreadId := "0"
	errorMessageCreateTred := ""
	currentTredBackstage := ""
	if r.Method == http.MethodPost {
		parentThreadId = r.FormValue("parentThreadId")

		tags := r.FormValue("tags")

		category := r.FormValue("category_select")

		content := r.FormValue("tredArea")

		client := openai.NewClient("Здесь токен сгенерированный в openAI")
		currentContent := fmt.Sprintf(`Ответь 1 если текст в кавычках содержит маты или 0 если нет "%s"`, content)

		resp, _ := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: currentContent,
					},
				},
			},
		)

		if resp.Choices[0].Message.Content == "0" {
			currentTredId, err := app.treds.Insert(id, parentThreadId, content, category)
			if err != nil {
				log.Println("Ошибка в insert-e:", err)
				return
			}
			err = app.Tag.Insert(currentTredId, tags)
			if err != nil {
				log.Println("Ошибка в insert-e тегов:", err)
				return
			}
			http.Redirect(w, r, fmt.Sprintf("/profile?id=%s", id), http.StatusSeeOther)
		} else {
			log.Println("Иная ошибка", resp.Choices[0].Message.Content)
			errorMessageCreateTred = `Пожалуйста, отредактируйте тред! Ваш текст содержит ненормативную лексику.`
			currentTredBackstage = content
		}

		// http.Redirect(w, r, "/profile", http.StatusSeeOther)
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
	categories, err := app.categories.Latest()
	if err != nil {
		log.Println("Ошибка при получении категорий:", err)
		http.NotFound(w, r)
		return
	}

	menu := `
	<div class="flex space-x-4">
	<!-- Current: "bg-gray-900 text-white", Default: "text-gray-300 hover:bg-gray-700 hover:text-white" -->
	<a href="http://127.0.0.1:4000/" class="text-gray-300 hover:text-gray-950 px-3 py-2 text-xl font-medium">Главная</a>
	<a href="http://127.0.0.1:4000/recs" class="text-gray-300 hover:text-gray-950 rounded-md px-3 py-2 text-xl font-medium">Рекомендовано</a>
	<a href="http://127.0.0.1:4000/events" class="text-gray-300 hover:text-gray-950 rounded-md px-3 py-2 text-xl font-medium">События</a>
	</div>
	`

	userID, err := app.getCurrentUserID(r)
	if err != nil {
		app.errorLog.Printf("Ошибка получения ID пользователя: %v\n", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	usersSubs, err := app.users.GetSubscrible(userID)
	if err != nil {
		app.errorLog.Printf("Ошибка получения пользователей: %v\n", err)
		http.Error(w, "Ошибка получения пользователей", http.StatusInternalServerError)
		return
	}

	// подписки
	userSubsTo, err := app.users.GetSubscribleCount(userID)
	if err != nil {
		app.errorLog.Printf("Ошибка получения подписок: %v\n", err)
		http.Error(w, "Ошибка получения подписок", http.StatusInternalServerError)
		return
	}

	// подписчики
	userSubsFrom, err := app.users.GetSubscribleToCount(userID)
	if err != nil {
		app.errorLog.Printf("Ошибка получения подписчиков: %v\n", err)
		http.Error(w, "Ошибка получения подписчиков", http.StatusInternalServerError)
		return
	}

	// Количество публикаций пользователя
	userTredsCount, err := app.treds.GetCount(userID)
	if err != nil {
		app.errorLog.Printf("Ошибка получения Количества тредов: %v\n", err)
		http.Error(w, "Ошибка получения количество тредов", http.StatusInternalServerError)
		return
	}
	// Создаем экземпляр структуры templateData, содержащей данные заметки.
	data := &templateData{
		User:                   s,
		Categories:             categories,
		CurrentMenu:            template.HTML(menu),
		ActiveValue:            "950",
		ErrorMessageCreateTred: errorMessageCreateTred,
		CurrentTredBackstage:   currentTredBackstage,
		ParentTredId:           parentThreadId,
		Users:                  usersSubs,
		UserSubsTo:             userSubsTo,
		UserSubsFrom:           userSubsFrom,
		UserTredsCount:         userTredsCount,
	}

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

func (app *application) createEvent(w http.ResponseWriter, r *http.Request) {

	id, err := app.getCurrentUserID(r)
	if err != nil {
		app.errorLog.Printf("Ошибка получения ID пользователя: %v\n", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		parentThreadId := r.FormValue("parentThreadId")

		category := r.FormValue("category_select")

		content := r.FormValue("tredArea")

		_, err := app.eventsObject.Insert(id, parentThreadId, content, category)
		if err != nil {
			log.Println("Ошибка в insert-e:", err)
			return
		}
		// http.Redirect(w, r, "/profile", http.StatusSeeOther)
		http.Redirect(w, r, fmt.Sprintf("/events"), http.StatusSeeOther)

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
	categories, err := app.eventsCategoy.Latest()
	if err != nil {
		log.Println("Ошибка при получении категорий:", err)
		http.NotFound(w, r)
		return
	}

	menu := `
	<div class="flex space-x-4">
	<!-- Current: "bg-gray-900 text-white", Default: "text-gray-300 hover:bg-gray-700 hover:text-white" -->
	<a href="http://127.0.0.1:4000/" class="text-gray-300 hover:text-gray-950 px-3 py-2 text-xl font-medium">Главная</a>
	<a href="http://127.0.0.1:4000/recs" class="text-gray-300 hover:text-gray-950 rounded-md px-3 py-2 text-xl font-medium">Рекомендовано</a>
	<a href="http://127.0.0.1:4000/events" class="text-gray-950 hover:text-gray-950 rounded-md px-3 py-2 text-xl font-medium">События</a>
	</div>
	`
	// Создаем экземпляр структуры templateData, содержащей данные заметки.
	data := &templateData{
		User:             s,
		EventsCategories: categories,
		CurrentMenu:      template.HTML(menu),
		ActiveValue:      "300",
	}

	files := []string{
		"..\\..\\ui\\html\\events.create.page.tmpl",
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

func (app *application) adminPage(w http.ResponseWriter, r *http.Request) {

	id, err := app.getCurrentUserID(r)
	if err != nil {
		app.errorLog.Printf("Ошибка получения ID пользователя: %v\n", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
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
	categories, err := app.eventsCategoy.Latest()
	if err != nil {
		log.Println("Ошибка при получении категорий:", err)
		http.NotFound(w, r)
		return
	}

	complaints, err := app.Complaint.GetAll()
	if err != nil {
		log.Println("Ошибка при получении жалоб:", err)
		http.NotFound(w, r)
		return
	}

	menu := `
	<div class="flex space-x-4">
	<!-- Current: "bg-gray-900 text-white", Default: "text-gray-300 hover:bg-gray-700 hover:text-white" -->
	<a href="http://127.0.0.1:4000/" class="text-gray-300 hover:text-gray-950 px-3 py-2 text-xl font-medium">Главная</a>
	<a href="http://127.0.0.1:4000/recs" class="text-gray-300 hover:text-gray-950 rounded-md px-3 py-2 text-xl font-medium">Рекомендовано</a>
	<a href="http://127.0.0.1:4000/events" class="text-gray-300 hover:text-gray-950 rounded-md px-3 py-2 text-xl font-medium">События</a>
	</div>
	`

	// Создаем экземпляр структуры templateData, содержащей данные заметки.
	data := &templateData{
		User:                 s,
		EventsCategories:     categories,
		CurrentMenu:          template.HTML(menu),
		ActiveValue:          "950",
		ComplaintWithDetails: complaints,
	}

	files := []string{
		"..\\..\\ui\\html\\admin.page.tmpl",
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

func (app *application) report(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id, err := app.getCurrentUserID(r)
	if err != nil {
		app.errorLog.Printf("Ошибка получения ID пользователя: %v\n", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
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

	var input struct {
		ID string `json:"id"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorLog.Printf("Ошибка декодирования JSON: %v\n", err)
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	_, err = app.Complaint.Insert(input.ID, s.ID)
	if err != nil {
		log.Println("Ошибка в insert-e:", err)
		return
	}

	app.infoLog.Printf("Получен отчет с ID: %s\n", input.ID)
	w.WriteHeader(http.StatusOK)
}

func (app *application) subscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		ID string `json:"userId"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorLog.Printf("Ошибка декодирования JSON: %v\n", err)
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	app.infoLog.Printf("Получен ID пользователя для подписки: %s\n", input.ID)

	id, err := app.getCurrentUserID(r)
	if err != nil {
		app.errorLog.Printf("Ошибка получения ID пользователя: %v\n", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
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

	_, err = app.Subscribe.Insert(s.ID, input.ID)
	if err != nil {
		log.Println("Ошибка в insert-e:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) moreuser(w http.ResponseWriter, r *http.Request) {
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

	tredsData, err := app.treds.GetChildTreds("all")
	if err != nil {
		app.errorLog.Printf("Ошибка получения тредов: %v\n", err)
		http.Error(w, "Ошибка получения тредов", http.StatusInternalServerError)
		return
	}

	usersSubs, err := app.users.GetSubscribleAll(userID)
	if err != nil {
		app.errorLog.Printf("Ошибка получения пользователей: %v\n", err)
		http.Error(w, "Ошибка получения пользователей", http.StatusInternalServerError)
		return
	}

	// подписки
	userSubsTo, err := app.users.GetSubscribleCount(userID)
	if err != nil {
		app.errorLog.Printf("Ошибка получения подписок: %v\n", err)
		http.Error(w, "Ошибка получения подписок", http.StatusInternalServerError)
		return
	}

	// подписчики
	userSubsFrom, err := app.users.GetSubscribleToCount(userID)
	if err != nil {
		app.errorLog.Printf("Ошибка получения подписчиков: %v\n", err)
		http.Error(w, "Ошибка получения подписчиков", http.StatusInternalServerError)
		return
	}

	// Количество публикаций пользователя
	userTredsCount, err := app.treds.GetCount(userID)
	if err != nil {
		app.errorLog.Printf("Ошибка получения Количества тредов: %v\n", err)
		http.Error(w, "Ошибка получения количество тредов", http.StatusInternalServerError)
		return
	}
	menu := `
	<div class="flex space-x-4">
	<!-- Current: "bg-gray-900 text-white", Default: "text-gray-300 hover:bg-gray-700 hover:text-white" -->
	<a href="http://127.0.0.1:4000/" class="text-gray-300 hover:text-gray-950 px-3 py-2 text-xl font-medium">Главная</a>
	<a href="http://127.0.0.1:4000/recs" class="text-gray-300 hover:text-gray-950 rounded-md px-3 py-2 text-xl font-medium">Рекомендовано</a>
	<a href="http://127.0.0.1:4000/events" class="text-gray-300 hover:text-gray-950 rounded-md px-3 py-2 text-xl font-medium">События</a>
	</div>
`

	data := &templateData{
		User:           userData,
		Treds:          tredsData,
		CurrentMenu:    template.HTML(menu),
		ActiveValue:    "950",
		Users:          usersSubs,
		UserSubsTo:     userSubsTo,
		UserSubsFrom:   userSubsFrom,
		UserTredsCount: userTredsCount,
	}

	// Ваш код для отображения профиля пользователя
	files := []string{
		"..\\..\\ui\\html\\moreuser.page.tmpl",
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
