package main

import (
	"fmt"
)

func main() {
	products := map[int]struct {
		name     string
		category string
		price    float64
	}{
		1: {"Apple", "Grocery", 0.5},
		2: {"Laptop", "Electronics", 999.99},
		3: {"T-Shirt", "Fashion", 19.99},
	}

	fmt.Println("Available Products:")
	for id, product := range products {
		fmt.Printf("ID: %d, Name: %s, Category: %s, Price: $%.2f\n", id, product.name, product.category, product.price)
	}

	var productID int
	var quantity int

	for {
		fmt.Print("Enter Product ID: ")
		_, err := fmt.Scanf("%d\n", &productID)
		if err != nil || products[productID].name == "" {
			fmt.Println("Invalid input! Enter a valid Product ID.")
			fmt.Scanln()
			continue
		}
		break
	}

	product := products[productID]

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

	totalPrice := float64(quantity) * product.price

	fmt.Printf("\nFinal Product Details:\n")
	fmt.Printf("ID: %d\n", productID)
	fmt.Printf("Name: %s\n", product.name)
	fmt.Printf("Category: %s\n", product.category)
	fmt.Printf("Price: $%.2f\n", product.price)
	fmt.Printf("Quantity: %d\n", quantity)
	fmt.Printf("Total Price: $%.2f\n", totalPrice)

	if quantity > 0 {
		fmt.Println("Product is in stock and ready for quick delivery!")
	} else {
		fmt.Println("Product is out of stock! Restocking soon.")
	}

	fmt.Println("\nProcessing Order...")
	for i := 0; i < quantity; i++ {
		fmt.Printf("Packing item %d\n", i+1)
	}

	switch product.category {
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
