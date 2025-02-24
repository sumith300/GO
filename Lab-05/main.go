package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
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

// ProductData represents the structure of the JSON file
type ProductData struct {
	Products []Product `json:"products"`
}

// initializeProductCatalog creates and returns a new product catalog
func initializeProductCatalog() (ProductCatalog, error) {
	// Read the JSON file
	data, err := os.ReadFile("products.json")
	if err != nil {
		return nil, fmt.Errorf("error reading products file: %v", err)
	}

	// Parse JSON data
	var productData ProductData
	if err := json.Unmarshal(data, &productData); err != nil {
		return nil, fmt.Errorf("error parsing products data: %v", err)
	}

	// Convert to catalog format
	catalog := make(ProductCatalog)
	for _, product := range productData.Products {
		catalog[product.ID] = product
	}

	return catalog, nil
}

// displayProducts shows all available products in the catalog
func displayProducts(catalog ProductCatalog) {
	fmt.Println("Available Products:")
	for _, product := range catalog {
		fmt.Printf("ID: %d, Name: %s, Category: %s, Price: ₹%.2f\n",
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
	fmt.Printf("Price: ₹%.2f\n", order.Product.Price)
	fmt.Printf("Quantity: %d\n", order.Quantity)
	fmt.Printf("Total Price: ₹%.2f\n", totalPrice)
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
    var grandTotal float64
    for i, order := range orders {
        orderTotal := calculateTotal(order)
        grandTotal += orderTotal
        fmt.Printf("Order %d: %s x%d - ₹%.2f\n",
            i+1, order.Product.Name, order.Quantity, orderTotal)
    }
    fmt.Printf("\nGrand Total: ₹%.2f\n", grandTotal)
}

// getValidProductID prompts for and validates a product ID
func getValidProductID(catalog ProductCatalog) (int, error) {
	var productID int
	fmt.Print("\nEnter Product ID: ")
	n, err := fmt.Scanf("%d", &productID)
	if err != nil || n != 1 {
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
    // Initialize product catalog
    catalog, err := initializeProductCatalog()
    if err != nil {
        fmt.Println("Error initializing product catalog:", err)
        return
    }
    
    // Initialize orders slice
    var orders []Order

    // Main loop for handling orders
    for {
        // Display available products
        displayProducts(catalog)
    
        // Get and validate product ID
        productID, err := getValidProductID(catalog)
        for err != nil {
            fmt.Println(err)
            bufio.NewReader(os.Stdin).ReadString('\n')
            productID, err = getValidProductID(catalog)
        }
    
        // Get and validate quantity
        quantity, err := getValidQuantity()
        for err != nil {
            fmt.Println(err)
            bufio.NewReader(os.Stdin).ReadString('\n')
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

        // Ask if user wants to continue shopping
        fmt.Print("\nDo you want to continue shopping? (y/n): ")
        var response string
        fmt.Scanf("%s", &response)
        bufio.NewReader(os.Stdin).ReadString('\n')

        if response != "y" && response != "Y" {
            fmt.Println("\nThank you for shopping with us!")
            displayAllOrders(orders)
            return
        }
    }
}