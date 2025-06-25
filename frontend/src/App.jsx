import React from 'react';
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate
} from 'react-router-dom';

import Login          from './components/Login';
import ReceitaList    from './components/ReceitaList';
import ReceitaForm    from './components/ReceitaForm';
import ReceitaDetail  from './components/ReceitaDetail';

function PrivateRoute({ children }) {
  const token = localStorage.getItem('token');
  return token ? children : <Navigate to="/login" />;
}

export default function App() {
  return (
    <Router>
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-6">Gerenciador de Receitas</h1>
        <h2 className="text-xl mb-4">Sistema de Gerenciamento de Receitas</h2>

        <Routes>
          <Route path="/login" element={<Login />} />

          <Route
            path="/"
            element={
              <PrivateRoute>
                <ReceitaList />
              </PrivateRoute>
            }
          />

          <Route
            path="/adicionar"
            element={
              <PrivateRoute>
                <ReceitaForm />
              </PrivateRoute>
            }
          />

          <Route
            path="/editar/:id"
            element={
              <PrivateRoute>
                <ReceitaForm />
              </PrivateRoute>
            }
          />

          <Route
            path="/detalhes/:id"
            element={
              <PrivateRoute>
                <ReceitaDetail />
              </PrivateRoute>
            }
          />
        </Routes>
      </div>
    </Router>
  );
}
