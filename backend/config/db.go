package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // _ porque nao esta sendo usado diretamente, mas precisa ser importado para registrar o driver
)

// Retorna uma instancia do banco de dados
// Cria a conexao com o banco de dados
// e retorna o enderco de memoria aonde foi salva a variel do endereco da conexao
func SetupDB() *sql.DB {
	err := godotenv.Load() // Carrega as variaveis de ambiente do arquivo .env
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSchema := os.Getenv("DB_SCHEMA")

	connectionStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSchema)

	dbConnection, err := sql.Open("postgres", connectionStr)

	if err != nil {
		log.Fatal(err)
	}

	err = dbConnection.Ping()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conexão com o banco de dados estabelecida com sucesso!")

	return dbConnection
}
