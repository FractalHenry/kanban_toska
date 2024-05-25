package main

import (
	"backend/handlers"
	"backend/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // Импортируем PostgreSQL драйвер
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()

	// Формируем строку подключения к PostgreSQL
	dsn := "host=localhost port=5433 user=postgres dbname=postgres password=web sslmode=disable"
	// Подключаемся к базе данных SQLite
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	// Инициализируем базу данных
	handlers.InitDatabase(db)

	// Определяем маршруты
	router.HandleFunc("/register", handlers.Register).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	//router.Handle("/protected", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProtectedEndpoint))).Methods("GET")

	router.Handle("/protected/{login}", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProtectedEndpointWithLogin))).Methods("POST")

	// Настройка CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Разрешенные источники
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		Debug:            true,
	})

	// Оборачиваем маршруты в обработчик CORS
	handler := c.Handler(router)

	// Запускаем сервер на порту 8000
	log.Fatal(http.ListenAndServe(":8000", handler))

}
