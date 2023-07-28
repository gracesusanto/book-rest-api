package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"grace/database"
	"grace/repositories"
	"grace/response"

	"github.com/gorilla/mux"
)

type BookHandler struct {
	Repository *repositories.BookRepository
}

func (h *BookHandler) CreateBook(r *http.Request) response.Response {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var book struct {
		database.Book
		PublishedAt string `json:"published_at"`
	}

	if err := dec.Decode(&book); err != nil {
		return response.BadRequest(fmt.Errorf("Invalid request payload: %w", err))
	}

	publishedAt, err := time.Parse(time.RFC3339, book.PublishedAt)
	if err != nil {
		publishedAt, err = time.Parse("2006-01-02", book.PublishedAt)
		if err != nil {
			return response.BadRequest(fmt.Errorf("Invalid date format for 'published_at': %w", err))
		}
	}

	book.Book.PublishedAt = publishedAt

	if err := h.Repository.Create(&book.Book); err != nil {
		return response.InternalError(fmt.Errorf("Failed to create book: %w", err))
	}

	return response.ResponseJsonCustomStatusCode(http.StatusCreated, book)
}

func (h *BookHandler) GetBook(r *http.Request) response.Response {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return response.BadRequest(fmt.Errorf("Invalid book ID: %w", err))
	}

	book, err := h.Repository.FindByID(uint(id))
	if err != nil {
		return response.NotFound(fmt.Errorf("Book not found: %w", err))
	}

	return response.ResponseJson(book)
}

func (h *BookHandler) UpdateBook(r *http.Request) response.Response {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return response.BadRequest(fmt.Errorf("Invalid book ID: %w", err))
	}

	_, err = h.Repository.FindByID(uint(id))
	if err != nil {
		return response.NotFound(fmt.Errorf("Book not found: %w", err))
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var newBook struct {
		database.Book
		PublishedAt string `json:"published_at"`
	}

	if err := dec.Decode(&newBook); err != nil {
		return response.BadRequest(fmt.Errorf("Invalid request payload: %w", err))
	}

	publishedAt, err := time.Parse(time.RFC3339, newBook.PublishedAt)
	if err != nil {
		publishedAt, err = time.Parse("2006-01-02", newBook.PublishedAt)
		if err != nil {
			return response.BadRequest(fmt.Errorf("Invalid date format for 'published_at': %w", err))
		}
	}

	newBook.Book.PublishedAt = publishedAt

	newBook.Book.ID = uint(id)
	if err := h.Repository.Update(&newBook.Book); err != nil {
		return response.InternalError(fmt.Errorf("Failed to update book: %w", err))
	}

	return response.ResponseJson(newBook)
}

func (h *BookHandler) DeleteBook(r *http.Request) response.Response {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return response.BadRequest(fmt.Errorf("Invalid book ID: %w", err))
	}

	book, err := h.Repository.FindByID(uint(id))
	if err != nil {
		return response.NotFound(fmt.Errorf("Book not found: %w", err))
	}

	if err := h.Repository.Delete(book); err != nil {
		return response.InternalError(fmt.Errorf("Failed to delete book: %w", err))
	}

	return response.ResponseNoContent(http.StatusNoContent)
}

func (h *BookHandler) PatchBook(r *http.Request) response.Response {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return response.BadRequest(fmt.Errorf("Invalid book ID: %w", err))
	}

	book, err := h.Repository.FindByID(uint(id))
	if err != nil {
		return response.NotFound(fmt.Errorf("Book not found: %w", err))
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var payload map[string]interface{}
	if err := dec.Decode(&payload); err != nil {
		return response.BadRequest(fmt.Errorf("Invalid request payload: %w", err))
	}

	// Dynamically update the fields that are present in the payload
	for key, value := range payload {
		switch strings.ToLower(key) {
		case "title":
			book.Title = value.(string)
		case "author":
			book.Author = value.(string)
		case "published_at":
			publishedAt, err := time.Parse(time.RFC3339, value.(string))
			if err != nil {
				publishedAt, err = time.Parse("2006-01-02", value.(string))
				if err != nil {
					return response.BadRequest(fmt.Errorf("Invalid date format for 'published_at': %w", err))
				}
			}
			book.PublishedAt = publishedAt
		case "edition":
			book.Edition = value.(string)
		case "description":
			book.Description = value.(string)
		case "genre":
			book.Genre = value.(string)
		}
	}

	if err := h.Repository.Update(book); err != nil {
		return response.InternalError(fmt.Errorf("Failed to update book: %w", err))
	}

	return response.ResponseJson(book)
}

func (h *BookHandler) GetBooks(r *http.Request) response.Response {
	var err error
	var start time.Time
	var end time.Time

	author := r.URL.Query().Get("author")
	genre := r.URL.Query().Get("genre")

	rstart, rend := r.URL.Query().Get("start"), r.URL.Query().Get("end")
	if rstart != "" {
		start, err = time.Parse("2006-01-02", rstart)
		if err != nil {
			return response.BadRequest(fmt.Errorf("Invalid start date: %w", err))
		}
	}

	if rend != "" {
		end, err = time.Parse("2006-01-02", rend)
		if err != nil {
			return response.BadRequest(fmt.Errorf("Invalid end date: %w", err))
		}
	}

	books, err := h.Repository.Find(author, genre, start, end)
	if err != nil {
		return response.InternalError(fmt.Errorf("Failed to get books: %w", err))
	}

	return response.ResponseJson(books)
}
