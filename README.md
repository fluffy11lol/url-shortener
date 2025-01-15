# URL Shortener

The **URL Shortener** is a web application designed to simplify 
the process of sharing long URLs by converting them into short, 
easy-to-remember aliases. Users can generate short aliases for 
lengthy links and use these aliases to redirect to the original 
URLs. The project is built using **Go (Golang)** and utilizes 
**SQLite** for efficient data storage.

## Key Features

- **Short URL Creation**: Users can submit a long URL and receive a short alias.
- **Redirection via Alias**: When accessing the short alias, users will be redirected to the original URL.
- **Link Deletion**: Users can delete existing aliases.
- **Basic Authentication**: Authentication is required to access the API.

## Technologies

- **Programming Language**: Go (Golang)
- **Database**: SQLite
- **Web Framework**: Chi (router and middleware)
- **Logging**: slog (structured logging)
- **Testing**: httpexpect, testify

## Project Structure

- **`cmd/url-shortener`**: Main package for running the application.
- **`internal/config`**: Application configuration (loaded from a YAML file).
- **`internal/http-server`**: HTTP server with request handlers.
- **`internal/storage`**: Database logic (SQLite).
- **`pkg/logger`**: Custom logger with support for different logging levels.
- **`pkg/random`**: Random alias generation.
- **`tests`**: Integration tests for the API.

## How to Run the Project

### Requirements

- Installed Go (version 1.23 or higher).
- SQLite3 (for data storage).

### Installation and Running

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/url-shortener.git
   cd url-shortener
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```
3. Create a configuration file config/local.yaml

   ```yaml
   env: "local"
   storage_path: "./storage/storage.db"
   http_server:
   host: "localhost"
   port: "8081"
   timeout: "5s"
   idle_timeout: "60s"
   user: "admin"
   password: "admin"
   ```

4. Run the application:

   ```bash
   go run cmd/url-shortener/main.go
   ```

5. The application will be available at: http://localhost:8081

#### You can also use Docker to run the application:

```bash
docker build -t url-shortener .
docker run -p 8081:8081 url-shortener
```

### Example Requests

Creating a Short URL

```bash
curl -X POST http://localhost:8081/url \
-H "Content-Type: application/json" \
-u admin:admin \
-d '{"url": "https://example.com", "alias": "example"}'
```

Response:

```json
{
  "status": "OK",
  "alias": "example"
} 
```

#### Access the link http://localhost:8081/example

#### You will be redirected to https://example.com.

Deleting a URL

```bash
curl -X DELETE http://localhost:8081/url/example \
-u admin:admin
```

Response:

```json
{
  "status": "OK"
}
```

## Logging

Logging is configured using the slog library. Depending on the environment (local, dev, prod), logs may be in text or
JSON format. In local mode, logs are pretty printed to the console.