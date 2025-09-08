document.addEventListener('DOMContentLoaded', function() {
    const orderIdInput = document.getElementById('order-id');
    const lookupBtn = document.getElementById('lookup-btn');
    const resultDiv = document.getElementById('result');

    // Function to fetch order information
    async function fetchOrderInfo(orderId) {
        try {
            // Clear previous results
            resultDiv.style.display = 'none';
            resultDiv.className = 'result';
            
            // Validate input
            if (!orderId) {
                showError('Please enter an order ID');
                return;
            }

            // Make API request
            const response = await fetch(`http://localhost:8080/order/${orderId}`);
            
            if (response.ok) {
                const data = await response.json();
                showSuccess(JSON.stringify(data, null, 2));
            } else if (response.status === 400) {
                showError('Bad request - please check the order ID format');
            } else if (response.status === 500) {
                showError('Server error - unable to retrieve order information');
            } else {
                showError(`Unexpected error: ${response.status} ${response.statusText}`);
            }
        } catch (error) {
            showError(`Network error: ${error.message}`);
        }
    }

    // Function to display success message
    function showSuccess(message) {
        resultDiv.innerHTML = `<pre>${message}</pre>`;
        resultDiv.classList.add('success');
        resultDiv.style.display = 'block';
    }

    // Function to display error message
    function showError(message) {
        resultDiv.innerHTML = `<strong>Error:</strong> ${message}`;
        resultDiv.classList.add('error');
        resultDiv.style.display = 'block';
    }

    // Event listener for button click
    lookupBtn.addEventListener('click', function() {
        const orderId = orderIdInput.value.trim();
        fetchOrderInfo(orderId);
    });

    // Event listener for Enter key in input field
    orderIdInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            const orderId = orderIdInput.value.trim();
            fetchOrderInfo(orderId);
        }
    });
});