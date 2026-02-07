import axios from 'axios';

const API_URL = 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
  getProfile: () => api.get('/users/me'),
};

// Products API
export const productsAPI = {
  getAll: () => api.get('/products'),
  getById: (id) => api.get(`/products/${id}`),
  getByCategory: (categoryId) => api.get(`/products/category/${categoryId}`),
  search: (query) => api.get(`/products/search?q=${query}`),
};

// Categories API
export const categoriesAPI = {
  getAll: () => api.get('/categories'),
};

// Cart API
export const cartAPI = {
  get: () => api.get('/cart'),
  addItem: (productId, quantity) => api.post('/cart/items', { product_id: productId, quantity }),
  updateItem: (productId, quantity) => api.put(`/cart/items/${productId}`, { quantity }),
  removeItem: (productId) => api.delete(`/cart/items/${productId}`),
  clear: () => api.delete('/cart'),
  addBulk: (items) => api.post('/cart/items/bulk', items),
};

// Recipes API
export const recipesAPI = {
  getAll: () => api.get('/recipes'),
  getById: (id) => api.get(`/recipes/${id}`),
  search: (query) => api.get(`/recipes/search?q=${query}`),
  calculate: (id, servings) => api.get(`/recipes/${id}/calculate?servings=${servings}`),
  addToCart: (id, servings) => api.post(`/recipes/${id}/add-to-cart`, { servings }),
};

// AI API
export const aiAPI = {
  dishToIngredients: (dishName, servings) => 
    api.post('/ai/dish-to-ingredients', { dish_name: dishName, servings }),
  cartToRecipes: () => api.get('/ai/cart-to-recipes'),
  productsToRecipes: (productIds) => 
    api.post('/ai/products-to-recipes', { product_ids: productIds }),
  addSuggestionToCart: (suggestion) => api.post('/ai/add-to-cart', suggestion),
};

// Admin API
export const adminAPI = {
  // Users
  getUsers: () => api.get('/admin/users'),
  updateUserRole: (id, role) => api.patch(`/admin/users/${id}/role`, { role }),
  deleteUser: (id) => api.delete(`/admin/users/${id}`),
  
  // Products
  createProduct: (data) => api.post('/admin/products', data),
  updateProduct: (id, data) => api.put(`/admin/products/${id}`, data),
  deleteProduct: (id) => api.delete(`/admin/products/${id}`),
  
  // Categories
  createCategory: (data) => api.post('/admin/categories', data),
  updateCategory: (id, name) => api.put(`/admin/categories/${id}`, { name }),
  deleteCategory: (id) => api.delete(`/admin/categories/${id}`),
  
  // Recipes
  createRecipe: (data) => api.post('/admin/recipes', data),
  updateRecipe: (id, data) => api.put(`/admin/recipes/${id}`, data),
  deleteRecipe: (id) => api.delete(`/admin/recipes/${id}`),
};

export default api;
