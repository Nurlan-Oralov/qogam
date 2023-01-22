package forms

// Определяет новый тип ошибок, который мы будем использовать для хранения ошибки проверки
// сообщения для форм. Имя поля формы будет использоваться в качестве ключа в
// this map.
type errors map[string][]string

// Реализовать метод Add() для добавления сообщений об ошибках для данного поля на карту.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Реализовать метод Get() для извлечения первого сообщения об ошибке для данного
// поле с карты.
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
