import React, { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import api from '../services/api';
import DeleteConfirmation from './DeleteConfirmation';

const ReceitaDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [receita, setReceita] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showDeleteModal, setShowDeleteModal] = useState(false);

  useEffect(() => {
    const fetchReceita = async () => {
      try {
        setLoading(true);
        const data = await api.getReceitaById(parseInt(id, 10));
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
      navigate('/'); // Redireciona para a página inicial após "excluir" o receita
    } catch (err) {
      console.error('Erro ao excluir receita:', err);
      setError('Erro ao excluir receita. Por favor, tente novamente.');
    }
  };

  const handleCancelDelete = () => {
    setShowDeleteModal(false);
  };

  if (loading) {
    return <div className="text-center py-8">Carregando detalhes do receita...</div>;
  }

  if (error) {
    return <div className="text-center text-red-600 py-8">{error}</div>;
  }

  if (!receita) {
    return <div className="text-center py-8">Receita não encontrado.</div>;
  }

  return (
    <div className="max-w-3xl mx-auto p-6 bg-white rounded-lg shadow-md">
      <div className="mb-6">
        <Link to="/" className="text-blue-600 hover:underline">
          &larr; Voltar para a lista
        </Link>
      </div>
      
      <h1 className="text-3xl font-bold mb-2">{receita.titulo}</h1>
      <h2 className="text-xl text-gray-700 mb-4">por {receita.autor}</h2>
      
      <div className="bg-gray-100 p-4 rounded-md mb-6">
        <p className="text-sm text-gray-600 mb-2">Ano de Publicação</p>
        <p className="font-semibold">{receita.anoPublicacao}</p>
      </div>
      
      <div className="mb-8">
        <h3 className="text-lg font-semibold mb-2">Sinopse</h3>
        <p className="text-gray-700 whitespace-pre-line">
          {receita.sinopse || "Nenhuma sinopse disponível."}
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