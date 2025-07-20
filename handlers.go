package main

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Auth(c *gin.Context) {
	var req AuthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Auth(&req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) Registration(c *gin.Context) {
	var req AuthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	matched, _ := regexp.MatchString(`[0-9]`, req.Password)

	if len(req.Login) > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login len must be < 20"})
		return
	}

	if strings.ContainsAny(req.Login, " \t\n") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login contains spaces"})
		return
	}

	if len(req.Password) < 8 || !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	user, err := h.service.Registration(&req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, user)
}

func (h *Handler) NewObj(c *gin.Context) {
	var obj ObjReq
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(obj.Header) > 25 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Header len must be < 25"})
		return
	}

	if len(obj.Body) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Body len must be < 500"})
		return
	}

	if obj.Price < 1 || obj.Price > 1000000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "1 < Price > 1000000"})
		return
	}

	if !strings.HasSuffix(obj.Image, ".jpg") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image must be .jpg"})
		return
	}

	withLogin := ObjReqWLogin{
		Header: obj.Header,
		Body:   obj.Body,
		Image:  obj.Image,
		Price:  obj.Price,
		Login:  c.GetString("user_login"),
	}

	export, err := h.service.NewObj(&withLogin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, export)

}

func (h *Handler) GetList(c *gin.Context) {
	login, exists := c.Get("user_login")
	if !exists {
		login = ""
	}

	var filters AdsFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if filters.MinPrice != nil && filters.MaxPrice != nil && *filters.MinPrice > *filters.MaxPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_price cannot be greater than max_price"})
		return
	}

	data, err := h.service.GetItems(&filters, login.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
