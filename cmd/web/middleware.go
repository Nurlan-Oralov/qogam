package main

import (
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Создает отложенную функцию (которая всегда будет выполняться в случае
		// паники, когда Go разматывает стек).
		defer func() {
			// Используется встроенную функцию восстановления, чтобы проверить, не было ли
			// паниковать или нет. Если так оно и было...
			if err := recover(); err != nil {
				// Устанавливаем заголовок "Подключение: закрыть" в ответе.
				w.Header().Set("Connection", "close")
				// Вызываем вспомогательный метод app.serverError, чтобы вернуть 500ь 500
				// Внутренний ответ сервера.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Если пользователь не прошел проверку подлинности,
		// перенаправит его на страницу входа в систему и вернет
		// из middleware chain, чтобы никакие последующие обработчики в цепочке не выполнялись.
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// В противном случае установит заголовок "Cache-Control: no-store",
		// чтобы страницы, требующие аутентификации, не сохранялись
		// в кэше браузера пользователя (или другом промежуточном кэше).
		w.Header().Add("Cache-Control", "no-store")
		// И вызовет следующий обработчик в цепочке.
		next.ServeHTTP(w, r)
	})
}

// Create a NoSurf middleware function which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly flags set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}
