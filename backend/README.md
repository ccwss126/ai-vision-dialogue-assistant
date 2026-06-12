# Go Backend

The backend exposes `POST /api/analyze` for the Vue frontend.

## Run

From this directory:

```bash
go run main.go
```

Or from the project root:

```bash
go run ./backend
```

The server listens on `PORT`, defaulting to `8080`.

## Environment

Create `.env` in the project root:

```env
GEMINI_API_KEY=
PORT=8080
HTTP_PROXY=http://127.0.0.1:7890
HTTPS_PROXY=http://127.0.0.1:7890
```

If `GEMINI_API_KEY` is empty, the backend returns Demo Mode responses. If Gemini network requests fail, it also falls back to Demo Mode.
