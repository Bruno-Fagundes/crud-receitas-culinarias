import React, { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import api from '../services/api';
// Verifique se o caminho para DeleteConfirmation está correto.
// Se você moveu DeleteConfirmation para ReceitaList.jsx, pode precisar ajustar ou remover esta linha.
import DeleteConfirmation from './DeleteConfirmation'; 

const ReceitaDetail = () => {
  const { id } = useParams(); // O ID já é uma string (UUID)
  const navigate = useNavigate();
  const [receita, setReceita] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showDeleteModal, setShowDeleteModal] = useState(false);

  useEffect(() => {
    const fetchReceita = async () => {
      try {
        setLoading(true);
        // Não use parseInt aqui, o ID é um UUID (string)
        const data = await api.getReceitaById(id); 
        setReceita(data);
      } catch (err) {
        console.error(`Erro ao buscar receita com id ${id}:`, err);
        setError('Erro ao carregar detalhes da receita. Por favor, tente novamente.');
      } finally {
        setLoading(false);
      }
    };

    fetchReceita();
  }, [id]);

  const handleDeleteClick = () => {
    setShowDeleteModal(true);
  };

  const handleConfirmDelete = async () => {
    try {
      await api.deleteReceita(receita.id);
      setShowDeleteModal(false);
      navigate('/'); // Redireciona para a página inicial após excluir a receita
    } catch (err) {
      console.error('Erro ao excluir receita:', err);
      setError('Erro ao excluir receita. Por favor, tente novamente.');
    }
  };

  const handleCancelDelete = () => {
    setShowDeleteModal(false);
  };

  if (loading) {
    return <div className="text-center py-8">Carregando detalhes da receita...</div>;
  }

  if (error) {
    return <div className="text-center text-red-600 py-8">{error}</div>;
  }

  if (!receita) {
    return <div className="text-center py-8">Receita não encontrada.</div>;
  }

  return (
    <div className="max-w-3xl mx-auto p-6 bg-white rounded-lg shadow-md">
      <div className="mb-6">
        <Link to="/" className="text-blue-600 hover:underline">
          &larr; Voltar para a lista
        </Link>
      </div>
      
      {/* Usando os nomes de campo corretos do seu modelo Receita */}
      <h1 className="text-3xl font-bold mb-2">{receita.nome}</h1>
      <h2 className="text-xl text-gray-700 mb-4">{receita.descricao}</h2> {/* Descrição como subtítulo */}
      
      <div className="bg-gray-100 p-4 rounded-md mb-6">
        <h3 className="text-lg font-semibold mb-2">Ingredientes</h3>
        {/* Verifica se ingredientes é um array antes de mapear */}
        {receita.ingredientes && Array.isArray(receita.ingredientes) && receita.ingredientes.length > 0 ? (
          <ul className="list-disc list-inside text-gray-700">
            {receita.ingredientes.map((ingrediente, index) => (
              <li key={index}>{ingrediente}</li>
            ))}
          </ul>
        ) : (
          <p className="text-gray-700">Nenhum ingrediente listado.</p>
        )}
      </div>
      
      <div className="mb-8">
        <h3 className="text-lg font-semibold mb-2">Instruções</h3>
        <p className="text-gray-700 whitespace-pre-line">
          {receita.instrucoes || "Nenhuma instrução disponível."}
        </p>
      </div>
      
      <div className="flex justify-between mt-6">
        <Link 
          to={`/editar/${receita.id}`}
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          Editar
        </Link>
        <button
          onClick={handleDeleteClick}
          className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
        >
          Excluir
        </button>
      </div>

      {showDeleteModal && (
        <DeleteConfirmation 
          onConfirm={handleConfirmDelete}
          onCancel={handleCancelDelete}
        />
      )}
    </div>
  );
};

export default ReceitaDetail;