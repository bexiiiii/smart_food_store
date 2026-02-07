import { useState, useEffect } from 'react';
import { productsAPI, categoriesAPI, cartAPI } from '../api';
import { useCartStore, useAuthStore } from '../store';
import { useSearchParams } from 'react-router-dom';
import ProductCard from '../components/ProductCard';
import toast from 'react-hot-toast';

export default function Products() {
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState(null);
  const [loading, setLoading] = useState(true);
  const [searchParams] = useSearchParams();
  const { setCart } = useCartStore();
  const { isAuthenticated } = useAuthStore();

  const searchQuery = (searchParams.get('q') || searchParams.get('search') || '').trim();

  useEffect(() => {
    loadCategories();
  }, []);

  useEffect(() => {
    if (searchQuery) {
      searchProducts(searchQuery);
    } else if (selectedCategory) {
      loadProductsByCategory(selectedCategory);
    } else {
      loadProducts();
    }
  }, [selectedCategory, searchQuery]);

  const loadCategories = async () => {
    try {
      const res = await categoriesAPI.getAll();
      setCategories(res.data);
    } catch (error) {
      console.error('Failed to load categories');
    }
  };

  const loadProducts = async () => {
    try {
      setLoading(true);
      const res = await productsAPI.getAll();
      setProducts(res.data);
    } catch (error) {
      toast.error('Failed to load products');
    } finally {
      setLoading(false);
    }
  };

  const loadProductsByCategory = async (categoryId) => {
    try {
      setLoading(true);
      const res = await productsAPI.getByCategory(categoryId);
      setProducts(res.data);
    } catch (error) {
      toast.error('Failed to load products');
    } finally {
      setLoading(false);
    }
  };

  const searchProducts = async (query) => {
    try {
      setLoading(true);
      const res = await productsAPI.search(query);
      setProducts(res.data);
    } catch (error) {
      toast.error('Search failed');
    } finally {
      setLoading(false);
    }
  };

  const refreshCart = async () => {
    if (isAuthenticated) {
      try {
        const res = await cartAPI.get();
        setCart(res.data);
      } catch (error) {
        console.error('Failed to refresh cart');
      }
    }
  };

  return (
    <div style={{ minHeight: '100vh', backgroundColor: '#f9fafb', padding: '32px 0' }}>
      <div style={{ maxWidth: '1280px', margin: '0 auto', padding: '0 24px' }}>
        {/* Header */}
        <div style={{ marginBottom: '32px' }}>
          <h1 style={{ fontSize: '1.875rem', fontWeight: 'bold', color: '#1f2937' }}>
            {searchQuery ? `Search: "${searchQuery}"` : 'All Products'}
          </h1>
          <p style={{ color: '#6b7280', marginTop: '8px' }}>
            {products.length} products found
          </p>
        </div>

        {/* Filters */}
        <div style={{ display: 'flex', flexDirection: 'row', gap: '24px' }}>
          {/* Sidebar - Categories */}
          <div style={{ width: '256px', flexShrink: 0 }}>
            <div style={{ backgroundColor: 'white', borderRadius: '12px', boxShadow: '0 1px 3px rgba(0,0,0,0.1)', padding: '24px', position: 'sticky', top: '96px' }}>
              <h3 style={{ fontWeight: '600', color: '#1f2937', marginBottom: '16px' }}>Categories</h3>
              <ul style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                <li>
                  <button
                    onClick={() => setSelectedCategory(null)}
                    style={{
                      width: '100%',
                      textAlign: 'left',
                      padding: '8px 12px',
                      borderRadius: '8px',
                      border: 'none',
                      cursor: 'pointer',
                      backgroundColor: selectedCategory === null ? '#dcfce7' : 'transparent',
                      color: selectedCategory === null ? '#15803d' : '#4b5563',
                      fontWeight: selectedCategory === null ? '500' : 'normal'
                    }}
                  >
                    All Products
                  </button>
                </li>
                {categories.map((category) => (
                  <li key={category.id}>
                    <button
                      onClick={() => setSelectedCategory(category.id)}
                      style={{
                        width: '100%',
                        textAlign: 'left',
                        padding: '8px 12px',
                        borderRadius: '8px',
                        border: 'none',
                        cursor: 'pointer',
                        backgroundColor: selectedCategory === category.id ? '#dcfce7' : 'transparent',
                        color: selectedCategory === category.id ? '#15803d' : '#4b5563',
                        fontWeight: selectedCategory === category.id ? '500' : 'normal'
                      }}
                    >
                      {category.name}
                    </button>
                  </li>
                ))}
              </ul>
            </div>
          </div>

          {/* Products Grid */}
          <div style={{ flex: 1 }}>
            {loading ? (
              <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '256px' }}>
                <div style={{
                  width: '48px',
                  height: '48px',
                  border: '3px solid #e5e7eb',
                  borderTopColor: '#22c55e',
                  borderRadius: '50%',
                  animation: 'spin 1s linear infinite'
                }}></div>
              </div>
            ) : products.length > 0 ? (
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '24px' }}>
                {products.map((product) => (
                  <ProductCard
                    key={product.id}
                    product={product}
                    onCartUpdate={refreshCart}
                  />
                ))}
              </div>
            ) : (
              <div style={{ textAlign: 'center', padding: '64px 0', backgroundColor: 'white', borderRadius: '12px' }}>
                <p style={{ color: '#6b7280', fontSize: '1.125rem' }}>No products found</p>
              </div>
            )}
          </div>
        </div>
      </div>

      <style>{`
        @keyframes spin {
          to { transform: rotate(360deg); }
        }
        @media (max-width: 768px) {
          .sidebar { display: none; }
        }
      `}</style>
    </div>
  );
}
