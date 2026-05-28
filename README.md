# In-Memory URL Shortener API

A lightweight Bitly-like URL shortener backend built with Go.

The service:
- shortens long URLs
- redirects users using shortcodes
- tracks click analytics
- stores data in memory
- handles concurrent requests safely

---

# Features

- URL shortening
- Random shortcode generation
- Redirect support
- Click analytics
- URL validation
- Thread-safe in-memory storage
- Middleware logging
- Unit tests

---

# Tech Stack

- Go
- net/http
- sync.RWMutex
- JSON APIs
- Go testing package

### Prerequisites
- Go 1.26 or higher

---

# Project Structure

```txt
.
├── cmd/
│   └── server/
│
├── internal/
│   ├── handler/
│   ├── middleware/
│   ├── models/
│   ├── shortener/
│   └── store/
│
├── go.mod
└── README.md
```

---

# API Endpoints

## Health Check

```http
GET /health
```

---

## Shorten URL

```http
POST /shorten
```

Request:

```json
{
  "url": "https://google.com"
}
```

Response:

```json
{
  "short_code": "abc123"
}
```

---

## Redirect

```http
GET /{shortCode}
```

Redirects to original URL.

---

## URL Analytics

```http
GET /stats/{shortCode}
```

Response:

```json
{
  "url": "https://google.com",
  "clicks": 3
}
```

---

# Concepts Learned

- HTTP servers
- Routing
- JSON APIs
- Middleware
- URL validation
- Redirect handling
- Shared application state
- Mutexes and concurrency safety
- Dependency injection
- Encapsulation
- Project structure
- Unit testing
- Backend architecture fundamentals

---

# Future Improvements

- Database persistence
- Docker support
- Graceful shutdown
- Rate limiting
- Request IDs
- Structured logging
- Custom shortcodes
- URL expiration
- Swagger/OpenAPI docs

---

