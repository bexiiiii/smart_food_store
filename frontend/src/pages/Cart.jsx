import { useState, useEffect } from 'react';
import { cartAPI } from '../api';
import { useCartStore, useAuthStore } from '../store';
import { FiTrash2, FiPlus, FiMinus, FiShoppingBag } from 'react-icons/fi';
import { Link, useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';

export default function Cart() {
  const { cart, setCart, loading } = useCartStore();
  const { isAuthenticated } = useAuthStore();
  const [updating, setUpdating] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }
    loadCart();
  }, [isAuthenticated]);

  const loadCart = async () => {
    try {
      const res = await cartAPI.get();
      setCart(res.data);
    } catch (error) {
      toast.error('Failed to load cart');
    }
  };

  const updateQuantity = async (productId, quantity) => {
    setUpdating(true);
    try {
      const res = await cartAPI.updateItem(productId, quantity);
      setCart(res.data);
    } catch (error) {
      toast.error('Failed to update quantity');
    } finally {
      setUpdating(false);
    }
  };

  const removeItem = async (productId) => {
    try {
      const res = await cartAPI.removeItem(productId);
      setCart(res.data);
      toast.success('Item removed');
    } catch (error) {
      toast.error('Failed to remove item');
    }
  };

  const clearCart = async () => {
    try {
      await cartAPI.clear();
      setCart({ items: [], total_price: 0, item_count: 0 });
      toast.success('Cart cleared');
    } catch (error) {
      toast.error('Failed to clear cart');
    }
  };

  if (!cart || cart.items?.length === 0) {
    return (
      <div style={{ minHeight: '100vh', backgroundColor: '#f9fafb', padding: '48px 0' }}>
        <div style={{ maxWidth: '896px', margin: '0 auto', padding: '0 16px', textAlign: 'center' }}>
          <div style={{ backgroundColor: 'white', borderRadius: '16px', boxShadow: '0 10px 15px -3px rgba(0,0,0,0.1)', padding: '48px' }}>
            <FiShoppingBag style={{ width: '80px', height: '80px', color: '#d1d5db', margin: '0 auto 24px' }} />
            <h2 style={{ fontSize: '24px', fontWeight: 'bold', color: '#1f2937', marginBottom: '16px' }}>Your cart is empty</h2>
            <p style={{ color: '#4b5563', marginBottom: '32px' }}>
              Looks like you haven't added any items yet.
            </p>
            <Link
              to="/products"
              style={{ display: 'inline-block', backgroundColor: '#22c55e', color: 'white', padding: '12px 32px', borderRadius: '9999px', fontWeight: '600', textDecoration: 'none' }}
            >
              Start Shopping
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div style={{ minHeight: '100vh', backgroundColor: '#f9fafb', padding: '32px 0' }}>
      <div style={{ maxWidth: '1152px', margin: '0 auto', padding: '0 16px' }}>
        <h1 style={{ fontSize: '30px', fontWeight: 'bold', color: '#1f2937', marginBottom: '32px' }}>Shopping Cart</h1>

        <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '32px' }}>
          {/* Cart Items */}
          <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
            {cart.items?.map((item) => (
              <div key={item.id} style={{ backgroundColor: 'white', borderRadius: '12px', boxShadow: '0 1px 3px rgba(0,0,0,0.1)', padding: '24px', display: 'flex', alignItems: 'center', gap: '24px' }}>
                {/* Product Image */}
                <div style={{ width: '96px', height: '96px', backgroundColor: '#f3f4f6', borderRadius: '8px', display: 'flex', alignItems: 'center', justifyContent: 'center', flexShrink: 0 }}>
                  <span style={{ fontSize: '36px' }}>🛒</span>
                </div>

                {/* Product Info */}
                <div style={{ flex: 1, minWidth: 0 }}>
                  <h3 style={{ fontWeight: '600', color: '#1f2937', fontSize: '18px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                    {item.product_name}
                  </h3>
                  <p style={{ color: '#6b7280', fontSize: '14px' }}>
                    ${item.price.toFixed(2)} per {item.unit}
                  </p>
                  <p style={{ color: '#22c55e', fontWeight: 'bold', marginTop: '4px' }}>
                    Subtotal: ${item.subtotal.toFixed(2)}
                  </p>
                </div>

                {/* Quantity Controls */}
                <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                  <button
                    onClick={() => updateQuantity(item.product_id, item.quantity - 1)}
                    disabled={updating || item.quantity <= 1}
                    style={{ padding: '8px', borderRadius: '9999px', backgroundColor: '#f3f4f6', border: 'none', cursor: 'pointer', opacity: updating || item.quantity <= 1 ? 0.5 : 1 }}
                  >
                    <FiMinus style={{ width: '16px', height: '16px' }} />
                  </button>
                  <span style={{ width: '48px', textAlign: 'center', fontWeight: '600' }}>{item.quantity}</span>
                  <button
                    onClick={() => updateQuantity(item.product_id, item.quantity + 1)}
                    disabled={updating}
                    style={{ padding: '8px', borderRadius: '9999px', backgroundColor: '#f3f4f6', border: 'none', cursor: 'pointer', opacity: updating ? 0.5 : 1 }}
                  >
                    <FiPlus style={{ width: '16px', height: '16px' }} />
                  </button>
                </div>

                {/* Remove Button */}
                <button
                  onClick={() => removeItem(item.product_id)}
                  style={{ padding: '8px', color: '#ef4444', background: 'none', border: 'none', cursor: 'pointer', borderRadius: '9999px' }}
                >
                  <FiTrash2 style={{ width: '20px', height: '20px' }} />
                </button>
              </div>
            ))}

            {/* Clear Cart Button */}
            <button
              onClick={clearCart}
              style={{ color: '#ef4444', fontWeight: '500', display: 'flex', alignItems: 'center', gap: '8px', background: 'none', border: 'none', cursor: 'pointer' }}
            >
              <FiTrash2 /> Clear Cart
            </button>
          </div>

          {/* Order Summary */}
          <div>
            <div style={{ backgroundColor: 'white', borderRadius: '12px', boxShadow: '0 1px 3px rgba(0,0,0,0.1)', padding: '24px', position: 'sticky', top: '96px' }}>
              <h2 style={{ fontSize: '20px', fontWeight: 'bold', color: '#1f2937', marginBottom: '24px' }}>Order Summary</h2>
              
              <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', marginBottom: '24px' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', color: '#4b5563' }}>
                  <span>Items ({cart.item_count})</span>
                  <span>${cart.total_price?.toFixed(2)}</span>
                </div>
                <div style={{ display: 'flex', justifyContent: 'space-between', color: '#4b5563' }}>
                  <span>Delivery</span>
                  <span style={{ color: '#22c55e' }}>Free</span>
                </div>
                <div style={{ borderTop: '1px solid #e5e7eb', paddingTop: '12px', display: 'flex', justifyContent: 'space-between', fontWeight: 'bold', fontSize: '18px' }}>
                  <span>Total</span>
                  <span style={{ color: '#22c55e' }}>${cart.total_price?.toFixed(2)}</span>
                </div>
              </div>

              <button style={{ width: '100%', backgroundColor: '#22c55e', color: 'white', padding: '12px', borderRadius: '9999px', fontWeight: '600', border: 'none', cursor: 'pointer' }}>
                Proceed to Checkout
              </button>

              <Link
                to="/ai-chef"
                style={{ display: 'block', textAlign: 'center', color: '#22c55e', fontWeight: '500', marginTop: '16px', textDecoration: 'none' }}
              >
                🤖 Get Recipe Ideas →
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
