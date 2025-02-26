package main

import (
    "bufio"
    "encoding/json"
    "errors"
    "fmt"
    "os"
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
        s.catalog[product.ID] = product
    }

    return nil
}

// GetProduct implements ProductManager interface
func (s *Store) GetProduct(id int) (Product, error) {
    product, exists := s.catalog[id]
    if !exists {
        return Product{}, errors.New("product not found")
    }
    return product, nil
}

// DisplayProducts implements ProductManager interface
func (s *Store) DisplayProducts() {
    fmt.Println("Available Products:")
    for _, product := range s.catalog {
        fmt.Printf("ID: %d, Name: %s, Category: %s, Price: ₹%.2f\n",
            product.ID, product.Name, product.Category, product.Price)
    }
}

// validateQuantity checks if the quantity is valid
func validateQuantity(quantity int) error {
    if quantity < 0 {
        return errors.New("quantity cannot be negative")
    }
    return nil
}

// CreateOrder implements OrderProcessor interface
func (s *Store) CreateOrder(product Product, quantity int) (Order, error) {
    if err := validateQuantity(quantity); err != nil {
        return Order{}, err
    }
    return Order{Product: product, Quantity: quantity}, nil
}

// CalculateTotal implements OrderProcessor interface
func (s *Store) CalculateTotal(order Order) float64 {
    return float64(order.Quantity) * order.Product.Price
}

// DisplayOrderDetails implements DisplayManager interface
func (s *Store) DisplayOrderDetails(order Order) {
    totalPrice := s.CalculateTotal(order)

    fmt.Printf("\nFinal Product Details:\n")
    fmt.Printf("ID: %d\n", order.Product.ID)
    fmt.Printf("Name: %s\n", order.Product.Name)
    fmt.Printf("Category: %s\n", order.Product.Category)
    fmt.Printf("Price: ₹%.2f\n", order.Product.Price)
    fmt.Printf("Quantity: %d\n", order.Quantity)
    fmt.Printf("Total Price: ₹%.2f\n", totalPrice)
}

// ProcessOrder implements OrderProcessor interface
func (s *Store) ProcessOrder(order Order) {
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

func main() {
    // Create new store instance
    store := NewStore()

    // Initialize product catalog
    if err := store.InitializeCatalog(); err != nil {
        fmt.Println("Error initializing product catalog:", err)
        return
    }

    // Initialize orders slice
    var orders []Order

    // Main loop for handling orders
    for {
        // Display available products
        store.DisplayProducts()

        // Get product ID
        productID, err := getValidProductID(store)
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }

        // Get product details
        product, err := store.GetProduct(productID)
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }

        // Get quantity
        quantity, err := getValidQuantity()
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }

        // Create and process order
        order, err := store.CreateOrder(product, quantity)
        if err != nil {
            fmt.Println("Error creating order:", err)
            continue
        }

        // Display order details
        store.DisplayOrderDetails(order)
        store.ProcessOrder(order)

        // Add to orders list
        orders = append(orders, order)

        // Display all orders
        store.DisplayAllOrders(orders)

        // Ask if user wants to continue
        fmt.Print("\nDo you want to place another order? (y/n): ")
        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        if scanner.Text() != "y" {
            break
        }
    }
}