package models

import "github.com/google/uuid"

type Receita struct {
	ID           uuid.UUID `json:"id"`
	Nome         string    `json:"nome"`
	Descricao    string    `json:"descricao"`
	Ingredientes []string  `json:"ingredientes"`
	Instrucoes   string    `json:"instrucoes"`
}

// Migration
const (
	TableName = "receitas"

	CreateTableQuery = `CREATE TABLE IF NOT EXISTS receitas (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		nome TEXT NOT NULL,
		descricao TEXT NOT NULL,
		ingredientes TEXT[] NOT NULL,
		instrucoes TEXT NOT NULL
	)`
)
