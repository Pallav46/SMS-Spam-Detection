# SMS Spam Detector

This project implements an SMS spam detection service using machine learning. It includes a Go backend to handle API requests and a simple frontend for user interaction.

## Folder Structure

```
/SMS-Spam-Detector
├── go-backend               # Contains the Go backend for the spam detection service
│   ├── main.go              # Main application file for the Go server
│   ├── go.mod               # Go module file
│   └── templates            # Contains HTML, CSS, and JavaScript files for the frontend
│       └── index.html       # Main HTML file for the frontend interface
├── python-api               # Contains the Python API for the spam detection model
│   └── ...                  # Python API files (if any)
├── spam.csv                 # Dataset for training the spam classification model
└── spam_classification_model.ipynb # Jupyter Notebook for training the spam classification model
```

## Overview

### Features

- **Spam Detection**: Users can enter a message to check whether it is spam or not.
- **User-Friendly Interface**: Simple and responsive design for easy interaction.

### Technologies Used

- **Backend**: Go (Gin framework)
- **Frontend**: HTML, CSS, JavaScript
- **Machine Learning**: Python, Scikit-learn for training the spam detection model

## Installation

1. Clone the repository:
   ```bash
   git clone <repository_url>
   cd SMS-Spam-Detector
   ```

2. Navigate to the `go-backend` directory and install Go if you haven't already:
   ```bash
   cd go-backend
   go mod tidy
   ```

3. Run the Go server:
   ```bash
   go run main.go
   ```

4. Open your web browser and navigate to `http://localhost:8080/` to access the spam detection interface.

## Running the Python API (Optional)

If you wish to run a separate Python API for spam detection, navigate to the `python-api` directory and run the necessary scripts.

## Usage

1. Enter the message you want to check in the text area provided.
2. Click on the "Check Spam" button.
3. The result will be displayed below the form, indicating whether the message is spam or not, along with the probability.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or raise an issue for improvements.

## License

This project is open-source and available under the [MIT License](LICENSE).