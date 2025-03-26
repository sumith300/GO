// Function to show popup message (reusing from script.js)
function showPopup(message, isError = false) {
    const overlay = document.querySelector('.popup-overlay');
    const content = document.querySelector('.popup-content');
    const messageEl = document.querySelector('.popup-message');
    
    messageEl.textContent = message;
    content.className = 'popup-content' + (isError ? ' popup-error' : ' popup-success');
    overlay.classList.add('active');
    
    document.querySelector('.popup-close').onclick = () => {
        overlay.classList.remove('active');
    };
}

// Function to show loading state
function setLoading(button, isLoading) {
    button.classList.toggle('loading', isLoading);
    button.disabled = isLoading;
}

// Toggle between login and signup forms
document.getElementById('showSignup').addEventListener('click', (e) => {
    e.preventDefault();
    document.getElementById('loginForm').style.display = 'none';
    document.getElementById('signupForm').style.display = 'block';
});

document.getElementById('showLogin').addEventListener('click', (e) => {
    e.preventDefault();
    document.getElementById('signupForm').style.display = 'none';
    document.getElementById('loginForm').style.display = 'block';
});

// Handle login form submission
document.getElementById('userLoginForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const submitButton = e.target.querySelector('button[type="submit"]');
    setLoading(submitButton, true);

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    try {
        const response = await fetch('http://localhost:8080/api/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password })
        });

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.error || 'Login failed');
        }

        // Store session token
        localStorage.setItem('sessionToken', data.token);
        
        // Redirect to main page
        window.location.href = 'index.html';
    } catch (error) {
        showPopup(error.message, true);
    } finally {
        setLoading(submitButton, false);
    }
});

// Handle signup form submission
document.getElementById('userSignupForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const submitButton = e.target.querySelector('button[type="submit"]');
    setLoading(submitButton, true);

    const username = document.getElementById('newUsername').value;
    const password = document.getElementById('newPassword').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    if (password !== confirmPassword) {
        showPopup('Passwords do not match', true);
        setLoading(submitButton, false);
        return;
    }

    try {
        const response = await fetch('http://localhost:8080/api/auth/signup', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password })
        });

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.error || 'Signup failed');
        }

        showPopup('Account created successfully! Please log in.');
        // Switch to login form
        document.getElementById('signupForm').style.display = 'none';
        document.getElementById('loginForm').style.display = 'block';
        document.getElementById('userSignupForm').reset();
    } catch (error) {
        showPopup(error.message, true);
    } finally {
        setLoading(submitButton, false);
    }
});

// Check if user is already logged in
window.addEventListener('load', () => {
    const token = localStorage.getItem('sessionToken');
    if (token && window.location.pathname.endsWith('login.html')) {
        window.location.href = 'index.html';
    }
});