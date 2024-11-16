package repository

import (
	"database/sql"
	"time"

	"github.com/anoying-kid/go-apps/blogAPI/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

// NewUserRepository returns a new instance of UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	now := time.Now()
	return r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.Password,
		now,
		now,
	).Scan(&user.ID)
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
    user := &models.User{}
    query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = $1`
    err := r.db.QueryRow(query, email).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.Password,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return user, err
}

func (r *UserRepository) UpdatePassword(userID int64, hashedPassword string) error {
    query := `
        UPDATE users 
        SET password = $1, updated_at = $2 
        WHERE id = $3`

    result, err := r.db.Exec(query, hashedPassword, time.Now(), userID)
    if err != nil {
        return err
    }

    // Check if any row was affected
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }

    return nil
}