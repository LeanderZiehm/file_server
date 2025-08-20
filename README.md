# File Server

This is a simple file server application with a Go backend and a React frontend.

## Features

*   Upload files
*   Download files
*   Delete files
*   List uploaded files

## Prerequisites

Before you begin, ensure you have the following installed:

*   [Go](https://golang.org/doc/install)
*   [Node.js](https://nodejs.org/en/download/) (which includes npm)

## Running the Application for Development

1.  **API Key**: The backend requires an API key to be set as an environment variable. You can create a `.env` file in the `backend` directory with the following content:

    ```
    API_KEY=your-secret-api-key
    ```

2.  **Start the application**: A convenience script is provided to start both the backend and frontend servers concurrently. Run the following command from the root of the project:

    ```bash
    ./start.sh
    ```

This will:
*   Start the Go backend server on `http://localhost:8080`.
*   Start the React frontend development server on `http://localhost:5173`.
