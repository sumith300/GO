# Shopping Cart System with Web Interface

A Go-based shopping cart system featuring a web interface, concurrent order processing, and RESTful API endpoints. The system demonstrates the use of interfaces, JSON operations, and web server implementation in Go.

## Features

### Product Management
- Product catalog initialization from JSON file
- Thread-safe product operations
- Stock management with concurrent access handling
- Product information retrieval and display

### Order Processing
- Concurrent order processing with worker pool
- Real-time stock updates
- Order validation and error handling
- Category-specific order handling

### Web Interface
- Modern responsive web interface
- Real-time product catalog display
- Shopping cart functionality
- Checkout process

## Installation

1. Clone the repository
2. Ensure Go is installed on your system
3. Navigate to the project directory
4. Run the following commands:
   ```bash
   go mod init shopping-cart
   go mod tidy
   ```

## Usage

1. Start the server:
   ```bash
   go run main.go
   ```
2. Access the web interface at `http://localhost:8080`

## API Endpoints

### Products

- `GET /api/products` - Get all products
- `GET /api/products/{id}` - Get a specific product
- `PUT /api/products/stock` - Update product stock

### Orders

- `POST /api/orders` - Create a new order
- `POST /api/checkout` - Process checkout

## Data Structure

### Product
```json
{
    "id": 1,
    "name": "Product Name",
    "category": "Category",
    "price": 99.99,
    "stock": 100
}
```

### Order
```json
{
    "productId": 1,
    "quantity": 5
}
```

## Features Implementation

### Interface-Based Design
- ProductManager interface for product operations
- OrderProcessor interface for order handling
- DisplayManager interface for output formatting

### Concurrency
- Thread-safe operations using mutex
- Worker pool for order processing
- Asynchronous order handling

### Error Handling
- Comprehensive error checking
- Stock validation
- Input validation

## Web Interface Features

- Real-time product updates
- Shopping cart management
- Order processing status
- Responsive design