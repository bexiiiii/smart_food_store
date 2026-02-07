import { FiHeart, FiEye, FiStar } from 'react-icons/fi';
import { useAuthStore } from '../store';
import { cartAPI } from '../api';
import toast from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';

export default function ProductCard({ product, onCartUpdate }) {
  const { isAuthenticated } = useAuthStore();
  const navigate = useNavigate();

  const handleAddToCart = async () => {
    if (!isAuthenticated) {
      toast.error('Please login to add items to cart');
      navigate('/login');
      return;
    }

    try {
      await cartAPI.addItem(product.id, 1);
      toast.success(`${product.name} added to cart!`);
      if (onCartUpdate) onCartUpdate();
    } catch (error) {
      toast.error(error.response?.data?.error || 'Failed to add to cart');
    }
  };

  // Generate a random rating for demo
  const rating = (Math.random() * 2 + 3).toFixed(1);
  const reviewCount = Math.floor(Math.random() * 10);

  // Check if on sale (random for demo)
  const onSale = Math.random() > 0.7;
  const originalPrice = onSale ? (product.price * 1.2).toFixed(2) : null;

  return (
    <div className="bg-white rounded-2xl shadow-sm hover:shadow-lg transition-shadow duration-300 overflow-hidden group">
      {/* Image Container */}
      <div className={`relative h-48 ${onSale ? 'bg-white' : 'bg-green-50'} flex items-center justify-center p-4`}>
        {onSale && (
          <span className="absolute top-3 left-3 bg-primary-500 text-white text-xs font-bold px-3 py-1 rounded-full">
            Sale!
          </span>
        )}
        
        {/* Product Image Placeholder */}
        <div className="w-32 h-32 rounded-full bg-gray-100 flex items-center justify-center overflow-hidden">
          {product.image_url ? (
            <img src={product.image_url} alt={product.name} className="w-full h-full object-cover" />
          ) : (
            <span className="text-5xl">
              {getCategoryEmoji(product.category?.name)}
            </span>
          )}
        </div>

        {/* Hover Actions */}
        <div className="absolute top-3 right-3 flex flex-col space-y-2 opacity-0 group-hover:opacity-100 transition-opacity">
          <button className="p-2 bg-white rounded-full shadow hover:bg-gray-50">
            <FiHeart className="w-4 h-4 text-gray-600" />
          </button>
          <button className="p-2 bg-white rounded-full shadow hover:bg-gray-50">
            <FiEye className="w-4 h-4 text-gray-600" />
          </button>
        </div>
      </div>

      {/* Content */}
      <div className="p-4">
        {/* Category */}
        <p className="text-xs text-gray-500 mb-1">
          {product.category?.name || 'General'}
        </p>

        {/* Name */}
        <h3 className="font-semibold text-gray-800 mb-2 truncate">
          {product.name}
        </h3>

        {/* Price & Rating Row */}
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2">
            {onSale && (
              <span className="text-gray-400 line-through text-sm">${originalPrice}</span>
            )}
            <span className="text-primary-600 font-bold text-lg">${product.price.toFixed(2)}</span>
          </div>
          
          <div className="flex items-center text-sm">
            <FiStar className="w-4 h-4 text-yellow-400 fill-current" />
            <span className="text-gray-600 ml-1">{rating}</span>
            {reviewCount > 0 && (
              <span className="text-gray-400 ml-1">({reviewCount})</span>
            )}
          </div>
        </div>

        {/* Add to Cart Button */}
        <div className="flex items-center gap-2">
          <button
            onClick={handleAddToCart}
            className="flex-1 bg-primary-500 hover:bg-primary-600 text-white py-2 px-4 rounded-full text-sm font-medium transition-colors flex items-center justify-center gap-2"
          >
            <span>Add to Cart →</span>
          </button>
        </div>

        {/* Stock Info */}
        <p className="text-xs text-gray-400 mt-2 text-center">
          {product.stock > 0 ? `${product.stock} ${product.unit} in stock` : 'Out of stock'}
        </p>
      </div>
    </div>
  );
}

function getCategoryEmoji(category) {
  const emojis = {
    'Vegetables': '🥬',
    'Fruits': '🍎',
    'Meat': '🥩',
    'Dairy': '🥛',
    'Grains': '🌾',
    'Spices': '🌶️',
    'Beverages': '🧃',
    'Bakery': '🍞',
    'Seafood': '🐟',
    'Frozen': '🧊',
    'Breakfast Cereals': '🥣',
    'Snacks': '🍿',
  };
  return emojis[category] || '🛒';
}
