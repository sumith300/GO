# Shopping Cart Program

## Summary

This Go program implements a simple but feature-rich shopping cart system with the following key features:

- **Product Catalog Management**: Uses Go maps to store and manage a catalog of products with their details (ID, Name, Category, Price)
- **Structured Data**: Implements custom structs (`Product` and `Order`) for organized data management
- **Multiple Orders**: Utilizes Go slices to handle multiple orders in a single session
- **Input Validation**: Robust error handling for user inputs including:
  - Product ID validation
  - Quantity validation (non-negative values)
- **Category-based Processing**: Special handling for different product categories (Grocery, Electronics, Fashion)
- **Order Processing**: Simulates order processing with item-by-item packing visualization
- **Price Calculation**: Automatic calculation of total price based on quantity and unit price

The program demonstrates the practical implementation of various Go programming concepts including:
- Structs and Custom Types
- Maps for Data Storage
- Slices for Dynamic Arrays
- Error Handling
- Control Structures (loops, switch statements)
- User Input Processing

## Output

Here's an example interaction with the shopping cart program:

```
Available Products:
ID: 1, Name: Apple, Category: Grocery, Price: $0.50
ID: 2, Name: Laptop, Category: Electronics, Price: $999.99
ID: 3, Name: T-Shirt, Category: Fashion, Price: $19.99

Enter Product ID: 2
Enter Quantity: 1

Final Product Details:
ID: 2
Name: Laptop
Category: Electronics
Price: $999.99
Quantity: 1
Total Price: $999.99
Product is in stock and ready for quick delivery!

Processing Order...
Packing item 1
This is an electronic item. Ensure safe packaging!
Order ready for dispatch!

All Orders:
Order 1: Laptop x1 - $999.99
```

This output demonstrates:
- Product catalog display
- Input validation
- Order processing
- Category-specific handling
- Multiple order tracking