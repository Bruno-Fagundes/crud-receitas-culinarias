import axios from 'axios';

const instance = axios.create({
  baseURL: 'http://localhost:8081/api',
});

// Adiciona o token antes de cada requisição
instance.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`;
  }
  return config;
});

export default {
  getAllReceitas: async () => {
    const res = await instance.get('/receitas');
    return res.data;
  },

  getReceitaById: async (id) => {
    const res = await instance.get(`/receitas/${id}`);
    return res.data;
  },

  createReceita: async (data) => {
    const res = await instance.post('/receitas', data);
    return res.data;
  },

  updateReceita: async (id, data) => {
    const res = await instance.put(`/receitas/${id}`, data);
    return res.data;
  },

  deleteReceita: async (id) => {
    const res = await instance.delete(`/receitas/${id}`);
    return res.data;
  },

  // Para login com endpoint que não exige token
  login: async (username, password) => {
    const res = await axios.post('http://localhost:8081/login', { username, password });
    return res.data;
  },
};
