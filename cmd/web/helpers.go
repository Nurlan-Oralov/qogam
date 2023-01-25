package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// Помощник serverError записывает сообщение об ошибке в errorLog и
// затем отправляет пользователю ответ 500 "Внутренняя ошибка сервера".
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Помощник clientError отправляет определенный код состояния и соответствующее описание
// пользователю. Мы будем использовать это в следующий уроках, чтобы отправлять ответы вроде 400 "Bad
// Request", когда есть проблема с пользовательским запросом.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Мы также реализуем помощник notFound. Это просто
// удобная оболочка вокруг clientError, которая отправляет пользователю ответ "404 Страница не найдена".
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Извлекаем соответствующий набор шаблонов из кэша в зависимости от названия страницы
	// (например, 'home.page.tmpl'). Если в кэше нет записи запрашиваемого шаблона, то
	// вызывается вспомогательный метод serverError(), который мы создали ранее.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("Шаблон %s не существует!", name))
		return
	}

	// Инициализировать новый буфер.
	buf := new(bytes.Buffer)
	// Запишем шаблон в буфер, а не прямо в
	// http.ResponseWriter. Если произошла ошибка, вызовите наш помощник serverError, а затем
	// return.
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Запишет содержимое буфера в http.ResponseWriter. Опять же, это
	// - это еще один случай, когда мы передаем наш http.ResponseWriter функции, которая
	// принимает ввод-вывод.Писатель.
	buf.WriteTo(w)
}

// Создает помощник по добавлению данных по умолчанию. Это принимает указатель на TemplateData
// struct, добавляет текущий год в поле currentYear, а затем возвращает
// указатель. Опять же, мы не используем *http. Запрашивает параметр в
// момент, но мы сделаем это позже.
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash")
	// Add the authentication status to the template data.
	td.IsAuthenticated = app.isAuthenticated(r)
	return td
}

// Возвращает значение true, если текущий запрос от аутентифицированного пользователя,
// в противном случае возвращает значение false.
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.session.Exists(r, "authenticatedUserID")
}
