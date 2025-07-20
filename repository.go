package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByLogin(req *AuthReq) (*User, error) {
	var user User

	query := `SELECT id, login, password FROM users WHERE login = $1`

	err := r.db.Get(&user, query, req.Login)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) InsertNewGuy(req *AuthReq) (*User, error) {
	var user User

	query := `
        INSERT INTO users (login, password) 
        VALUES ($1, $2)
        RETURNING id, login, password
    `

	err := r.db.Get(&user, query, req.Login, req.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) InsertObj(obj *ObjReqWLogin) (*ObjExport, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var id int
	if err := tx.Get(&id, `SELECT id FROM users WHERE login = $1`, obj.Login); err != nil {
		return nil, err
	}

	var fullObj ObjExport
	query := `
        INSERT INTO objects (header, body, image, price, user_id, date)
        VALUES ($1, $2, $3, $4, $5, NOW())
        RETURNING *
    `
	if err := tx.Get(&fullObj, query, obj.Header, obj.Body, obj.Image, obj.Price, id); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &fullObj, nil
}

func (r *Repository) GetFilteredItems(filters *AdsFilters, login string) ([]Ad, error) {
	query := `
        SELECT 
            o.id,
            o.header,
            o.body,
            o.image,
            o.price,
            o.date,
            u.login as author_login,
            (u.login = $1) as is_mine
        FROM objects o
        JOIN users u ON o.user_id = u.id
        WHERE 1=1
    `

	args := []interface{}{login}
	paramPos := 2

	if filters.MinPrice != nil {
		query += fmt.Sprintf(" AND o.price >= $%d", paramPos)
		args = append(args, *filters.MinPrice)
		paramPos++
	}
	if filters.MaxPrice != nil {
		query += fmt.Sprintf(" AND o.price <= $%d", paramPos)
		args = append(args, *filters.MaxPrice)
		paramPos++
	}

	switch filters.SortBy {
	case "price_asc":
		query += " ORDER BY o.price ASC"
	case "price_desc":
		query += " ORDER BY o.price DESC"
	case "date_asc":
		query += " ORDER BY o.date ASC"
	default:
		query += " ORDER BY o.date DESC"
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramPos, paramPos+1)
	args = append(args, filters.PerPage, (filters.Page-1)*filters.PerPage)

	var ads []Ad
	err := r.db.Select(&ads, query, args...)

	if err != nil {
		return nil, fmt.Errorf("failed to get ads: %w", err)
	}

	return ads, nil
}
