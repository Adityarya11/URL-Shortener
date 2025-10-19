// frontend/app.js
// IMPORTANT: set API_BASE_URL to your deployed Go backend (Render) URL.
// For local testing use: "http://localhost:8000"
const API_BASE_URL = "https://url-shortener-5qu1.onrender.com";

const form = document.getElementById('shorten-form');
const longUrlInput = document.getElementById('long-url');
const customCodeInput = document.getElementById('custom-code');
const resultDiv = document.getElementById('result');
const errorDiv = document.getElementById('error');

form.addEventListener('submit', async (event) => {
    event.preventDefault();
    resultDiv.textContent = '';
    errorDiv.textContent = '';

    const longUrl = longUrlInput.value.trim();
    const customCode = customCodeInput.value.trim();

    if (!longUrl) {
        errorDiv.textContent = 'Please enter a URL.';
        return;
    }

    // Build body matching backend expectation: { url, customCode }
    const body = { url: longUrl };
    if (customCode) body.customCode = customCode;

    try {
        const res = await fetch(`${API_BASE_URL}/shorten`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(body)
        });

        const data = await res.json();

        if (!res.ok) {
            // backend should return reasonable error JSON or text
            const msg = data.error || data.message || (data && (data.error || JSON.stringify(data))) || 'Server error';
            throw new Error(msg);
        }

        // backend returns { shortCode, originalUrl, ... }
        const shortCode = data.shortCode || data.short_code || data.short; // tolerant
        const shortUrl = `${API_BASE_URL}/${shortCode}`;
        resultDiv.innerHTML = `Success! Short URL: <a href="${shortUrl}" target="_blank" rel="noopener noreferrer">${shortUrl}</a>`;
        longUrlInput.value = '';
        customCodeInput.value = '';
    } catch (err) {
        errorDiv.textContent = `Error: ${err.message}`;
        console.error(err);
    }
});
