# API Search

A simple API that uses Meilisearch as the engine for "List" endpoints, it uses an event driven backend to asynchronously update the index

<img width="741" alt="Screenshot 2025-04-23 at 2 17 01 PM" src="https://github.com/user-attachments/assets/038a468c-37c0-42f6-bb3a-645ba01f8e98" />

## Requirements
- Go
- Docker
- Make

## Running the Project

_make sure you have a `.env` file setup (see `.envrc` for an example)_

```sh
make
```

## Endpoints

```http
GET    http://localhost:8080/vehicles       HTTP/1.1
GET    http://localhost:8080/vehicles/:id   HTTP/1.1
POST   http://localhost:8080/vehicles       HTTP/1.1
PUT    http://localhost:8080/vehicles/:id   HTTP/1.1
DELETE http://localhost:8080/vehicles/:id   HTTP/1.1
```

## List Endpoint query params

_*Note: make sure all query params are properly URL encoded*_

| Param | Description | Example |
|-|-|-|
| search | a general purpose search | `search=Bugatti`  |
| order  | a key-value pair to order the results by | `order=created_at:desc`  |
| filter  | filter by a specific query (see meilisearch docs for details | `filter=horse_power>500`  |
| size  | the maximum allowed size per page | `size=50`  |
| page  | the page to query into  | `page=3`  |
