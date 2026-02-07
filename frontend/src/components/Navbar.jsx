import { Link, useNavigate } from 'react-router-dom';
import { FiShoppingCart, FiUser, FiSearch, FiLogOut, FiSettings } from 'react-icons/fi';
import { useAuthStore, useCartStore } from '../store';
import { useState } from 'react';

export default function Navbar() {
  const { isAuthenticated, user, logout, isAdmin } = useAuthStore();
  const { getItemCount } = useCartStore();
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  const handleSearch = (e) => {
    e.preventDefault();
    const query = searchQuery.trim();

    if (query) {
      navigate(`/products?q=${encodeURIComponent(query)}`);
    } else {
      navigate('/products');
    }
  };

  return (
    <nav style={{ 
      backgroundColor: 'white', 
      boxShadow: '0 1px 3px rgba(0,0,0,0.1)', 
      position: 'sticky', 
      top: 0, 
      zIndex: 50 
    }}>
      <div style={{ maxWidth: '1280px', margin: '0 auto', padding: '0 24px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', height: '64px' }}>
          {/* Logo */}
          <Link to="/" style={{ display: 'flex', alignItems: 'center', gap: '8px', textDecoration: 'none' }}>
            <span style={{ fontSize: '28px' }}>🍕</span>
            <span style={{ fontWeight: 'bold', fontSize: '20px', color: '#1f2937' }}>Smart Food Store</span>
          </Link>

          {/* Search Bar */}
          <form onSubmit={handleSearch} style={{ flex: 1, maxWidth: '400px', margin: '0 32px', display: 'flex' }}>
            <div style={{ position: 'relative', width: '100%' }}>
              <input
                type="text"
                placeholder="Search products..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                style={{
                  width: '100%',
                  paddingLeft: '40px',
                  paddingRight: '16px',
                  paddingTop: '10px',
                  paddingBottom: '10px',
                  border: '1px solid #d1d5db',
                  borderRadius: '9999px',
                  outline: 'none',
                  fontSize: '14px'
                }}
              />
              <FiSearch style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: '#9ca3af' }} />
            </div>
          </form>

          {/* Navigation Links */}
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            <Link to="/products" style={{ color: '#4b5563', fontWeight: '500', padding: '8px 16px', textDecoration: 'none' }}>
              Products
            </Link>
            <Link to="/recipes" style={{ color: '#4b5563', fontWeight: '500', padding: '8px 16px', textDecoration: 'none' }}>
              Recipes
            </Link>
            <Link to="/ai-chef" style={{ color: '#4b5563', fontWeight: '500', padding: '8px 16px', textDecoration: 'none' }}>
              🤖 AI Chef
            </Link>

            {isAuthenticated ? (
              <div style={{ display: 'flex', alignItems: 'center', gap: '16px', marginLeft: '16px' }}>
                <Link to="/cart" style={{ position: 'relative', padding: '8px', color: '#4b5563' }}>
                  <FiShoppingCart style={{ width: '24px', height: '24px' }} />
                  {getItemCount() > 0 && (
                    <span style={{
                      position: 'absolute', top: '-4px', right: '-4px',
                      backgroundColor: '#22c55e', color: 'white', fontSize: '12px',
                      borderRadius: '9999px', width: '20px', height: '20px',
                      display: 'flex', alignItems: 'center', justifyContent: 'center'
                    }}>{getItemCount()}</span>
                  )}
                </Link>
                <div style={{ position: 'relative' }} className="group">
                  <button style={{ display: 'flex', alignItems: 'center', gap: '8px', color: '#4b5563', background: 'none', border: 'none', cursor: 'pointer' }}>
                    <FiUser style={{ width: '24px', height: '24px' }} />
                    <span>{user?.name}</span>
                  </button>
                  <div className="hidden group-hover:block" style={{
                    position: 'absolute', right: 0, marginTop: '8px', width: '192px',
                    backgroundColor: 'white', borderRadius: '8px', boxShadow: '0 10px 15px -3px rgba(0,0,0,0.1)', padding: '8px 0'
                  }}>
                    <Link to="/profile" style={{ display: 'block', padding: '8px 16px', color: '#374151', textDecoration: 'none' }}>Profile</Link>
                    {isAdmin() && (
                      <Link to="/admin" style={{ display: 'flex', alignItems: 'center', padding: '8px 16px', color: '#374151', textDecoration: 'none' }}>
                        <FiSettings style={{ marginRight: '8px' }} /> Admin Panel
                      </Link>
                    )}
                    <button onClick={handleLogout} style={{ display: 'flex', alignItems: 'center', width: '100%', padding: '8px 16px', color: '#374151', background: 'none', border: 'none', cursor: 'pointer', textAlign: 'left' }}>
                      <FiLogOut style={{ marginRight: '8px' }} /> Logout
                    </button>
                  </div>
                </div>
              </div>
            ) : (
              <div style={{ display: 'flex', alignItems: 'center', marginLeft: '16px', gap: '12px' }}>
                <Link to="/login" style={{ padding: '8px 16px', color: '#22c55e', fontWeight: '500', textDecoration: 'none' }}>Login</Link>
                <Link to="/register" style={{ padding: '10px 20px', backgroundColor: '#22c55e', color: 'white', borderRadius: '9999px', fontWeight: '500', textDecoration: 'none', boxShadow: '0 1px 3px rgba(0,0,0,0.1)' }}>Sign Up</Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}
