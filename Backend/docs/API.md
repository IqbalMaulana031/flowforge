# API Starter

## Base URL

```text
http://localhost:8080
```

## Endpoints

| Method | Path | Description |
|---|---|---|
| GET | `/` | API root status |
| GET | `/v1/health/ping` | Liveness ping |
| GET | `/v1/health` | API, database, dan Redis health status |

## Response Format

Success:

```json
{
  "success": true,
  "data": {}
}
```

Error:

```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "route not found"
  },
  "requestId": "uuid"
}
```
