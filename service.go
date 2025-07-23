package main

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IService interface {
	Auth(req *AuthReq) (string, error)
	Registration(req *AuthReq) (*User, error)
	NewObj(obj *ObjReqWLogin) (*ObjExport, error)
	GetItems(filters *AdsFilters, login string) ([]Ad, error)
}

type Service struct {
	repo   *Repository
	secret string
}

func NewService(repo *Repository, secret string) *Service {
	return &Service{repo: repo, secret: secret}
}

func (s *Service) Auth(req *AuthReq) (string, error) {
	user, err := s.repo.GetByLogin(req)
	if err != nil {
		return "", err
	}

	if user.Password != req.Password {
		return "", errors.New("invalid credintials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_login": user.Login,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(s.secret))
}

func (s *Service) Registration(req *AuthReq) (*User, error) {
	existing, _ := s.repo.GetByLogin(&AuthReq{Login: req.Login})
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	user, err := s.repo.InsertNewGuy(req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) NewObj(obj *ObjReqWLogin) (*ObjExport, error) {
	newobj, err := s.repo.InsertObj(obj)
	if err != nil {
		return nil, err
	}
	return newobj, nil
}

func (s *Service) GetItems(filters *AdsFilters, login string) ([]Ad, error) {
	data, err := s.repo.GetFilteredItems(filters, login)
	if err != nil {
		return nil, err
	}

	return data, nil
}
