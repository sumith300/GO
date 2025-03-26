// Global state to store current account
let currentAccount = null;

// Check authentication on page load
function checkAuth() {
    const token = localStorage.getItem('sessionToken');
    if (!token) {
        window.location.href = 'login.html';
        return;
    }
}

// Add authorization header to fetch requests
function getAuthHeaders() {
    return {
        'Content-Type': 'application/json',
        'Authorization': localStorage.getItem('sessionToken')
    };
}

// Function to show popup message
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

// Event listener for form submission
document.getElementById('createAccountForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const submitButton = e.target.querySelector('button[type="submit"]');
    setLoading(submitButton, true);

    const accountType = document.getElementById('accountType').value;
    const initialBalance = parseFloat(document.getElementById('initialBalance').value);
    const accountId = document.getElementById('accountId').value;

    try {
        const response = await fetch('http://localhost:8080/api/accounts', {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify({
                accountType: parseInt(accountType),
                initialBalance,
                accountId
            })
        });

        if (response.status === 401) {
            window.location.href = 'login.html';
            return;
        }

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.error || 'Failed to create account');
        }

        showPopup('Account created successfully!');
        document.getElementById('createAccountForm').reset();
        updateAccountList();
    } catch (error) {
        showPopup(error.message, true);
    } finally {
        setLoading(submitButton, false);
    }
});

// Function to update the list of accounts
async function updateAccountList() {
    try {
        const response = await fetch('http://localhost:8080/api/accounts', {
            headers: getAuthHeaders()
        });

        if (response.status === 401) {
            window.location.href = 'login.html';
            return;
        }

        const accounts = await response.json();
        const accountListContent = document.getElementById('accountListContent');
        accountListContent.innerHTML = '';

        accounts.forEach(account => {
            const accountDiv = document.createElement('div');
            accountDiv.textContent = `${account.id} (${account.type}): ₹${account.balance.toFixed(2)}`;
            accountDiv.onclick = () => selectAccount(account);
            accountListContent.appendChild(accountDiv);
        });
    } catch (error) {
        console.error('Error fetching accounts:', error);
        if (error.message.includes('Unauthorized')) {
            window.location.href = 'login.html';
        }
    }
}

// Function to select and display account details
function selectAccount(account) {
    currentAccount = account;
    document.getElementById('accountForm').style.display = 'none';
    document.getElementById('accountManagement').style.display = 'block';
    document.getElementById('displayAccountType').textContent = account.type;
    document.getElementById('displayBalance').textContent = account.balance.toFixed(2);
    viewTransactionHistory();
}

// Function to handle deposits
async function handleDeposit() {
    const amount = parseFloat(document.getElementById('amount').value);
    if (!amount || amount <= 0) {
        showPopup('Please enter a valid amount', true);
        return;
    }

    const depositButton = document.querySelector('button[onclick="handleDeposit()"]');
    setLoading(depositButton, true);

    try {
        const response = await fetch(`http://localhost:8080/api/accounts/${currentAccount.id}/deposit`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify({ amount })
        });

        if (response.status === 401) {
            window.location.href = 'login.html';
            return;
        }

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.error || 'Failed to deposit');
        }

        document.getElementById('displayBalance').textContent = data.balance.toFixed(2);
        document.getElementById('amount').value = '';
        viewTransactionHistory();
        updateAccountList();
        showPopup('Deposit successful!');
    } catch (error) {
        showPopup(error.message, true);
    } finally {
        setLoading(depositButton, false);
    }
}

// Function to handle withdrawals
async function handleWithdraw() {
    const amount = parseFloat(document.getElementById('amount').value);
    if (!amount || amount <= 0) {
        showPopup('Please enter a valid amount', true);
        return;
    }

    const withdrawButton = document.querySelector('button[onclick="handleWithdraw()"]');
    setLoading(withdrawButton, true);

    try {
        const response = await fetch(`http://localhost:8080/api/accounts/${currentAccount.id}/withdraw`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify({ amount })
        });

        if (response.status === 401) {
            window.location.href = 'login.html';
            return;
        }

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.error || 'Failed to withdraw');
        }

        document.getElementById('displayBalance').textContent = data.balance.toFixed(2);
        document.getElementById('amount').value = '';
        viewTransactionHistory();
        updateAccountList();
        showPopup('Withdrawal successful!');
    } catch (error) {
        showPopup(error.message, true);
    } finally {
        setLoading(withdrawButton, false);
    }
}

// Function to view transaction history
async function viewTransactionHistory() {
    try {
        const response = await fetch(`http://localhost:8080/api/accounts/${currentAccount.id}/history`, {
            headers: getAuthHeaders()
        });

        if (response.status === 401) {
            window.location.href = 'login.html';
            return;
        }

        const transactions = await response.json();
        const transactionList = document.getElementById('transactionList');
        transactionList.innerHTML = '';

        transactions.forEach(transaction => {
            const li = document.createElement('li');
            li.textContent = transaction;
            transactionList.appendChild(li);
        });
    } catch (error) {
        console.error('Error fetching transaction history:', error);
        if (error.message.includes('Unauthorized')) {
            window.location.href = 'login.html';
        }
    }
}

// Function to go back to main menu
function backToMain() {
    document.getElementById('accountForm').style.display = 'block';
    document.getElementById('accountManagement').style.display = 'none';
    currentAccount = null;
}

// Check authentication and load account list
checkAuth();
updateAccountList();