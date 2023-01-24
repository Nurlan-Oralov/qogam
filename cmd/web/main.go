package main

import (
	"database/sql" // Новый импорт
	"flag"
	"github.com/golangcollege/sessions"
	"golangify.com/snippetbox/pkg/models/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // Новый импорт
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "Сетевой адрес веб-сервера")
	// Определение нового флага из командной строки для настройки MySQL подключения.
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "Название MySQL источника данных")
	flag.Parse()

	// Определяем новый флаг командной строки для секрета сеанса (случайный ключ, который
	// будет использоваться для шифрования и аутентификации сеансовых файлов cookie).
	// Это должно быть 32 байт длиной.
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Используем sessions.New() функция для инициализации нового диспетчера сеансов,
	// передавая секретный ключ в качестве параметра.
	//Затем мы настраиваем его так, чтобы сеансы всегда истекали через 12 часов.
	//Note: The sessions.New() function returns a Session struct which holds the
	//configuration settings for the session. In the code above we’ve set the Lifetime field
	//of this struct so that sessions expire after 12 hours, but there’s a range of other fields
	//that you can and should configure depending on your application’s needs.
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Запуск сервера на https://localhost%s/", *addr)
	// Используем метод ListenAndServeTLS() для запуска HTTPS-сервера. Мы
	// передаем пути к tls-сертификату и соответствующему секретному ключу в качестве
	// двух параметров.
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// Функция openDB() обертывает sql.Open() и возвращает пул соединений sql.DB
// для заданной строки подключения (DSN).
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
