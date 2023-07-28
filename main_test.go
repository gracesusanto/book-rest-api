package main_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"grace/api"
	"grace/database"
	"grace/server"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDatabase() (*gorm.DB, string) {
	// Generate a unique name for the test database
	testDBName := fmt.Sprintf("test_%s", strings.Replace(uuid.New().String(), "-", "", -1))

	// Set up connection to test database
	mainDB := setupDatabase("app")

	// Create the test database
	if err := mainDB.Exec(fmt.Sprintf("CREATE DATABASE %s;", testDBName)).Error; err != nil {
		log.Fatalf("Failed to create test database: %v", err)
	}

	// Set up connection to test database
	testDB := setupDatabase(testDBName)

	// Run migrations
	server.MigrateDatabase(testDB)

	return testDB, testDBName
}

func cleanupTestDatabase(testDB *gorm.DB, dbName string) {
	// Close all connections to the test database
	sqlDB, err := testDB.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying SQL database: %v", err)
	}

	sqlDB.Close()

	// Connect to the main database (replace 'mainDB' with your actual main database name)
	mainDB := setupDatabase("app")

	// Drop the test database
	if err := mainDB.Exec(fmt.Sprintf("DROP DATABASE %s;", dbName)).Error; err != nil {
		log.Fatalf("Failed to drop test database: %v", err)
	}
}

func TestAddBook(t *testing.T) {
	// Setup the test database
	db, dbName := setupTestDatabase()
	defer cleanupTestDatabase(db, dbName)

	// Start the REST server with the test database
	httpServer := server.StartRESTServer(db)

	// Create new book JSON
	bookJSON := `{"title": "New Test Book", "author": "New Test Author", "published_at": "2015-01-01"}`

	// Perform the POST operation
	req, _ := http.NewRequest("POST", "/v1/book", strings.NewReader(bookJSON))
	rr := httptest.NewRecorder()
	httpServer.Handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusCreated, rr.Code)

	var result map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal("Failed to decode the response body: ", err)
	}

	message, ok := result["message"].(map[string]interface{})
	if !ok {
		t.Fatal("Could not convert 'message' to map[string]interface{}")
	}

	title, ok := message["title"].(string)
	if !ok {
		t.Fatal("Could not convert 'title' to string")
	}

	author, ok := message["author"].(string)
	if !ok {
		t.Fatal("Could not convert 'author' to string")
	}

	require.Equal(t, "New Test Book", title)
	require.Equal(t, "New Test Author", author)

	var newBook database.Book
	if err := db.Where("title = ? AND author = ?", "New Test Book", "New Test Author").First(&newBook).Error; err != nil {
		t.Fatal("Failed to find the new book in the database: ", err)
	}
}

func TestGetBook(t *testing.T) {
	// Setup the test database
	db, dbName := setupTestDatabase()
	defer cleanupTestDatabase(db, dbName)

	// Create a new book directly in the database
	book := &database.Book{Title: "Test Book", Author: "Test Author"} // Replace with your actual Book model and values
	db.Save(book)

	// Start the REST server with the test database
	httpServer := server.StartRESTServer(db)

	// Perform the GET operation
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/book/%d", book.ID), nil)
	rr := httptest.NewRecorder()
	httpServer.Handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var result map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		fmt.Println("Failed to decode the response body: ", err)
		return
	}

	message, ok := result["message"].(map[string]interface{})
	if !ok {
		fmt.Println("Could not convert 'message' to map[string]interface{}")
		return
	}

	title, ok := message["title"].(string)
	if !ok {
		fmt.Println("Could not convert 'title' to string")
		return
	}

	author, ok := message["author"].(string)
	if !ok {
		fmt.Println("Could not convert 'author' to string")
		return
	}
	require.Equal(t, title, "Test Book")
	require.Equal(t, author, "Test Author")
}

func TestGetBooks(t *testing.T) {
	// Setup the test database
	db, dbName := setupTestDatabase()
	defer cleanupTestDatabase(db, dbName)

	// Create new books directly in the database
	books := []*database.Book{
		{Title: "Title 1", Author: "Author 1"},
		{Title: "Title 2", Author: "Author 2"},
		{Title: "Title 3", Author: "Author 3"},
	}

	for _, book := range books {
		db.Save(book)
	}

	// Start the REST server with the test database
	httpServer := server.StartRESTServer(db)

	// Perform the GET operation
	req, _ := http.NewRequest("GET", "/v1/book", nil)
	rr := httptest.NewRecorder()
	httpServer.Handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var result map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode the response body: %v", err)
	}

	// Extract 'message' array from response
	message, ok := result["message"].([]interface{})
	if !ok {
		t.Fatalf("Could not convert 'message' to []interface{}")
	}

	// Check number of books returned
	require.Equal(t, len(books), len(message))

	// Check returned book details
	for i, bookInterface := range message {
		bookMap, ok := bookInterface.(map[string]interface{})
		if !ok {
			t.Fatalf("Could not convert bookInterface to map[string]interface{}")
		}

		// Extract 'title' and 'author'
		title, ok := bookMap["title"].(string)
		if !ok {
			t.Fatalf("Could not convert 'title' to string")
		}

		author, ok := bookMap["author"].(string)
		if !ok {
			t.Fatalf("Could not convert 'author' to string")
		}

		// Check if returned details match created book details
		require.Equal(t, books[i].Title, title)
		require.Equal(t, books[i].Author, author)
	}
}

func TestUpdateBook(t *testing.T) {
	// Setup the test database
	db, dbName := setupTestDatabase()
	defer cleanupTestDatabase(db, dbName)

	// Create a new book directly in the database
	book := &database.Book{Title: "Test Book", Author: "Test Author", PublishedAt: time.Now()}
	db.Save(book)

	// Start the REST server with the test database
	httpServer := server.StartRESTServer(db)

	// Perform the UPDATE operation
	bookJSON := `{"title": "Updated Test Book","description": "Updated Test Description"}`
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/v1/book/%d", book.ID), strings.NewReader(bookJSON))
	rr := httptest.NewRecorder()
	httpServer.Handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var result map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		fmt.Println("Failed to decode the response body: ", err)
		return
	}

	message, ok := result["message"].(map[string]interface{})
	if !ok {
		fmt.Println("Could not convert 'message' to map[string]interface{}")
		return
	}

	title, ok := message["title"].(string)
	if !ok {
		fmt.Println("Could not convert 'title' to string")
		return
	}

	author, ok := message["author"].(string)
	if !ok {
		fmt.Println("Could not convert 'author' to string")
		return
	}

	description, ok := message["description"].(string)
	if !ok {
		fmt.Println("Could not convert 'author' to string")
		return
	}
	require.Equal(t, title, "Updated Test Book")
	require.Equal(t, author, "Test Author")
	require.Equal(t, description, "Updated Test Description")

	var updatedBook database.Book
	db.First(&updatedBook, book.ID)
	require.Equal(t, updatedBook.Title, "Updated Test Book")
	require.Equal(t, updatedBook.Author, "Test Author")
	require.Equal(t, updatedBook.Description, "Updated Test Description")
}

func TestDeleteBook(t *testing.T) {
	// Setup the test database
	db, dbName := setupTestDatabase()
	defer cleanupTestDatabase(db, dbName)

	// Create a new book directly in the database
	book := &database.Book{Title: "Test Book", Author: "Test Author"} // Replace with your actual Book model and values
	db.Save(book)

	// Start the REST server with the test database
	httpServer := server.StartRESTServer(db)

	// Perform the DELETE operation
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/book/%d", book.ID), nil)
	rr := httptest.NewRecorder()
	httpServer.Handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusNoContent, rr.Code)

	expected, err := buildSuccessResponse(nil)
	if err != nil {
		t.Fatal(err)
	}
	require.JSONEq(t, expected, rr.Body.String())

	var deletedBook database.Book
	err = db.First(&deletedBook, 1).Error
	// Assert that the error is a "record not found" error
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestCreateCollection(t *testing.T) {
	db, dbName := setupTestDatabase()
	defer cleanupTestDatabase(db, dbName)

	httpServer := server.StartRESTServer(db)

	collectionJSON := `{"name": "Classic Novels"}`

	req, _ := http.NewRequest("POST", "/v1/collection", strings.NewReader(collectionJSON))
	rr := httptest.NewRecorder()
	httpServer.Handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusCreated, rr.Code)

	var result map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal("Failed to decode the response body: ", err)
	}

	message, ok := result["message"].(map[string]interface{})
	if !ok {
		t.Fatal("Could not convert 'message' to map[string]interface{}")
	}

	name, ok := message["name"].(string)
	if !ok {
		t.Fatal("Could not convert 'name' to string")
	}

	require.Equal(t, "Classic Novels", name)

	var newCollection database.Collection
	if err := db.Where("name = ?", "Classic Novels").First(&newCollection).Error; err != nil {
		t.Fatal("Failed to find the new collection in the database: ", err)
	}
	require.Equal(t, newCollection.Name, "Classic Novels")
}

func TestGetCollection(t *testing.T) {
	db, dbName := setupTestDatabase()
	defer cleanupTestDatabase(db, dbName)

	collection := &database.Collection{Name: "Classic Novels"}
	db.Save(collection)

	httpServer := server.StartRESTServer(db)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/collection/%d", collection.ID), nil)
	rr := httptest.NewRecorder()
	httpServer.Handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var result map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal("Failed to decode the response body: ", err)
	}

	message, ok := result["message"].(map[string]interface{})
	if !ok {
		t.Fatal("Could not convert 'message' to map[string]interface{}")
	}

	name, ok := message["name"].(string)
	if !ok {
		t.Fatal("Could not convert 'name' to string")
	}

	require.Equal(t, name, collection.Name)
}

func TestUpdateCollection(t *testing.T) {
	db, dbName := setupTestDatabase()
	defer cleanupTestDatabase(db, dbName)

	collection := &database.Collection{Name: "Classic Novels"}
	db.Save(collection)

	httpServer := server.StartRESTServer(db)

	collectionJSON := `{"name": "Updated Collection Name"}`
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/v1/collection/%d", collection.ID), strings.NewReader(collectionJSON))
	rr := httptest.NewRecorder()
	httpServer.Handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var result map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal("Failed to decode the response body: ", err)
	}

	message, ok := result["message"].(map[string]interface{})
	if !ok {
		t.Fatal("Could not convert 'message' to map[string]interface{}")
	}

	name, ok := message["name"].(string)
	if !ok {
		t.Fatal("Could not convert 'name' to string")
	}

	require.Equal(t, "Updated Collection Name", name)

	var updatedCollection database.Collection
	db.First(&updatedCollection, collection.ID)
	require.Equal(t, updatedCollection.Name, "Updated Collection Name")
}

func TestDeleteCollection(t *testing.T) {
	db, dbName := setupTestDatabase()
	defer cleanupTestDatabase(db, dbName)

	collection := &database.Collection{Name: "Classic Novels"}
	db.Save(collection)

	httpServer := server.StartRESTServer(db)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/collection/%d", collection.ID), nil)
	rr := httptest.NewRecorder()
	httpServer.Handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusNoContent, rr.Code)

	expected, err := buildSuccessResponse(nil)
	if err != nil {
		t.Fatal(err)
	}
	require.JSONEq(t, expected, rr.Body.String())

	var deletedCollection database.Collection
	err = db.First(&deletedCollection, collection.ID).Error
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestAddAndUpdateBookToCollection(t *testing.T) {
	db, dbName := setupTestDatabase()
	defer cleanupTestDatabase(db, dbName)

	book := &database.Book{Title: "Test Book", Author: "Test Author"}
	collection := &database.Collection{Name: "Classic Novels"}
	db.Save(book)
	db.Save(collection)

	httpServer := server.StartRESTServer(db)

	bookIDJSON := fmt.Sprintf(`{"book_id": %d}`, book.ID)
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://localhost:3000/v1/collection/%d/book", collection.ID), strings.NewReader(bookIDJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	httpServer.Handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusAccepted, rr.Code)

	// Update the book
	book.Title = "Updated Title"
	db.Save(book)

	// Fetch the book from the API
	req, _ = http.NewRequest("GET", fmt.Sprintf("http://localhost:3000/v1/collection/%d/book", collection.ID), nil)
	rr = httptest.NewRecorder()
	httpServer.Handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	// Unmarshal the response body
	var result map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode the response body: %v", err)
	}

	messages, ok := result["message"].([]interface{})
	if !ok {
		t.Fatal("Could not convert 'message' to []interface{}")
	}

	// assuming the book we want is the first in the list
	firstMessage, ok := messages[0].(map[string]interface{})
	if !ok {
		t.Fatal("Could not convert first 'message' to map[string]interface{}")
	}

	title, ok := firstMessage["title"].(string)
	if !ok {
		t.Fatal("Could not convert 'title' to string")
	}

	require.Equal(t, book.Title, title)
}

func TestRemoveBookFromCollection(t *testing.T) {
	db, dbName := setupTestDatabase()
	defer cleanupTestDatabase(db, dbName)

	book := &database.Book{Title: "Test Book", Author: "Test Author"}
	collection := &database.Collection{Name: "Classic Novels"}
	db.Save(book)
	db.Save(collection)

	bookCollection := &database.CollectionBook{BookID: book.ID, CollectionID: collection.ID}
	db.Save(bookCollection)

	httpServer := server.StartRESTServer(db)

	bookIDJSON := fmt.Sprintf(`{"book_id": %d}`, book.ID)
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/collection/%d/book", collection.ID), strings.NewReader(bookIDJSON))
	rr := httptest.NewRecorder()
	httpServer.Handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusNoContent, rr.Code)

	var deletedBookCollection database.CollectionBook
	err := db.Where("book_id = ? AND collection_id = ?", book.ID, collection.ID).First(&deletedBookCollection).Error
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func buildSuccessResponse(message interface{}) (string, error) {
	resp := api.Response{
		Type:       api.SyncResponse,
		Status:     "Success",
		StatusCode: 200,
		Code:       0,
		Error:      "",
	}

	if message != nil {
		resp.Message = message
	}

	// marshal Response instance to JSON
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(jsonResp), nil
}

func setupDatabase(db_name string) *gorm.DB {
	dsn := fmt.Sprintf("host=localhost user=app password=password dbname=%s port=5432 sslmode=disable", db_name)
	newLogger := logger.Default.LogMode(logger.Silent)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}
