package models

import (
	"database/sql"
	"github.com/pkg/errors"
	"html"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Role struct {
	ID        int       `sql:"primary_key" json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `sql:"index" json:"deleted_at,omitempty"`
	Name      string    `json:"name,omitempty"`
}

func (r *Role) Prepare() {
	r.ID = int(rand.Uint32())
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	r.Name = html.EscapeString(strings.TrimSpace(r.Name))
}

func (r *Role) Validate() error {
	if r.Name == "" {
		return errors.New("Required Role Name")
	}
	return nil
}

func (r Role) SaveRole(db *sql.DB) (*Role, error) {
	r.Prepare()
	stmt, err := db.Prepare("INSERT INTO role VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return &r, err
	}
	_, err = stmt.Exec(&r.ID, &r.CreatedAt, &r.UpdatedAt, nil, &r.Name)
	if err != nil {
		return &r, err
	}
	defer stmt.Close()
	return &r, err
}

func (r Role) FindAllRole(db *sql.DB) (*[]Role, error) {
	var roles []Role
	rows, err := db.Query("SELECT id, created_at, updated_at, name FROM role WHERE deleted_at IS NULL")
	if err != nil {
		return &roles, err
	}
	for rows.Next() {
		var role Role
		err := rows.Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name)
		if err != nil {
			return &roles, err
		}
		roles = append(roles, role)
	}
	err = rows.Err()
	if err != nil {
		return &roles, err
	}
	defer rows.Close()
	return &roles, err
}

func (r Role) FindRoleById(db *sql.DB, id int) (*Role, error) {
	err := db.QueryRow("SELECT id, created_at, updated_at, name FROM role WHERE id=$1", id).Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt, &r.Name)
	if err != nil {
		return &r, errors.New("Data Not Found By Id" + strconv.Itoa(id))
	}
	return &r, err
}

func (r Role) UpdateRoleById(db *sql.DB, id int) (*Role, error) {
	role, err := r.FindRoleById(db, id)
	if err != nil {
		return role, err
	}
	stmt, err := db.Prepare("UPDATE role SET updated_at=$1, name=$2 WHERE id=$3")
	if err != nil {
		return role, err
	}
	_, err = stmt.Exec(r.UpdatedAt, r.Name, id)
	if err != nil {
		return role, err
	}
	defer stmt.Close()
	return &r, err
}

func (r Role) SoftDeleteRoleById(db *sql.DB, id int) (*Role, error) {
	role, err := r.FindRoleById(db, id)
	if err != nil {
		return role, err
	}
	stmt, err := db.Prepare("UPDATE role SET deleted_at=$1 WHERE id=$2")
	if err != nil {
		return role, err
	}
	_, err = stmt.Exec(&r.DeletedAt, id)
	if err != nil {
		return role, err
	}
	defer stmt.Close()
	return &r, err
}
