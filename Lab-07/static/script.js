let cart = [];
let products = [];

// Fetch products from the server
async function fetchProducts() {
    try {
        const response = await fetch('/api/products');
        products = await response.json();
        displayProducts();
    } catch (error) {
        console.error('Error fetching products:', error);
    }
}

// Display products in the grid
function displayProducts() {
    const productGrid = document.getElementById('productGrid');
    productGrid.innerHTML = products.map(product => `
        <div class="col">
            <div class="card product-card position-relative">
                <span class="badge bg-primary category-badge">${product.Category}</span>
                <span class="badge ${product.Stock > 0 ? 'bg-success' : 'bg-danger'} stock-badge">
                    ${product.Stock > 0 ? 'In Stock' : 'Out of Stock'}
                </span>
                <div class="card-body d-flex flex-column">
                    <h5 class="card-title">${product.Name}</h5>
                    <p class="card-text">₹${product.Price.toFixed(2)}</p>
                    <p class="card-text">Stock: ${product.Stock}</p>
                    <div class="quantity-control mb-3">
                        <button class="btn btn-sm btn-outline-secondary" onclick="event.stopPropagation(); updateCardQuantity(${product.ID}, 'decrease')">-</button>
                        <span id="quantity-${product.ID}">1</span>
                        <button class="btn btn-sm btn-outline-secondary" onclick="event.stopPropagation(); updateCardQuantity(${product.ID}, 'increase')">+</button>
                    </div>
                    <button onclick="addToCartWithQuantity(${product.ID})" 
                            class="btn btn-primary mt-auto" 
                            ${product.Stock === 0 ? 'disabled' : ''}>
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
    const product = products.find(p => p.ID === productId);
    let currentQuantity = parseInt(quantityElement.textContent);

    if (action === 'increase' && currentQuantity < product.Stock) {
        currentQuantity++;
    } else if (action === 'decrease' && currentQuantity > 1) {
        currentQuantity--;
    }

    quantityElement.textContent = currentQuantity;
}

function addToCartWithQuantity(productId) {
    const product = products.find(p => p.ID === productId);
    if (!product) return;

    const quantityElement = document.getElementById(`quantity-${productId}`);
    const quantity = parseInt(quantityElement.textContent);

    const existingItem = cart.find(item => item.product.ID === productId);
    if (existingItem) {
        if (existingItem.quantity + quantity <= product.Stock) {
            existingItem.quantity += quantity;
        } else {
            alert('Not enough stock available!');
            return;
        }
    } else {
        if (quantity <= product.Stock) {
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
    cart = cart.filter(item => item.product.ID !== productId);
    updateCart();
}

// Update cart display
function updateCart() {
    const cartItems = document.getElementById('cartItems');
    const cartCount = document.getElementById('cartCount');
    const cartTotal = document.getElementById('cartTotal');

    cartCount.textContent = cart.reduce((total, item) => total + item.quantity, 0);

    cartItems.innerHTML = cart.map(item => `
        <div class="cart-item">
            <div class="d-flex justify-content-between align-items-center mb-2">
                <h6 class="mb-0">${item.product.Name}</h6>
                <button class="btn btn-sm btn-danger" onclick="removeFromCart(${item.product.ID})">&times;</button>
            </div>
            <div class="d-flex justify-content-between align-items-center">
                <div class="quantity-control">
                    <button class="btn btn-sm btn-outline-secondary" onclick="updateQuantity(${item.product.ID}, ${item.quantity - 1})">-</button>
                    <span>${item.quantity}</span>
                    <button class="btn btn-sm btn-outline-secondary" onclick="updateQuantity(${item.product.ID}, ${item.quantity + 1})">+</button>
                </div>
                <span>₹${(item.product.Price * item.quantity).toFixed(2)}</span>
            </div>
        </div>
    `).join('');

    cartTotal.textContent = cart.reduce((total, item) => 
        total + (item.product.Price * item.quantity), 0).toFixed(2);
}

// Update quantity of cart item
function updateQuantity(productId, newQuantity) {
    const item = cart.find(item => item.product.ID === productId);
    if (!item) return;

    const product = products.find(p => p.ID === productId);
    if (!product) return;

    if (newQuantity <= 0) {
        removeFromCart(productId);
    } else if (newQuantity <= product.Stock) {
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
        const response = await fetch('/api/checkout', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(cart)
        });

        if (response.ok) {
            alert('Order placed successfully!');
            cart = [];
            updateCart();
            fetchProducts(); // Refresh products to update stock
        } else {
            const error = await response.text();
            alert(`Checkout failed: ${error}`);
        }
    } catch (error) {
        console.error('Error during checkout:', error);
        alert('An error occurred during checkout. Please try again.');
    }
}

// Initialize the page
fetchProducts();