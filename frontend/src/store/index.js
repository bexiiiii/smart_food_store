import { create } from 'zustand';
import { persist } from 'zustand/middleware';

// Auth Store
export const useAuthStore = create(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      
      setAuth: (user, token) => {
        localStorage.setItem('token', token);
        set({ user, token, isAuthenticated: true });
      },
      
      logout: () => {
        localStorage.removeItem('token');
        set({ user: null, token: null, isAuthenticated: false });
      },
      
      checkAuth: () => {
        const token = localStorage.getItem('token');
        const state = get();
        if (token && state.user) {
          set({ isAuthenticated: true });
        } else if (!token) {
          set({ user: null, token: null, isAuthenticated: false });
        }
      },
      
      isAdmin: () => get().user?.role === 'admin',
    }),
    {
      name: 'auth-storage',
    }
  )
);

// Cart Store
export const useCartStore = create((set, get) => ({
  cart: null,
  loading: false,
  
  setCart: (cart) => set({ cart }),
  setLoading: (loading) => set({ loading }),
  
  getItemCount: () => get().cart?.item_count || 0,
  getTotalPrice: () => get().cart?.total_price || 0,
}));

// Products Store
export const useProductsStore = create((set) => ({
  products: [],
  categories: [],
  selectedCategory: null,
  searchQuery: '',
  loading: false,
  
  setProducts: (products) => set({ products }),
  setCategories: (categories) => set({ categories }),
  setSelectedCategory: (category) => set({ selectedCategory: category }),
  setSearchQuery: (query) => set({ searchQuery: query }),
  setLoading: (loading) => set({ loading }),
}));
