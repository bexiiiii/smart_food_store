package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bexiiiii/smart_food_store/internal/config"
	"github.com/bexiiiii/smart_food_store/internal/models"
	"github.com/bexiiiii/smart_food_store/internal/repository"
)

type AIService struct {
	productRepo *repository.ProductRepository
	recipeRepo  *repository.RecipeRepository
	config      *config.Config
	apiKey      string
}

// Gemini REST API structures
type GeminiRequest struct {
	Contents         []GeminiContent        `json:"contents"`
	GenerationConfig GeminiGenerationConfig `json:"generationConfig,omitempty"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiGenerationConfig struct {
	Temperature float64 `json:"temperature,omitempty"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error,omitempty"`
}

func NewAIService(productRepo *repository.ProductRepository, recipeRepo *repository.RecipeRepository, cfg *config.Config) (*AIService, error) {
	return &AIService{
		productRepo: productRepo,
		recipeRepo:  recipeRepo,
		config:      cfg,
		apiKey:      cfg.GeminiAPIKey,
	}, nil
}

// callGeminiAPI - вызов Gemini REST API напрямую
func (s *AIService) callGeminiAPI(ctx context.Context, prompt string) (string, error) {
	if s.apiKey == "" || s.apiKey == "your-gemini-api-key-here" {
		return "", errors.New("Gemini API not configured. Please set GEMINI_API_KEY in .env file")
	}


	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/gemini-2.5-flash:generateContent?key=%s", s.apiKey)

	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: GeminiGenerationConfig{
			Temperature: 0.7,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Gemini API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %v (body: %s)", err, string(body))
	}

	if geminiResp.Error != nil {
		return "", fmt.Errorf("Gemini API error: %s (code: %d)", geminiResp.Error.Message, geminiResp.Error.Code)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("empty response from Gemini API")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

// GetIngredientsForDish - вводишь название блюда, AI подбирает продукты из магазина
func (s *AIService) GetIngredientsForDish(ctx context.Context, dishName string, servings int) (*models.DishIngredientsResponse, error) {
	// Get all available products from store
	products, err := s.productRepo.GetAllWithStock()
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %v", err)
	}

	// Build product list for AI
	productList := s.buildProductListString(products)

	prompt := fmt.Sprintf(`You are a cooking assistant for a food store. A user wants to make "%s" for %d servings.

Here are the available products in our store:
%s

Based on these available products, please provide ingredients needed with the following JSON format (ONLY JSON, no other text):
{
  "dish_name": "%s",
  "description": "Brief description of the dish",
  "servings": %d,
  "required_ingredients": [
    {"name": "Ingredient Name", "quantity": 500, "unit": "g"}
  ],
  "matched_products": [
    {"id": 1, "name": "Product Name", "price": 5.99, "unit": "kg"}
  ],
  "cooking_tips": "Brief cooking tips"
}

Important rules:
1. In matched_products, ONLY include products from the store list above
2. Use the exact product ID from the store list
3. required_ingredients lists what you need for the recipe
4. matched_products lists which store products to buy (with real IDs and prices from the list)
5. Units: "g", "kg", "l", "ml", "pcs"
6. Return ONLY valid JSON, no markdown code blocks`, dishName, servings, productList, dishName, servings)

	responseText, err := s.callGeminiAPI(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// Parse JSON response
	var response models.DishIngredientsResponse
	cleanedResponse := s.cleanJSONResponse(responseText)
	if err := json.Unmarshal([]byte(cleanedResponse), &response); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %v (response: %s)", err, cleanedResponse[:min(200, len(cleanedResponse))])
	}

	// Verify and enrich matched products with real data
	var validProducts []models.MatchedProduct
	var totalPrice float64
	for _, mp := range response.MatchedProducts {
		product, err := s.productRepo.GetByID(mp.ID)
		if err == nil && product != nil {
			validProducts = append(validProducts, models.MatchedProduct{
				ID:    product.ID,
				Name:  product.Name,
				Price: product.Price,
				Unit:  string(product.Unit),
			})
			totalPrice += product.Price
		}
	} 
	response.MatchedProducts = validProducts
	response.TotalPrice = totalPrice

	return &response, nil
}

// GetRecipesFromCart - на основе продуктов в корзине AI предлагает что можно приготовить
func (s *AIService) GetRecipesFromCart(ctx context.Context, productIDs []uint) ([]models.AIRecipeSuggestion, error) {
	// Get products from cart
	products, err := s.productRepo.GetByIDs(productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %v", err)
	}

	if len(products) == 0 {
		return nil, errors.New("no products found in cart")
	}

	// Build product list
	var productNames []string
	productMap := make(map[uint]models.Product)
	for _, p := range products {
		productNames = append(productNames, fmt.Sprintf("- %s (ID: %d, Unit: %s, Price: %.2f)", p.Name, p.ID, p.Unit, p.Price))
		productMap[p.ID] = p
	}

	prompt := fmt.Sprintf(`You are a creative cooking assistant. A user has the following products in their cart:

%s

Please suggest 3 different recipes they can make with these ingredients. Return ONLY a JSON array (no other text):
[
  {
    "name": "Recipe Name",
    "description": "Brief description",
    "instructions": "Step-by-step cooking instructions",
    "prep_time": 15,
    "cook_time": 30,
    "servings": 4,
    "ingredients": [
      {
        "product_id": 1,
        "product_name": "Product Name",
        "quantity": 200,
        "unit": "g"
      }
    ],
    "confidence": 0.9
  }
]

Rules:
1. ONLY use products from the cart (use exact product_id)
2. Units: "g", "kg", "l", "ml", "pcs"
3. Order recipes by how well they use the available ingredients (best first)
4. confidence reflects how complete the recipe is with available ingredients
5. Return ONLY valid JSON array`, strings.Join(productNames, "\n"))

	responseText, err := s.callGeminiAPI(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// Parse JSON response
	var suggestions []models.AIRecipeSuggestion
	cleanedResponse := s.cleanJSONResponse(responseText)
	if err := json.Unmarshal([]byte(cleanedResponse), &suggestions); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %v", err)
	}

	// Calculate prices
	for i := range suggestions {
		suggestions[i].TotalPrice = 0
		for j := range suggestions[i].Ingredients {
			if product, ok := productMap[suggestions[i].Ingredients[j].ProductID]; ok {
				suggestions[i].Ingredients[j].Available = true
				suggestions[i].Ingredients[j].Price = s.calculatePrice(&product, suggestions[i].Ingredients[j].Quantity, suggestions[i].Ingredients[j].Unit)
				suggestions[i].TotalPrice += suggestions[i].Ingredients[j].Price
			}
		}
	}

	return suggestions, nil
}

// Helper functions
func (s *AIService) buildProductListString(products []models.Product) string {
	var lines []string
	for _, p := range products {
		lines = append(lines, fmt.Sprintf("- ID: %d, Name: %s, Price: %.2f per %s, Stock: %.0f %s", 
			p.ID, p.Name, p.Price, p.Unit, p.Stock, p.Unit))
	}
	return strings.Join(lines, "\n")
}

func (s *AIService) cleanJSONResponse(response string) string {
	// Remove markdown code blocks
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	return strings.TrimSpace(response)
}

func (s *AIService) calculatePrice(product *models.Product, quantity float64, unit models.Unit) float64 {
	// Simple price calculation - assumes product price is per base unit
	// In real scenario, you'd need more sophisticated unit conversion
	return product.Price * quantity / 1000 // Assuming price is per kg/liter and quantity in g/ml
}

func (s *AIService) Close() error {
	return nil
}
