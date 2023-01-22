package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// Создает пользовательскую структуру формы, которая анонимно вставляет объект url.Values
// (для хранения данных формы) и поле ошибок для хранения любых ошибок проверки
// для данных формы.
type Form struct {
	url.Values
	Errors errors
}

// Определяет новую функцию для инициализации структуры пользовательской формы. Обратите внимание, что
// это принимает данные формы в качестве параметра?
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Реализуем требуемый метод для проверки наличия определенных полей в форме
// данные присутствуют, а не пустые.
//Если какие-либо поля не проходят эту проверку, добавьте соответствующее сообщение
//к ошибкам формы.

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Вызываем  метод MaxLength для проверки того, что определенное поле в форме
// содержит максимальное количество символов.
// Если проверка завершится неудачей, добавьте соответствующее сообщение
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

// Реализуем метод PermittedValues для проверки того, что определенное поле в форме
// соответствует одному из набора определенных допустимых значений. Если проверка завершится неудачей
// затем добавьте соответствующее сообщение в форму ошибок.
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

// Реализуйте допустимый метод, который возвращает значение true, если ошибок нет.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
