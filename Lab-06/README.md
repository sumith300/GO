# Shopping Cart System with Web Interface

A Go-based shopping cart system featuring a web interface, concurrent order processing, and RESTful API endpoints. The system demonstrates the use of interfaces, worker pools, and thread-safe operations in Go.

## Features

### Product Management
- Thread-safe product catalog operations
- JSON-based product initialization
- Real-time stock updates
- Concurrent access handling

### Order Processing
- Worker pool implementation for concurrent order processing
- Automatic worker scaling based on load
- Category-specific order handling
- Real-time order status updates

### Web Interface
- Modern responsive design using Bootstrap
- Real-time cart management
- Dynamic product display
- Offcanvas shopping cart

## Implementation Details

### Interface-Based Design
- `ProductManager`: Handles product catalog operations
- `OrderProcessor`: Manages order creation and processing
- `DisplayManager`: Handles display formatting

### Concurrency Features
- Worker pool pattern for order processing
- Mutex-based thread safety for catalog operations
- Channel-based communication between components
- Automatic worker scaling with monitoring

### Web Components
- Bootstrap 5.1.3 for responsive design
- Dynamic product grid layout
- Real-time cart updates
- Checkout functionality

## Project Structure
```
├── main.go           # Main application logic
├── products.json     # Product catalog data
└── static/
    ├── index.html    # Web interface
    ├── script.js     # Frontend logic
    └── styles.css    # Custom styling
```

## Usage

1. Start the server:
   ```bash
   go run main.go
   ```

2. Access the web interface at `http://localhost:8080`

## API Endpoints

### Products
- `GET /products` - Retrieve all products
- `GET /products/{id}` - Get specific product

### Orders
- `POST /order` - Create and process a new order

## Technical Features

### Thread Safety
- Mutex-protected catalog operations
- Atomic counters for worker management
- Thread-safe order processing

### Worker Pool
- Dynamic worker pool with auto-scaling
- Idle timeout management
- Worker monitoring and statistics

### Error Handling
- Comprehensive error checking
- Graceful error responses
- Thread-safe error reporting