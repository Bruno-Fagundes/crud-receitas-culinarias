package main

import (
	"log"
	"net/http"
	"os" // Certifique-se que 'os' está importado!

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

	// *** ADICIONE ESTA LINHA AQUI ***
	log.Println("DEBUG: JWT_SECRET carregado no backend:", os.Getenv("JWT_SECRET"))
	// ******************************

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
	api.HandleFunc("/receitas/{id}", receitaHandler.ReadReceitasById).Methods("GET")
	api.HandleFunc("/receitas", receitaHandler.CreateReceitas).Methods("POST")
	api.HandleFunc("/receitas/{id}", receitaHandler.DeleteReceitas).Methods("DELETE")
	api.HandleFunc("/receitas/{id}", receitaHandler.UpdateReceitas).Methods("PUT")

	// Configurações de CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handlerWithCORS := c.Handler(router)
port := os.Getenv("PORT")
if port == "" {
    port = "5555" // valor de fallback local
}

log.Printf("Servidor rodando em http://0.0.0.0:%s", port)
log.Fatal(http.ListenAndServe(":"+port, handlerWithCORS))

}
