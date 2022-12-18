package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/varjangn/urlsweetner/models"
)

type SQLiteRepository struct {
	db *sql.DB
}

var (
	DbRepo *SQLiteRepository
)

func NewSQLiteRepository(dbFileName string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return nil, err
	}
	return &SQLiteRepository{
		db: db,
	}, nil
}

func (r *SQLiteRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
		uuid TEXT NOT NULL,
        firstname TEXT,
		lastname TEXT
    );`
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) AddUserToDB(u *models.User) error {
	res, err := r.db.Exec("INSERT INTO users(email, password, uuid, firstname, lastname) values(?,?,?,?,?)", u.Email, u.Password, u.UUID, u.FirstName, u.LastName)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.Id = id
	return nil
}

func (r *SQLiteRepository) GetUser(email string) (*models.User, error) {
	row := r.db.QueryRow(fmt.Sprintf("SELECT * from users WHERE email='%s'", email))
	var user models.User
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.UUID, &user.FirstName, &user.LastName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
