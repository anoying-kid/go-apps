package repository

import (
	"database/sql"
	"time"

	"github.com/anoying-kid/go-apps/blogAPI/internal/models"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *models.Post) error {
	query := `
		INSERT INTO posts (title, body, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	now := time.Now()
	return r.db.QueryRow(
		query,
		post.Title,
		post.Body,
		post.AuthorID,
		now,
		now,
	).Scan(&post.ID)
}

func (r *PostRepository) GetByID(id int64) (*models.Post, error) {
    post := &models.Post{}
    query := `
        SELECT p.id, p.title, p.body, p.author_id, p.created_at, p.updated_at,
               u.username, u.email
        FROM posts p
        JOIN users u ON p.author_id = u.id
        WHERE p.id = $1`
    
    var author models.User
    err := r.db.QueryRow(query, id).Scan(
		&post.ID,
        &post.Title,
        &post.Body,
		&post.AuthorID,
        &post.CreatedAt,
        &post.UpdatedAt,
        &author.Username,
        &author.Email,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    author.ID = post.AuthorID
    post.Author = &author
    return post, nil
}

func (r *PostRepository) List(limit, offset int) ([]*models.Post, error) {
    query := `
        SELECT p.id, p.title, p.body, p.author_id, p.created_at, p.updated_at,
               u.username, u.email
        FROM posts p
        JOIN users u ON p.author_id = u.id
        ORDER BY p.created_at DESC
        LIMIT $1 OFFSET $2`
    
    rows, err := r.db.Query(query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []*models.Post
    for rows.Next() {
        post := &models.Post{}
        author := &models.User{}
        
        err := rows.Scan(
            &post.ID,
            &post.Title,
            &post.Body,
            &post.AuthorID,
            &post.CreatedAt,
            &post.UpdatedAt,
            &author.Username,
            &author.Email,
        )
        if err != nil {
            return nil, err
        }
        
        author.ID = post.AuthorID
        post.Author = author
        posts = append(posts, post)
    }
    
    return posts, nil
}

func (r *PostRepository) Update(post *models.Post) error {
    query := `
        UPDATE posts 
        SET title = $1, body = $2, updated_at = $3
        WHERE id = $4 AND author_id = $5`
    
    result, err := r.db.Exec(
        query,
        post.Title,
        post.Body,
        time.Now(),
        post.ID,
        post.AuthorID,
    )
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }
    
    return nil
}