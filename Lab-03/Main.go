package main

import (
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

func main() {
	// Initialize product catalog using a map of Product structs
	productCatalog := map[int]Product{
		1: {ID: 1, Name: "Apple", Category: "Grocery", Price: 0.5},
		2: {ID: 2, Name: "Laptop", Category: "Electronics", Price: 999.99},
		3: {ID: 3, Name: "T-Shirt", Category: "Fashion", Price: 19.99},
	}

	// Slice to store multiple orders
	var orders []Order

	// Loop for multiple orders
	for {

	// Display available products
	fmt.Println("Available Products:")
	for _, product := range productCatalog {
		fmt.Printf("ID: %d, Name: %s, Category: %s, Price: $%.2f\n",
			product.ID, product.Name, product.Category, product.Price)
	}

	var productID int
	var quantity int

	// Get product ID with input validation
	for {
		fmt.Print("\nEnter Product ID: ")
		_, err := fmt.Scanf("%d\n", &productID)
		_, exists := productCatalog[productID]
		if err != nil || !exists {
			fmt.Println("Invalid input! Enter a valid Product ID.")
			fmt.Scanln()
			continue
		}
		break
	}

	// Get quantity with input validation
	for {
		fmt.Print("Enter Quantity: ")
		_, err := fmt.Scanf("%d\n", &quantity)
		if err != nil {
			fmt.Println("Invalid input! Enter a valid integer for Quantity.")
			fmt.Scanln()
			continue
		}
		if quantity < 0 {
			fmt.Println("Quantity cannot be negative.")
			continue
		}
		break
	}

	// Create new order and append to orders slice
	newOrder := Order{
		Product:  productCatalog[productID],
		Quantity: quantity,
	}
	orders = append(orders, newOrder)

	// Calculate total price
	totalPrice := float64(quantity) * productCatalog[productID].Price

	// Display order details
	fmt.Printf("\nFinal Product Details:\n")
	fmt.Printf("ID: %d\n", productID)
	fmt.Printf("Name: %s\n", productCatalog[productID].Name)
	fmt.Printf("Category: %s\n", productCatalog[productID].Category)
	fmt.Printf("Price: $%.2f\n", productCatalog[productID].Price)
	fmt.Printf("Quantity: %d\n", quantity)
	fmt.Printf("Total Price: $%.2f\n", totalPrice)

	// Check stock status
	if quantity > 0 {
		fmt.Println("Product is in stock and ready for quick delivery!")
	} else {
		fmt.Println("Product is out of stock! Restocking soon.")
	}

	// Process order
	fmt.Println("\nProcessing Order...")
	for i := 0; i < quantity; i++ {
		fmt.Printf("Packing item %d\n", i+1)
	}

	// Handle different categories
	switch productCatalog[productID].Category {
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

		// Ask if user wants to place another order
		fmt.Print("\nDo you want to place another order? (y/n): ")
		var response string
		fmt.Scanf("%s\n", &response)
		if response != "y" && response != "Y" {
			break
		}
	}

	// Display all orders with total amount
	fmt.Println("\nAll Orders:")
	var grandTotal float64
	for i, order := range orders {
		orderTotal := float64(order.Quantity) * order.Product.Price
		grandTotal += orderTotal
		fmt.Printf("Order %d: %s x%d - $%.2f\n",
			i+1, order.Product.Name, order.Quantity, orderTotal)
	}
	fmt.Printf("\nGrand Total: $%.2f\n", grandTotal)
}