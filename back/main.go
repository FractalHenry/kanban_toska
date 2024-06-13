package main

import (
	"backend/handlers"
	"backend/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	router := mux.NewRouter()

	// Формируем строку подключения к PostgreSQL
	dsn := "host=localhost port=5433 user=postgres dbname=postgres password=web sslmode=disable"
	// Подключаемся к базе данных SQLite
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Инициализируем базу данных
	handlers.InitHandlers(db)

	// Определяем маршруты
	router.HandleFunc("/register", handlers.Register).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")

	// User
	router.Handle("/user", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserLogin))).Methods("GET")
	router.Handle("/user/{login}", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserInfo))).Methods("GET")
	router.Handle("/description", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateUserDescription))).Methods("PUT")

	// Board
	router.Handle("/boards", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserBoards))).Methods("GET")
	router.Handle("/board", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateBoardHandler))).Methods("POST")
	router.Handle("/board/{boardId}", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetBoardDetailsHandler))).Methods("GET")

	// Card
	router.Handle("/board/{boardId}/card", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateCardHandler))).Methods("POST")
	router.Handle("/removeCard/{cardID}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteCardHandler))).Methods("DELETE")

	router.Handle("/card/{cardID}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateCardHandler))).Methods("PUT")

	// Task
	router.Handle("/board/{boardId}/{cardId}/task", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateTaskHandler))).Methods("POST")

	router.Handle("/task/{taskId}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateTaskHandler))).Methods("PUT")
	router.Handle("/removeTask/{TaskID}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteTaskHandler))).Methods("DELETE")

	// Infoblock
	router.Handle("/board/{boardId}/Infoblock", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateInfoblockHandler))).Methods("POST")
	router.Handle("/removeInfoBlock/{InfoBlockID}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteInfoblockHandler))).Methods("DELETE")
	router.Handle("/updateInfoBlock/{InfoBlockID}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateInfoblockHandler))).Methods("PUT")

	// Space
	router.Handle("/spaces", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserSpaces))).Methods("GET")

	router.Handle("/space/{spaceId}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateSpaceHandler))).Methods("PUT")
	router.Handle("/space/{spaceId}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteSpaceHandler))).Methods("DELETE")
	router.Handle("/space", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateSpaceHandler))).Methods("POST")

	// Protected
	router.Handle("/protected/{name}", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProtectedEndpointWithLogin))).Methods("GET")

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
