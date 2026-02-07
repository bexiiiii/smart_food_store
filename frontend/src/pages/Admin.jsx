import { useCallback, useEffect, useMemo, useState } from 'react';
import {
  FiBookOpen,
  FiGrid,
  FiLoader,
  FiPackage,
  FiPlus,
  FiShield,
  FiTrash2,
  FiUsers,
} from 'react-icons/fi';
import toast from 'react-hot-toast';
import api from '../api/axios';

const tabs = [
  { id: 'products', label: 'Products', icon: FiPackage, accent: '#22c55e', tint: '#dcfce7' },
  { id: 'categories', label: 'Categories', icon: FiGrid, accent: '#f59e0b', tint: '#fef3c7' },
  { id: 'users', label: 'Users', icon: FiUsers, accent: '#0ea5e9', tint: '#e0f2fe' },
  { id: 'recipes', label: 'Recipes', icon: FiBookOpen, accent: '#8b5cf6', tint: '#ede9fe' },
];

const inputStyle = {
  width: '100%',
  padding: '12px 14px',
  borderRadius: '10px',
  border: '1px solid #d1d5db',
  outline: 'none',
  fontSize: '14px',
  color: '#1f2937',
  boxSizing: 'border-box',
};

const panelStyle = {
  backgroundColor: 'white',
  borderRadius: '16px',
  border: '1px solid #e5e7eb',
  boxShadow: '0 12px 30px -24px rgba(2, 6, 23, 0.75)',
};

function createEmptyRecipeIngredient() {
  return {
    product_id: '',
    quantity: '',
    unit: 'g',
    notes: '',
  };
}

function createDefaultRecipeForm() {
  return {
    name: '',
    description: '',
    instructions: '',
    servings: '2',
    prep_time: '',
    cook_time: '',
    image_url: '',
    ingredients: [createEmptyRecipeIngredient()],
  };
}

function formatPrice(value) {
  const number = Number(value);
  if (Number.isNaN(number)) {
    return '$0.00';
  }
  return `$${number.toFixed(2)}`;
}

function getErrorMessage(error, fallback) {
  return error?.response?.data?.error || fallback;
}

export default function Admin() {
  const [activeTab, setActiveTab] = useState('products');
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [users, setUsers] = useState([]);
  const [recipes, setRecipes] = useState([]);
  const [loading, setLoading] = useState(false);

  const [productForm, setProductForm] = useState({
    name: '',
    description: '',
    price: '',
    stock: '',
    unit: 'g',
    category_id: '',
    image_url: '',
  });

  const [categoryForm, setCategoryForm] = useState({ name: '' });
  const [recipeForm, setRecipeForm] = useState(createDefaultRecipeForm());

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      if (activeTab === 'products') {
        const [prodRes, catRes] = await Promise.all([api.get('/products'), api.get('/categories')]);
        setProducts(prodRes.data);
        setCategories(catRes.data);
      } else if (activeTab === 'categories') {
        const res = await api.get('/categories');
        setCategories(res.data);
      } else if (activeTab === 'users') {
        const res = await api.get('/admin/users');
        setUsers(res.data);
      } else if (activeTab === 'recipes') {
        const [recipeRes, productRes] = await Promise.all([api.get('/recipes'), api.get('/products')]);
        setRecipes(recipeRes.data);
        setProducts(productRes.data);
      }
    } catch (error) {
      console.error('Error loading data:', error);
      toast.error(getErrorMessage(error, 'Failed to load data'));
    } finally {
      setLoading(false);
    }
  }, [activeTab]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const stats = useMemo(
    () => [
      { label: 'Products', value: products.length, icon: FiPackage, color: '#15803d', tint: '#dcfce7' },
      { label: 'Categories', value: categories.length, icon: FiGrid, color: '#b45309', tint: '#fef3c7' },
      { label: 'Users', value: users.length, icon: FiUsers, color: '#0369a1', tint: '#e0f2fe' },
      { label: 'Recipes', value: recipes.length, icon: FiBookOpen, color: '#6d28d9', tint: '#ede9fe' },
    ],
    [products.length, categories.length, users.length, recipes.length]
  );

  const activeTabConfig = tabs.find((tab) => tab.id === activeTab) || tabs[0];

  const handleProductSubmit = async (event) => {
    event.preventDefault();

    const price = Number(productForm.price);
    const stock = Number(productForm.stock);
    const categoryId = Number(productForm.category_id);

    if (Number.isNaN(price) || price <= 0) {
      toast.error('Price must be greater than 0');
      return;
    }

    if (Number.isNaN(stock) || stock < 0) {
      toast.error('Stock cannot be negative');
      return;
    }

    if (!categoryId) {
      toast.error('Please select a category');
      return;
    }

    try {
      await api.post('/admin/products', {
        ...productForm,
        price,
        stock,
        category_id: categoryId,
      });

      toast.success('Product created successfully');
      setProductForm({
        name: '',
        description: '',
        price: '',
        stock: '',
        unit: 'g',
        category_id: '',
        image_url: '',
      });
      loadData();
    } catch (error) {
      console.error('Error creating product:', error);
      toast.error(getErrorMessage(error, 'Failed to create product'));
    }
  };

  const handleDeleteProduct = async (id) => {
    if (!window.confirm('Delete this product?')) {
      return;
    }

    try {
      await api.delete(`/admin/products/${id}`);
      toast.success('Product deleted');
      loadData();
    } catch (error) {
      console.error('Error deleting product:', error);
      toast.error(getErrorMessage(error, 'Failed to delete product'));
    }
  };

  const handleCategorySubmit = async (event) => {
    event.preventDefault();
    try {
      await api.post('/admin/categories', categoryForm);
      toast.success('Category created successfully');
      setCategoryForm({ name: '' });
      loadData();
    } catch (error) {
      console.error('Error creating category:', error);
      toast.error(getErrorMessage(error, 'Failed to create category'));
    }
  };

  const handleDeleteCategory = async (id) => {
    if (!window.confirm('Delete this category?')) {
      return;
    }

    try {
      await api.delete(`/admin/categories/${id}`);
      toast.success('Category deleted');
      loadData();
    } catch (error) {
      console.error('Error deleting category:', error);
      toast.error(getErrorMessage(error, 'Failed to delete category'));
    }
  };

  const handleRecipeFieldChange = (field, value) => {
    setRecipeForm((prev) => ({ ...prev, [field]: value }));
  };

  const handleRecipeIngredientChange = (index, field, value) => {
    setRecipeForm((prev) => ({
      ...prev,
      ingredients: prev.ingredients.map((ingredient, ingredientIndex) =>
        ingredientIndex === index ? { ...ingredient, [field]: value } : ingredient
      ),
    }));
  };

  const handleAddIngredient = () => {
    setRecipeForm((prev) => ({
      ...prev,
      ingredients: [...prev.ingredients, createEmptyRecipeIngredient()],
    }));
  };

  const handleRemoveIngredient = (index) => {
    setRecipeForm((prev) => {
      if (prev.ingredients.length <= 1) {
        return prev;
      }

      return {
        ...prev,
        ingredients: prev.ingredients.filter((_, ingredientIndex) => ingredientIndex !== index),
      };
    });
  };

  const handleRecipeSubmit = async (event) => {
    event.preventDefault();

    const servings = Number(recipeForm.servings);
    const prepTime = recipeForm.prep_time === '' ? 0 : Number(recipeForm.prep_time);
    const cookTime = recipeForm.cook_time === '' ? 0 : Number(recipeForm.cook_time);

    if (!recipeForm.name.trim()) {
      toast.error('Recipe name is required');
      return;
    }

    if (!recipeForm.instructions.trim()) {
      toast.error('Instructions are required');
      return;
    }

    if (Number.isNaN(servings) || servings < 1) {
      toast.error('Servings must be at least 1');
      return;
    }

    if (Number.isNaN(prepTime) || prepTime < 0 || Number.isNaN(cookTime) || cookTime < 0) {
      toast.error('Prep/Cook time cannot be negative');
      return;
    }

    const ingredients = [];
    for (let index = 0; index < recipeForm.ingredients.length; index += 1) {
      const ingredient = recipeForm.ingredients[index];
      const productId = Number(ingredient.product_id);
      const quantity = Number(ingredient.quantity);

      if (!productId) {
        toast.error(`Select a product for ingredient ${index + 1}`);
        return;
      }

      if (Number.isNaN(quantity) || quantity <= 0) {
        toast.error(`Quantity for ingredient ${index + 1} must be greater than 0`);
        return;
      }

      if (!ingredient.unit) {
        toast.error(`Select unit for ingredient ${index + 1}`);
        return;
      }

      ingredients.push({
        product_id: productId,
        quantity,
        unit: ingredient.unit,
        notes: ingredient.notes?.trim() || '',
      });
    }

    try {
      await api.post('/admin/recipes', {
        name: recipeForm.name.trim(),
        description: recipeForm.description.trim(),
        instructions: recipeForm.instructions.trim(),
        servings,
        prep_time: prepTime,
        cook_time: cookTime,
        image_url: recipeForm.image_url.trim(),
        ingredients,
      });

      toast.success('Recipe created successfully');
      setRecipeForm(createDefaultRecipeForm());
      loadData();
    } catch (error) {
      console.error('Error creating recipe:', error);
      toast.error(getErrorMessage(error, 'Failed to create recipe'));
    }
  };

  const handleUpdateRole = async (userId, newRole) => {
    try {
      await api.patch(`/admin/users/${userId}/role`, { role: newRole });
      toast.success('Role updated');
      loadData();
    } catch (error) {
      console.error('Error updating role:', error);
      toast.error(getErrorMessage(error, 'Failed to update role'));
    }
  };

  const handleDeleteUser = async (id) => {
    if (!window.confirm('Delete this user?')) {
      return;
    }

    try {
      await api.delete(`/admin/users/${id}`);
      toast.success('User deleted');
      loadData();
    } catch (error) {
      console.error('Error deleting user:', error);
      toast.error(getErrorMessage(error, 'Failed to delete user'));
    }
  };

  const handleDeleteRecipe = async (id) => {
    if (!window.confirm('Delete this recipe?')) {
      return;
    }

    try {
      await api.delete(`/admin/recipes/${id}`);
      toast.success('Recipe deleted');
      loadData();
    } catch (error) {
      console.error('Error deleting recipe:', error);
      toast.error(getErrorMessage(error, 'Failed to delete recipe'));
    }
  };

  return (
    <div
      style={{
        minHeight: '100vh',
        background: 'linear-gradient(160deg, #f0fdf4 0%, #ffffff 45%, #fef3c7 100%)',
        padding: '32px 0 56px',
      }}
    >
      <div style={{ maxWidth: '1280px', margin: '0 auto', padding: '0 24px' }}>
        <section
          style={{
            ...panelStyle,
            padding: '28px',
            marginBottom: '24px',
            background: 'linear-gradient(120deg, #14532d 0%, #15803d 42%, #22c55e 100%)',
            color: 'white',
            border: 'none',
            boxShadow: '0 22px 35px -26px rgba(21, 128, 61, 0.8)',
          }}
        >
          <div className="admin-hero-content">
            <div>
              <p style={{ fontSize: '13px', opacity: 0.85, letterSpacing: '0.08em', textTransform: 'uppercase' }}>
                Smart Food Store
              </p>
              <h1 style={{ fontSize: '34px', fontWeight: 700, marginTop: '8px', marginBottom: '10px' }}>
                Admin Control Center
              </h1>
              <p style={{ opacity: 0.9, maxWidth: '680px' }}>
                Manage products, categories, users and recipes from one place with the same visual style as the
                customer-facing app.
              </p>
            </div>
            <div
              style={{
                width: '74px',
                height: '74px',
                borderRadius: '18px',
                backgroundColor: 'rgba(255,255,255,0.2)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                flexShrink: 0,
              }}
            >
              <FiShield style={{ width: '36px', height: '36px' }} />
            </div>
          </div>
        </section>

        <section className="admin-stats-grid" style={{ marginBottom: '24px' }}>
          {stats.map((item) => {
            const Icon = item.icon;
            return (
              <div key={item.label} style={{ ...panelStyle, padding: '18px 16px' }}>
                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                  <div>
                    <p style={{ color: '#6b7280', fontSize: '13px' }}>{item.label}</p>
                    <p style={{ fontSize: '30px', fontWeight: 700, color: '#111827', marginTop: '6px' }}>{item.value}</p>
                  </div>
                  <div
                    style={{
                      width: '42px',
                      height: '42px',
                      borderRadius: '12px',
                      backgroundColor: item.tint,
                      color: item.color,
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                    }}
                  >
                    <Icon style={{ width: '20px', height: '20px' }} />
                  </div>
                </div>
              </div>
            );
          })}
        </section>

        <section style={{ ...panelStyle, padding: '10px', marginBottom: '24px' }}>
          <div className="admin-tabs">
            {tabs.map((tab) => {
              const Icon = tab.icon;
              const isActive = activeTab === tab.id;

              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  style={{
                    display: 'inline-flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    gap: '8px',
                    border: 'none',
                    cursor: 'pointer',
                    borderRadius: '10px',
                    padding: '11px 15px',
                    fontWeight: 600,
                    fontSize: '14px',
                    backgroundColor: isActive ? tab.tint : 'transparent',
                    color: isActive ? tab.accent : '#4b5563',
                    boxShadow: isActive ? 'inset 0 0 0 1px rgba(255,255,255,0.35)' : 'none',
                    transition: 'all 0.2s ease',
                  }}
                >
                  <Icon style={{ width: '17px', height: '17px' }} />
                  {tab.label}
                </button>
              );
            })}
          </div>
        </section>

        <section style={{ marginBottom: '16px' }}>
          <h2 style={{ fontSize: '24px', fontWeight: 700, color: '#111827' }}>{activeTabConfig.label}</h2>
          <p style={{ color: '#4b5563', marginTop: '5px' }}>
            {activeTab === 'products' && 'Create products and control inventory details.'}
            {activeTab === 'categories' && 'Organize your catalog into clear categories.'}
            {activeTab === 'users' && 'Manage customer access and admin roles.'}
            {activeTab === 'recipes' && 'Review and moderate recipe content.'}
          </p>
        </section>

        {loading ? (
          <div style={{ ...panelStyle, padding: '42px', textAlign: 'center', color: '#4b5563' }}>
            <FiLoader
              style={{
                width: '28px',
                height: '28px',
                margin: '0 auto 12px',
                animation: 'admin-spin 1s linear infinite',
                color: activeTabConfig.accent,
              }}
            />
            Loading {activeTabConfig.label.toLowerCase()}...
          </div>
        ) : (
          <>
            {activeTab === 'products' && (
              <div style={{ display: 'grid', gap: '18px' }}>
                <div style={{ ...panelStyle, padding: '22px' }}>
                  <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '16px' }}>
                    <h3 style={{ fontSize: '20px', fontWeight: 700, color: '#111827' }}>Add New Product</h3>
                    <span
                      style={{
                        display: 'inline-flex',
                        alignItems: 'center',
                        gap: '6px',
                        fontSize: '13px',
                        fontWeight: 600,
                        color: '#15803d',
                        backgroundColor: '#dcfce7',
                        borderRadius: '999px',
                        padding: '6px 12px',
                      }}
                    >
                      <FiPlus style={{ width: '14px', height: '14px' }} />
                      Quick create
                    </span>
                  </div>

                  <form onSubmit={handleProductSubmit} className="admin-product-form">
                    <input
                      type="text"
                      placeholder="Product name"
                      value={productForm.name}
                      onChange={(event) => setProductForm({ ...productForm, name: event.target.value })}
                      style={inputStyle}
                      required
                    />
                    <select
                      value={productForm.category_id}
                      onChange={(event) =>
                        setProductForm({
                          ...productForm,
                          category_id: event.target.value,
                        })
                      }
                      style={inputStyle}
                      required
                    >
                      <option value="">Select category</option>
                      {categories.map((category) => (
                        <option key={category.id} value={category.id}>
                          {category.name}
                        </option>
                      ))}
                    </select>

                    <input
                      type="number"
                      step="0.01"
                      min="0"
                      placeholder="Price"
                      value={productForm.price}
                      onChange={(event) => setProductForm({ ...productForm, price: event.target.value })}
                      style={inputStyle}
                      required
                    />
                    <input
                      type="number"
                      step="0.01"
                      min="0"
                      placeholder="Stock"
                      value={productForm.stock}
                      onChange={(event) => setProductForm({ ...productForm, stock: event.target.value })}
                      style={inputStyle}
                      required
                    />

                    <select
                      value={productForm.unit}
                      onChange={(event) => setProductForm({ ...productForm, unit: event.target.value })}
                      style={inputStyle}
                    >
                      <option value="g">g</option>
                      <option value="kg">kg</option>
                      <option value="l">l</option>
                      <option value="ml">ml</option>
                      <option value="pcs">pcs</option>
                    </select>
                    <input
                      type="text"
                      placeholder="Image URL"
                      value={productForm.image_url}
                      onChange={(event) => setProductForm({ ...productForm, image_url: event.target.value })}
                      style={inputStyle}
                    />

                    <textarea
                      placeholder="Description"
                      value={productForm.description}
                      onChange={(event) => setProductForm({ ...productForm, description: event.target.value })}
                      rows={3}
                      style={{ ...inputStyle, resize: 'vertical', gridColumn: '1 / -1' }}
                    />

                    <button
                      type="submit"
                      style={{
                        gridColumn: '1 / -1',
                        backgroundColor: '#22c55e',
                        color: 'white',
                        border: 'none',
                        borderRadius: '10px',
                        padding: '12px 16px',
                        fontWeight: 600,
                        cursor: 'pointer',
                      }}
                    >
                      Add Product
                    </button>
                  </form>
                </div>

                <div style={{ ...panelStyle, padding: '6px 0 0' }}>
                  <div style={{ padding: '0 22px 12px' }}>
                    <h3 style={{ fontSize: '20px', fontWeight: 700, color: '#111827' }}>Product List</h3>
                  </div>

                  {products.length === 0 ? (
                    <div style={{ padding: '8px 22px 22px', color: '#6b7280' }}>No products yet.</div>
                  ) : (
                    <div className="admin-table-wrap">
                      <table style={{ width: '100%', minWidth: '860px', borderCollapse: 'collapse' }}>
                        <thead>
                          <tr style={{ backgroundColor: '#f8fafc', borderTop: '1px solid #e5e7eb', borderBottom: '1px solid #e5e7eb' }}>
                            <th style={{ textAlign: 'left', padding: '12px 22px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                              Product
                            </th>
                            <th style={{ textAlign: 'left', padding: '12px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                              Category
                            </th>
                            <th style={{ textAlign: 'left', padding: '12px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                              Price
                            </th>
                            <th style={{ textAlign: 'left', padding: '12px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                              Stock
                            </th>
                            <th style={{ textAlign: 'right', padding: '12px 22px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                              Actions
                            </th>
                          </tr>
                        </thead>
                        <tbody>
                          {products.map((product) => (
                            <tr key={product.id} style={{ borderBottom: '1px solid #e5e7eb' }}>
                              <td style={{ padding: '14px 22px' }}>
                                <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                                  <div
                                    style={{
                                      width: '38px',
                                      height: '38px',
                                      borderRadius: '10px',
                                      backgroundColor: '#f0fdf4',
                                      display: 'flex',
                                      alignItems: 'center',
                                      justifyContent: 'center',
                                      color: '#15803d',
                                      fontSize: '18px',
                                      flexShrink: 0,
                                    }}
                                  >
                                    {product.image_url ? '🖼️' : '📦'}
                                  </div>
                                  <div>
                                    <p style={{ color: '#111827', fontWeight: 600 }}>{product.name}</p>
                                    <p style={{ color: '#6b7280', fontSize: '13px' }}>{product.id}</p>
                                  </div>
                                </div>
                              </td>
                              <td style={{ padding: '14px 12px', color: '#334155' }}>
                                {product.category?.name ||
                                  categories.find((category) => Number(category.id) === Number(product.category_id))
                                    ?.name ||
                                  'Uncategorized'}
                              </td>
                              <td style={{ padding: '14px 12px', color: '#111827', fontWeight: 600 }}>{formatPrice(product.price)}</td>
                              <td style={{ padding: '14px 12px', color: '#475569' }}>
                                {Number(product.stock || 0).toFixed(2)} {product.unit}
                              </td>
                              <td style={{ padding: '14px 22px', textAlign: 'right' }}>
                                <button
                                  onClick={() => handleDeleteProduct(product.id)}
                                  style={{
                                    display: 'inline-flex',
                                    alignItems: 'center',
                                    gap: '8px',
                                    color: '#dc2626',
                                    fontWeight: 600,
                                    border: 'none',
                                    background: 'none',
                                    cursor: 'pointer',
                                  }}
                                >
                                  <FiTrash2 style={{ width: '15px', height: '15px' }} />
                                  Delete
                                </button>
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  )}
                </div>
              </div>
            )}

            {activeTab === 'categories' && (
              <div style={{ display: 'grid', gap: '18px' }}>
                <div style={{ ...panelStyle, padding: '22px' }}>
                  <h3 style={{ fontSize: '20px', fontWeight: 700, color: '#111827', marginBottom: '16px' }}>Add Category</h3>
                  <form onSubmit={handleCategorySubmit} className="admin-category-form">
                    <input
                      type="text"
                      placeholder="Category name"
                      value={categoryForm.name}
                      onChange={(event) => setCategoryForm({ name: event.target.value })}
                      style={inputStyle}
                      required
                    />
                    <button
                      type="submit"
                      style={{
                        minWidth: '160px',
                        backgroundColor: '#f59e0b',
                        color: 'white',
                        border: 'none',
                        borderRadius: '10px',
                        padding: '12px 16px',
                        fontWeight: 600,
                        cursor: 'pointer',
                      }}
                    >
                      Add Category
                    </button>
                  </form>
                </div>

                <div style={{ ...panelStyle, padding: '6px 0 0' }}>
                  <div style={{ padding: '0 22px 12px' }}>
                    <h3 style={{ fontSize: '20px', fontWeight: 700, color: '#111827' }}>Category List</h3>
                  </div>

                  {categories.length === 0 ? (
                    <div style={{ padding: '8px 22px 22px', color: '#6b7280' }}>No categories yet.</div>
                  ) : (
                    <div className="admin-table-wrap">
                      <table style={{ width: '100%', minWidth: '620px', borderCollapse: 'collapse' }}>
                        <thead>
                          <tr style={{ backgroundColor: '#f8fafc', borderTop: '1px solid #e5e7eb', borderBottom: '1px solid #e5e7eb' }}>
                            <th style={{ textAlign: 'left', padding: '12px 22px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                              ID
                            </th>
                            <th style={{ textAlign: 'left', padding: '12px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                              Name
                            </th>
                            <th style={{ textAlign: 'right', padding: '12px 22px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                              Actions
                            </th>
                          </tr>
                        </thead>
                        <tbody>
                          {categories.map((category) => (
                            <tr key={category.id} style={{ borderBottom: '1px solid #e5e7eb' }}>
                              <td style={{ padding: '14px 22px', color: '#334155' }}>{category.id}</td>
                              <td style={{ padding: '14px 12px', color: '#111827', fontWeight: 600 }}>{category.name}</td>
                              <td style={{ padding: '14px 22px', textAlign: 'right' }}>
                                <button
                                  onClick={() => handleDeleteCategory(category.id)}
                                  style={{
                                    display: 'inline-flex',
                                    alignItems: 'center',
                                    gap: '8px',
                                    color: '#dc2626',
                                    fontWeight: 600,
                                    border: 'none',
                                    background: 'none',
                                    cursor: 'pointer',
                                  }}
                                >
                                  <FiTrash2 style={{ width: '15px', height: '15px' }} />
                                  Delete
                                </button>
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  )}
                </div>
              </div>
            )}

            {activeTab === 'users' && (
              <div style={{ ...panelStyle, padding: '6px 0 0' }}>
                <div style={{ padding: '0 22px 12px' }}>
                  <h3 style={{ fontSize: '20px', fontWeight: 700, color: '#111827' }}>User Management</h3>
                </div>

                {users.length === 0 ? (
                  <div style={{ padding: '8px 22px 22px', color: '#6b7280' }}>No users found.</div>
                ) : (
                  <div className="admin-table-wrap">
                    <table style={{ width: '100%', minWidth: '780px', borderCollapse: 'collapse' }}>
                      <thead>
                        <tr style={{ backgroundColor: '#f8fafc', borderTop: '1px solid #e5e7eb', borderBottom: '1px solid #e5e7eb' }}>
                          <th style={{ textAlign: 'left', padding: '12px 22px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                            User
                          </th>
                          <th style={{ textAlign: 'left', padding: '12px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                            Email
                          </th>
                          <th style={{ textAlign: 'left', padding: '12px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                            Role
                          </th>
                          <th style={{ textAlign: 'right', padding: '12px 22px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                            Actions
                          </th>
                        </tr>
                      </thead>
                      <tbody>
                        {users.map((user) => (
                          <tr key={user.id} style={{ borderBottom: '1px solid #e5e7eb' }}>
                            <td style={{ padding: '14px 22px' }}>
                              <p style={{ color: '#111827', fontWeight: 600 }}>{user.name}</p>
                              <p style={{ color: '#6b7280', fontSize: '13px' }}>ID: {user.id}</p>
                            </td>
                            <td style={{ padding: '14px 12px', color: '#334155' }}>{user.email}</td>
                            <td style={{ padding: '14px 12px' }}>
                              <select
                                value={user.role}
                                onChange={(event) => handleUpdateRole(user.id, event.target.value)}
                                style={{
                                  ...inputStyle,
                                  padding: '8px 10px',
                                  maxWidth: '150px',
                                  borderColor: user.role === 'admin' ? '#93c5fd' : '#d1d5db',
                                  backgroundColor: user.role === 'admin' ? '#eff6ff' : '#fff',
                                }}
                              >
                                <option value="user">User</option>
                                <option value="admin">Admin</option>
                              </select>
                            </td>
                            <td style={{ padding: '14px 22px', textAlign: 'right' }}>
                              <button
                                onClick={() => handleDeleteUser(user.id)}
                                style={{
                                  display: 'inline-flex',
                                  alignItems: 'center',
                                  gap: '8px',
                                  color: '#dc2626',
                                  fontWeight: 600,
                                  border: 'none',
                                  background: 'none',
                                  cursor: 'pointer',
                                }}
                              >
                                <FiTrash2 style={{ width: '15px', height: '15px' }} />
                                Delete
                              </button>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                )}
              </div>
            )}

            {activeTab === 'recipes' && (
              <div style={{ display: 'grid', gap: '18px' }}>
                <div style={{ ...panelStyle, padding: '22px' }}>
                  <h3 style={{ fontSize: '20px', fontWeight: 700, color: '#111827', marginBottom: '16px' }}>Add New Recipe</h3>
                  {products.length === 0 ? (
                    <p style={{ color: '#6b7280' }}>Create products first, then you can add recipe ingredients.</p>
                  ) : (
                    <form onSubmit={handleRecipeSubmit} className="admin-recipe-form">
                      <input
                        type="text"
                        placeholder="Recipe name"
                        value={recipeForm.name}
                        onChange={(event) => handleRecipeFieldChange('name', event.target.value)}
                        style={inputStyle}
                        required
                      />
                      <input
                        type="number"
                        min="1"
                        placeholder="Servings"
                        value={recipeForm.servings}
                        onChange={(event) => handleRecipeFieldChange('servings', event.target.value)}
                        style={inputStyle}
                        required
                      />

                      <input
                        type="number"
                        min="0"
                        placeholder="Prep time (min)"
                        value={recipeForm.prep_time}
                        onChange={(event) => handleRecipeFieldChange('prep_time', event.target.value)}
                        style={inputStyle}
                      />
                      <input
                        type="number"
                        min="0"
                        placeholder="Cook time (min)"
                        value={recipeForm.cook_time}
                        onChange={(event) => handleRecipeFieldChange('cook_time', event.target.value)}
                        style={inputStyle}
                      />

                      <input
                        type="text"
                        placeholder="Image URL"
                        value={recipeForm.image_url}
                        onChange={(event) => handleRecipeFieldChange('image_url', event.target.value)}
                        style={{ ...inputStyle, gridColumn: '1 / -1' }}
                      />

                      <textarea
                        placeholder="Short description"
                        value={recipeForm.description}
                        onChange={(event) => handleRecipeFieldChange('description', event.target.value)}
                        rows={3}
                        style={{ ...inputStyle, resize: 'vertical', gridColumn: '1 / -1' }}
                      />

                      <textarea
                        placeholder="Instructions (one step per line)"
                        value={recipeForm.instructions}
                        onChange={(event) => handleRecipeFieldChange('instructions', event.target.value)}
                        rows={5}
                        style={{ ...inputStyle, resize: 'vertical', gridColumn: '1 / -1' }}
                        required
                      />

                      <div className="admin-recipe-form-full">
                        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '10px' }}>
                          <h4 style={{ color: '#111827', fontWeight: 700, fontSize: '16px' }}>Ingredients</h4>
                          <button
                            type="button"
                            onClick={handleAddIngredient}
                            style={{
                              border: 'none',
                              backgroundColor: '#ede9fe',
                              color: '#6d28d9',
                              borderRadius: '999px',
                              padding: '8px 12px',
                              fontWeight: 600,
                              cursor: 'pointer',
                              display: 'inline-flex',
                              alignItems: 'center',
                              gap: '6px',
                            }}
                          >
                            <FiPlus style={{ width: '14px', height: '14px' }} />
                            Ingredient
                          </button>
                        </div>

                        {recipeForm.ingredients.map((ingredient, index) => (
                          <div key={`ingredient-${index}`} className="admin-ingredient-grid">
                            <select
                              value={ingredient.product_id}
                              onChange={(event) => handleRecipeIngredientChange(index, 'product_id', event.target.value)}
                              style={inputStyle}
                              required
                            >
                              <option value="">Select product</option>
                              {products.map((product) => (
                                <option key={product.id} value={product.id}>
                                  {product.name}
                                </option>
                              ))}
                            </select>

                            <input
                              type="number"
                              min="0.01"
                              step="0.01"
                              placeholder="Qty"
                              value={ingredient.quantity}
                              onChange={(event) => handleRecipeIngredientChange(index, 'quantity', event.target.value)}
                              style={inputStyle}
                              required
                            />

                            <select
                              value={ingredient.unit}
                              onChange={(event) => handleRecipeIngredientChange(index, 'unit', event.target.value)}
                              style={inputStyle}
                            >
                              <option value="g">g</option>
                              <option value="kg">kg</option>
                              <option value="l">l</option>
                              <option value="ml">ml</option>
                              <option value="pcs">pcs</option>
                            </select>

                            <input
                              type="text"
                              placeholder="Notes (optional)"
                              value={ingredient.notes}
                              onChange={(event) => handleRecipeIngredientChange(index, 'notes', event.target.value)}
                              style={inputStyle}
                            />

                            <button
                              type="button"
                              onClick={() => handleRemoveIngredient(index)}
                              disabled={recipeForm.ingredients.length <= 1}
                              style={{
                                width: '42px',
                                height: '42px',
                                borderRadius: '10px',
                                border: 'none',
                                backgroundColor: recipeForm.ingredients.length <= 1 ? '#e5e7eb' : '#fee2e2',
                                color: recipeForm.ingredients.length <= 1 ? '#94a3b8' : '#dc2626',
                                cursor: recipeForm.ingredients.length <= 1 ? 'not-allowed' : 'pointer',
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                              }}
                            >
                              <FiTrash2 style={{ width: '16px', height: '16px' }} />
                            </button>
                          </div>
                        ))}
                      </div>

                      <button
                        type="submit"
                        style={{
                          gridColumn: '1 / -1',
                          backgroundColor: '#8b5cf6',
                          color: 'white',
                          border: 'none',
                          borderRadius: '10px',
                          padding: '12px 16px',
                          fontWeight: 600,
                          cursor: 'pointer',
                        }}
                      >
                        Add Recipe
                      </button>
                    </form>
                  )}
                </div>

                <div style={{ ...panelStyle, padding: '6px 0 0' }}>
                  <div style={{ padding: '0 22px 12px' }}>
                    <h3 style={{ fontSize: '20px', fontWeight: 700, color: '#111827' }}>Recipe Management</h3>
                  </div>

                  {recipes.length === 0 ? (
                    <div style={{ padding: '8px 22px 22px', color: '#6b7280' }}>No recipes found.</div>
                  ) : (
                    <div className="admin-table-wrap">
                      <table style={{ width: '100%', minWidth: '860px', borderCollapse: 'collapse' }}>
                        <thead>
                          <tr style={{ backgroundColor: '#f8fafc', borderTop: '1px solid #e5e7eb', borderBottom: '1px solid #e5e7eb' }}>
                            <th style={{ textAlign: 'left', padding: '12px 22px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                              Recipe
                            </th>
                            <th style={{ textAlign: 'left', padding: '12px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                              Servings
                            </th>
                            <th style={{ textAlign: 'left', padding: '12px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                              Time
                            </th>
                            <th style={{ textAlign: 'left', padding: '12px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                              Ingredients
                            </th>
                            <th style={{ textAlign: 'right', padding: '12px 22px', fontSize: '12px', color: '#64748b', textTransform: 'uppercase' }}>
                              Actions
                            </th>
                          </tr>
                        </thead>
                        <tbody>
                          {recipes.map((recipe) => (
                            <tr key={recipe.id} style={{ borderBottom: '1px solid #e5e7eb' }}>
                              <td style={{ padding: '14px 22px' }}>
                                <p style={{ color: '#111827', fontWeight: 600 }}>{recipe.name}</p>
                                <p style={{ color: '#6b7280', fontSize: '13px' }}>ID: {recipe.id}</p>
                              </td>
                              <td style={{ padding: '14px 12px', color: '#334155' }}>{recipe.servings || 0}</td>
                              <td style={{ padding: '14px 12px', color: '#334155' }}>
                                {(Number(recipe.prep_time || 0) + Number(recipe.cook_time || 0)) || 0} min
                              </td>
                              <td style={{ padding: '14px 12px', color: '#334155' }}>{recipe.ingredients?.length || 0}</td>
                              <td style={{ padding: '14px 22px', textAlign: 'right' }}>
                                <button
                                  onClick={() => handleDeleteRecipe(recipe.id)}
                                  style={{
                                    display: 'inline-flex',
                                    alignItems: 'center',
                                    gap: '8px',
                                    color: '#dc2626',
                                    fontWeight: 600,
                                    border: 'none',
                                    background: 'none',
                                    cursor: 'pointer',
                                  }}
                                >
                                  <FiTrash2 style={{ width: '15px', height: '15px' }} />
                                  Delete
                                </button>
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  )}
                </div>
              </div>
            )}
          </>
        )}
      </div>

      <style>{`
        @keyframes admin-spin {
          to {
            transform: rotate(360deg);
          }
        }

        .admin-hero-content {
          display: flex;
          align-items: center;
          justify-content: space-between;
          gap: 18px;
        }

        .admin-stats-grid {
          display: grid;
          grid-template-columns: repeat(4, minmax(0, 1fr));
          gap: 12px;
        }

        .admin-tabs {
          display: grid;
          grid-template-columns: repeat(4, minmax(0, 1fr));
          gap: 8px;
        }

        .admin-product-form {
          display: grid;
          grid-template-columns: repeat(2, minmax(0, 1fr));
          gap: 14px;
        }

        .admin-category-form {
          display: flex;
          gap: 12px;
        }

        .admin-recipe-form {
          display: grid;
          grid-template-columns: repeat(2, minmax(0, 1fr));
          gap: 14px;
        }

        .admin-recipe-form-full {
          grid-column: 1 / -1;
        }

        .admin-ingredient-grid {
          display: grid;
          grid-template-columns: minmax(190px, 2fr) minmax(90px, 1fr) minmax(90px, 0.8fr) minmax(180px, 1.3fr) auto;
          gap: 10px;
          align-items: center;
          margin-bottom: 10px;
        }

        .admin-table-wrap {
          width: 100%;
          overflow-x: auto;
        }

        @media (max-width: 1024px) {
          .admin-stats-grid {
            grid-template-columns: repeat(2, minmax(0, 1fr));
          }

          .admin-tabs {
            grid-template-columns: repeat(2, minmax(0, 1fr));
          }
        }

        @media (max-width: 768px) {
          .admin-hero-content {
            flex-direction: column;
            align-items: flex-start;
          }

          .admin-product-form {
            grid-template-columns: 1fr;
          }

          .admin-recipe-form {
            grid-template-columns: 1fr;
          }

          .admin-product-form textarea,
          .admin-product-form button {
            grid-column: auto !important;
          }

          .admin-recipe-form textarea,
          .admin-recipe-form button {
            grid-column: auto !important;
          }

          .admin-ingredient-grid {
            grid-template-columns: 1fr;
          }

          .admin-ingredient-grid button {
            width: 100%;
          }

          .admin-category-form {
            flex-direction: column;
          }

          .admin-category-form button {
            width: 100%;
          }
        }

        @media (max-width: 560px) {
          .admin-stats-grid {
            grid-template-columns: 1fr;
          }

          .admin-tabs {
            grid-template-columns: 1fr;
          }
        }
      `}</style>
    </div>
  );
}
