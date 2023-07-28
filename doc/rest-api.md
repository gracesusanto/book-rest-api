# REST API

This documentation provides an overview of the RESTful API for interactions between the client and the server over HTTP. 

## API versioning

The current version of our API is `v1`.

## Return values

The API generally returns two types of values:

* Standard return value
* Error

### Standard return value

For a standard synchronous operation, the following JSON object is returned:

```js
{
    "type": "sync",
    "status": "Success",
    "status_code": 200,
    "error_code": 0,
    "error": "",
    "message": {}    // Extra resource/action specific data
}
```

HTTP status code is in the range of 200-299.

### List of success status codes

Code  | Meaning   | Methods
:-----| :---------| :-----
200   | OK        | GET, PUT, PATCH
201   | Created   | POST
202   | Accepted  | POST
204   | No Content| DELETE

### Error

There are various situations in which an error may occur, in those cases, the following return value is used:

```js
{
    "type": "error",
    "status": "",
    "status_code": 0,
    "error_code": 400,
    "error": "Error message",
    "message": null
}
```
| Code| Meaning         | 
|-----|-----------------|
| 400 | Bad Request     |
| 404 | Not Found       |
| 500 | Internal Error  |
| 501 | Not Implemented |

## Status codes

The codes are always 3 digits, with the following ranges:

200 to 299: positive action result
400 to 599: negative action result

### List of current status codes

Code  | Meaning
:---  | :------
200   | Success
400   | Failure


## PUT vs PATCH

Our API supports both PUT and PATCH methods to modify existing objects:

- **PUT**: This method replaces the entire object with the new data that you provide. It's typically used after you've retrieved the current object state with a GET request. The entire object data needs to be provided.
  
- **PATCH**: This method allows you to modify a specific part of the object by only specifying the part you want to change. This is beneficial when you only need to update a single field or a subset of fields of an object.

In either case, it's important to handle possible race conditions. If another client modifies the object after you've retrieved it but before you update it, your update might overwrite their changes.

## Filtering

Our API supports filtering on GET requests, which allows you to limit the data returned by the API based on specific criteria. You can filter on several different fields, including:

- `author`: Allows you to filter books based on the author's name.
- `start` and `end`: Allows you to filter books based on the publication date.

Filtering is done by providing the filter criteria in the query parameters of the GET request. The format of a filter is: `?fieldname=value`.

For example, to filter books by author, you could use: `GET /books?author=JohnDoe`.

## API Structure

Our API endpoints are documented in Postman collections which can be found in [`rest-api.postman_collection.json`](https://github.com/gracesusanto/book-rest-api/blob/main/doc/rest-api.postman_collection.json). You can import the Postman collection into your Postman application to explore the API interactively.