package main

import (
	"log"
	"net/http"

	"github.com/Bruno-Fagundes/crud-receitas-culinarias/config"
	_ "github.com/Bruno-Fagundes/crud-receitas-culinarias/docs"
	"github.com/Bruno-Fagundes/crud-receitas-culinarias/handlers"
	"github.com/Bruno-Fagundes/crud-receitas-culinarias/middleware"
	"github.com/Bruno-Fagundes/crud-receitas-culinarias/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Carrega .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	db := config.SetupDB()
	defer db.Close()
	if _, err := db.Exec(models.CreateTableQuery); err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	receitaHandler := handlers.NewReceitaHandler(db)

	router := mux.NewRouter()
	// Public
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Protegidas
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware)
	api.HandleFunc("/receitas", receitaHandler.ReadReceitas).Methods("GET")
	api.HandleFunc("/receitas", receitaHandler.CreateReceitas).Methods("POST")
	api.HandleFunc("/receitas/{id}", receitaHandler.DeleteReceitas).Methods("DELETE")
	api.HandleFunc("/receitas/{id}", receitaHandler.UpdateReceitas).Methods("PUT")

	// Configurações de CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5174"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handlerWithCORS := c.Handler(router)

	log.Println("Servidor rodando em http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", handlerWithCORS))
}
