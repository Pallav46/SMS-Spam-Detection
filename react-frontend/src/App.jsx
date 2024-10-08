import { useState } from 'react';
import axios from 'axios';
import './App.css';

function App() {
  const [count, setCount] = useState(0);
  const [message, setMessage] = useState('');
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await axios.post('http://localhost:8080/predict', {
        text: message
      });
      setResult(response.data);
    } catch (err) {
      setError('An error occurred while processing the request.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container">
      <h1>Spam Detection</h1>
      <form onSubmit={handleSubmit}>
        <textarea
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder="Enter text to check for spam"
          required
        />
        <button type="submit" disabled={loading}>
          {loading ? 'Checking...' : 'Check Spam'}
        </button>
      </form>
      {result && (
        <div className={`result ${result.prediction === 'Spam' ? 'spam' : 'success'}`}>
          Prediction: <strong>{result.prediction}</strong>
          <br />
          Probability: {result.probability}
        </div>
      )}
      {error && <div className="error">{error}</div>}
    </div>
  );
}

export default App;
