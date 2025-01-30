# Gopher Chat 💬

A simple WebSocket-based chat application powered by OpenAI's GPT-4o model. The backend is built with Go and Gorilla WebSocket, while the frontend is a minimal HTML, CSS, and JavaScript interface.

## Features 🚀

- WebSocket-based real-time communication.
- Chatbot powered by OpenAI GPT-4o Mini.
- Simple frontend with dynamic message display.
- Uses environment variables for secure API key management.

## Prerequisites ⚙️

Make sure you have the following installed:

- [Go](https://go.dev/doc/install) (1.18+ recommended)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)
- [godotenv](https://github.com/joho/godotenv)

## Installation 🛠️

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/gopher-chat.git
   cd gopher-chat
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

3. Create a `.env` file and add your OpenAI API key:

   ```sh
   echo "OPENAI_API_KEY=your-api-key-here" > .env
   ```

## Running the Application ▶️

1. Start the Go server:

   ```sh
   go run main.go
   ```

2. Open a browser and navigate to:

   ```
   http://localhost:8080
   ```

3. Start chatting! The server processes messages via WebSockets and returns AI-generated responses.

## Project Structure 📂

```
/gopher-chat
│── main.go          # Go WebSocket server
│── .env             # Environment variables (API key)
│── styles.css       # Frontend styles
│── index.html       # Chat interface
│── README.md        # Project documentation
```

## API Endpoints 🌐

| Method | Endpoint   | Description                        |
|--------|-----------|------------------------------------|
| GET    | `/`       | Serves the chat UI (index.html)   |
| WS     | `/chat`   | WebSocket endpoint for messaging  |

## Troubleshooting 🔍

### WebSocket Connection Issues
- Ensure the Go server is running on `localhost:8080`.
- Check browser console (`F12 → Console`) for errors.

### OpenAI API Errors
- Ensure your `.env` file contains the correct API key.
- Check if your OpenAI API usage limits have been reached.

## License 📜

This project is open-source and available under the MIT License.
