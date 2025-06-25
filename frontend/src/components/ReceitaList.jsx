import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import api from '../services/api';
import '../styles.css';

// O componente DeleteConfirmation pode ser mantido aqui ou em um arquivo separado.
// Para este exemplo, vou mantê-lo aqui como você o forneceu.
const DeleteConfirmation = ({ onConfirm, onCancel }) => (
  <div className="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center">
    <div className="bg-white p-6 rounded-lg shadow-lg">
      <p className="mb-4">Tem certeza que deseja excluir esta receita? Esta ação é permanente.</p>
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
        setReceitas(data); 
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

      // Remove a receita da lista após a exclusão bem-sucedida
      setReceitas(prevReceitas => prevReceitas.filter(receita => receita.id !== selectedReceitaId));
      setShowConfirmation(false);
    } catch (err) {
      console.error("Erro ao excluir receita:", err);
      alert("Erro ao excluir a receita.");
    }
  };

  if (loading) return <p className="text-center py-4">Carregando...</p>;
  if (error) return <p className="text-center text-red-500 py-4">{error}</p>;

  return (
    <div className="max-w-6xl mx-auto px-6 py-8">
      <h2 className="text-3xl font-semibold mb-6 text-gray-800">Lista de Receitas</h2>

      <Link
        to="/adicionar"
        className="mb-6 inline-block px-5 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition-colors"
      >
        Adicionar Nova Receita
      </Link>

      <div className="overflow-x-auto">
        <table className="min-w-full bg-white shadow-lg rounded-lg overflow-hidden text-sm md:text-base">
          <thead className="bg-gray-200 text-gray-800">
            <tr>
              <th className="px-6 py-4 text-left">Nome da Receita</th>
              <th className="px-6 py-4 text-left">Descrição</th>
              {/* <th className="px-6 py-4 text-left">Ingredientes</th> Removido */}
              {/* <th className="px-6 py-4 text-left">Instruções</th> Removido */}
              <th className="px-6 py-4 text-left">Ações</th> {/* Manter Ações */}
            </tr>
          </thead>
          <tbody>
            {receitas.length === 0 ? (
              <tr>
                {/* Ajuste o colSpan para o número de colunas visíveis (Nome, Descrição, Ações = 3) */}
                <td colSpan="3" className="text-center px-6 py-6 text-gray-500">
                  Nenhuma receita encontrada.
                </td>
              </tr>
            ) : (
              receitas.map((receita) => (
                <tr key={receita.id} className="hover:bg-gray-50 border-t border-gray-200">
                  <td className="px-6 py-4">{receita.nome}</td>
                  <td className="px-6 py-4">{receita.descricao}</td>
                  {/* <td className="px-6 py-4">{receita.ingredientes ? receita.ingredientes.join(', ') : ''}</td> Removido */}
                  {/* <td className="px-6 py-4">{receita.instrucoes}</td> Removido */}
                  <td className="px-6 py-4">
                    <div className="flex flex-wrap gap-2">
                      <Link
                        to={`/detalhes/${receita.id}`}
                        className="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded"
                      >
                        Detalhes
                      </Link>
                      <Link
                        onClick={() => handleEdit(receita.id)}
                        className="bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-1 rounded"
                      >
                        Editar
                      </Link>
                      <Link
                        onClick={() => initiateDelete(receita.id)}
                        className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded"
                      >
                        Excluir
                      </Link>
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