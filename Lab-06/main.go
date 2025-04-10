package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "os"
    "sync"
    "sync/atomic"
    "runtime"
    "time"
)

// ProductManager interface defines product-related operations
type ProductManager interface {
    InitializeCatalog() error
    GetProduct(id int) (Product, error)
    DisplayProducts()
}

// OrderProcessor interface defines order-related operations
type OrderProcessor interface {
    CreateOrder(product Product, quantity int) (Order, error)
    ProcessOrder(order Order)
    CalculateTotal(order Order) float64
}

// DisplayManager interface defines display-related operations
type DisplayManager interface {
    DisplayOrderDetails(order Order)
    DisplayAllOrders(orders []Order)
}

// Product struct to define the structure of a product
type Product struct {
    ID       int
    Name     string
    Category string
    Price    float64
    Stock    int
}

// Order struct to store order details
type Order struct {
    Product  Product
    Quantity int
}

// ProductCatalog represents the store's product inventory
type ProductCatalog map[int]Product

// ProductData represents the structure of the JSON file
type ProductData struct {
    Products []Product `json:"products"`
}

// Store struct implements all interfaces with concurrency support
type Store struct {
    catalog     ProductCatalog
    mu          sync.RWMutex
    orderChan   chan Order
    resultChan  chan string
    workerCount int
    workerDone  chan bool
    activeWorkers int32
}

// NewStore creates a new store instance with concurrent features
func NewStore() *Store {
    store := &Store{
        catalog:     make(ProductCatalog),
        orderChan:   make(chan Order, 100),
        resultChan:  make(chan string, 100),
        workerCount: 3, // Number of concurrent workers
        workerDone:  make(chan bool),
        activeWorkers: 0,
    }
    // Start the worker pool and monitoring
    store.startWorkerPool()
    go store.monitorGoroutines()
    return store
}

// startWorkerPool initializes a pool of workers to process orders concurrently
func (s *Store) startWorkerPool() {
    for i := 0; i < s.workerCount; i++ {
        go func(workerID int) {
            atomic.AddInt32(&s.activeWorkers, 1)
            defer atomic.AddInt32(&s.activeWorkers, -1)
            
            idleTimeout := time.NewTimer(30 * time.Second)
            defer idleTimeout.Stop()
            
            for {
                select {
                case order, ok := <-s.orderChan:
                    if !ok {
                        return
                    }
                    idleTimeout.Reset(30 * time.Second)
                    s.processOrderAsync(order, workerID)
                case <-idleTimeout.C:
                    select {
                    case s.workerDone <- true:
                        return
                    default:
                        // Continue if channel is full
                    }
                }
            }
        }(i + 1)
    }
}

// monitorGoroutines tracks and manages goroutine count
func (s *Store) monitorGoroutines() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            currentGoroutines := runtime.NumGoroutine()
            activeWorkers := atomic.LoadInt32(&s.activeWorkers)
            fmt.Printf("Current goroutines: %d, Active workers: %d\n", currentGoroutines, activeWorkers)
        case <-s.workerDone:
            if atomic.LoadInt32(&s.activeWorkers) > 1 { // Keep at least one worker
                continue
            }
            // Start a new worker if needed
            go func() {
                atomic.AddInt32(&s.activeWorkers, 1)
                s.startWorkerPool()
            }()
        }
    }
}

// processOrderAsync handles order processing asynchronously
func (s *Store) processOrderAsync(order Order, workerID int) {
    fmt.Printf("Worker %d processing order for %s\n", workerID, order.Product.Name)

    // Update stock quantity with thread safety
    s.mu.Lock()
    if product, exists := s.catalog[order.Product.ID]; exists {
        product.Stock -= order.Quantity
        s.catalog[order.Product.ID] = product
    }
    s.mu.Unlock()

    if order.Quantity > 0 {
        s.resultChan <- fmt.Sprintf("Worker %d: Product is in stock and ready for quick delivery!", workerID)
    } else {
        s.resultChan <- fmt.Sprintf("Worker %d: Product is out of stock! Restocking soon.", workerID)
    }

    // Process order
    s.resultChan <- fmt.Sprintf("\nWorker %d: Processing Order...", workerID)
    for i := 0; i < order.Quantity; i++ {
        s.resultChan <- fmt.Sprintf("Worker %d: Packing item %d", workerID, i+1)
    }

    // Handle different categories
    switch order.Product.Category {
    case "Grocery":
        s.resultChan <- fmt.Sprintf("Worker %d: This is a grocery item. Perishable and needs fast delivery!", workerID)
    case "Electronics":
        s.resultChan <- fmt.Sprintf("Worker %d: This is an electronic item. Ensure safe packaging!", workerID)
    case "Fashion":
        s.resultChan <- fmt.Sprintf("Worker %d: This is a fashion item. Speed and presentation matter!", workerID)
    default:
        s.resultChan <- fmt.Sprintf("Worker %d: Unknown category. Classify properly for quick commerce.", workerID)
    }

    s.resultChan <- fmt.Sprintf("Worker %d: Order ready for dispatch!", workerID)
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
        s.catalog[product.ID] = product
    }

    return nil
}

// GetProduct implements ProductManager interface with thread safety
func (s *Store) GetProduct(id int) (Product, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    product, exists := s.catalog[id]
    if !exists {
        return Product{}, errors.New("product not found")
    }
    return product, nil
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
func (s *Store) CreateOrder(product Product, quantity int) (Order, error) {
    if err := validateQuantity(quantity); err != nil {
        return Order{}, err
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    // Check if enough stock is available
    if product.Stock < quantity {
        return Order{}, fmt.Errorf("insufficient stock: only %d items available", product.Stock)
    }

    order := Order{Product: product, Quantity: quantity}
    
    // Send order to processing channel
    go func() {
        s.orderChan <- order
    }()

    return order, nil
}

// CalculateTotal implements OrderProcessor interface
func (s *Store) CalculateTotal(order Order) float64 {
    return float64(order.Quantity) * order.Product.Price
}

// DisplayOrderDetails implements DisplayManager interface with concurrent result handling
func (s *Store) DisplayOrderDetails(order Order) {
    totalPrice := s.CalculateTotal(order)

    fmt.Printf("\nFinal Product Details:\n")
    fmt.Printf("ID: %d\n", order.Product.ID)
    fmt.Printf("Name: %s\n", order.Product.Name)
    fmt.Printf("Category: %s\n", order.Product.Category)
    fmt.Printf("Price: ₹%.2f\n", order.Product.Price)
    fmt.Printf("Quantity: %d\n", order.Quantity)
    fmt.Printf("Total Price: ₹%.2f\n", totalPrice)

    // Display processing results
    go func() {
        for result := range s.resultChan {
            fmt.Println(result)
        }
    }()
}

// ProcessOrder implements OrderProcessor interface with concurrent processing
func (s *Store) ProcessOrder(order Order) {
    // Send order to processing channel
    s.orderChan <- order

    // Display processing results
    go func() {
        for result := range s.resultChan {
            fmt.Println(result)
        }
    }()
}

// DisplayAllOrders implements DisplayManager interface
func (s *Store) DisplayAllOrders(orders []Order) {
    fmt.Println("\nAll Orders:")
    var grandTotal float64
    for i, order := range orders {
        orderTotal := s.CalculateTotal(order)
        grandTotal += orderTotal
        fmt.Printf("Order %d: %s x%d - ₹%.2f\n",
            i+1, order.Product.Name, order.Quantity, orderTotal)
    }
    fmt.Printf("\nGrand Total: ₹%.2f\n", grandTotal)
}

// getValidProductID prompts for and validates a product ID
func getValidProductID(store *Store) (int, error) {
    var productID int
    fmt.Print("\nEnter Product ID: ")
    n, err := fmt.Scanf("%d", &productID)
    if err != nil || n != 1 {
        return 0, errors.New("invalid input: please enter a valid number")
    }
    
    _, err = store.GetProduct(productID)
    if err != nil {
        return 0, errors.New("invalid product ID: product not found")
    }
    
    return productID, nil
}

// getValidQuantity prompts for and validates quantity input
func getValidQuantity() (int, error) {
    var quantity int
    fmt.Print("Enter Quantity: ")
    n, err := fmt.Scanf("%d", &quantity)
    if err != nil || n != 1 {
        return 0, errors.New("invalid input: please enter a valid number")
    }
    
    if err := validateQuantity(quantity); err != nil {
        return 0, err
    }
    
    return quantity, nil
}

// CartItem represents an item in the shopping cart
type CartItem struct {
    Product  Product `json:"product"`
    Quantity int     `json:"quantity"`
}

// handleGetProducts returns the product catalog as JSON
func (s *Store) handleGetProducts(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    // Convert map to array
    products := make([]Product, 0, len(s.catalog))
    for _, product := range s.catalog {
        products = append(products, product)
    }
    json.NewEncoder(w).Encode(products)
}

// handleCheckout processes the checkout from the web interface
func (s *Store) handleCheckout(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var cartItems []CartItem
    if err := json.NewDecoder(r.Body).Decode(&cartItems); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Process each cart item
    for _, item := range cartItems {
        order, err := s.CreateOrder(item.Product, item.Quantity)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        s.ProcessOrder(order)
    }

    w.WriteHeader(http.StatusOK)
}

func main() {
    // Create new store instance
    store := NewStore()

    // Initialize product catalog
    if err := store.InitializeCatalog(); err != nil {
        fmt.Println("Error initializing product catalog:", err)
        return
    }

    // Set up HTTP routes
    http.HandleFunc("/api/products", store.handleGetProducts)
    http.HandleFunc("/api/checkout", store.handleCheckout)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "static/index.html")
    })

    // Start the server
    fmt.Println("Server starting on http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Server error:", err)
    }
}