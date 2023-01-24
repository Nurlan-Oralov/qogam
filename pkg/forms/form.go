package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailRX Используйте функцию regexp.MustCompile() для анализа шаблона
// и компиляции регулярного выражения для проверки формата адреса электронной почты на работоспособность.
// Это возвращает объект *regexp.Regexp или вызывает панику в случае ошибки.
// Выполнение этого один раз во время выполнения и сохранение скомпилированного
// объекта регулярного выражения в переменной является более производительным,
// чем повторная компиляция шаблона при каждом запросе.
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Form Создает пользовательскую структуру формы, которая анонимно вставляет объект url.Values
// (для хранения данных формы) и поле ошибок для хранения любых ошибок проверки
// для данных формы.
type Form struct {
	url.Values
	Errors errors
}

// New определяет новую функцию для инициализации структуры пользовательской формы.
// Обратите внимание, что это принимает данные формы в качестве параметра?
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required метод для проверки наличия определенных полей в форме
// данные присутствуют, а не пустые.
// Если какие-либо поля не проходят эту проверку, добавит соответствующее сообщение
// к ошибкам формы.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MaxLength Вызываем метод для проверки того, что определенное поле в форме
// содержит максимальное количество символов.
// Если проверка завершится неудачей, добавит соответствующее сообщение
// к ошибкам формы.
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

// MinLength метод для проверки того, что определенное поле в форме содержит
// минимальное количество символов. Если проверка завершится неудачей,
// добавит соответствующее сообщение в форму ошибок.
func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
	}
}

// PermittedValues Реализуем метод для проверки того, что определенное поле в форме
// соответствует одному из набора определенных допустимых значений. Если проверка завершится неудачей,
// то добавит соответствующее сообщение в форму ошибок.
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// MatchesPattern метод проверки того,
// что определенное поле в форме соответствует регулярному выражению.
// Если проверка завершится неудачей, добавит соответствующее сообщение в форму ошибок.
func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

// Valid Реализуем допустимый метод, который возвращает значение true, если ошибок нет.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
