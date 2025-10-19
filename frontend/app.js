// *** IMPORTANT ***
// Change this to your backend's URL once it's deployed!
// For local testing: const API_BASE_URL = "http://localhost:8080"; 
// For production: const API_BASE_URL = "https://your-backend-url.onrender.com";

const API_BASE_URL = "https://your-backend-url-goes-here.onrender.com";

// Get the elements from the HTML
const form = document.getElementById('shorten-form');
const longUrlInput = document.getElementById('long-url');
const customCodeInput = document.getElementById('custom-code');
const resultDiv = document.getElementById('result');
const errorDiv = document.getElementById('error');

form.addEventListener('submit', async (event) => {
    event.preventDefault(); // Stop the form from reloading the page

    // Clear previous results
    resultDiv.textContent = '';
    errorDiv.textContent = '';

    const longUrl = longUrlInput.value;
    const customCode = customCodeInput.value;

    // Prepare the data to send
    const body = {
        long_url: longUrl,
    };

    if (customCode) {
        body.custom_code = customCode;
    }

    try {
        // Use the fetch() API to send a POST request
        const response = await fetch(`${API_BASE_URL}/shorten`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(body),
        });

        const data = await response.json();

        if (!response.ok) {
            // If the server returns an error (like 400 or 500)
            throw new Error(data.error || 'Something went wrong');
        }

        // Success! Display the short URL
        const shortUrl = `${API_BASE_URL}/${data.short_code}`;
        resultDiv.innerHTML = `Success! Your short URL is: <a href="${shortUrl}" target="_blank">${shortUrl}</a>`;

        // Clear the inputs
        longUrlInput.value = '';
        customCodeInput.value = '';

    } catch (err) {
        // Show the error message
        errorDiv.textContent = `Error: ${err.message}`;
    }
});