package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "net/http"
    "os"
    "sync"
)

// ProductManager interface defines product-related operations
type ProductManager interface {
    InitializeCatalog() error
    GetProduct(id int) (*Product, error)
    DisplayProducts()
    UpdateStock(id int, quantity int) error
}

// OrderProcessor interface defines order-related operations
type OrderProcessor interface {
    CreateOrder(product *Product, quantity int) (*Order, error)
    ProcessOrder(order *Order)
    CalculateTotal(order *Order) float64
}

// DisplayManager interface defines display-related operations
type DisplayManager interface {
    DisplayOrderDetails(order *Order)
    DisplayAllOrders(orders []*Order)
}

// Product struct to define the structure of a product
type Product struct {
    ID       int     `json:"id"`
    Name     string  `json:"name"`
    Category string  `json:"category"`
    Price    float64 `json:"price"`
    Stock    int     `json:"stock"`
}

// Order struct to store order details
type Order struct {
    Product  *Product
    Quantity int
    Status   string
}

// ProductCatalog represents the store's product inventory
type ProductCatalog map[int]*Product

// ProductData represents the structure of the JSON file
type ProductData struct {
    Products []Product `json:"products"`
}

// Store struct implements all interfaces with concurrency support
type Store struct {
    catalog     ProductCatalog
    mu          sync.RWMutex
    orderChan   chan *Order
    resultChan  chan string
    workerCount int
}

// NewStore creates a new store instance with concurrent features
func NewStore() *Store {
    store := &Store{
        catalog:     make(ProductCatalog),
        orderChan:   make(chan *Order, 100),
        resultChan:  make(chan string, 100),
        workerCount: 3, // Number of concurrent workers
    }
    // Start the worker pool
    store.startWorkerPool()
    return store
}

// startWorkerPool initializes a pool of workers to process orders concurrently
func (s *Store) startWorkerPool() {
    for i := 0; i < s.workerCount; i++ {
        go func(workerID int) {
            for order := range s.orderChan {
                s.processOrderAsync(order, workerID)
            }
        }(i + 1)
    }
}

// InitializeCatalog implements ProductManager interface with thread safety
func (s *Store) InitializeCatalog() error {
    data, err := os.ReadFile("products.json")
    if err != nil {
        return fmt.Errorf("error reading products file: %v", err)
    }

    var productData ProductData
    if err := json.Unmarshal(data, &productData); err != nil {
        return fmt.Errorf("error parsing products data: %v", err)
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    s.catalog = make(ProductCatalog)
    for _, product := range productData.Products {
        newProduct := product
        s.catalog[product.ID] = &newProduct
    }

    return nil
}

// GetProduct implements ProductManager interface with thread safety
func (s *Store) GetProduct(id int) (*Product, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    product, exists := s.catalog[id]
    if !exists {
        return nil, errors.New("product not found")
    }
    return product, nil
}

// UpdateStock implements ProductManager interface with thread safety
func (s *Store) UpdateStock(id int, quantity int) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, exists := s.catalog[id]; !exists {
        return errors.New("product not found")
    }
    if quantity < 0 {
        return errors.New("quantity cannot be negative")
    }
    s.catalog[id].Stock = quantity
    return nil
}

// DisplayProducts implements ProductManager interface with thread safety
func (s *Store) DisplayProducts() {
    s.mu.RLock()
    defer s.mu.RUnlock()

    fmt.Println("Available Products:")
    for _, product := range s.catalog {
        fmt.Printf("ID: %d, Name: %s, Category: %s, Price: ₹%.2f, Stock: %d\n",
            product.ID, product.Name, product.Category, product.Price, product.Stock)
    }
}

// validateQuantity checks if the quantity is valid
func validateQuantity(quantity int) error {
    if quantity < 0 {
        return errors.New("quantity cannot be negative")
    }
    return nil
}

// CreateOrder implements OrderProcessor interface with thread safety
func (s *Store) CreateOrder(product *Product, quantity int) (*Order, error) {
    if err := validateQuantity(quantity); err != nil {
        return nil, err
    }
    if product == nil {
        return nil, errors.New("product cannot be nil")
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    if quantity > 0 {
        if product.Stock < quantity {
            return nil, fmt.Errorf("insufficient stock: only %d items available", product.Stock)
        }
        s.catalog[product.ID].Stock -= quantity
    }

    order := &Order{Product: product, Quantity: quantity, Status: "Created"}
    // Send order to processing channel
    go func() {
        s.orderChan <- order
    }()

    return order, nil
}

// processOrderAsync handles order processing asynchronously
func (s *Store) processOrderAsync(order *Order, workerID int) {
    fmt.Printf("Worker %d processing order for %s\n", workerID, order.Product.Name)

    if order.Quantity > 0 {
        fmt.Printf("Worker %d: Product is in stock and ready for quick delivery!\n", workerID)
    } else {
        fmt.Printf("Worker %d: Product is out of stock! Restocking soon.\n", workerID)
    }

    for i := 0; i < order.Quantity; i++ {
        fmt.Printf("Worker %d: Packing item %d\n", workerID, i+1)
    }

    switch order.Product.Category {
    case "Grocery":
        fmt.Printf("Worker %d: This is a grocery item. Perishable and needs fast delivery!\n", workerID)
    case "Electronics":
        fmt.Printf("Worker %d: This is an electronic item. Ensure safe packaging!\n", workerID)
    case "Fashion":
        fmt.Printf("Worker %d: This is a fashion item. Speed and presentation matter!\n", workerID)
    default:
        fmt.Printf("Worker %d: Unknown category. Classify properly for quick commerce.\n", workerID)
    }

    order.Status = "Processed"
    s.resultChan <- fmt.Sprintf("Worker %d: Order for %s has been processed successfully!", workerID, order.Product.Name)
}

// CalculateTotal implements OrderProcessor interface
func (s *Store) CalculateTotal(order *Order) float64 {
    return float64(order.Quantity) * order.Product.Price
}

// DisplayOrderDetails implements DisplayManager interface
func (s *Store) DisplayOrderDetails(order *Order) {
    totalPrice := s.CalculateTotal(order)

    fmt.Printf("\nOrder Details:\n")
    fmt.Printf("ID: %d\n", order.Product.ID)
    fmt.Printf("Name: %s\n", order.Product.Name)
    fmt.Printf("Category: %s\n", order.Product.Category)
    fmt.Printf("Price: ₹%.2f\n", order.Product.Price)
    fmt.Printf("Quantity: %d\n", order.Quantity)
    fmt.Printf("Total Price: ₹%.2f\n", totalPrice)
    fmt.Printf("Status: %s\n", order.Status)
}

// DisplayAllOrders implements DisplayManager interface
func (s *Store) DisplayAllOrders(orders []*Order) {
    fmt.Println("\nAll Orders:")
    for i, order := range orders {
        fmt.Printf("\nOrder %d:\n", i+1)
        s.DisplayOrderDetails(order)
    }
}

// OrderRequest represents the structure of order requests from the frontend
type OrderRequest struct {
    ProductID int `json:"productId"`
    Quantity  int `json:"quantity"`
}

func main() {
    store := NewStore()
    if err := store.InitializeCatalog(); err != nil {
        log.Fatalf("Error initializing catalog: %v\n", err)
    }

    // Serve static files
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/", fs)

    // API endpoints
    http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
        store.mu.RLock()
        products := make([]Product, 0, len(store.catalog))
        for _, p := range store.catalog {
            products = append(products, *p)
        }
        store.mu.RUnlock()

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(products)
    })

    http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var orderReq OrderRequest
        if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        product, err := store.GetProduct(orderReq.ProductID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        order, err := store.CreateOrder(product, orderReq.Quantity)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(order)
    })

    // Start the server
    fmt.Println("Server starting on http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}