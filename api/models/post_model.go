package models

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"html"
	"strings"
	"time"
)

type Post struct {
	ID        uuid.UUID `sql:"primary_key" json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `sql:"index" json:"deleted_at,omitempty"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	AuthorID  uuid.UUID `json:"author_id,omitempty"`
	Author    User      `json:"author,omitempty"`
}

func (p *Post) Prepare() {
	p.ID, _ = uuid.NewUUID()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
}

func (p *Post) Validate() error {
	if p.Title == "" {
		return errors.New("Required Title")
	}
	if p.Content == "" {
		return errors.New("Required Content")
	}
	if p.AuthorID == uuid.Nil {
		return errors.New("Required Author")
	}
	return nil
}

func (p Post) SavePost(db *sql.DB) (*Post, error) {
	p.Prepare()
	stmt, err := db.Prepare("INSERT INTO post VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		return &p, err
	}
	_, err = stmt.Exec(p.ID, p.CreatedAt, p.UpdatedAt, nil, p.Title, p.Content, p.AuthorID)
	if err != nil {
		return &p, err
	}
	defer stmt.Close()
	return &p, err
}

func (p Post) FindAllPost(db *sql.DB) (posts []Post, err error) {
	rows, err := db.Query("SELECT id, created_at, updated_at, title, content, author_id FROM post WHERE deleted_at IS NULL")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.Title, &post.Content, &post.AuthorID)
		if err != nil {
			return posts, err
		}
		var user User
		row := db.QueryRow("SELECT id, created_at, updated_at, full_name, username, password, email FROM users WHERE id=$1", &post.AuthorID)
		err = row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.FullName, &user.Username, &user.Password, &user.Email)
		if err != nil {
			return posts, err
		}
		post.Author = user
		posts = append(posts, post)
	}
	err = rows.Err()
	if err != nil {
		return
	}
	defer rows.Close()
	return
}

func (p Post) FindPostByID(db *sql.DB, uuid uuid.UUID) (*Post, error) {
	err := db.QueryRow("SELECT id, created_at, updated_at, title, content, author_id FROM post WHERE id=$1", uuid).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.Title, &p.Content, &p.AuthorID)
	if err != nil {
		return &p, errors.New("Data Not Found By Id " + uuid.String())
	}
	var user User
	err = db.QueryRow("SELECT id, created_at, updated_at, full_name, username, password, email FROM users WHERE id=$1", &p.AuthorID).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.FullName, &user.Username, &user.Password, &user.Email)
	if err != nil {
		return &p, err
	}
	p.Author = user
	return &p, err
}

func (p Post) UpdatePostById(db *sql.DB, uuid uuid.UUID) (*Post, error) {
	p.UpdatedAt = time.Now().Local()
	stmt, err := db.Prepare("UPDATE post SET updated_at=$1, tittle=$2, content=$3 WHERE id=$4")
	if err != nil {
		return &p, err
	}
	_, err = stmt.Exec(p.UpdatedAt, p.Title, p.Content, p.ID)
	if err != nil {
		return &p, err
	}
	defer stmt.Close()
	return &p, err
}

func (p Post) SoftDeletePostById(db *sql.DB, uuid uuid.UUID) (*Post, error) {
	p.DeletedAt = time.Now().Local()
	stmt, err := db.Prepare("UPDATE post SET deleted_at=$1 WHERE id=$2")
	if err != nil {
		return &p, err
	}
	_, err = stmt.Exec(p.DeletedAt, uuid)
	if err != nil {
		return &p, err
	}
	defer stmt.Close()
	return &p, err
}
