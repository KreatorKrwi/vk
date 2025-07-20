package main

import (
	"log"
	"test-vk/config"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	connStr :=
		"host=" + cfg.DB.Host +
			" port=" + cfg.DB.Port +
			" user=" + cfg.DB.User +
			" password=" + cfg.DB.Password +
			" dbname=" + cfg.DB.Name +
			" sslmode=" + cfg.DB.SSLMode

	db, err := sqlx.Connect("postgres", connStr)

	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping DB:", err)
	}
	log.Println("Successfully connected to database")

	repo := NewRepo(db)
	service := NewService(repo, cfg.Secret.Secret)
	handler := NewHandler(service)

	r := gin.Default()

	r.Use(JWTAuth())

	r.POST("/login", handler.Auth)
	r.POST("/registration", handler.Registration)
	r.POST("/new", handler.NewObj)
	r.GET("/list", handler.GetList)
	r.Run(":" + cfg.Server.Port)
}
