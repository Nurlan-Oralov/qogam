package mysql

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"golangify.com/snippetbox/pkg/models"
	"strings"
)

type UserModel struct {
	DB *sql.DB
}

// Insert Мы будем использовать этот метод, чтобы добавить новую запись в таблицу users.
func (m *UserModel) Insert(name, email, password string) error {
	// Создает хэш bcrypt для пароля в виде обычного текста.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created)
VALUES(?, ?, ?, UTC_TIMESTAMP())`
	// Метод Exec(), чтобы вставить данные
	// пользователя и хэшированный пароль в таблицу users.
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// Если это возвращает ошибку, мы используем errors.As() для проверки того,
		// имеет ли ошибка тип *mysql.MySQLError.  Если это произойдет, то ошибка
		// будет присвоена переменной mySQLError. Затем мы можем проверить, относится
		// ли ошибка к нашему ключу users_uc_email или нет, проверив содержимое строки сообщения.
		// Если это произойдет, мы вернем ErrDuplicateEmail error.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
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
