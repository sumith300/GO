let cart = [];
let products = [];

// Fetch products from the server
async function fetchProducts() {
    try {
        const response = await fetch('/api/products');
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        products = await response.json();
        console.log('Fetched products:', products); // Debug log
        displayProducts();
    } catch (error) {
        console.error('Error fetching products:', error);
        alert('Failed to load products. Please try again later.');
    }
}

// Display products in the grid
function displayProducts() {
    const productGrid = document.getElementById('productGrid');
    if (!productGrid) {
        console.error('Product grid element not found');
        return;
    }

    productGrid.innerHTML = products.map(product => `
        <div class="col">
            <div class="card product-card position-relative">
                <span class="badge bg-primary category-badge">${product.category}</span>
                <span class="badge ${product.stock > 0 ? 'bg-success' : 'bg-danger'} stock-badge">
                    ${product.stock > 0 ? 'In Stock' : 'Out of Stock'}
                </span>
                <div class="card-body d-flex flex-column">
                    <h5 class="card-title">${product.name}</h5>
                    <p class="card-text">₹${product.price.toFixed(2)}</p>
                    <p class="card-text">Stock: ${product.stock}</p>
                    <div class="quantity-control mb-3">
                        <button class="btn btn-sm btn-outline-secondary" onclick="event.stopPropagation(); updateCardQuantity(${product.id}, 'decrease')">-</button>
                        <span id="quantity-${product.id}">1</span>
                        <button class="btn btn-sm btn-outline-secondary" onclick="event.stopPropagation(); updateCardQuantity(${product.id}, 'increase')">+</button>
                    </div>
                    <button onclick="addToCartWithQuantity(${product.id})" 
                            class="btn btn-primary mt-auto" 
                            ${product.stock === 0 ? 'disabled' : ''}>
                        Add to Cart
                    </button>
                </div>
            </div>
        </div>
    `).join('');
}

// Add product to cart
function updateCardQuantity(productId, action) {
    const quantityElement = document.getElementById(`quantity-${productId}`);
    if (!quantityElement) {
        console.error(`Quantity element not found for product ${productId}`);
        return;
    }

    const product = products.find(p => p.id === productId);
    if (!product) {
        console.error(`Product not found with ID ${productId}`);
        return;
    }

    let currentQuantity = parseInt(quantityElement.textContent);
    if (isNaN(currentQuantity)) {
        currentQuantity = 1;
    }

    if (action === 'increase' && currentQuantity < product.stock) {
        currentQuantity++;
    } else if (action === 'decrease' && currentQuantity > 1) {
        currentQuantity--;
    }

    quantityElement.textContent = currentQuantity;
}

// Add product to cart
function addToCartWithQuantity(productId) {
    const product = products.find(p => p.id === productId);
    if (!product) {
        console.error(`Product not found with ID ${productId}`);
        return;
    }

    const quantityElement = document.getElementById(`quantity-${productId}`);
    if (!quantityElement) {
        console.error(`Quantity element not found for product ${productId}`);
        return;
    }

    const quantity = parseInt(quantityElement.textContent);
    if (isNaN(quantity) || quantity < 1) {
        alert('Invalid quantity');
        return;
    }

    const existingItem = cart.find(item => item.product.id === productId);
    if (existingItem) {
        if (existingItem.quantity + quantity <= product.stock) {
            existingItem.quantity += quantity;
        } else {
            alert('Not enough stock available!');
            return;
        }
    } else {
        if (quantity <= product.stock) {
            cart.push({ product, quantity });
        } else {
            alert('Not enough stock available!');
            return;
        }
    }
    
    quantityElement.textContent = '1';
    updateCart();
}

// Remove product from cart
function removeFromCart(productId) {
    cart = cart.filter(item => item.product.id !== productId);
    updateCart();
}

// Update cart display
function updateCart() {
    const cartItems = document.getElementById('cartItems');
    const cartCount = document.getElementById('cartCount');
    const cartTotal = document.getElementById('cartTotal');

    if (!cartItems || !cartCount || !cartTotal) {
        console.error('Cart elements not found');
        return;
    }

    cartCount.textContent = cart.reduce((total, item) => total + item.quantity, 0);

    cartItems.innerHTML = cart.map(item => `
        <div class="cart-item">
            <div class="d-flex justify-content-between align-items-center mb-2">
                <h6 class="mb-0">${item.product.name}</h6>
                <button class="btn btn-sm btn-danger" onclick="removeFromCart(${item.product.id})">&times;</button>
            </div>
            <div class="d-flex justify-content-between align-items-center">
                <div class="quantity-control">
                    <button class="btn btn-sm btn-outline-secondary" onclick="updateQuantity(${item.product.id}, ${item.quantity - 1})">-</button>
                    <span>${item.quantity}</span>
                    <button class="btn btn-sm btn-outline-secondary" onclick="updateQuantity(${item.product.id}, ${item.quantity + 1})">+</button>
                </div>
                <span>₹${(item.product.price * item.quantity).toFixed(2)}</span>
            </div>
        </div>
    `).join('');

    cartTotal.textContent = cart.reduce((total, item) => 
        total + (item.product.price * item.quantity), 0).toFixed(2);
}

// Update quantity of cart item
function updateQuantity(productId, newQuantity) {
    const item = cart.find(item => item.product.id === productId);
    if (!item) {
        console.error(`Cart item not found for product ${productId}`);
        return;
    }

    const product = products.find(p => p.id === productId);
    if (!product) {
        console.error(`Product not found with ID ${productId}`);
        return;
    }

    if (newQuantity <= 0) {
        removeFromCart(productId);
    } else if (newQuantity <= product.stock) {
        item.quantity = newQuantity;
        updateCart();
    } else {
        alert('Not enough stock available!');
    }
}

// Checkout function
async function checkout() {
    if (cart.length === 0) {
        alert('Your cart is empty!');
        return;
    }

    try {
        console.log('Sending cart items:', cart); // Debug log
        const response = await fetch('/api/checkout', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(cart)
        });

        if (response.ok) {
            alert('Order placed successfully!');
            // Clear the cart
            cart = [];
            updateCart();
            // Refresh products to update stock
            await fetchProducts();
            // Close the cart offcanvas
            const cartOffcanvas = document.getElementById('cartOffcanvas');
            if (cartOffcanvas) {
                const bsOffcanvas = bootstrap.Offcanvas.getInstance(cartOffcanvas);
                if (bsOffcanvas) {
                    bsOffcanvas.hide();
                }
            }
        } else {
            const error = await response.text();
            console.error('Checkout error:', error); // Debug log
            alert(`Checkout failed: ${error}`);
        }
    } catch (error) {
        console.error('Error during checkout:', error);
        alert('An error occurred during checkout. Please try again.');
    }
}

// Initialize the page
document.addEventListener('DOMContentLoaded', () => {
    fetchProducts();
});

// Test case functions
async function runTest(testName) {
    const testOutput = document.getElementById('testOutput');
    const testButton = event.target;
    const originalText = testButton.innerHTML;
    
    // Disable button and show loading state
    testButton.disabled = true;
    testButton.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> Running...';
    
    try {
        testOutput.innerHTML += `<div class="alert alert-info mt-2">
            <i class="bi bi-info-circle-fill"></i> Starting test: ${testName}
        </div>`;
        
        switch (testName) {
            case 'initializeCatalog':
                await testInitializeCatalog();
                break;
            case 'getProduct':
                await testGetProduct();
                break;
            case 'updateStock':
                await testUpdateStock();
                break;
            case 'createOrder':
                await testCreateOrder();
                break;
            case 'calculateTotal':
                await testCalculateTotal();
                break;
            case 'invalidProduct':
                await testInvalidProduct();
                break;
            case 'negativeStock':
                await testNegativeStock();
                break;
            case 'exceedStock':
                await testExceedStock();
                break;
            case 'tableDriven':
                await testTableDriven();
                break;
        }
    } catch (error) {
        testOutput.innerHTML += `<div class="alert alert-danger mt-2">
            <i class="bi bi-exclamation-triangle-fill"></i> Test failed: ${testName}<br>
            Error: ${error.message}
        </div>`;
        console.error(`Test failed: ${testName}`, error);
    } finally {
        // Restore button state
        testButton.disabled = false;
        testButton.innerHTML = originalText;
    }
}

async function runAllTests() {
    const testOutput = document.getElementById('testOutput');
    const allTestButtons = document.querySelectorAll('.test-cases button');
    testOutput.innerHTML = '<div class="alert alert-info">Running all tests...</div>';
    
    for (const button of allTestButtons) {
        const testName = button.getAttribute('onclick').match(/'([^']+)'/)[1];
        await runTest(testName);
    }
    
    testOutput.innerHTML += '<div class="alert alert-success mt-2">All tests completed!</div>';
}

async function testInitializeCatalog() {
    const testOutput = document.getElementById('testOutput');
    const response = await fetch('/api/products');
    if (!response.ok) throw new Error(`Failed to initialize catalog: ${response.status}`);
    const products = await response.json();
    if (products.length === 0) throw new Error('Catalog is empty');
    testOutput.innerHTML += `<div class="alert alert-success mt-2">
        <i class="bi bi-check-circle-fill"></i> Catalog initialized successfully with ${products.length} products
    </div>`;
}

async function testGetProduct() {
    const testOutput = document.getElementById('testOutput');
    const productId = document.getElementById('productId').value;
    if (!productId) throw new Error('Product ID is required');
    
    const response = await fetch(`/api/products/${productId}`);
    if (!response.ok) throw new Error(`Failed to get product: ${response.status}`);
    const product = await response.json();
    if (!product) throw new Error('Product is null');
    testOutput.innerHTML += `<div class="alert alert-success mt-2">
        <i class="bi bi-check-circle-fill"></i> Product retrieved successfully: ${product.name}
    </div>`;
}

async function testUpdateStock() {
    const testOutput = document.getElementById('testOutput');
    const productId = document.getElementById('productId').value;
    const newStock = document.getElementById('stockUpdate').value;
    
    if (!productId || !newStock) throw new Error('Product ID and new stock value are required');
    
    // First get the current product to verify it exists
    const response = await fetch(`/api/products/${productId}`);
    if (!response.ok) throw new Error(`Failed to get product: ${response.status}`);
    const product = await response.json();
    
    // Update the stock
    const updateResponse = await fetch(`/api/products/${productId}/stock`, {
        method: 'PUT',
        headers: { 
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        },
        body: JSON.stringify({ stock: parseInt(newStock) })
    });
    
    if (!updateResponse.ok) {
        const errorText = await updateResponse.text();
        throw new Error(`Failed to update stock: ${updateResponse.status} - ${errorText}`);
    }
    
    // Verify the update
    const verifyResponse = await fetch(`/api/products/${productId}`);
    if (!verifyResponse.ok) throw new Error(`Failed to verify update: ${verifyResponse.status}`);
    const updatedProduct = await verifyResponse.json();
    
    if (updatedProduct.stock !== parseInt(newStock)) {
        throw new Error(`Stock update verification failed. Expected ${newStock}, got ${updatedProduct.stock}`);
    }
    
    testOutput.innerHTML += `<div class="alert alert-success mt-2">
        <i class="bi bi-check-circle-fill"></i> Stock updated successfully from ${product.stock} to ${newStock}
    </div>`;
}

async function testCreateOrder() {
    const testOutput = document.getElementById('testOutput');
    const productId = document.getElementById('orderProductId').value;
    const quantity = document.getElementById('orderQuantity').value;
    
    if (!productId || !quantity) throw new Error('Product ID and quantity are required');
    
    const response = await fetch(`/api/products/${productId}`);
    if (!response.ok) throw new Error(`Failed to get product: ${response.status}`);
    const product = await response.json();
    
    const orderResponse = await fetch('/api/orders', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ productId: parseInt(productId), quantity: parseInt(quantity) })
    });
    if (!orderResponse.ok) throw new Error(`Failed to create order: ${orderResponse.status}`);
    
    const order = await orderResponse.json();
    testOutput.innerHTML += `<div class="alert alert-success mt-2">
        <i class="bi bi-check-circle-fill"></i> Order created successfully with quantity ${order.quantity}
    </div>`;
}

async function testCalculateTotal() {
    const testOutput = document.getElementById('testOutput');
    const productId = document.getElementById('totalProductId').value;
    const quantity = document.getElementById('totalQuantity').value;
    
    if (!productId || !quantity) throw new Error('Product ID and quantity are required');
    
    // Get the product details
    const response = await fetch(`/api/products/${productId}`);
    if (!response.ok) throw new Error(`Failed to get product: ${response.status}`);
    const product = await response.json();
    
    // Calculate total directly from product price and quantity
    const total = product.price * parseInt(quantity);
    
    testOutput.innerHTML += `<div class="alert alert-success mt-2">
        <i class="bi bi-check-circle-fill"></i> Total calculated successfully: ₹${total.toFixed(2)}<br>
        (Price: ₹${product.price.toFixed(2)} × Quantity: ${quantity})
    </div>`;
}

async function testInvalidProduct() {
    const testOutput = document.getElementById('testOutput');
    const productId = document.getElementById('invalidProductId').value;
    if (!productId) throw new Error('Invalid product ID is required');
    
    const response = await fetch(`/api/products/${productId}`);
    if (response.ok) throw new Error('Expected error for invalid product');
    testOutput.innerHTML += `<div class="alert alert-success mt-2">
        <i class="bi bi-check-circle-fill"></i> Invalid product test passed
    </div>`;
}

async function testNegativeStock() {
    const testOutput = document.getElementById('testOutput');
    const productId = document.getElementById('productId').value;
    const stockValue = document.getElementById('negativeStockValue').value;
    
    if (!productId || !stockValue) throw new Error('Product ID and stock value are required');
    
    const response = await fetch(`/api/products/${productId}/stock`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ stock: parseInt(stockValue) })
    });
    if (response.ok) throw new Error('Expected error for negative stock');
    testOutput.innerHTML += `<div class="alert alert-success mt-2">
        <i class="bi bi-check-circle-fill"></i> Negative stock test passed
    </div>`;
}

async function testExceedStock() {
    const testOutput = document.getElementById('testOutput');
    const productId = document.getElementById('productId').value;
    const quantity = document.getElementById('exceedStockQuantity').value;
    
    if (!productId || !quantity) throw new Error('Product ID and quantity are required');
    
    const response = await fetch(`/api/products/${productId}`);
    if (!response.ok) throw new Error(`Failed to get product: ${response.status}`);
    const product = await response.json();
    
    const orderResponse = await fetch('/api/orders', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ productId: parseInt(productId), quantity: parseInt(quantity) })
    });
    if (orderResponse.ok) throw new Error('Expected error for exceeding stock');
    testOutput.innerHTML += `<div class="alert alert-success mt-2">
        <i class="bi bi-check-circle-fill"></i> Exceed stock test passed
    </div>`;
}

async function testTableDriven() {
    const testOutput = document.getElementById('testOutput');
    const productId = document.getElementById('tableProductId').value;
    if (!productId) throw new Error('Product ID is required');
    
    const tests = [
        { quantity: 1, expected: 40.00 },
        { quantity: 3, expected: 120.00 },
        { quantity: 0, expected: 0.00 }
    ];

    for (const test of tests) {
        const response = await fetch(`/api/products/${productId}`);
        if (!response.ok) throw new Error(`Failed to get product: ${response.status}`);
        const product = await response.json();
        
        const orderResponse = await fetch('/api/orders', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ productId: parseInt(productId), quantity: test.quantity })
        });
        if (!orderResponse.ok) throw new Error(`Failed to create order: ${orderResponse.status}`);
        
        const order = await orderResponse.json();
        const total = order.quantity * product.price;
        
        if (Math.abs(total - test.expected) > 0.01) {
            throw new Error(`Expected ${test.expected}, got ${total}`);
        }
    }
    testOutput.innerHTML += `<div class="alert alert-success mt-2">
        <i class="bi bi-check-circle-fill"></i> Table driven tests passed
    </div>`;
}