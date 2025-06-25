package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Bruno-Fagundes/crud-receitas-culinarias/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type ReceitaHandler struct {
	DBConnection *sql.DB
}

// Construtor de ReceitaHandler
func NewReceitaHandler(dbConnection *sql.DB) *ReceitaHandler {
	return &ReceitaHandler{DBConnection: dbConnection}
}

// ReadReceitas godoc
// @Summary Lista todas as receitas
// @Description Retorna todas as receitas cadastradas no banco de dados
// @Tags receitas
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Receita
// @Failure 500 {object} map[string]string
// @Router /api/receitas [get]
func (receitaHandler *ReceitaHandler) ReadReceitas(w http.ResponseWriter, r *http.Request) {

	rows, err := receitaHandler.DBConnection.Query("SELECT * FROM receitas")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var receitas []models.Receita

	for rows.Next() {
		var receita models.Receita
		err := rows.Scan(&receita.ID, &receita.Nome, &receita.Descricao, pq.Array(&receita.Ingredientes), &receita.Instrucoes)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		receitas = append(receitas, receita)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receitas)
}

// CreateReceitas godoc
// @Summary Cria uma nova receita
// @Description Adiciona uma nova receita com nome, descrição, ingredientes e instruções
// @Tags receitas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param receita body models.Receita true "Dados da nova receita"
// @Success 200 {object} models.Receita
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/receitas [post]
func (receitaHandler *ReceitaHandler) CreateReceitas(w http.ResponseWriter, r *http.Request) {
	var receita models.Receita

	err := json.NewDecoder(r.Body).Decode(&receita)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO receitas (nome, descricao, ingredientes, instrucoes) VALUES ($1, $2, $3, $4) RETURNING id`
	err = receitaHandler.DBConnection.QueryRow(query, receita.Nome, receita.Descricao, pq.Array(receita.Ingredientes), receita.Instrucoes).Scan(&receita.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receita)
}

// DeleteReceitas godoc
// @Summary Deleta uma receita
// @Description Remove uma receita do banco de dados pelo ID
// @Tags receitas
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da receita (UUID)"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/receitas/{id} [delete]
func (receitaHandler *ReceitaHandler) DeleteReceitas(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	result, err := receitaHandler.DBConnection.Exec("DELETE FROM receitas WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Receita não encontrada", http.StatusNotFound)
		return
	}
}

// UpdateReceitas godoc
// @Summary Atualiza uma receita
// @Description Atualiza os dados de uma receita existente
// @Tags receitas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da receita (UUID)"
// @Param receita body models.Receita true "Dados atualizados da receita"
// @Success 200 {object} models.Receita
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/receitas/{id} [put]
func (receitaHandler *ReceitaHandler) UpdateReceitas(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)

	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var receita models.Receita
	err = json.NewDecoder(r.Body).Decode(&receita)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE receitas SET nome = $1, descricao = $2, ingredientes = $3, instrucoes = $4 WHERE id = $5`
	result, err := receitaHandler.DBConnection.Exec(query, receita.Nome, receita.Descricao, pq.Array(receita.Ingredientes), receita.Instrucoes, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Receita não encontrada", http.StatusNotFound)
		return
	}

	receita.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receita)
}
