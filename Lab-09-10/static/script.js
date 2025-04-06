let products = [];
let cart = [];
let wishlist = [];
let filteredProducts = [];

// Search and filter state
let currentSort = 'default';
let searchQuery = '';

// Search function
function searchProducts(query) {
    searchQuery = query.toLowerCase();
    filterAndSortProducts();
}

// Sort function
function sortProducts(sortType) {
    currentSort = sortType;
    filterAndSortProducts();
}

// Filter and sort products
function filterAndSortProducts() {
    filteredProducts = products.filter(product =>
        product.name.toLowerCase().includes(searchQuery) ||
        product.category.toLowerCase().includes(searchQuery)
    );

    switch(currentSort) {
        case 'price-asc':
            filteredProducts.sort((a, b) => a.price - b.price);
            break;
        case 'price-desc':
            filteredProducts.sort((a, b) => b.price - a.price);
            break;
        case 'category':
            filteredProducts.sort((a, b) => a.category.localeCompare(b.category));
            break;
        default:
            filteredProducts = [...filteredProducts];
    }

    displayProducts();
}

// Fetch products from the server
async function fetchProducts() {
    try {
        const response = await fetch('/products');
        products = await response.json();
        displayProducts();
    } catch (error) {
        console.error('Error fetching products:', error);
    }
}

// Display products in the grid
function displayProducts() {
    const productGrid = document.getElementById('productGrid');
    productGrid.innerHTML = '';
    
    const productsToDisplay = filteredProducts.length > 0 ? filteredProducts : products;

    productsToDisplay.forEach(product => {
        const col = document.createElement('div');
        col.className = 'col-12 col-md-6 col-lg-4 mb-4';

        const stockStatus = product.stock > 0 ? 'In Stock' : 'Out of Stock';
        const stockBadgeClass = product.stock > 0 ? 'bg-success' : 'bg-danger';
        const isInWishlist = wishlist.some(item => item.id === product.id);

        col.innerHTML = `
            <div class="card product-card">
                <div class="product-image">
                    <img src="${product.image}" alt="${product.name}" class="product-img">
                </div>
                <span class="badge ${stockBadgeClass} stock-badge">${stockStatus}</span>
                <span class="badge bg-primary category-badge">${product.category}</span>
                <button class="btn btn-link wishlist-btn ${isInWishlist ? 'active' : ''}" onclick="toggleWishlist(${product.id})">
                    <i class="bi ${isInWishlist ? 'bi-heart-fill' : 'bi-heart'}"></i>
                </button>
                <div class="card-body">
                    <h5 class="card-title">${product.name}</h5>
                    <p class="card-text">₹${product.price.toFixed(2)}</p>
                    <div class="quantity-control">
                        <input type="number" class="form-control" value="1" min="1" max="${product.stock}" id="quantity-${product.id}">
                        <button class="btn btn-primary" onclick="addToCart(${product.id})" ${product.stock === 0 ? 'disabled' : ''}>
                            Add to Cart
                        </button>
                    </div>
                </div>
            </div>
        `;

        productGrid.appendChild(col);
    });
}

// Add product to cart
async function addToCart(productId) {
    const product = products.find(p => p.id === productId);
    const quantityInput = document.getElementById(`quantity-${productId}`);
    const quantity = parseInt(quantityInput.value);

    if (!product || quantity <= 0) {
        alert('Invalid product or quantity!');
        return;
    }

    if (quantity > product.stock) {
        alert('Not enough stock available!');
        return;
    }

    const existingItem = cart.find(item => item.product.id === productId);
    if (existingItem) {
        if (existingItem.quantity + quantity > product.stock) {
            alert('Not enough stock available!');
            return;
        }
    }

    try {
        // Verify stock availability from server
        const response = await fetch('/check-stock', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                productId: productId,
                quantity: quantity
            })
        });

        if (!response.ok) {
            const error = await response.text();
            alert(error || 'Failed to verify stock availability');
            return;
        }

        // Update cart
        if (existingItem) {
            existingItem.quantity += quantity;
        } else {
            cart.push({
                product: product,
                quantity: quantity
            });
        }

        // Update product stock locally
        product.stock -= quantity;
        
        // Reset quantity input
        quantityInput.value = 1;
        quantityInput.max = product.stock;
        
        // Update display
        displayProducts();
        updateCartDisplay();
    } catch (error) {
        console.error('Error adding to cart:', error);
        alert('Failed to add item to cart. Please try again.');
    }
}

// Wishlist management functions
function toggleWishlist(productId) {
    const product = products.find(p => p.id === productId);
    const index = wishlist.findIndex(item => item.id === productId);
    
    if (index === -1) {
        wishlist.push(product);
    } else {
        wishlist.splice(index, 1);
    }
    
    displayProducts();
    updateWishlistDisplay();
}

function updateWishlistDisplay() {
    const wishlistItems = document.getElementById('wishlistItems');
    if (!wishlistItems) return;

    wishlistItems.innerHTML = wishlist.map(item => `
        <div class="wishlist-item">
            <div class="wishlist-item-details">
                <h6>${item.name}</h6>
                <p>₹${item.price.toFixed(2)}</p>
            </div>
            <div class="wishlist-item-actions">
                <button class="btn btn-primary btn-sm" onclick="addToCart(${item.id})" ${item.stock === 0 ? 'disabled' : ''}>
                    Add to Cart
                </button>
                <button class="btn btn-danger btn-sm" onclick="toggleWishlist(${item.id})">
                    Remove
                </button>
            </div>
        </div>
    `).join('');
}

// Update cart display with tax calculations
function updateCartDisplay() {
    const cartItems = document.getElementById('cartItems');
    const cartCount = document.getElementById('cartCount');
    const cartTotal = document.getElementById('cartTotal');
    const cartSubtotal = document.getElementById('cartSubtotal');
    const cartTax = document.getElementById('cartTax');
    
    let subtotal = 0;

    if (cart.length === 0) {
        cartItems.innerHTML = `
            <div class="text-center py-5">
                <i class="bi bi-cart-x" style="font-size: 3rem; color: var(--shadow-color);"></i>
                <h5 class="mt-3">Your cart is empty</h5>
                <p class="text-muted">Add some products to your cart to see them here.</p>
            </div>
        `;
    } else {
        cartItems.innerHTML = '';
        cart.forEach(item => {
            const itemTotal = item.product.price * item.quantity;
            subtotal += itemTotal;

            cartItems.innerHTML += `
                <div class="cart-item">
                    <div class="cart-item-image" style="width: 60px; height: 60px; overflow: hidden; border-radius: 8px;">
                        <img src="${item.product.image}" alt="${item.product.name}" style="width: 100%; height: 100%; object-fit: cover;">
                    </div>
                    <div class="cart-item-details">
                        <h6>${item.product.name}</h6>
                        <p class="text-muted mb-0">₹${item.product.price.toFixed(2)} × ${item.quantity}</p>
                    </div>
                    <div class="cart-item-price">₹${itemTotal.toFixed(2)}</div>
                    <button class="btn btn-outline-danger btn-sm ms-2" onclick="removeFromCart(${item.product.id})">
                        <i class="bi bi-trash"></i>
                    </button>
                </div>
            `;
        });
    }

    const tax = subtotal * 0.18; // 18% GST
    const total = subtotal + tax;

    cartCount.textContent = cart.reduce((sum, item) => sum + item.quantity, 0);
    if (cartSubtotal) cartSubtotal.textContent = subtotal.toFixed(2);
    if (cartTax) cartTax.textContent = tax.toFixed(2);
    cartTotal.textContent = total.toFixed(2);
}

// Remove item from cart
function removeFromCart(productId) {
    cart = cart.filter(item => item.product.id !== productId);
    updateCartDisplay();
}

// Process checkout
async function processCheckout() {
    try {
        for (const item of cart) {
            const response = await fetch('/order', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    productId: item.product.id,
                    quantity: item.quantity
                })
            });

            if (!response.ok) {
                throw new Error('Failed to process order');
            }
        }

        alert('Order processed successfully!');
        cart = [];
        updateCartDisplay();
        fetchProducts(); // Refresh product list to update stock
    } catch (error) {
        console.error('Error processing order:', error);
        alert('Failed to process order. Please try again.');
    }
}

// Event listeners
document.addEventListener('DOMContentLoaded', () => {
    fetchProducts();
    document.getElementById('checkoutBtn').addEventListener('click', processCheckout);
    
    // Search functionality
    const searchInput = document.getElementById('searchInput');
    const searchBtn = document.getElementById('searchBtn');
    
    searchInput.addEventListener('keyup', (e) => {
        if (e.key === 'Enter') {
            searchProducts(searchInput.value);
        }
    });
    
    searchBtn.addEventListener('click', () => {
        searchProducts(searchInput.value);
    });
    
    // Sort functionality
    const sortOptions = document.querySelectorAll('.sort-option');
    sortOptions.forEach(option => {
        option.addEventListener('click', (e) => {
            e.preventDefault();
            const sortType = option.getAttribute('data-sort');
            
            // Update active class
            sortOptions.forEach(opt => opt.classList.remove('active'));
            option.classList.add('active');
            
            sortProducts(sortType);
        });
    });
});