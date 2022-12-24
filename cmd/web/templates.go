package main

import (
	"golangify.com/snippetbox/pkg/models"
	"html/template" // новый импорт
	"path/filepath" // новый импорт
)

// Эту структуру не трогаем, она тут из прошлого урока...
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
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

	// Перебираем файл шаблона от каждой страницы.
	for _, page := range pages {
		// Извлечение конечное названия файла (например, 'home.page.tmpl') из полного пути к файлу
		// и присваивание его переменной name.
		name := filepath.Base(page)

		// Обрабатываем итерируемый файл шаблона.
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Используем метод ParseGlob для добавления всех каркасных шаблонов.
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
