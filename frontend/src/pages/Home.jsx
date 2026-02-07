import { useState, useEffect } from 'react';
import { productsAPI, categoriesAPI, cartAPI } from '../api';
import { useCartStore, useAuthStore } from '../store';
import ProductCard from '../components/ProductCard';
import toast from 'react-hot-toast';

export default function Home() {
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState(null);
  const [loading, setLoading] = useState(true);
  const { setCart } = useCartStore();
  const { isAuthenticated } = useAuthStore();

  useEffect(() => {
    loadData();
  }, []);

  useEffect(() => {
    if (selectedCategory) {
      loadProductsByCategory(selectedCategory);
    } else {
      loadProducts();
    }
  }, [selectedCategory]);

  const loadData = async () => {
    try {
      const [productsRes, categoriesRes] = await Promise.all([
        productsAPI.getAll(),
        categoriesAPI.getAll(),
      ]);
      setProducts(productsRes.data);
      setCategories(categoriesRes.data);
    } catch (error) {
      toast.error('Failed to load products');
    } finally {
      setLoading(false);
    }
  };

  const loadProducts = async () => {
    try {
      const res = await productsAPI.getAll();
      setProducts(res.data);
    } catch (error) {
      toast.error('Failed to load products');
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
    <div style={{ minHeight: '100vh', backgroundColor: '#f9fafb' }}>
      {/* Hero Section */}
      <div style={{
        background: 'linear-gradient(to right, #22c55e, #16a34a)',
        color: 'white',
        padding: '64px 0'
      }}>
        <div style={{
          maxWidth: '1280px',
          margin: '0 auto',
          padding: '0 24px',
          textAlign: 'center'
        }}>
          <h1 style={{ fontSize: '3rem', fontWeight: 'bold', marginBottom: '16px' }}>
            🍕 Smart Food Store
          </h1>
          <p style={{ fontSize: '1.25rem', color: 'rgba(255,255,255,0.8)', marginBottom: '32px' }}>
            Fresh groceries delivered to your door. Let AI plan your meals!
          </p>
          <div style={{ display: 'flex', justifyContent: 'center', flexWrap: 'wrap', gap: '16px' }}>
            <a
              href="/ai-chef"
              style={{
                backgroundColor: 'white',
                color: '#22c55e',
                padding: '12px 32px',
                borderRadius: '9999px',
                fontWeight: '600',
                textDecoration: 'none',
                boxShadow: '0 4px 6px rgba(0,0,0,0.1)'
              }}
            >
              🤖 Try AI Chef
            </a>
            <a
              href="/products"
              style={{
                border: '2px solid white',
                color: 'white',
                padding: '12px 32px',
                borderRadius: '9999px',
                fontWeight: '600',
                textDecoration: 'none',
                boxShadow: '0 4px 6px rgba(0,0,0,0.1)',
                backgroundColor: 'transparent'
              }}
            >
              Browse Products
            </a>
          </div>
        </div>
      </div>

      {/* Category Tabs */}
      <div style={{ maxWidth: '1280px', margin: '0 auto', padding: '48px 24px' }}>
        <div style={{ display: 'flex', flexWrap: 'wrap', justifyContent: 'center', gap: '12px', marginBottom: '32px' }}>
          <button
            onClick={() => setSelectedCategory(null)}
            style={{
              padding: '8px 24px',
              borderRadius: '9999px',
              fontWeight: '500',
              cursor: 'pointer',
              border: selectedCategory === null ? 'none' : '1px solid #e5e7eb',
              backgroundColor: selectedCategory === null ? '#22c55e' : 'white',
              color: selectedCategory === null ? 'white' : '#4b5563',
              boxShadow: '0 1px 2px rgba(0,0,0,0.05)'
            }}
          >
            All Products
          </button>
          {categories.slice(0, 5).map((category) => (
            <button
              key={category.id}
              onClick={() => setSelectedCategory(category.id)}
              style={{
                padding: '8px 24px',
                borderRadius: '9999px',
                fontWeight: '500',
                cursor: 'pointer',
                border: selectedCategory === category.id ? 'none' : '1px solid #e5e7eb',
                backgroundColor: selectedCategory === category.id ? '#22c55e' : 'white',
                color: selectedCategory === category.id ? 'white' : '#4b5563',
                boxShadow: '0 1px 2px rgba(0,0,0,0.05)'
              }}
            >
              {category.name}
            </button>
          ))}
        </div>

        {/* Products Grid */}
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
        ) : (
          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
            gap: '24px'
          }}>
            {products.map((product) => (
              <ProductCard
                key={product.id}
                product={product}
                onCartUpdate={refreshCart}
              />
            ))}
          </div>
        )}

        {products.length === 0 && !loading && (
          <div style={{ textAlign: 'center', padding: '64px 0' }}>
            <p style={{ color: '#6b7280', fontSize: '1.125rem' }}>No products found</p>
          </div>
        )}
      </div>

      {/* Features Section */}
      <div style={{ backgroundColor: 'white', padding: '64px 0', marginTop: '32px' }}>
        <div style={{ maxWidth: '1280px', margin: '0 auto', padding: '0 24px' }}>
          <h2 style={{ fontSize: '1.875rem', fontWeight: 'bold', textAlign: 'center', color: '#1f2937', marginBottom: '48px' }}>
            Why Choose Smart Food Store?
          </h2>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: '32px' }}>
            <div style={{ textAlign: 'center', padding: '32px', backgroundColor: '#f9fafb', borderRadius: '16px', boxShadow: '0 1px 3px rgba(0,0,0,0.1)' }}>
              <div style={{ fontSize: '3rem', marginBottom: '16px' }}>🤖</div>
              <h3 style={{ fontSize: '1.25rem', fontWeight: '600', marginBottom: '12px', color: '#1f2937' }}>AI Recipe Planner</h3>
              <p style={{ color: '#6b7280' }}>
                Tell us what you want to cook, and our AI will find all ingredients for you
              </p>
            </div>
            <div style={{ textAlign: 'center', padding: '32px', backgroundColor: '#f9fafb', borderRadius: '16px', boxShadow: '0 1px 3px rgba(0,0,0,0.1)' }}>
              <div style={{ fontSize: '3rem', marginBottom: '16px' }}>🛒</div>
              <h3 style={{ fontSize: '1.25rem', fontWeight: '600', marginBottom: '12px', color: '#1f2937' }}>Smart Cart</h3>
              <p style={{ color: '#6b7280' }}>
                Based on your cart, get recipe suggestions for what you can cook
              </p>
            </div>
            <div style={{ textAlign: 'center', padding: '32px', backgroundColor: '#f9fafb', borderRadius: '16px', boxShadow: '0 1px 3px rgba(0,0,0,0.1)' }}>
              <div style={{ fontSize: '3rem', marginBottom: '16px' }}>🚚</div>
              <h3 style={{ fontSize: '1.25rem', fontWeight: '600', marginBottom: '12px', color: '#1f2937' }}>Fresh Delivery</h3>
              <p style={{ color: '#6b7280' }}>
                Get fresh groceries delivered straight to your doorstep
              </p>
            </div>
          </div>
        </div>
      </div>

      <style>{`
        @keyframes spin {
          to { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  );
}
