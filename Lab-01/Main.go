package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var productID int
	var quantity int
	var price float64
	var productName string
	var category string

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter Product ID: ")
		_, err := fmt.Scanf("%d\n", &productID)
		if err != nil {
			fmt.Println("Invalid input! Enter a valid integer for Product ID.")
			fmt.Scanln()
			continue
		}
		break
	}

	fmt.Print("Enter Product Name: ")
	productName, _ = reader.ReadString('\n')
	productName = strings.TrimSpace(productName)

	fmt.Print("Enter Product Category: ")
	category, _ = reader.ReadString('\n')
	category = strings.TrimSpace(category)

	for {
		fmt.Print("Enter Price: ")
		_, err := fmt.Scanf("%f\n", &price)
		if err != nil {
			fmt.Println("Invalid input! Enter a valid float for Price.")
			fmt.Scanln()
			continue
		}
		if price <= 0 {
			fmt.Println("Price must be greater than zero.")
			continue
		}
		break
	}

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

	if quantity > 0 {
		fmt.Println("Product is in stock and ready for quick delivery!")
	} else {
		fmt.Println("Product is out of stock! Restocking soon.")
	}

	fmt.Println("\nProcessing Order...")
	for i := 0; i < quantity; i++ {
		fmt.Printf("Packing item %d\n", i+1)
	}

	switch category {
	case "Grocery":
		fmt.Println("This is a grocery item. Perishable and needs fast delivery!")
	case "Electronics":
		fmt.Println("This is an electronic item. Ensure safe packaging!")
	case "Fashion":
		fmt.Println("This is a fashion item. Speed and presentation matter!")
	default:
		fmt.Println("Unknown category. Classify properly for quick commerce.")
	}

	fmt.Printf("\nFinal Product Details:\n")
	fmt.Printf("ID: %d\n", productID)
	fmt.Printf("Name: %s\n", productName)
	fmt.Printf("Category: %s\n", category)
	fmt.Printf("Price: $%.2f\n", price)
	fmt.Printf("Quantity: %d\n", quantity)
	fmt.Println("Order ready for dispatch!")
}
