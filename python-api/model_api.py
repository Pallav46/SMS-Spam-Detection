from flask import Flask, request, jsonify
import pickle
import nltk
from nltk.stem import PorterStemmer
from nltk.corpus import stopwords
import string
import logging
from datetime import datetime

app = Flask(__name__)

# Configure logging
logging.basicConfig(filename='api.log', level=logging.INFO, format='%(asctime)s - %(message)s')

nltk.download('punkt')
nltk.download('stopwords')

# Load the model and vectorizer once, globally
with open('model.pkl', 'rb') as model_file, open('vectorizer.pkl', 'rb') as vectorizer_file:
    model = pickle.load(model_file)
    vectorizer = pickle.load(vectorizer_file)

def preprocess_text(text):
    ps = PorterStemmer()
    text = text.lower()
    text = nltk.word_tokenize(text)

    y = []
    for i in text:
        if i.isalnum():
            y.append(i)

    text = y[:]
    y.clear()

    for i in text:
        if i not in stopwords.words('english') and i not in string.punctuation:
            y.append(ps.stem(i))

    return " ".join(y)

@app.route('/predict', methods=['POST'])
def predict():
    data = request.get_json()
    
    if not data or 'text' not in data:
        return jsonify({'error': 'No text provided'}), 400

    text = data.get('text')

    # Log the received text
    logging.info(f"Received text: {text}")

    # Preprocess the text
    processed_text = preprocess_text(text)

    # Vectorize the processed text
    vector_input = vectorizer.transform([processed_text])

    # Predict the result
    prediction = model.predict(vector_input)[0]
    probability = model.predict_proba(vector_input)[0][1]

    # Log the prediction and probability
    logging.info(f"Prediction: {'Spam' if prediction == 1 else 'Not Spam'}, Probability: {probability:.4f}")

    # Return the result as JSON
    return jsonify({
        'prediction': "Spam" if prediction == 1 else "Not Spam",
        'probability': round(probability, 4)
    })

if __name__ == '__main__':
    app.run(debug=True)
