import React from 'react';
import { Link } from 'react-router-dom';

const ReceitaCard = ({ receita, onDeleteClick }) => {
  return (
    <div className="bg-white rounded-lg shadow-md p-4 mb-4">
      <h3 className="text-xl font-semibold mb-2">{receita.nome}</h3>
      <p className="text-gray-700 mb-2">Descrição: {receita.descricao}</p>
      <p className="text-gray-600 mb-2">Ingredientes: {receita.ingredientes}</p>
      <p className="text-gray-600 mb-2">Instruções: {receita.instrucoes}</p>

      <div className="flex justify-between mt-4">
        <Link 
          to={`/receita/${receita.id}`} 
          className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
        >
          Ver Detalhes
        </Link>
        <button 
          onClick={() => onDeleteClick(receita.id)} 
          className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600"
        >
          Excluir
        </button>
      </div>
    </div>
  );
};

export default ReceitaCard;