package database

import (
	"context"
	"database/sql"
	"log"

	"golang-rest-api-websockets/models"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

// Constructor funcion of the repository
func NewPostgresrepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users (email, password) VALUES ($1, $2)", user.Password, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = $1", id)

	//defer behind a function, tells us that this is a function
	// that will execute after the whole function (GetUserById) ends
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}() // <--- this () means that we are executing the anonymous function (func(){})

	user := models.User{}
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email); err == nil {
			return &user, nil
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}
