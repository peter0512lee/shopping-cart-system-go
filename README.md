# Shopping Cart System

A simple e-commerce shopping cart system built with Go (backend) and React (frontend), featuring real-time cart updates and a responsive design.

## Features

- Product listing with dynamic price and stock display
- Shopping cart management (add, update, remove items)
- Real-time cart total calculation
- Responsive design for mobile and desktop
- MongoDB for data persistence
- Docker containerization

## Tech Stack

### Backend
- Go
- Gin (Web Framework)
- MongoDB
- Docker

### Frontend
- React
- TypeScript
- Tailwind CSS

## Prerequisites

- Docker and Docker Compose
- Node.js (for local development)
- Go (for local development)

## Getting Started

1. Clone the repository
```bash
git clone [your-repository-url]
cd shopping-cart-system
```

2. Start the application using Docker Compose
```bash
docker-compose up --build
```

The application will be available at:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- MongoDB: localhost:27017

## API Endpoints

### Products
- `GET /api/v1/products` - Get all products
- `POST /api/v1/products` - Create a new product
- `POST /api/v1/products/bulk` - Create multiple products
- `DELETE /api/v1/products` - Clear all products

### Cart
- `GET /api/v1/cart/:user_id` - Get user's cart
- `POST /api/v1/cart` - Add item to cart
- `PUT /api/v1/cart/:user_id` - Update cart item quantity

## Sample API Usage

### Create Products
```bash
curl -X POST http://localhost:8080/api/v1/products/bulk \
  -H "Content-Type: application/json" \
  -d '[
    {
      "name": "Gaming Laptop",
      "price": 1299.99,
      "stock": 10,
      "emoji": "💻"
    },
    {
      "name": "Wireless Headphones",
      "price": 199.99,
      "stock": 15,
      "emoji": "🎧"
    }
  ]'
```

### Add to Cart
```bash
curl -X POST http://localhost:8080/api/v1/cart \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "product_id": "YOUR_PRODUCT_ID",
    "quantity": 1
  }'
```

## Project Structure

```
shopping-cart-system/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── models/
│   │   ├── product.go
│   │   └── cart.go
│   └── handlers/
│       ├── product_handler.go
│       └── cart_handler.go
├── pkg/
│   └── database/
│       └── mongodb.go
├── web/
│   └── frontend/
│       ├── src/
│       │   ├── components/
│       │   └── App.tsx
│       └── package.json
├── docker-compose.yml
└── README.md
```

## Development

### Backend Development
```bash
# Run backend locally
go run cmd/api/main.go
```

### Frontend Development
```bash
cd web/frontend
npm install
npm start
```

## Testing

Currently implementing tests...

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details