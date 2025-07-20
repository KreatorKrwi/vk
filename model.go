package main

import "time"

type AuthReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	Id       int    `db:"id" json:"id"`
	Login    string `db:"login" json:"login"`
	Password string `db:"password" json:"password"`
}

type ObjReq struct {
	Header string `json:"header"`
	Body   string `json:"body"`
	Image  string `json:"image"`
	Price  int    `json:"price"`
}

type ObjReqWLogin struct {
	Header string
	Body   string
	Image  string
	Price  int
	Login  string
}

type ObjExport struct {
	Id      int       `json:"id" db:"id"`
	Header  string    `json:"header" db:"header"`
	Body    string    `json:"body" db:"body"`
	Image   string    `json:"image" db:"image"`
	Price   int       `json:"price" db:"price"`
	User_id int       `json:"user_id" db:"user_id"`
	Date    time.Time `json:"date" db:"date"`
}

type AdsFilters struct {
	Page     int    `form:"page" binding:"min=1"`
	PerPage  int    `form:"per_page" binding:"min=1,max=100"`
	SortBy   string `form:"sort_by,omitempty" binding:"omitempty,oneof=date_desc date_asc price_desc price_asc"`
	MinPrice *int   `form:"min_price" binding:"omitempty,min=0"`
	MaxPrice *int   `form:"max_price" binding:"omitempty,min=1"`
}

type Ad struct {
	ID     int       `json:"id" db:"id"`
	Header string    `json:"header" db:"header"`
	Body   string    `json:"body" db:"body"`
	Image  string    `json:"image" db:"image"`
	Price  int       `json:"price" db:"price"`
	Date   time.Time `json:"date" db:"date"`
	Author string    `json:"author" db:"author_login"`
	IsMine bool      `json:"is_mine" db:"is_mine"`
}
