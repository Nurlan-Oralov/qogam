package main

import (
	"golangify.com/snippetbox/pkg/forms"
	"golangify.com/snippetbox/pkg/models"
	"html/template" // новый импорт
	"path/filepath" // новый импорт
	"time"
)

// Update the templateData fields, removing the individual FormData and
// FormErrors fields and replacing them with a single Form field.
type templateData struct {
	CurrentYear     int
	Flash           string
	Form            *forms.Form
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	IsAuthenticated bool
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Инициализируем новую карту, которая будет хранить кэш.
	cache := map[string]*template.Template{}

	// Используем функцию filepath.Glob, чтобы получить срез всех файловых путей с
	// расширением '.page.tmpl'. По сути, мы получим список всех файлов шаблонов для страниц
	// нашего веб-приложения.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Перебор файла шаблона от каждой страницы.
	for _, page := range pages {
		// Извлечение конечное названия файла (например, 'home.page.tmpl') из полного пути к файлу
		// и присваивание его переменной name.
		name := filepath.Base(page)

		// Обрабатываем итерируемый файл шаблона.
		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method. This means we have to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// метод ParseGlob для добавления всех каркасных шаблонов.
		// В нашем случае это только файл base.layout.tmpl (основная структура шаблона).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Используем метод ParseGlob для добавления всех вспомогательных шаблонов.
		// В нашем случае это footer.partial.tmpl "подвал" нашего шаблона.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Добавляем полученный набор шаблонов в кэш, используя название страницы
		// (например, home.page.tmpl) в качестве ключа для нашей карты.
		cache[name] = ts
	}

	// Возвращаем полученную карту.
	return cache, nil
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 13:04")
}
