package models

import (
	"database/sql"
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"html"
	"strings"
	"time"
)

type User struct {
	ID          uuid.UUID `sql:"primary_key" json:"id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	DeletedAt   time.Time `sql:"index" json:"deleted_at,omitempty"`
	FullName    string    `json:"full_name,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Username    string    `json:"username,omitempty"`
	Password    string    `json:"password,omitempty"`
	Email       string    `json:"email,omitempty"`
	Posts       []Post    `json:"posts,omitempty"`
	Role        Role      `json:"role,omitempty"`
	RoleId      int       `json:"role_id,omitempty"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID, _ = uuid.NewUUID()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.FullName = html.EscapeString(strings.TrimSpace(u.FullName))
	u.PhoneNumber = html.EscapeString(strings.TrimSpace(u.PhoneNumber))
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Password = html.EscapeString(strings.TrimSpace(u.Password))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
}

func Validate(u *User) error {
	if u.FullName == "" {
		return errors.New("Required Full Name")
	}
	if u.Username == "" {
		return errors.New("Required Username")
	}
	if u.Password == "" {
		return errors.New("Required Password")
	}
	if u.Email == "" {
		return errors.New("Required Email")
	}
	err := checkmail.ValidateFormat(u.Email)
	if err != nil {
		return errors.New("Invalid Email")
	}
	return nil
}

func (u *User) ValidateUser(action string) error {
	switch strings.ToLower(action) {
	case "update":
		return Validate(u)
	case "login":
		if u.Email == "" && u.Username == "" {
			return errors.New("Required Email or Username")
		}
		if u.Email != "" {
			err := checkmail.ValidateFormat(u.Email)
			if err != nil {
				return errors.New("Invalid Email")
			}
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
	default:
		return Validate(u)
	}
	return nil
}

func (u User) SaveUser(db *sql.DB) (*User, error) {
	u.Prepare()
	stmt, err := db.Prepare("INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)")
	if err != nil {
		return &u, err
	}
	_, err = stmt.Exec(u.ID, u.CreatedAt, u.UpdatedAt, nil, u.FullName, u.PhoneNumber, u.Username, u.Password, u.Email, u.RoleId)
	if err != nil {
		return &u, err
	}
	defer stmt.Close()
	return &u, err
}

func (u User) FindAllUser(db *sql.DB) (users []User, err error) {
	rows, err := db.Query("SELECT id, created_at, updated_at, full_name, phone_number, username, password, email, role_id FROM users WHERE deleted_at IS NULL")
	if err != nil {
		return
	}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.FullName, &user.PhoneNumber, &user.Username, &user.Password, &user.Email, &user.RoleId)
		if err != nil {
			return users, err
		}
		var role Role
		err = db.QueryRow("SELECT id, created_at, updated_at, name FROM role WHERE id=$1", &user.RoleId).
			Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name)
		if err != nil {
			fmt.Println(err)
		}
		user.Role = role
		rowsPost, err := db.Query("SELECT id, created_at, updated_at, title, content, author_id FROM post WHERE author_id=$1", &user.ID)
		if err != nil {
			fmt.Println(err)
		}
		for rowsPost.Next() {
			var post Post
			err := rowsPost.Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.Title, &post.Content, &post.AuthorID)
			if err != nil {
				fmt.Println(err.Error())
			}
			user.Posts = append(user.Posts, post)
		}
		users = append(users, user)
	}
	defer rows.Close()
	return
}

func (u User) FindUserById(db *sql.DB, uuid uuid.UUID) (*User, error) {
	err := db.QueryRow("SELECT id, created_at, updated_at, full_name, phone_number, username, password, email, role_id FROM users WHERE id=$1", uuid).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &u.FullName, &u.Username, &u.PhoneNumber, &u.Password, &u.Email, &u.RoleId)
	if err != nil {
		fmt.Println(err)
		return &u, err
	}
	rows, err := db.Query("SELECT id, created_at, updated_at, title, content, author_id FROM post WHERE author_id=$1", u.ID)
	if err != nil {
		return &u, err
	}
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.Title, &post.Content, &post.AuthorID)
		if err != nil {
			return &u, err
		}
		u.Posts = append(u.Posts, post)
	}
	var role Role
	err = db.QueryRow("SELECT id, created_at, updated_at, name FROM role WHERE id=$1", &u.RoleId).
		Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name)
	if err != nil {
		fmt.Println(err)
		return &u, err
	}
	u.Role = role
	return &u, err
}

func (u User) UpdateUserById(db *sql.DB, uuid uuid.UUID) (*User, error) {
	u.UpdatedAt = time.Now()
	stmt, err := db.Prepare("UPDATE users SET updated_at=$1, full_name=$2, phone_number=$3, username=$4, password=$5, email=$6 WHERE id=$7")
	if err != nil {
		return &u, err
	}
	_, err = stmt.Exec(u.UpdatedAt, u.FullName, u.PhoneNumber, u.Username, u.Password, u.Email, uuid)
	if err != nil {
		return &u, err
	}
	defer stmt.Close()
	return &u, err
}

func (u User) SoftDeleteUserById(db *sql.DB, uuid uuid.UUID) (*User, error) {
	u.DeletedAt = time.Now()
	stmt, err := db.Prepare("UPDATE users SET deleted_at=$1 WHERE id=$2")
	if err != nil {
		return &u, err
	}
	_, err = stmt.Exec(u.DeletedAt, uuid)
	if err != nil {
		return &u, err
	}
	defer stmt.Close()
	return &u, err
}
