import { useState } from 'react';
import { aiAPI, productsAPI, cartAPI } from '../api';
import { useCartStore, useAuthStore } from '../store';
import { FiSearch, FiShoppingCart, FiArrowRight, FiPlus, FiCheck } from 'react-icons/fi';
import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';

export default function AIChef() {
  const [activeTab, setActiveTab] = useState('dish');
  const [dishQuery, setDishQuery] = useState('');
  const [dishResult, setDishResult] = useState(null);
  const [cartSuggestions, setCartSuggestions] = useState(null);
  const [loading, setLoading] = useState(false);
  const [addedProducts, setAddedProducts] = useState(new Set());
  
  const { cart, setCart } = useCartStore();
  const { isAuthenticated } = useAuthStore();
  const navigate = useNavigate();

  const handleDishSearch = async (e) => {
    e.preventDefault();
    if (!dishQuery.trim()) {
      toast.error('Please enter a dish name');
      return;
    }

    setLoading(true);
    try {
      const res = await aiAPI.dishToIngredients(dishQuery);
      setDishResult(res.data);
      setAddedProducts(new Set());
    } catch (error) {
      toast.error(error.response?.data?.error || 'Failed to get ingredients');
    } finally {
      setLoading(false);
    }
  };

  const handleCartSuggestions = async () => {
    if (!isAuthenticated) {
      toast.error('Please login to use this feature');
      navigate('/login');
      return;
    }

    if (!cart || cart.items?.length === 0) {
      toast.error('Your cart is empty');
      return;
    }

    setLoading(true);
    try {
      const res = await aiAPI.cartToRecipes();
      setCartSuggestions(res.data);
    } catch (error) {
      toast.error(error.response?.data?.error || 'Failed to get suggestions');
    } finally {
      setLoading(false);
    }
  };

  const addToCart = async (productId) => {
    if (!isAuthenticated) {
      toast.error('Please login to add items to cart');
      navigate('/login');
      return;
    }

    try {
      await cartAPI.addItem(productId, 1);
      setAddedProducts(prev => new Set([...prev, productId]));
      toast.success('Added to cart!');
      
      // Refresh cart
      const res = await cartAPI.get();
      setCart(res.data);
    } catch (error) {
      toast.error('Failed to add to cart');
    }
  };

  const addAllToCart = async () => {
    if (!dishResult?.matched_products) return;
    
    for (const product of dishResult.matched_products) {
      if (!addedProducts.has(product.id)) {
        await addToCart(product.id);
      }
    }
    toast.success('All ingredients added to cart!');
  };

  return (
    <div style={{ minHeight: '100vh', background: 'linear-gradient(to bottom right, #f0fdf4, white, #fef3c7)', padding: '32px 0' }}>
      <div style={{ maxWidth: '896px', margin: '0 auto', padding: '0 24px' }}>
        {/* Header */}
        <div style={{ textAlign: 'center', marginBottom: '48px' }}>
          <div style={{
            width: '80px',
            height: '80px',
            background: 'linear-gradient(to bottom right, #4ade80, #facc15)',
            borderRadius: '50%',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            margin: '0 auto 24px',
            boxShadow: '0 4px 6px rgba(0,0,0,0.1)'
          }}>
            <span style={{ fontSize: '2.5rem' }}>🤖</span>
          </div>
          <h1 style={{ fontSize: '2.25rem', fontWeight: 'bold', color: '#1f2937', marginBottom: '16px' }}>AI Chef Assistant</h1>
          <p style={{ color: '#6b7280', maxWidth: '576px', margin: '0 auto' }}>
            Let our AI help you find ingredients for any dish or suggest recipes based on your cart
          </p>
        </div>

        {/* Tabs */}
        <div style={{ display: 'flex', gap: '16px', marginBottom: '32px', justifyContent: 'center' }}>
          <button
            onClick={() => setActiveTab('dish')}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              padding: '12px 24px',
              borderRadius: '12px',
              fontWeight: '500',
              border: 'none',
              cursor: 'pointer',
              backgroundColor: activeTab === 'dish' ? '#22c55e' : 'white',
              color: activeTab === 'dish' ? 'white' : '#4b5563',
              boxShadow: activeTab === 'dish' ? '0 4px 6px rgba(34,197,94,0.3)' : '0 1px 3px rgba(0,0,0,0.1)'
            }}
          >
            <FiSearch style={{ width: '20px', height: '20px' }} />
            Find Ingredients
          </button>
          <button
            onClick={() => {
              setActiveTab('cart');
              if (activeTab !== 'cart') {
                handleCartSuggestions();
              }
            }}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              padding: '12px 24px',
              borderRadius: '12px',
              fontWeight: '500',
              border: 'none',
              cursor: 'pointer',
              backgroundColor: activeTab === 'cart' ? '#f59e0b' : 'white',
              color: activeTab === 'cart' ? 'white' : '#4b5563',
              boxShadow: activeTab === 'cart' ? '0 4px 6px rgba(245,158,11,0.3)' : '0 1px 3px rgba(0,0,0,0.1)'
            }}
          >
            <FiShoppingCart style={{ width: '20px', height: '20px' }} />
            Recipe Suggestions
          </button>
        </div>

        {/* Dish to Ingredients Tab */}
        {activeTab === 'dish' && (
          <div style={{ backgroundColor: 'white', borderRadius: '16px', boxShadow: '0 10px 15px rgba(0,0,0,0.1)', padding: '32px' }}>
            <form onSubmit={handleDishSearch} style={{ marginBottom: '32px' }}>
              <label style={{ display: 'block', color: '#374151', fontWeight: '500', marginBottom: '8px' }}>
                What do you want to cook?
              </label>
              <div style={{ display: 'flex', gap: '16px' }}>
                <input
                  type="text"
                  value={dishQuery}
                  onChange={(e) => setDishQuery(e.target.value)}
                  placeholder="e.g., Pasta Carbonara, Caesar Salad..."
                  style={{
                    flex: 1,
                    padding: '12px 16px',
                    border: '1px solid #e5e7eb',
                    borderRadius: '12px',
                    outline: 'none',
                    fontSize: '1rem'
                  }}
                />
                <button
                  type="submit"
                  disabled={loading}
                  style={{
                    padding: '12px 24px',
                    backgroundColor: '#22c55e',
                    color: 'white',
                    borderRadius: '12px',
                    border: 'none',
                    cursor: 'pointer',
                    fontWeight: '500',
                    display: 'flex',
                    alignItems: 'center',
                    gap: '8px',
                    opacity: loading ? 0.5 : 1
                  }}
                >
                  {loading ? (
                    <div style={{
                      width: '20px',
                      height: '20px',
                      border: '2px solid rgba(255,255,255,0.3)',
                      borderTopColor: 'white',
                      borderRadius: '50%',
                      animation: 'spin 1s linear infinite'
                    }}></div>
                  ) : (
                    <>
                      <FiSearch style={{ width: '20px', height: '20px' }} />
                      Search
                    </>
                  )}
                </button>
              </div>
            </form>

            {/* Dish Result */}
            {dishResult && (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
                <div style={{ background: 'linear-gradient(to right, #f0fdf4, #fef3c7)', borderRadius: '12px', padding: '24px' }}>
                  <h3 style={{ fontSize: '1.25rem', fontWeight: '600', color: '#1f2937', marginBottom: '8px' }}>
                    {dishResult.dish_name}
                  </h3>
                  <p style={{ color: '#6b7280' }}>{dishResult.description}</p>
                </div>

                {/* Required Ingredients */}
                {dishResult.required_ingredients?.length > 0 && (
                  <div>
                    <h4 style={{ fontWeight: '500', color: '#1f2937', marginBottom: '12px' }}>Required Ingredients:</h4>
                    <ul style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '8px' }}>
                      {dishResult.required_ingredients.map((ing, idx) => (
                        <li key={idx} style={{ display: 'flex', alignItems: 'center', gap: '8px', color: '#6b7280', backgroundColor: '#f9fafb', padding: '8px 12px', borderRadius: '8px' }}>
                          <span style={{ width: '8px', height: '8px', backgroundColor: '#4ade80', borderRadius: '50%' }}></span>
                          {ing.name} - {ing.quantity} {ing.unit}
                        </li>
                      ))}
                    </ul>
                  </div>
                )}

                {/* Matched Products */}
                {dishResult.matched_products?.length > 0 && (
                  <div>
                    <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '16px' }}>
                      <h4 style={{ fontWeight: '500', color: '#1f2937' }}>
                        Available in Store ({dishResult.matched_products.length})
                      </h4>
                      <button
                        onClick={addAllToCart}
                        style={{
                          fontSize: '0.875rem',
                          backgroundColor: '#dcfce7',
                          color: '#15803d',
                          padding: '8px 16px',
                          borderRadius: '8px',
                          border: 'none',
                          cursor: 'pointer',
                          display: 'flex',
                          alignItems: 'center',
                          gap: '8px'
                        }}
                      >
                        <FiPlus style={{ width: '16px', height: '16px' }} />
                        Add All to Cart
                      </button>
                    </div>
                    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))', gap: '16px' }}>
                      {dishResult.matched_products.map((product) => (
                        <div
                          key={product.id}
                          style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', backgroundColor: '#f9fafb', padding: '16px', borderRadius: '12px' }}
                        >
                          <div>
                            <p style={{ fontWeight: '500', color: '#1f2937' }}>{product.name}</p>
                            <p style={{ color: '#22c55e', fontWeight: '600' }}>${product.price.toFixed(2)}</p>
                          </div>
                          <button
                            onClick={() => addToCart(product.id)}
                            disabled={addedProducts.has(product.id)}
                            style={{
                              width: '40px',
                              height: '40px',
                              borderRadius: '50%',
                              display: 'flex',
                              alignItems: 'center',
                              justifyContent: 'center',
                              border: 'none',
                              cursor: 'pointer',
                              backgroundColor: addedProducts.has(product.id) ? '#dcfce7' : '#22c55e',
                              color: addedProducts.has(product.id) ? '#22c55e' : 'white'
                            }}
                          >
                            {addedProducts.has(product.id) ? (
                              <FiCheck style={{ width: '20px', height: '20px' }} />
                            ) : (
                              <FiPlus style={{ width: '20px', height: '20px' }} />
                            )}
                          </button>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Cooking Tips */}
                {dishResult.cooking_tips && (
                  <div style={{ backgroundColor: '#fef3c7', borderRadius: '12px', padding: '24px' }}>
                    <h4 style={{ fontWeight: '500', color: '#1f2937', marginBottom: '8px' }}>💡 Cooking Tips</h4>
                    <p style={{ color: '#6b7280' }}>{dishResult.cooking_tips}</p>
                  </div>
                )}
              </div>
            )}

            {/* Placeholder when no search */}
            {!dishResult && !loading && (
              <div style={{ textAlign: 'center', padding: '48px 0', color: '#9ca3af' }}>
                <FiSearch style={{ width: '64px', height: '64px', margin: '0 auto 16px', strokeWidth: 1 }} />
                <p>Enter a dish name to find all the ingredients you need</p>
              </div>
            )}
          </div>
        )}

        {/* Cart to Recipes Tab */}
        {activeTab === 'cart' && (
          <div style={{ backgroundColor: 'white', borderRadius: '16px', boxShadow: '0 10px 15px rgba(0,0,0,0.1)', padding: '32px' }}>
            {loading ? (
              <div style={{ display: 'flex', justifyContent: 'center', padding: '48px 0' }}>
                <div style={{
                  width: '48px',
                  height: '48px',
                  border: '3px solid #e5e7eb',
                  borderTopColor: '#f59e0b',
                  borderRadius: '50%',
                  animation: 'spin 1s linear infinite'
                }}></div>
              </div>
            ) : cartSuggestions ? (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
                <div style={{ background: 'linear-gradient(to right, #fef3c7, #f0fdf4)', borderRadius: '12px', padding: '24px' }}>
                  <h3 style={{ fontSize: '1.25rem', fontWeight: '600', color: '#1f2937', marginBottom: '8px' }}>
                    Recipes You Can Make
                  </h3>
                  <p style={{ color: '#6b7280' }}>
                    Based on {cartSuggestions.cart_items?.length || 0} items in your cart
                  </p>
                </div>

                {/* Cart Items */}
                {cartSuggestions.cart_items?.length > 0 && (
                  <div>
                    <h4 style={{ fontWeight: '500', color: '#1f2937', marginBottom: '12px' }}>Your Ingredients:</h4>
                    <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
                      {cartSuggestions.cart_items.map((item, idx) => (
                        <span
                          key={idx}
                          style={{ backgroundColor: '#f3f4f6', color: '#374151', padding: '4px 12px', borderRadius: '9999px', fontSize: '0.875rem' }}
                        >
                          {item}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Suggested Recipes */}
                {cartSuggestions.recipes?.length > 0 && (
                  <div>
                    <h4 style={{ fontWeight: '500', color: '#1f2937', marginBottom: '16px' }}>Suggested Recipes:</h4>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                      {cartSuggestions.recipes.map((recipe, idx) => (
                        <div
                          key={idx}
                          style={{ border: '1px solid #f3f4f6', borderRadius: '12px', padding: '24px' }}
                        >
                          <h5 style={{ fontSize: '1.125rem', fontWeight: '600', color: '#1f2937', marginBottom: '8px' }}>
                            {recipe.name}
                          </h5>
                          <p style={{ color: '#6b7280', fontSize: '0.875rem', marginBottom: '16px' }}>{recipe.description}</p>
                          
                          {recipe.instructions && (
                            <div style={{ backgroundColor: '#f9fafb', borderRadius: '8px', padding: '16px' }}>
                              <h6 style={{ fontWeight: '500', color: '#374151', marginBottom: '8px' }}>How to make:</h6>
                              <p style={{ color: '#6b7280', fontSize: '0.875rem', whiteSpace: 'pre-line' }}>
                                {recipe.instructions}
                              </p>
                            </div>
                          )}

                          {recipe.missing_ingredients?.length > 0 && (
                            <div style={{ marginTop: '16px', backgroundColor: '#fff7ed', borderRadius: '8px', padding: '16px' }}>
                              <h6 style={{ fontWeight: '500', color: '#c2410c', marginBottom: '8px' }}>
                                Missing Ingredients:
                              </h6>
                              <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
                                {recipe.missing_ingredients.map((ing, i) => (
                                  <span
                                    key={i}
                                    style={{ backgroundColor: '#ffedd5', color: '#c2410c', padding: '4px 8px', borderRadius: '4px', fontSize: '0.875rem' }}
                                  >
                                    {ing}
                                  </span>
                                ))}
                              </div>
                            </div>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ) : (
              <div style={{ textAlign: 'center', padding: '48px 0' }}>
                <FiShoppingCart style={{ width: '64px', height: '64px', margin: '0 auto 16px', color: '#d1d5db' }} />
                <p style={{ color: '#6b7280', marginBottom: '16px' }}>
                  {isAuthenticated
                    ? 'Add items to your cart to get recipe suggestions'
                    : 'Please login to use this feature'}
                </p>
                <button
                  onClick={() => navigate(isAuthenticated ? '/products' : '/login')}
                  style={{
                    display: 'inline-flex',
                    alignItems: 'center',
                    gap: '8px',
                    backgroundColor: '#f59e0b',
                    color: 'white',
                    padding: '12px 24px',
                    borderRadius: '12px',
                    border: 'none',
                    cursor: 'pointer',
                    fontWeight: '500'
                  }}
                >
                  {isAuthenticated ? 'Browse Products' : 'Login'}
                  <FiArrowRight />
                </button>
              </div>
            )}
          </div>
        )}
      </div>

      <style>{`
        @keyframes spin {
          to { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  );
}
