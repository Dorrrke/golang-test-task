package storage

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/Dorrrke/golang-test-task/internal/domain/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrDataNotFound = errors.New("no data")
)

type Storage struct {
	DB  *pgxpool.Pool
	log *slog.Logger
}

func New(ctx context.Context, db *pgxpool.Pool, log *slog.Logger) (*Storage, error) {
	createTablestr := `CREATE TABLE IF NOT EXISTS users
	(
		uid serial PRIMARY KEY,
		name character(50) NOT NULL,
		age integer NOT NULL,
		occupation character(50) NOT NULL,
		salary numeric(5,2) NOT NULL
	)`

	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, createTablestr)
	if err != nil {
		return nil, err
	}

	return &Storage{
		DB:  db,
		log: log,
	}, tx.Commit(ctx)
}

func (stor *Storage) InsertUser(ctx context.Context, user models.User) (int, error) {
	row := stor.DB.QueryRow(ctx, "insert into users (name, age, occupation, salary) values ($1, $2, $3, $4) RETURNING uid;", user.Name, user.Age, user.Occupation, user.Salary)
	var uID int
	if err := row.Scan(&uID); err != nil {
		return -1, err
	}
	return uID, nil
}

func (stor *Storage) GetUserByID(ctx context.Context, userID int) (models.User, error) {
	row := stor.DB.QueryRow(ctx, "SELECT name, age, occupation, salary FROM users WHERE uid = $1", userID)
	var user models.User
	if err := row.Scan(&user.Name, &user.Age, &user.Occupation, &user.Salary); err != nil {
		stor.log.Debug("Check error", slog.Bool("Boolean", errors.Is(err, pgx.ErrNoRows)))
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{
				Name:       "",
				Age:        0,
				Salary:     0,
				Occupation: "",
			}, ErrUserNotFound
		}
		return models.User{
			Name:       "",
			Age:        0,
			Salary:     0,
			Occupation: "",
		}, err
	}
	user.Name = strings.TrimSpace(user.Name)
	user.Occupation = strings.TrimSpace(user.Occupation)
	return user, nil
}

func (stor *Storage) GetAllUsers(ctx context.Context) ([]models.User, error) {
	rows, err := stor.DB.Query(ctx, "SELECT name, age, occupation, salary FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User

	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.Name, &user.Age, &user.Occupation, &user.Salary)
		if err != nil {
			return nil, err
		}
		user.Name = strings.TrimSpace(user.Name)
		user.Occupation = strings.TrimSpace(user.Occupation)
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrDataNotFound
	}
	return users, nil
}
