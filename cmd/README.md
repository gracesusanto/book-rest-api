## Book CLI Documentation 

`bookcli` is a command-line interface (CLI) for managing a book database. The CLI has commands for adding, getting, updating, and deleting books. Here is a description of each command:

## Installation

As this is a Go application, you can build it using the Go command:

```shell
go build -o bookcli
```

This will create an executable called `bookcli`.

## Available Commands

The following commands are available:

- `add` - Adds a new book.
- `get` - Retrieves existing books.
- `update` - Updates an existing book.
- `delete` - Deletes an existing book.

### Add a Book 

`add` command is used to add a book to the database.

**Usage:**

```bash
$ ./bookcli add --title [title] --author [author] --published-at [published_date] [--edition [edition]] [--description [description]] [--genre [genre]]
```

**Example:**

```bash
# Add the book "To Kill a Mockingbird" by Harper Lee, published on 1960-07-11
$ ./bookcli add --title "To Kill a Mockingbird" --author "Harper Lee" --published-at "1960-07-11"
Book added successfully
```

### Get Books

`get` command is used to retrieve books from the database.

**Usage:**

```bash
$ ./bookcli get --author [author]
```

**Example:**

```bash
# Get books authored by "Harper Lee"
$ ./bookcli get --author "Harper Lee"
+----+-----------------------+------------+--------------+---------+-------------+-------+
| ID | TITLE                 | AUTHOR     | PUBLISHED AT | EDITION | DESCRIPTION | GENRE |
+----+-----------------------+------------+--------------+---------+-------------+-------+
|  1 | To Kill a Mockingbird | Harper Lee | 1960-07-10   |         |             |       |
+----+-----------------------+------------+--------------+---------+-------------+-------+
```

### Update a Book 

The update command is used to selectively update the information of a book in the database. You only need to provide the fields that you want to update.

**Usage:**

```bash
$ ./bookcli update --id [id] [--title [title]] [--author [author]] [--published-at [published_date]] [--edition [edition]] [--description [description]] [--genre [genre]]
```

**Example:**

```bash
# Update the book with id 1 to add edition and genre information
$ ./bookcli update --id 1 --edition "First Edition" --genre "Fiction"
Book updated successfully
```

### Delete a Book 

`delete` command is used to remove a book from the database.

**Usage:**

```bash
$ ./bookcli delete --id [id]
```

**Example:**

```bash
# Deleting the book "To Kill a Mockingbird"
$ ./bookcli delete --id 1
Book deleted successfully
```

In these examples, we first add the book "To Kill a Mockingbird" to the database. We then retrieve all books authored by "Harper Lee". We update the edition and genre information for the book and finally, we delete the book from the database.

## Error Handling

In the event of an error, the CLI will print a message to standard error and exit with a non-zero status code.
