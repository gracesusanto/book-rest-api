# API Documentation

## Overview

This document describes the API endpoints related to book and book collections management. There are two main sections of this API:

1. **Book API** which handles:
   - Creating a new book
   - Retrieving a specific book
   - Updating an existing book
   - Deleting a book
   - Patching (partially updating) a book
   - Retrieving all books with optional filtering
   
2. **Collection API** which handles:
   - Creating a new collection
   - Retrieving a specific collection
   - Updating an existing collection
   - Deleting a collection
   - Getting books in a collection
   - Adding a book to a collection
   - Removing a book from a collection

---

## Book API

### Create a Book

- **URL:** /v1/book
- **Method:** POST
- **Data Params:**
  - `title`: Title of the book [string]
  - `author`: Author of the book [string]
  - `published_at`: Publication date of the book [datetime]. It must be in one of the following formats:
    - Date only: "YYYY-MM-DD" (e.g., "2016-01-02"). This is interpreted as midnight (start of the day) in the UTC.
    - Full date-time in UTC: "YYYY-MM-DDTHH:MM:SSZ" (e.g., "2023-07-07T00:00:00Z"). The "T" is a separator indicating the start of the time component, and the "Z" signifies that the date-time is in Coordinated Universal Time (UTC).
  - `edition`: Edition of the book [integer]
  - `description`: Description of the book [string]
  - `genre`: Genre of the book [string]

#### Sample Call:

```bash
curl -X POST -H "Content-Type: application/json" -d "{ book data }" localhost:3000/v1/book
```
#### Response:
Status code: 201
```json
{
    "type": "sync",
    "status": "Success",
    "status_code": 200,
    "error_code": 0,
    "error": "",
    "message": {
        "id": 1,
        "title": "Book Title",
        "author": "Book Author",
        "published_at": "2023-07-07T00:00:00Z",
        "edition": "1",
        "description": "Book description",
        "genre": "Book genre"
    }
}
```

### Get a Book

- **URL:** /v1/book/:id
- **Method:** GET
- **URL Params:** 
  - `id`: ID of the book to retrieve [integer]

#### Sample Call:

```bash
curl -X GET localhost:3000/v1/book/1
```

#### Response:
Status code: 200
```json
{
    "type": "sync",
    "status": "Success",
    "status_code": 200,
    "error_code": 0,
    "error": "",
    "message": {
        "id": 1,
        "title": "Book Title",
        "author": "Book Author",
        "published_at": "2023-07-07T00:00:00Z",
        "edition": 1,
        "description": "Book description",
        "genre": "Book genre"
    }
}
```

### Update a Book

- **URL:** /v1/book/:id
- **Method:** PUT
- **URL Params:** 
  - `id`: ID of the book to update [integer]
- **Data Params:**
  - `title`: Title of the book [string]
  - `author`: Author of the book [string]
  - `published_at`: Publication date of the book [datetime]. It must be in one of the following formats:
    - Date only: "YYYY-MM-DD" (e.g., "2016-01-02"). This is interpreted as midnight (start of the day) in the UTC.
    - Full date-time in UTC: "YYYY-MM-DDTHH:MM:SSZ" (e.g., "2023-07-07T00:00:00Z"). The "T" is a separator indicating the start of the time component, and the "Z" signifies that the date-time is in Coordinated Universal Time (UTC).
  - `edition`: Edition of the book [integer]
  - `description`: Description of the book [string]
  - `genre`: Genre of the book [string]

#### Sample Call:

```bash
curl -X PUT -H "Content-Type: application/json" -d "{ updated book data }" localhost:3000/v1/book/1
```
#### Response:
Status code: 200
```json
{
    "type": "sync",
    "status": "Success",
    "status_code": 200,
    "error_code": 0,
    "error": "",
    "message": {
        "id": 1,
        "title": "Updated Book Title",
        "author": "Updated Book Author",
        "published_at": "2023-07-07T00:00:00Z",
        "edition": 1,
        "description": "Updated Book description",
        "genre": "Updated Book genre"
    }
}
```

### Delete a Book

- **URL:** /v1/book/:id
- **Method:** DELETE
- **URL Params:** 
  - `id`: ID of the book to delete [integer]

#### Sample Call:

```bash
curl -X DELETE localhost:3000/v1/book/1
```

#### Response:
Status code: 204

### Patch a Book

- **URL:** /v1/book/:id
- **Method:** PATCH
- **URL Params:** 
  - `id`: ID of the book to update [integer]
- **Data Params:** 
  - `title` (optional): Title of the book [string]
  - `author` (optional): Author of the book [string]
  - `published_at` (optional): Publication date of the book [datetime]. It must be in one of the following formats:
    - Date only: "YYYY-MM-DD" (e.g., "2016-01-02"). This is interpreted as midnight (start of the day) in the UTC.
    - Full date-time in UTC: "YYYY-MM-DDTHH:MM:SSZ" (e.g., "2023-07-07T00:00:00Z"). The "T" is a separator indicating the start of the time component, and the "Z" signifies that the date-time is in Coordinated Universal Time (UTC).
  - `edition` (optional): Edition of the book [integer]
  - `description` (optional): Description of the book [string]
  - `genre` (optional): Genre of the book [string]

#### Sample Call:

```bash
curl -X PATCH -H "Content-Type: application/json" -d "{ "author" : "Updated Book Author" }" localhost:3000/v1/book/1
```

#### Response:
Status code: 200
```json
{
    "type": "sync",
    "status": "Success",
    "status_code": 200,
    "error_code": 0,
    "error": "",
    "message": {
        "id": 1,
        "title": "Book Title",
        "author": "Updated Book Author",
        "published_at": "2023-07-07T00:00:00Z",
        "edition": 1,
        "description": "Book description",
        "genre": "Book genre"
    }
}
```
### Get All Books

- **URL:** /v1/book
- **Method:** GET
- **URL Params:** 
  - `author` (optional): Author of the books to retrieve [string]
  - `genre` (optional): Genre of the books to retrieve [string]
  - `start` (optional): Start date of the books to retrieve [datetime]
  - `end` (optional): End date of the books to retrieve [datetime]

#### Sample Call:

```bash
curl -X GET localhost:3000/v1/book?author=John&genre=Fiction&start=2022-01-01&end=2022-12-31
```

#### Response:
Status code: 200
```json
{
    "type": "sync",
    "status": "Success",
    "status_code": 200,
    "error_code": 0,
    "error": "",
    "message": [
      {
        "id": 1,
        "title": "Book Title",
        "author": "John",
        "published_at": "2022-07-07T00:00:00Z",
        "edition": 1,
        "description": "Book description",
        "genre": "Book genre"
      },
      {
        "id": 2,
        "title": "Book Title",
        "author": "John Doe",
        "published_at": "2022-10-10T00:00:00Z",
        "edition": 1,
        "description": "Book description",
        "genre": "Book genre"
      },
    ]
}
```

Please note that all datetime parameters should follow the "2006-01-02" format.

---

## Collections API

### Create a Collection

- **URL:** /v1/collection
- **Method:** POST
- **Data Params:**
  - `name`: Name of the collection [string]

#### Sample Call:

```bash
curl -X POST -H "Content-Type: application/json" -d "{ 'name': 'Classic Novels' }" localhost:3000/v1/collection
```

#### Response:
Status code: 201
```json
{
    "type": "sync",
    "status": "Success",
    "status_code": 200,
    "error_code": 0,
    "error": "",
    "message": {
        "id": 1,
        "name": "Classic Novels",
        "books": null,
    }
}
```

### Get a Collection

- **URL:** /v1/collection/:id
- **Method:** GET
- **URL Params:** 
  - `id`: ID of the collection to retrieve [integer]

#### Sample Call:

```bash
curl -X GET localhost:3000/v1/collection/1
```
#### Response:
Status code: 200
```json
{
    "type": "sync",
    "status": "Success",
    "status_code": 200,
    "error_code": 0,
    "error": "",
    "message": {
        "id": 1,
        "name": "Classic Novels",
        "books": [
            {
              "id": 1,
              "title": "Book Title",
              "author": "Book Author",
              "published_at": "2022-07-07T00:00:00Z",
              "edition": 1,
              "description": "Book description",
              "genre": "Book genre"
            },
        ],
    }
}
```

### Get All Collection

- **URL:** /v1/collection
- **Method:** GET

#### Sample Call:

```bash
curl -X GET localhost:3000/v1/collection
```
#### Response:
Status code: 200
```json
{
    "type": "sync",
    "status": "Success",
    "status_code": 200,
    "error_code": 0,
    "error": "",
    "message": [
      {
        "id": 1,
        "name": "Classic Novels",
        "books": [
            {
              "id": 1,
              "title": "Book Title",
              "author": "Book Author",
              "published_at": "2022-07-07T00:00:00Z",
              "edition": 1,
              "description": "Book description",
              "genre": "Book genre"
            },
        ],
      },
      {
        "id": 2,
        "name": "Non Fiction",
        "books": [],
      }
    ]
}
```

### Update a Collection

- **URL:** /v1/collection/:id
- **Method:** PUT
- **URL Params:** 
  - `id`: ID of the collection to update [integer]
- **Data Params:**
  - `name`: Name of the collection [string]

#### Sample Call:

```bash
curl -X PUT -H "Content-Type: application/json" -d "{ 'name': 'Updated Collection Name' }" localhost:3000/v1/collection/1
```
#### Response:
Status code: 200
```json
{
    "type": "sync",
    "status": "Success",
    "status_code": 200,
    "error_code": 0,
    "error": "",
    "message": {
        "id": 1,
        "name": "Updated Collection Name",
        "books": [
            {
              "id": 1,
              "title": "Book Title",
              "author": "Book Author",
              "published_at": "2022-07-07T00:00:00Z",
              "edition": 1,
              "description": "Book description",
              "genre": "Book genre"
            },
        ],
    }
}
```

### Delete a Collection

- **URL:** /v1/collection/:id
- **Method:** DELETE
- **URL Params:** 
  - `id`: ID of the collection to delete [integer]

#### Sample Call:

```bash
curl -X DELETE localhost:3000/v1/collection/1
```

#### Response:
Status code: 204

### Add a Book to a Collection

- **URL:** /v1/collection/:id/book
- **Method:** POST
- **URL Params:** 
  - `id`: ID of the collection to update [integer]
- **Data Params:**
  - `book_id`: ID of the book to add [integer]

#### Sample Call:

```bash
curl -X POST -H "Content-Type: application/json" -d "{ 'book_id': 1 }" localhost:3000/v1/collection/1/book
```

#### Response:
Status code: 202
```json
{
    "type": "sync",
    "status": "Success",
    "status_code": 200,
    "error_code": 0,
    "error": "",
    "message": null
}
```

#### Response:
Status code: 204

### Remove a Book from a Collection

- **URL:** /v1/collection/:id/book/:book_id
- **Method:** DELETE
- **URL Params:** 
  - `id`: ID of the collection [integer]
  - `book_id`: ID of the book to remove [integer]

#### Sample Call:

```bash
curl -X DELETE localhost:3000/v1/collection/1/book/1
```

#### Response:
Status code: 204

### Get Books Collection

- **URL:** /v1/collection/:id/book
- **Method:** GET
- **URL Params:** 
  - `id`: ID of the collection to update [integer]
- **Data Params:**
  - `book_id`: ID of the book to add [integer]

#### Sample Call:

```bash
curl -X GET -H "Content-Type: application/json" -d "{ 'book_id': 1 }" localhost:3000/v1/collection/1/book
```

#### Response:
Status code: 200
```json
{
    "type": "sync",
    "status": "Success",
    "status_code": 200,
    "error_code": 0,
    "error": "",
    "message": [
        {
          "id": 1,
          "title": "Book Title",
          "author": "Book Author",
          "published_at": "2022-07-07T00:00:00Z",
          "edition": 1,
          "description": "Book description",
          "genre": "Book genre"
        },
    ],
}
```

---
## API Structure

All API responses, regardless of status, will follow the same basic structure:

```json
{
    "type": string,
    "status": string,
    "status_code": integer,
    "error_code": integer,
    "error": string,
    "message": object
}
```

### Fields

- `type`: The type of response, either "sync" for successful responses or "error" for responses indicating an error.
- `status`: This will be "Success" if the operation was successful, and "Failure" if there was an error.
- `status_code`: The HTTP status code for the response. This will be 2xx for successful responses. For error responses, the HTTP status code will match the `error_code`.
- `error_code`: If there was an error, this will be the HTTP status code indicating the type of error. If there was no error, this will be 0.
- `error`: If there was an error, this will be a string describing the error. If there was no error, this will be an empty string ("").
- `message`: If the operation was successful, this will be an object containing the relevant data. If there was an error, this will be `null`.

### Success Response

Here's an example of a successful response :

```json
{
    "type": "sync",
    "status": "Success",
    "status_code": 200,
    "message": {}
}
```

### Error Response

Here's an example of an error response:

```json
{
    "type": "error",
    "error_code": 400,
    "error": "Error message",
}
```

Please note that error messages will vary based on the specifics of the error that occurred.

---

## Errors

For all endpoints, if the request is invalid (e.g., a required parameter is missing, the provided ID does not exist, etc.), a `400 Bad Request` error is returned. If an operation fails on the server, a `500 Internal Server Error` is returned.

For `GET`, `PATCH`, `DELETE`, `POST`, `DELETE` operations, if the specified book or collection does not exist, a `404 Not Found` error is returned.
