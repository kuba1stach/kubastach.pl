# kubastach.pl

Site about me.

## Backend API

Implements endpoints defined in `openapi.yaml` backed by Azure Cosmos DB.

### Endpoints (base path `/api/v1`)

* `GET /posts` - list posts
* `GET /posts/{date}` - post for given date (YYYY-MM-DD)
* `GET /posts/dates` - list distinct dates with posts

### Azure Cosmos Data Assumptions

Documents live in a single container and have a `category` field equal to `post` to be included. Each document contains:

```jsonc
{
	"category": "post",
	"content": "markdown or text",
	"date": "2025-09-16",
	"media": [ { "name": "", "type": "" } ]
}
```

### Environment Variables

Set these before running the server:

* `COSMOS_ENDPOINT` - Cosmos DB endpoint URL
* `COSMOS_KEY` - Primary key / auth key
* `COSMOS_DB` - Database name
* `COSMOS_CONTAINER` - Container (collection) name
* `PORT` (optional) - Port to listen on (default 8080)

### Run

```sh
cd backend/cmd/server
go run .
```

Or build:

```sh
go build -o server ./backend/cmd/server
./server
```

### OpenAPI

See `openapi.yaml` for the contract.
