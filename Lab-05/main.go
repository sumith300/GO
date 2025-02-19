package main

import (
	"errors"
	"fmt"
)

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

// initializeProductCatalog creates and returns a new product catalog
func initializeProductCatalog() ProductCatalog {
	return ProductCatalog{
		1: {ID: 1, Name: "Apple", Category: "Grocery", Price: 0.5},
		2: {ID: 2, Name: "Laptop", Category: "Electronics", Price: 999.99},
		3: {ID: 3, Name: "T-Shirt", Category: "Fashion", Price: 19.99},
	}
}

// displayProducts shows all available products in the catalog
func displayProducts(catalog ProductCatalog) {
	fmt.Println("Available Products:")
	for _, product := range catalog {
		fmt.Printf("ID: %d, Name: %s, Category: %s, Price: $%.2f\n",
			product.ID, product.Name, product.Category, product.Price)
	}
}

// getProductByID retrieves a product from the catalog by its ID
func getProductByID(catalog ProductCatalog, id int) (Product, error) {
	product, exists := catalog[id]
	if !exists {
		return Product{}, errors.New("product not found")
	}
	return product, nil
}

// validateQuantity checks if the quantity is valid
func validateQuantity(quantity int) error {
	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}
	return nil
}

// createOrder creates a new order with the given product and quantity
func createOrder(product Product, quantity int) (Order, error) {
	if err := validateQuantity(quantity); err != nil {
		return Order{}, err
	}
	return Order{Product: product, Quantity: quantity}, nil
}

// calculateTotal calculates the total price for an order
func calculateTotal(order Order) float64 {
	return float64(order.Quantity) * order.Product.Price
}

// displayOrderDetails shows the details of a processed order
func displayOrderDetails(order Order) {
	totalPrice := calculateTotal(order)

	fmt.Printf("\nFinal Product Details:\n")
	fmt.Printf("ID: %d\n", order.Product.ID)
	fmt.Printf("Name: %s\n", order.Product.Name)
	fmt.Printf("Category: %s\n", order.Product.Category)
	fmt.Printf("Price: $%.2f\n", order.Product.Price)
	fmt.Printf("Quantity: %d\n", order.Quantity)
	fmt.Printf("Total Price: $%.2f\n", totalPrice)
}

// processOrder handles the order processing and displays appropriate messages
func processOrder(order Order) {
	// Check stock status
	if order.Quantity > 0 {
		fmt.Println("Product is in stock and ready for quick delivery!")
	} else {
		fmt.Println("Product is out of stock! Restocking soon.")
	}

	// Process order
	fmt.Println("\nProcessing Order...")
	for i := 0; i < order.Quantity; i++ {
		fmt.Printf("Packing item %d\n", i+1)
	}

	// Handle different categories
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

// displayAllOrders shows all orders in the current session
func displayAllOrders(orders []Order) {
	fmt.Println("\nAll Orders:")
	for i, order := range orders {
		fmt.Printf("Order %d: %s x%d - $%.2f\n",
			i+1, order.Product.Name, order.Quantity, calculateTotal(order))
	}
}

// getValidProductID prompts for and validates a product ID
func getValidProductID(catalog ProductCatalog) (int, error) {
	var productID int
	fmt.Print("\nEnter Product ID: ")
	_, err := fmt.Scanf("%d\n", &productID)
	if err != nil {
		return 0, errors.New("invalid input: please enter a valid number")
	}
	
	_, exists := catalog[productID]
	if !exists {
		return 0, errors.New("invalid product ID: product not found")
	}
	
	return productID, nil
}

// getValidQuantity prompts for and validates quantity input
func getValidQuantity() (int, error) {
	var quantity int
	fmt.Print("Enter Quantity: ")
	_, err := fmt.Scanf("%d\n", &quantity)
	if err != nil {
		return 0, errors.New("invalid input: please enter a valid number")
	}
	
	if err := validateQuantity(quantity); err != nil {
		return 0, err
	}
	
	return quantity, nil
}

func main() {
	// Initialize product catalog
	catalog := initializeProductCatalog()
	
	// Initialize orders slice
	var orders []Order

	// Display available products
	displayProducts(catalog)

	// Get and validate product ID
	productID, err := getValidProductID(catalog)
	for err != nil {
		fmt.Println(err)
		fmt.Scanln() // Clear input buffer
		productID, err = getValidProductID(catalog)
	}

	// Get and validate quantity
	quantity, err := getValidQuantity()
	for err != nil {
		fmt.Println(err)
		fmt.Scanln() // Clear input buffer
		quantity, err = getValidQuantity()
	}

	// Get product from catalog
	product, err := getProductByID(catalog, productID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Create new order
	newOrder, err := createOrder(product, quantity)
	if err != nil {
		fmt.Println("Error creating order:", err)
		return
	}

	// Add order to orders slice
	orders = append(orders, newOrder)

	// Display and process order
	displayOrderDetails(newOrder)
	processOrder(newOrder)

	// Display all orders
	displayAllOrders(orders)
}