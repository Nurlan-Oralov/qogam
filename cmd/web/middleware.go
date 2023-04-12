package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/justinas/nosurf"
	"golangify.com/snippetbox/pkg/models"
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

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверит, существует ли в сеансе значение  authenticatedUserID.
		// Если *isn't present*, вызовет следующий обработчик в цепочке в обычном режиме.
		exists := app.session.Exists(r, "authenticatedUserID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}
		// Извлекает сведения о текущем пользователе из базы данных.
		// Если соответствующая запись не найдена или текущий пользователь был деактивирован,
		// удалит (недопустимое) значение идентификатора пользователя, прошедшего проверку подлинности,
		// из их сеанса и вызовет следующий обработчик в цепочке в обычном режиме.
		user, err := app.users.Get(app.session.GetInt(r, "authenticatedUserID"))
		if errors.Is(err, models.ErrNoRecord) || !user.Active {
			app.session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}
		// Otherwise, we know that the request is coming from a active, authenticated,
		// user. We create a new copy of the request, with a true boolean value
		// added to the request context to indicate this, and call the next handler
		// in the chain *using this new copy of the request*.
		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
