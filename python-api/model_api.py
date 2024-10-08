from flask import Flask, request, jsonify
import pickle
import nltk
from nltk.stem import PorterStemmer
from nltk.corpus import stopwords

# Initialize the Flask app
app = Flask(__name__)

# Download necessary NLTK resources
nltk.download('punkt')
nltk.download('stopwords')

# Load the model and vectorizer globally at the start
try:
    with open('model.pkl', 'rb') as model_file:
        model = pickle.load(model_file)
    with open('vectorizer.pkl', 'rb') as vectorizer_file:
        vectorizer = pickle.load(vectorizer_file)
except Exception as e:
    print(f"Error loading model or vectorizer: {e}")
    raise e

# Preprocess function to clean and stem the input text
def preprocess_text(text):
    ps = PorterStemmer()
    text = text.lower()
    text = nltk.word_tokenize(text)

    # Filter out non-alphanumeric characters and stopwords
    y = [ps.stem(i) for i in text if i.isalnum() and i not in stopwords.words('english')]

    return " ".join(y)

# Define the prediction route
@app.route('/predict', methods=['POST'])
def predict():
    data = request.get_json()

    # Check if the text is provided in the request body
    if not data or 'text' not in data:
        return jsonify({'error': 'No text provided'}), 400

    text = data.get('text')

    # Preprocess the input text
    try:
        processed_text = preprocess_text(text)
    except Exception as e:
        return jsonify({'error': 'Error during text preprocessing'}), 500

    # Vectorize the preprocessed text
    try:
        vector_input = vectorizer.transform([processed_text])
    except Exception as e:
        return jsonify({'error': 'Error during text vectorization'}), 500

    # Make the prediction and calculate the probability
    try:
        prediction = model.predict(vector_input)[0]
        probability = model.predict_proba(vector_input)[0][1]  # Probability of being spam
    except Exception as e:
        return jsonify({'error': 'Error during model prediction'}), 500

    # Return the prediction and probability in JSON format
    return jsonify({
        'prediction': "Spam" if prediction == 1 else "Not Spam",
        'probability': round(probability, 4)
    })

# Run the Flask app
if __name__ == '__main__':
    app.run(debug=True)
