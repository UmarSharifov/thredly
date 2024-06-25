package main

import (
	"database/sql"
	"flag"
	"log"

	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"thredly.com/thredly/pkg/models/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	users         *mysql.UserModel
	treds         *mysql.TredModel
	sessionStore  *sessions.CookieStore
	categories    *mysql.CategoriesMoodel
	eventsObject  *mysql.EventsMoodel
	eventsCategoy *mysql.EventsCategoriesMoodel
	Complaint     *mysql.ComplaintModel
	Subscribe     *mysql.SubscribeModel
	Tag           *mysql.TagsModel
}

func main() {
	addr := flag.String("addr", ":4000", "Сетевой адрес веб-сервера")

	// Определение нового флага из командной строки для настройки MySQL подключения.
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "Название MySQL источника данных")
	secret := flag.String("secret", "your-secret-key", "Ключ для шифрования сессий") // Добавьте ключ для шифрования сессий
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Чтобы функция main() была более компактной, мы поместили код для создания
	// пула соединений в отдельную функцию openDB(). Мы передаем в нее полученный
	// источник данных (DSN) из флага командной строки.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// Мы также откладываем вызов db.Close(), чтобы пул соединений был закрыт
	// до выхода из функции main().
	// Подробнее про defer: https://golangs.org/errors#defer
	defer db.Close()
	sessionStore := sessions.NewCookieStore([]byte(*secret))
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		users:         &mysql.UserModel{DB: db},
		treds:         &mysql.TredModel{DB: db},
		categories:    &mysql.CategoriesMoodel{DB: db},
		eventsObject:  &mysql.EventsMoodel{DB: db},
		eventsCategoy: &mysql.EventsCategoriesMoodel{DB: db},
		Complaint:     &mysql.ComplaintModel{DB: db},
		Subscribe:     &mysql.SubscribeModel{DB: db},
		Tag:           &mysql.TagsModel{DB: db},
		sessionStore:  sessionStore,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Запуск сервера на %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
