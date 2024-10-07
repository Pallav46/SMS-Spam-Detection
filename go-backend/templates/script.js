document.getElementById('spamForm').addEventListener('submit', function(event) {
    event.preventDefault();

    // Get the input text
    const message = document.getElementById('message').value;

    // Fake API call for demo (replace this with actual API call)
    const result = fakeSpamCheck(message);

    // Display result
    const resultDiv = document.getElementById('result');
    resultDiv.style.display = 'block';

    if (result.isSpam) {
        resultDiv.classList.remove('success');
        resultDiv.classList.add('spam');
        resultDiv.innerHTML = `Prediction: <strong>Spam</strong><br>Probability: ${result.probability}`;
    } else {
        resultDiv.classList.remove('spam');
        resultDiv.classList.add('success');
        resultDiv.innerHTML = `Prediction: <strong>Not Spam</strong><br>Probability: ${result.probability}`;
    }
});

// Fake API function (replace with actual API integration)
function fakeSpamCheck(message) {
    // For demo purposes, if the message contains "lottery" it's considered spam
    if (message.toLowerCase().includes('lottery')) {
        return { isSpam: true, probability: 0.95 };
    } else {
        return { isSpam: false, probability: 0.10 };
    }
}
