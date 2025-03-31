package main

import (
    "testing"
)

func TestInitializeCatalog(t *testing.T) {
    store := NewStore()
    err := store.InitializeCatalog()
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if len(store.catalog) == 0 {
        t.Errorf("Expected catalog to be initialized, got empty catalog")
    }
}

func TestGetProduct(t *testing.T) {
    store := NewStore()
    store.InitializeCatalog()
    product, err := store.GetProduct(1)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if product == nil {
        t.Errorf("Expected product, got nil")
    }
}

func TestUpdateStock(t *testing.T) {
    store := NewStore()
    store.InitializeCatalog()
    err := store.UpdateStock(1, 5)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    product, _ := store.GetProduct(1)
    if product.Stock != 5 {
        t.Errorf("Expected stock to be 5, got %d", product.Stock)
    }
}

func TestCreateOrder(t *testing.T) {
    store := NewStore()
    store.InitializeCatalog()
    product, _ := store.GetProduct(1)
    order, err := store.CreateOrder(product, 2)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if order.Quantity != 2 {
        t.Errorf("Expected quantity to be 2, got %d", order.Quantity)
    }
}

func TestCalculateTotal(t *testing.T) {
    store := NewStore()
    store.InitializeCatalog()
    product, _ := store.GetProduct(1)
    order, _ := store.CreateOrder(product, 2)
    total := store.CalculateTotal(order)
    expectedTotal := product.Price * 2
    if total != expectedTotal {
        t.Errorf("Expected total to be %v, got %v", expectedTotal, total)
    }
}

func TestGetProductInvalidID(t *testing.T) {
	store := NewStore()
	store.InitializeCatalog()
	
	tests := []struct {
		name string
		id   int
	}{
		{"non-existent id", 999},
		{"negative id", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := store.GetProduct(tt.id)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", tt.name)
			}
			if product != nil {
				t.Errorf("Expected no product for %s, got %v", tt.name, product)
			}
		})
	}
}

func TestUpdateStockNegative(t *testing.T) {
	store := NewStore()
	store.InitializeCatalog()
	err := store.UpdateStock(1, -5)
	if err == nil {
		t.Errorf("Expected error for negative stock update, got nil")
	}
}

func TestCreateOrderExceedsStock(t *testing.T) {
	store := NewStore()
	store.InitializeCatalog()
	store.UpdateStock(1, 5) // Set known stock

	product, _ := store.GetProduct(1)
	_, err := store.CreateOrder(product, 6)
	if err == nil {
		t.Errorf("Expected error for quantity exceeding stock, got nil")
	}
}

func TestProductsJSONLoading(t *testing.T) {
	store := NewStore()
	err := store.InitializeCatalog()
	if err != nil {
		t.Fatalf("Failed to load products: %v", err)
	}

	if len(store.catalog) == 0 {
		t.Fatal("Catalog failed to load from JSON")
	}

	// Verify at least one product has non-zero fields
	product, _ := store.GetProduct(1)
	if product.Name == "" || product.Price <= 0 {
		t.Errorf("Invalid product data loaded from JSON: %+v", product)
	}
}

func TestCalculateTotalTableDriven(t *testing.T) {
	tests := []struct {
		name     string
		quantity int
		expected float64
	}{
		{"single item", 1, 40.00},
		{"multiple items", 3, 120.00},
		{"zero quantity", 0, 0.00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new store instance for each test case
			store := NewStore()
			if err := store.InitializeCatalog(); err != nil {
				t.Fatalf("Failed to initialize catalog: %v", err)
			}

			product, err := store.GetProduct(1)
			if err != nil {
				t.Fatalf("Failed to get product: %v", err)
			}
			
			order, err := store.CreateOrder(product, tt.quantity)
			if err != nil {
				t.Fatalf("Failed to create order: %v", err)
			}
			total := store.CalculateTotal(order)
			if total != tt.expected {
				t.Errorf("Expected %.2f, got %.2f", tt.expected, total)
			}
		})
	}
}