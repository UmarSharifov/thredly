package main

import (
	"database/sql" // Новый импорт
	"flag"
	"log"

	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"thredly.com/thredly/pkg/models/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
	users    *mysql.UserModel
	treds    *mysql.TredModel
}

func main() {
	addr := flag.String("addr", ":4000", "Сетевой адрес веб-сервера")
	// Определение нового флага из командной строки для настройки MySQL подключения.
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "Название MySQL источника данных")
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

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{DB: db},
		users:    &mysql.UserModel{DB: db},
		treds:    &mysql.TredModel{DB: db},
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
