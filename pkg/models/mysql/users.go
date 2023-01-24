package mysql

import (
	"database/sql"
	"golangify.com/snippetbox/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

// Insert Мы будем использовать этот метод, чтобы добавить новую запись в таблицу users.
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate Мы будем использовать этот метод, чтобы проверить,
// существует ли пользователь с указанным адресом электронной почты и пароль.
// Это вернет соответствующий идентификатор пользователя, если они это сделают.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Get Мы будем использовать этот метод для получения сведений о конкретном пользователе на основе
// on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
