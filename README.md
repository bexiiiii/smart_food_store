# 🍕 Smart Food Store + Recipe Planner

A modern food store API with AI-powered recipe suggestions using **Gemini API**.

## 🌟 Features

### Core Store Features
- User registration and authentication (JWT)
- Browse products by category
- Search products
- Shopping cart management
- Role-based access (User/Admin)

### 🤖 AI Features (Gemini API)
1. **Dish to Ingredients**: Enter a dish name → AI suggests products from the store with exact quantities
2. **Cart to Recipes**: Based on products in your cart → AI suggests what you can cook

### Admin Features
- Product management (CRUD)
- Category management
- Recipe management
- User management

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 14+
- Gemini API Key (optional, for AI features)

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/bexiiiii/smart_food_store.git
cd smart_food_store
```

2. **Setup environment**
```bash
cp .env.example .env
# Edit .env with your database credentials and Gemini API key
```

3. **Install dependencies**
```bash
go mod tidy
```

4. **Create PostgreSQL database**
```bash
createdb smart_food_store
```

5. **Run the application**
```bash
go run cmd/main.go
```

The server will start at `http://localhost:8080`

## 📚 API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login user |

### Products (Public)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/products` | Get all products |
| GET | `/api/v1/products/:id` | Get product by ID |
| GET | `/api/v1/products/category/:id` | Get products by category |
| GET | `/api/v1/products/search?q=query` | Search products |
| GET | `/api/v1/categories` | Get all categories |

### Recipes (Public)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/recipes` | Get all recipes |
| GET | `/api/v1/recipes/:id` | Get recipe by ID |
| GET | `/api/v1/recipes/search?q=query` | Search recipes |
| GET | `/api/v1/recipes/:id/calculate?servings=4` | Calculate ingredients for servings |

### Cart (Protected - requires JWT)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/cart` | Get user's cart |
| POST | `/api/v1/cart/items` | Add item to cart |
| POST | `/api/v1/cart/items/bulk` | Add multiple items |
| PUT | `/api/v1/cart/items/:product_id` | Update item quantity |
| DELETE | `/api/v1/cart/items/:product_id` | Remove item |
| DELETE | `/api/v1/cart` | Clear cart |

### AI Features
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/ai/dish-to-ingredients` | Get ingredients for a dish | Public |
| POST | `/api/v1/ai/products-to-recipes` | Get recipes from products | Public |
| GET | `/api/v1/ai/cart-to-recipes` | Get recipes from cart items | Protected |
| POST | `/api/v1/ai/add-to-cart` | Add AI suggestion to cart | Protected |

### Admin (Protected - requires Admin role)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/users` | Get all users |
| PATCH | `/api/v1/admin/users/:id/role` | Update user role |
| POST | `/api/v1/admin/products` | Create product |
| PUT | `/api/v1/admin/products/:id` | Update product |
| DELETE | `/api/v1/admin/products/:id` | Delete product |
| POST | `/api/v1/admin/recipes` | Create recipe |
| PUT | `/api/v1/admin/recipes/:id` | Update recipe |
| DELETE | `/api/v1/admin/recipes/:id` | Delete recipe |

## 🔐 Authentication

Use JWT token in Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## 📝 Example Requests

### Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "password": "password123"}'
```

### Get Ingredients for Dish (AI)
```bash
curl -X POST http://localhost:8080/api/v1/ai/dish-to-ingredients \
  -H "Content-Type: application/json" \
  -d '{"dish_name": "Pasta Carbonara", "servings": 4}'
```

### Get Recipe Suggestions from Cart (AI)
```bash
curl http://localhost:8080/api/v1/ai/cart-to-recipes \
  -H "Authorization: Bearer <your-token>"
```

## 🏗️ Project Structure

```
smart_food_store/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/              # Configuration
│   ├── database/            # Database connection & migrations
│   ├── handlers/            # HTTP handlers
│   ├── middleware/          # Auth middleware
│   ├── models/              # Data models
│   ├── repository/          # Database operations
│   └── services/            # Business logic
├── .env.example             # Environment template
├── go.mod                   # Go modules
└── README.md
```

## 👥 Team

- Tukezhan Kausar
- Ginayat Yerassyl
- Behruz Tokhtamishov

**Group: SE-2422**

## 📄 License

This project is part of the Advanced Programming 1 course at Astana IT University.
