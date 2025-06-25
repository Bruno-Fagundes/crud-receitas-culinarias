package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
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

// ReadReceitaByID godoc
// @Summary Busca uma receita por ID
// @Description Retorna uma única receita com base no ID fornecido
// @Tags receitas
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da receita (UUID)"
// @Success 200 {object} models.Receita
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/receitas/{id} [get]
func (receitaHandler *ReceitaHandler) ReadReceitasById(w http.ResponseWriter, r *http.Request) {
	log.Println("ReadReceitaByID: Recebendo requisição para detalhes de receita.")
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Valida se o ID é um UUID válido
	_, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("ReadReceitaByID: ID inválido '%s': %v\n", idStr, err)
		http.Error(w, "ID da receita inválido", http.StatusBadRequest)
		return
	}

	var receita models.Receita
	var ingredientesArray pq.StringArray // Para escanear o array do PostgreSQL

	// Consulta a receita pelo ID
	query := `SELECT id, nome, descricao, ingredientes, instrucoes FROM receitas WHERE id = $1`
	err = receitaHandler.DBConnection.QueryRow(query, idStr).Scan(&receita.ID, &receita.Nome, &receita.Descricao, &ingredientesArray, &receita.Instrucoes)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("ReadReceitaByID: Receita não encontrada para o ID %s.\n", idStr)
			http.Error(w, "Receita não encontrada", http.StatusNotFound)
		} else {
			log.Printf("ReadReceitaByID: Erro ao buscar receita por ID: %v\n", err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		}
		return
	}

	receita.Ingredientes = []string(ingredientesArray) // Converte para []string para o modelo

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receita)
	log.Printf("ReadReceitaByID: Receita '%s' carregada com sucesso.\n", receita.Nome)
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
	log.Println("DeleteReceitas: Recebendo requisição para deletar receita.") // Log de início da função

	vars := mux.Vars(r)
	idStr := vars["id"] // Pega o ID da URL

	// 1. Validação do ID: Garante que é um UUID válido
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("DeleteReceitas: Erro de ID inválido '%s': %v\n", idStr, err)
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// 2. Executa a deleção no banco de dados
	result, err := receitaHandler.DBConnection.Exec("DELETE FROM receitas WHERE id = $1", id)
	if err != nil {
		log.Printf("DeleteReceitas: Erro ao executar DELETE no banco para ID %s: %v\n", idStr, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Verifica se alguma linha foi afetada (se a receita existia)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("DeleteReceitas: Erro ao verificar RowsAffected para ID %s: %v\n", idStr, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Printf("DeleteReceitas: Receita com ID %s não encontrada para exclusão.\n", idStr)
		http.Error(w, "Receita não encontrada", http.StatusNotFound) // 404 Not Found se não encontrou
		return
	}

	// 4. Se chegou até aqui, a exclusão foi bem-sucedida
	w.WriteHeader(http.StatusNoContent)                                            // <-- ESTA LINHA É CRUCIAL! Envia o status 204 No Content
	log.Printf("DeleteReceitas: Receita com ID %s deletada com sucesso.\n", idStr) // Log de sucesso
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
