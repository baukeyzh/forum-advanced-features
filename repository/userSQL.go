package repository

import (
	"database/sql"
	"fmt"
	"time"

	"forum/models"

	_ "github.com/mattn/go-sqlite3"
)

type usersSQL struct {
	db *sql.DB
}

// NewAuthSQL create new database struct.
func NewAuthSQL(db *sql.DB) *usersSQL {
	return &usersSQL{db: db}
}

// CreateUser in users table | INSERT
func (r *usersSQL) CreateUser(User models.User) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s ( email, username, password_hash, expire_at) values (?,?,?,?) RETURNING id`, usersTable)

	row := r.db.QueryRow(query, User.Email, User.UserName, User.PassHash, time.Now())
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// GetUser by email from users table | SELECT
func (r *usersSQL) GetUser(Email string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE email=?", usersTable)
	var token sql.NullString
	err := r.db.QueryRow(query, Email).Scan(
		&user.Id,
		&user.Email,
		&user.UserName,
		&user.PassHash,
		&token,
		&user.ExpireAt,
	)
	if err == sql.ErrNoRows {
		return user, err
	}
	user.Token = token.String
	if err != nil {
		return user, fmt.Errorf("can't get user: %w", err)
	}

	return user, nil
}

// AddToken in users table | UPDATE
func (r *usersSQL) AddToken(User models.User) error {
	query := fmt.Sprintf(`UPDATE %s SET token = ?, expire_at = ?  WHERE id = ?`, usersTable)

	if _, err := r.db.Exec(query, User.Token, User.ExpireAt, User.Id); err != nil {
		return fmt.Errorf("can't add token: %w", err)
	}

	return nil
}

// DeleteToken in users table | UPDATE
func (r *usersSQL) DeleteToken(User models.User) error {
	query := fmt.Sprintf(`UPDATE %s SET token = ?, expire_at = ?  WHERE id = ?`, usersTable)

	if _, err := r.db.Exec(query, nil, time.Now(), User.Id); err != nil {
		return fmt.Errorf("can't delete token: %w", err)
	}

	return nil
}

// GetUserByToken from users table  | SELECT
func (r *usersSQL) GetUserByToken(Token string) (models.User, error) {
	var user models.User
	var token sql.NullString
	query := fmt.Sprintf("SELECT * FROM %s WHERE token=?", usersTable)
	err := r.db.QueryRow(query, Token).Scan(&user.Id, &user.Email, &user.UserName, &user.PassHash, &token, &user.ExpireAt)
	if err == sql.ErrNoRows {
		return user, models.ErrorUnauthorized
	}
	user.Token = token.String

	if err != nil {
		return user, fmt.Errorf("can't get user by token: %w", err)
	}

	return user, nil
}
