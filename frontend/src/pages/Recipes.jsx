import { useState, useEffect } from 'react';
import { recipesAPI } from '../api';
import { FiClock, FiUsers, FiChevronRight } from 'react-icons/fi';
import toast from 'react-hot-toast';

function getRecipeTotalTime(recipe) {
  const prepTime = Number(recipe?.prep_time || 0);
  const cookTime = Number(recipe?.cook_time || 0);
  const totalTime = prepTime + cookTime;

  if (totalTime > 0) {
    return totalTime;
  }

  const legacyCookTime = Number(recipe?.cooking_time || 0);
  return legacyCookTime > 0 ? legacyCookTime : 0;
}

function getIngredientName(ingredient, index) {
  return ingredient?.product?.name || ingredient?.product_name || ingredient?.name || `Ingredient ${index + 1}`;
}

export default function Recipes() {
  const [recipes, setRecipes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selectedRecipe, setSelectedRecipe] = useState(null);

  useEffect(() => {
    loadRecipes();
  }, []);

  const loadRecipes = async () => {
    try {
      setLoading(true);
      const res = await recipesAPI.getAll();
      setRecipes(res.data);
    } catch {
      toast.error('Failed to load recipes');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ minHeight: '100vh', backgroundColor: '#f9fafb', padding: '32px 0' }}>
      <div style={{ maxWidth: '1280px', margin: '0 auto', padding: '0 24px' }}>
        {/* Header */}
        <div style={{ textAlign: 'center', marginBottom: '48px' }}>
          <h1 style={{ fontSize: '2.25rem', fontWeight: 'bold', color: '#1f2937', marginBottom: '16px' }}>Recipes</h1>
          <p style={{ color: '#6b7280', maxWidth: '672px', margin: '0 auto' }}>
            Discover delicious recipes and find all the ingredients you need in our store
          </p>
        </div>

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
        ) : recipes.length > 0 ? (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(320px, 1fr))', gap: '32px' }}>
            {recipes.map((recipe) => (
              <div
                key={recipe.id}
                style={{
                  backgroundColor: 'white',
                  borderRadius: '16px',
                  boxShadow: '0 4px 6px rgba(0,0,0,0.1)',
                  overflow: 'hidden',
                  cursor: 'pointer',
                  transition: 'box-shadow 0.2s'
                }}
                onClick={() => setSelectedRecipe(recipe)}
              >
                <div style={{
                  height: '192px',
                  background: 'linear-gradient(to bottom right, #4ade80, #22c55e)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center'
                }}>
                  <span style={{ fontSize: '4rem' }}>🍽️</span>
                </div>
                <div style={{ padding: '24px' }}>
                  <h3 style={{ fontSize: '1.25rem', fontWeight: '600', color: '#1f2937', marginBottom: '8px' }}>
                    {recipe.name}
                  </h3>
                  <p style={{ color: '#6b7280', fontSize: '0.875rem', marginBottom: '16px', overflow: 'hidden', display: '-webkit-box', WebkitLineClamp: 2, WebkitBoxOrient: 'vertical' }}>
                    {recipe.description}
                  </p>
                  <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', fontSize: '0.875rem', color: '#6b7280' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
                      <span style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                        <FiClock style={{ width: '16px', height: '16px' }} />
                        {getRecipeTotalTime(recipe)} min
                      </span>
                      <span style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                        <FiUsers style={{ width: '16px', height: '16px' }} />
                        {recipe.servings} servings
                      </span>
                    </div>
                    <FiChevronRight style={{ width: '20px', height: '20px', color: '#22c55e' }} />
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div style={{ textAlign: 'center', padding: '64px 0', backgroundColor: 'white', borderRadius: '12px' }}>
            <span style={{ fontSize: '4rem', display: 'block', marginBottom: '16px' }}>📖</span>
            <p style={{ color: '#6b7280', fontSize: '1.125rem' }}>No recipes available yet</p>
            <p style={{ color: '#9ca3af', marginTop: '8px' }}>Check back later for new recipes!</p>
          </div>
        )}
      </div>

      {/* Recipe Modal */}
      {selectedRecipe && (
        <div style={{
          position: 'fixed',
          inset: 0,
          backgroundColor: 'rgba(0,0,0,0.5)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          padding: '16px',
          zIndex: 50
        }}>
          <div style={{
            backgroundColor: 'white',
            borderRadius: '16px',
            maxWidth: '672px',
            width: '100%',
            maxHeight: '90vh',
            overflowY: 'auto'
          }}>
            <div style={{
              height: '192px',
              background: 'linear-gradient(to bottom right, #4ade80, #22c55e)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              position: 'relative'
            }}>
              <span style={{ fontSize: '6rem' }}>🍽️</span>
              <button
                onClick={() => setSelectedRecipe(null)}
                style={{
                  position: 'absolute',
                  top: '16px',
                  right: '16px',
                  width: '40px',
                  height: '40px',
                  backgroundColor: 'rgba(255,255,255,0.2)',
                  borderRadius: '50%',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  color: 'white',
                  border: 'none',
                  cursor: 'pointer',
                  fontSize: '18px'
                }}
              >
                ✕
              </button>
            </div>
            <div style={{ padding: '24px' }}>
              <h2 style={{ fontSize: '1.5rem', fontWeight: 'bold', color: '#1f2937', marginBottom: '8px' }}>
                {selectedRecipe.name}
              </h2>
              <p style={{ color: '#6b7280', marginBottom: '16px' }}>{selectedRecipe.description}</p>
              
              <div style={{ display: 'flex', gap: '16px', marginBottom: '24px' }}>
                <span style={{ display: 'flex', alignItems: 'center', gap: '8px', backgroundColor: '#f3f4f6', padding: '8px 16px', borderRadius: '8px' }}>
                  <FiClock style={{ color: '#22c55e' }} />
                  {getRecipeTotalTime(selectedRecipe)} min
                </span>
                <span style={{ display: 'flex', alignItems: 'center', gap: '8px', backgroundColor: '#f3f4f6', padding: '8px 16px', borderRadius: '8px' }}>
                  <FiUsers style={{ color: '#22c55e' }} />
                  {selectedRecipe.servings} servings
                </span>
              </div>

              {selectedRecipe.ingredients && (
                <div style={{ marginBottom: '24px' }}>
                  <h3 style={{ fontWeight: '600', color: '#1f2937', marginBottom: '12px' }}>Ingredients</h3>
                  <ul style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                    {selectedRecipe.ingredients.map((ing, idx) => (
                      <li key={idx} style={{ display: 'flex', alignItems: 'center', gap: '8px', color: '#6b7280' }}>
                        <span style={{ width: '8px', height: '8px', backgroundColor: '#22c55e', borderRadius: '50%' }}></span>
                        {getIngredientName(ing, idx)} - {ing.quantity} {ing.unit}
                      </li>
                    ))}
                  </ul>
                </div>
              )}

              {selectedRecipe.instructions && (
                <div>
                  <h3 style={{ fontWeight: '600', color: '#1f2937', marginBottom: '12px' }}>Instructions</h3>
                  <ol style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                    {selectedRecipe.instructions.split('\n').map((step, idx) => (
                      <li key={idx} style={{ display: 'flex', gap: '12px', color: '#6b7280' }}>
                        <span style={{
                          width: '24px',
                          height: '24px',
                          backgroundColor: '#dcfce7',
                          color: '#15803d',
                          borderRadius: '50%',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          fontSize: '0.875rem',
                          fontWeight: '500',
                          flexShrink: 0
                        }}>
                          {idx + 1}
                        </span>
                        {step}
                      </li>
                    ))}
                  </ol>
                </div>
              )}

              <button
                onClick={() => setSelectedRecipe(null)}
                style={{
                  marginTop: '24px',
                  width: '100%',
                  backgroundColor: '#22c55e',
                  color: 'white',
                  padding: '12px',
                  borderRadius: '12px',
                  border: 'none',
                  cursor: 'pointer',
                  fontWeight: '500',
                  fontSize: '1rem'
                }}
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}

      <style>{`
        @keyframes spin {
          to { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  );
}
