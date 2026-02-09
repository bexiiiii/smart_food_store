# Project Proposal - Smart Food Store + AI Recipe Planner

## Overview
This project builds a full-stack web application that combines a digital food store with an AI-powered cooking assistant.  
Users can register, browse products, manage a cart, explore recipes, calculate ingredient quantities by servings, and get AI suggestions based on a dish name or cart contents.  
The backend is implemented in Go (Gin + GORM + PostgreSQL), and the frontend is implemented in React + Vite with role-based flows for users and admins.

## a) Project relevance (why this topic)
Food shopping and meal planning are usually split across different apps: one for buying products and another for finding recipes.  
This causes friction, overbuying, and product waste. The project addresses this by connecting inventory-aware shopping with recipe planning in one workflow:
- users see what they can cook from available products,
- ingredient quantities can be calculated for exact servings,
- AI suggestions are grounded in real products from the store catalog.

## b) Competitor context
Most grocery platforms focus on product delivery, while recipe platforms focus on cooking content without direct product matching.  
Smart Food Store combines both layers in a single system and adds technical advantages:
- **Inventory-aware suggestions:** recipe suggestions are based on real products in stock.
- **Unified flow:** browse -> cart -> recipe planning -> AI suggestions -> add ingredients to cart.
- **Role-based management:** admins can manage users, categories, products, and recipes from the same API.
- **Modular architecture:** handlers/services/repositories allow future expansion (orders, payments, recommendation analytics).

## c) Target audience
- **Primary users:** students and busy professionals who need fast meal planning and shopping.
- **Secondary users:** families who want predictable grocery planning and fewer extra purchases.
- **Admin users:** store operators maintaining catalog quality, stock, and recipe content.

## d) Project features
- **Authentication and RBAC:** JWT login/register, protected routes, admin-only endpoints.
- **Product catalog:** category filtering, search, product details.
- **Cart management:** add/update/remove/clear with stock-aware validation.
- **Recipe module:** recipe search and dynamic ingredient calculation per serving.
- **AI assistant (Gemini):**
  - dish name -> required ingredients + matched store products,
  - selected products/cart -> recipe suggestions,
  - one-click add available AI ingredients to cart.
- **Admin module:** CRUD for products, categories, recipes; user role updates and user management.
- **Seeded startup data:** initial categories and products for immediate testing/demo.
