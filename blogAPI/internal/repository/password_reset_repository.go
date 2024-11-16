package repository

import (
	"database/sql"
	"time"

	"github.com/anoying-kid/go-apps/blogAPI/internal/models"
)

type PasswordResetRepository struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
    return &PasswordResetRepository{db: db}
}

func (r *PasswordResetRepository) Create(reset *models.PasswordResetToken) error {
    query := `
        INSERT INTO password_reset_tokens (user_id, token, expired_at, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

    return r.db.QueryRow(
        query,
        reset.UserID,
        reset.Token,
        reset.ExpiredAt,
        time.Now(),
    ).Scan(&reset.ID)
}

func (r *PasswordResetRepository) GetByToken(token string) (*models.PasswordResetToken, error) {
	reset := &models.PasswordResetToken{}
	query := `
	SELECT id, user_id, token, expired_at, used, created_at
	FROM password_reset_tokens
	WHERE token = $1`
	err := r.db.QueryRow(query, token).Scan(
		&reset.ID,
		&reset.UserID,
		&reset.Token,
		&reset.ExpiredAt,
		&reset.Used,
		&reset.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return reset, err
}

func (r *PasswordResetRepository) MarkAsUsed(id int64) error {
	query := `UPDATE password_reset_tokens SET used = true WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}