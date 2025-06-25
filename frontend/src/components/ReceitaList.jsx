  import React, { useState, useEffect } from 'react';
  import { Link, useNavigate } from 'react-router-dom';
  import api from '../services/api';
  import '../styles.css';

  const DeleteConfirmation = ({ onConfirm, onCancel }) => (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center">
      <div className="bg-white p-6 rounded-lg shadow-lg">
        <p className="mb-4">Tem certeza que deseja desativar este receita? Ele será marcado como inativo.</p>
        <div className="flex gap-4">
          <button
            onClick={onConfirm}
            className="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded"
          >
            Confirmar
          </button>
          <button
            onClick={onCancel}
            className="bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded"
          >
            Cancelar
          </button>
        </div>
      </div>
    </div>
  );

  function ReceitaList() {
    const [receitas, setReceitas] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [selectedReceitaId, setSelectedReceitaId] = useState(null);
    const [showConfirmation, setShowConfirmation] = useState(false);
    const navigate = useNavigate();

    useEffect(() => {
      const fetchReceitas = async () => {
        try {
          setLoading(true);
          const data = await api.getAllReceitas();
          // Filtra apenas receitas ativos no frontend
          const receitasAtivos = data.filter(receita => receita.ativo === true);
          setReceitas(receitasAtivos);
          setError(null);
        } catch (err) {
          console.error("Erro ao carregar receitas:", err);
          setError("Erro ao carregar receitas. Tente novamente.");
        } finally {
          setLoading(false);
        }
      };

      fetchReceitas();
    }, []);

    const handleEdit = (id) => {
      navigate(`/editar/${id}`);
    };

    const initiateDelete = (id) => {
      setSelectedReceitaId(id);
      setShowConfirmation(true);
    };

    const handleDelete = async () => {
      try {
        await api.deleteReceita(selectedReceitaId);

        setReceitas(prevReceitas => prevReceitas.filter(receita => receita.id !== selectedReceitaId));
        setShowConfirmation(false);
      } catch (err) {
        console.error("Erro ao desativar receita:", err);
        alert("Erro ao desativar o receita.");
      }
    };

    if (loading) return <p className="text-center py-4">Carregando...</p>;
    if (error) return <p className="text-center text-red-500 py-4">{error}</p>;

    return (
      <div className="max-w-6xl mx-auto px-6 py-8">
        <h2 className="text-3xl font-semibold mb-6 text-gray-800">Lista de receitas</h2>

        <Link
          to="/adicionar"
          className="mb-6 inline-block px-5 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition-colors"
        >
          Adicionar Novo receita
        </Link>

        <div className="overflow-x-auto">
          <table className="min-w-full bg-white shadow-lg rounded-lg overflow-hidden text-sm md:text-base">
            <thead className="bg-gray-200 text-gray-800">
              <tr>
                <th className="px-6 py-4 text-left">Título</th>
                <th className="px-6 py-4 text-left">Autor</th>
                <th className="px-6 py-4 text-left">Ano</th>
                <th className="px-6 py-4 text-left">Ações</th>
              </tr>
            </thead>
            <tbody>
              {receitas.length === 0 ? (
                <tr>
                  <td colSpan="4" className="text-center px-6 py-6 text-gray-500">
                    Nenhuma receita encontrada.
                  </td>
                </tr>
              ) : (
                receitas.map((receita) => (
                  <tr key={receita.id} className="hover:bg-gray-50 border-t border-gray-200">
                    <td className="px-6 py-4">{receita.nome}</td>
                    <td className="px-6 py-4">{receita.descricao}</td>
                    <td className="px-6 py-4">{receita.ingredientes}</td>
                  <td className="px-6 py-4">{receita.ingredientes}</td>

                    <td className="px-6 py-4">
                      <div className="flex flex-wrap gap-2">
                        <Link
                          to={`/detalhes/${receita.id}`}
                          className="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded"
                        >
                          Detalhes
                        </Link>
                        <button
                          onClick={() => handleEdit(receita.id)}
                          className="bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-1 rounded"
                        >
                          Editar
                        </button>
                        <button
                          onClick={() => initiateDelete(receita.id)}
                          className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded"
                        >
                          Excluir
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {showConfirmation && (
          <DeleteConfirmation
            onConfirm={handleDelete}
            onCancel={() => {
              setSelectedReceitaId(null);
              setShowConfirmation(false);
            }}
          />
        )}
      </div>
    );
  }

  export default ReceitaList;