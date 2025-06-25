import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import api from '../services/api';

const ReceitaForm = () => {
  const navigate = useNavigate();
  const { id } = useParams(); // Captura o ID da URL se existir
  const [formData, setFormData] = useState({
    nome: '',
    descricao: '',
    ingredientes: '',
    instrucoes: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (id) {
      const fetchReceita = async () => {
        try {
          setLoading(true);
          const receita = await api.getReceitaById(id);
          setFormData({
            nome: receita.nome,
            descricao: receita.descricao || '',
            ingredientes: receita.ingredientes || '',
            instrucoes: receita.instrucoes || ''
          });
        } catch (err) {
          console.error('Erro ao buscar receita:', err);
          setError('Receita não encontrada!');
          navigate('/');
        } finally {
          setLoading(false);
        }
      };
      fetchReceita();
    }
  }, [id, navigate]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.nome || !formData.descricao || !formData.ingredientes || !formData.instrucoes) {
      setError('Por favor, preencha todos os campos obrigatórios.');
      return;
    }
    
    try {
      setLoading(true);
      
      if (id) {
        // Modo edição: PUT
        await api.updateReceita(id, formData);
      } else {
        // Modo criação: POST
        await api.createReceita(formData);
      }
      
      navigate('/');
    } catch (err) {
      console.error('Erro ao salvar receita:', err);
      setError(`Erro ao ${id ? 'atualizar' : 'adicionar'} receita. Tente novamente.`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto p-6 bg-white rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-6">
        {id ? 'Editar Receita' : 'Adicionar Nova Receita'}
      </h2>
      
      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}
      
      <form onSubmit={handleSubmit}>
        {/* Campos mantidos iguais */}
        <div className="mb-4">
          <label htmlFor="nome" className="block text-gray-700 font-semibold mb-2">
            Nome *
          </label>
          <input
            type="text"
            id="nome"
            name="nome"
            value={formData.nome}
            onChange={e => setFormData({ ...formData, nome: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
        
        <div className="mb-4">
          <label htmlFor="descricao" className="block text-gray-700 font-semibold mb-2">
            Descrição *
          </label>
          <input
            type="text"
            id="descricao"
            name="descricao"
            value={formData.descricao}
            onChange={e => setFormData({ ...formData, descricao: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
        
        <div className="mb-4">
          <label htmlFor="ingredientes" className="block text-gray-700 font-semibold mb-2">
            Ingredientes *
          </label>
          <input
            type="text"
            id="ingredientes"
            name="ingredientes"
            value={formData.ingredientes}
            onChange={e => setFormData({ ...formData, ingredientes: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
        
        <div className="mb-4">
          <label htmlFor="instrucoes" className="block text-gray-700 font-semibold mb-2">
            Instruções *
          </label>
          <textarea
            id="instrucoes"
            name="instrucoes"
            value={formData.instrucoes}
            onChange={e => setFormData({ ...formData, instrucoes: e.target.value })}
            rows="4"
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          ></textarea>
        </div>
        
        <div className="flex justify-between mt-6">
          <button
            type="button"
            onClick={() => navigate('/')}
            className="px-4 py-2 bg-gray-300 text-gray-800 rounded hover:bg-gray-400"
          >
            Cancelar
          </button>
          <button
            type="submit"
            disabled={loading}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-blue-300"
          >
            {loading 
              ? (id ? 'Atualizando...' : 'Salvando...') 
              : (id ? 'Atualizar Receita' : 'Salvar Receita')}
          </button>
        </div>
      </form>
    </div>
  );
};

export default ReceitaForm;