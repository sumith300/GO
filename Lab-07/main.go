package main

import (
    "bufio"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "os"
    "strings"
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
    ID       int
    Name     string
    Category string
    Price    float64
    Stock    int
}

// Order struct to store order details
type Order struct {
    Product  *Product
    Quantity int
}

// ProductCatalog represents the store's product inventory
type ProductCatalog map[int]*Product

// ProductData represents the structure of the JSON file
type ProductData struct {
    Products []Product `json:"products"`
}

// Store struct implements all interfaces
type Store struct {
    catalog ProductCatalog
}

// NewStore creates a new store instance
func NewStore() *Store {
    return &Store{
        catalog: make(ProductCatalog),
    }
}

// InitializeCatalog implements ProductManager interface
func (s *Store) InitializeCatalog() error {
    data, err := os.ReadFile("products.json")
    if err != nil {
        return fmt.Errorf("error reading products file: %v", err)
    }

    var productData ProductData
    if err := json.Unmarshal(data, &productData); err != nil {
        return fmt.Errorf("error parsing products data: %v", err)
    }

    s.catalog = make(ProductCatalog)
    for _, product := range productData.Products {
        // Create a new product pointer for each product
        newProduct := product // Copy the product
        s.catalog[product.ID] = &newProduct
    }

    return nil
}

// GetProduct implements ProductManager interface (Call by Reference)
func (s *Store) GetProduct(id int) (*Product, error) {
    product, exists := s.catalog[id]
    if !exists {
        return nil, errors.New("product not found")
    }
    return product, nil
}

// UpdateStock implements ProductManager interface (Call by Reference)
func (s *Store) UpdateStock(id int, quantity int) error {
    product, exists := s.catalog[id]
    if !exists {
        return errors.New("product not found")
    }
    if product.Stock < quantity {
        return fmt.Errorf("insufficient stock: only %d items available", product.Stock)
    }
    product.Stock -= quantity
    return nil
}

// DisplayProducts implements ProductManager interface (Call by Value)
func (s *Store) DisplayProducts() {
    fmt.Println("Available Products:")
    for _, product := range s.catalog {
        fmt.Printf("ID: %d, Name: %s, Category: %s, Price: ₹%.2f, Stock: %d\n",
            product.ID, product.Name, product.Category, product.Price, product.Stock)
    }
}

// validateQuantity checks if the quantity is valid (Call by Value)
func validateQuantity(quantity int) error {
    if quantity < 0 {
        return errors.New("quantity cannot be negative")
    }
    return nil
}

// CreateOrder implements OrderProcessor interface (Call by Reference)
func (s *Store) CreateOrder(product *Product, quantity int) (*Order, error) {
    if err := validateQuantity(quantity); err != nil {
        return nil, err
    }
    if product.Stock < quantity {
        return nil, fmt.Errorf("insufficient stock: only %d items available", product.Stock)
    }
    return &Order{Product: product, Quantity: quantity}, nil
}

// CalculateTotal implements OrderProcessor interface (Call by Reference)
func (s *Store) CalculateTotal(order *Order) float64 {
    return float64(order.Quantity) * order.Product.Price
}

// DisplayOrderDetails implements DisplayManager interface (Call by Reference)
func (s *Store) DisplayOrderDetails(order *Order) {
    totalPrice := s.CalculateTotal(order)

    fmt.Printf("\nFinal Product Details:\n")
    fmt.Printf("ID: %d\n", order.Product.ID)
    fmt.Printf("Name: %s\n", order.Product.Name)
    fmt.Printf("Category: %s\n", order.Product.Category)
    fmt.Printf("Price: ₹%.2f\n", order.Product.Price)
    fmt.Printf("Quantity: %d\n", order.Quantity)
    fmt.Printf("Total Price: ₹%.2f\n", totalPrice)
}

// ProcessOrder implements OrderProcessor interface (Call by Reference)
func (s *Store) ProcessOrder(order *Order) {
    // Update stock quantity using UpdateStock
    if err := s.UpdateStock(order.Product.ID, order.Quantity); err != nil {
        fmt.Printf("Error updating stock: %v\n", err)
        return
    }

    if order.Quantity > 0 {
        fmt.Println("Product is in stock and ready for quick delivery!")
    } else {
        fmt.Println("Product is out of stock! Restocking soon.")
    }

    fmt.Println("\nProcessing Order...")
    for i := 0; i < order.Quantity; i++ {
        fmt.Printf("Packing item %d\n", i+1)
    }

    switch order.Product.Category {
    case "Grocery":
        fmt.Println("This is a grocery item. Perishable and needs fast delivery!")
    case "Electronics":
        fmt.Println("This is an electronic item. Ensure safe packaging!")
    case "Fashion":
        fmt.Println("This is a fashion item. Speed and presentation matter!")
    default:
        fmt.Println("Unknown category. Classify properly for quick commerce.")
    }

    fmt.Println("Order ready for dispatch!")
}

// DisplayAllOrders implements DisplayManager interface (Call by Reference)
func (s *Store) DisplayAllOrders(orders []*Order) {
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

// getValidProductID prompts for and validates a product ID (Call by Reference)
func getValidProductID(store *Store) (int, error) {
    reader := bufio.NewReader(os.Stdin)
    for {
        var productID int
        fmt.Print("\nEnter Product ID: ")
        
        input, err := reader.ReadString('\n')
        if err != nil {
            return 0, errors.New("input error")
        }
        input = strings.TrimSpace(input)
        
        _, err = fmt.Sscanf(input, "%d", &productID)
        if err != nil {
            fmt.Println("Invalid input: please enter a valid number")
            continue
        }
        
        _, err = store.GetProduct(productID)
        if err != nil {
            fmt.Println("Invalid product ID: product not found")
            continue
        }
        
        return productID, nil
    }
}

// getValidQuantity prompts for and validates quantity input (Call by Value)
func getValidQuantity() (int, error) {
    reader := bufio.NewReader(os.Stdin)
    for {
        var quantity int
        fmt.Print("Enter Quantity: ")
        
        input, err := reader.ReadString('\n')
        if err != nil {
            return 0, errors.New("input error")
        }
        input = strings.TrimSpace(input)
        
        _, err = fmt.Sscanf(input, "%d", &quantity)
        if err != nil {
            fmt.Println("Invalid input: please enter a valid number")
            continue
        }
        
        if err := validateQuantity(quantity); err != nil {
            fmt.Println(err.Error())
            continue
        }
        
        return quantity, nil
    }
}

// CartItem represents an item in the shopping cart
type CartItem struct {
    Product  *Product `json:"product"`
    Quantity int      `json:"quantity"`
}

// handleGetProducts returns the product catalog as JSON
func (s *Store) handleGetProducts(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    // Convert map to array
    products := make([]Product, 0, len(s.catalog))
    for _, product := range s.catalog {
        products = append(products, *product) // Dereference pointer
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