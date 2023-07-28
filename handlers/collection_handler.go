package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"grace/database"
	"grace/repositories"
	"grace/response"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slices"
)

type CollectionHandler struct {
	Repository         *repositories.CollectionRepository
	BookRepository     *repositories.BookRepository
	CollectionBookRepo *repositories.CollectionBookRepository
}

func (h *CollectionHandler) CreateCollection(r *http.Request) response.Response {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var collection database.Collection
	if err := dec.Decode(&collection); err != nil {
		return response.BadRequest(fmt.Errorf("Invalid request payload: %w", err))
	}

	if err := h.Repository.Create(&collection); err != nil {
		return response.InternalError(fmt.Errorf("Failed to create collection: %w", err))
	}

	return response.ResponseJsonCustomStatusCode(http.StatusCreated, collection)
}

func (h *CollectionHandler) GetCollection(r *http.Request) response.Response {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return response.BadRequest(fmt.Errorf("Invalid collection ID: %w", err))
	}

	collection, err := h.Repository.FindByID(uint(id))
	if err != nil {
		return response.NotFound(fmt.Errorf("Collection not found: %w", err))
	}

	return response.ResponseJson(collection)
}

func (h *CollectionHandler) UpdateCollection(r *http.Request) response.Response {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return response.BadRequest(fmt.Errorf("Invalid collection ID: %w", err))
	}

	collection, err := h.Repository.FindByID(uint(id))
	if err != nil {
		return response.NotFound(fmt.Errorf("Collection not found: %w", err))
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&collection); err != nil {
		return response.BadRequest(fmt.Errorf("Invalid request payload: %w", err))
	}

	if err := h.Repository.Update(collection); err != nil {
		return response.InternalError(fmt.Errorf("Failed to update collection: %w", err))
	}

	return response.ResponseJson(collection)
}

func (h *CollectionHandler) GetCollections(r *http.Request) response.Response {
	var err error

	collections, err := h.Repository.FindAll()
	if err != nil {
		return response.InternalError(fmt.Errorf("Failed to get collections: %w", err))
	}

	return response.ResponseJson(collections)
}

func (h *CollectionHandler) DeleteCollection(r *http.Request) response.Response {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return response.BadRequest(fmt.Errorf("Invalid collection ID: %w", err))
	}

	collection, err := h.Repository.FindByID(uint(id))
	if err != nil {
		return response.NotFound(fmt.Errorf("Collection not found: %w", err))
	}

	if err := h.Repository.Delete(collection); err != nil {
		return response.InternalError(fmt.Errorf("Failed to delete collection: %w", err))
	}

	return response.ResponseNoContent(http.StatusNoContent)
}

func (h *CollectionHandler) GetCollectionBooks(r *http.Request) response.Response {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return response.BadRequest(fmt.Errorf("Invalid collection ID: %w", err))
	}

	collection, err := h.Repository.FindByID(uint(id))
	if err != nil {
		return response.NotFound(fmt.Errorf("Collection not found: %w", err))
	}

	return response.ResponseJson(collection.Books)
}

func (h *CollectionHandler) AddBookToCollection(r *http.Request) response.Response {
	collectionID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return response.BadRequest(fmt.Errorf("Invalid collection ID: %w", err))
	}

	collection, err := h.Repository.FindByID(uint(collectionID))
	if err != nil {
		return response.NotFound(fmt.Errorf("Collection not found: %w", err))
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var bookIDPayload struct {
		BookID uint `json:"book_id"`
	}
	if err := dec.Decode(&bookIDPayload); err != nil {
		return response.BadRequest(fmt.Errorf("Invalid request payload: %w", err))
	}

	book, err := h.BookRepository.FindByID(bookIDPayload.BookID)
	if err != nil {
		return response.NotFound(fmt.Errorf("Book not found: %w", err))
	}

	var collection_book_id []uint
	for _, collection_book := range collection.Books {
		collection_book_id = append(collection_book_id, collection_book.ID)
	}

	if slices.Contains(collection_book_id, book.ID) {
		return response.ResponseNoContent(http.StatusAccepted)
	}

	if err := h.CollectionBookRepo.AddBookToCollection(book, collection); err != nil {
		return response.InternalError(fmt.Errorf("Failed to add book to collection: %w", err))
	}

	return response.ResponseNoContent(http.StatusAccepted)
}

func (h *CollectionHandler) RemoveBookFromCollection(r *http.Request) response.Response {
	collectionID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return response.BadRequest(fmt.Errorf("Invalid collection ID: %w", err))
	}

	collection, err := h.Repository.FindByID(uint(collectionID))
	if err != nil {
		return response.NotFound(fmt.Errorf("Collection not found: %w", err))
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var bookIDPayload struct {
		BookID uint `json:"book_id"`
	}
	if err := dec.Decode(&bookIDPayload); err != nil {
		return response.BadRequest(fmt.Errorf("Invalid request payload: %w", err))
	}

	book, err := h.BookRepository.FindByID(bookIDPayload.BookID)
	if err != nil {
		return response.NotFound(fmt.Errorf("Book not found: %w", err))
	}

	if err := h.CollectionBookRepo.RemoveBookFromCollection(book, collection); err != nil {
		return response.InternalError(fmt.Errorf("Failed to remove book from collection: %w", err))
	}

	return response.ResponseNoContent(http.StatusNoContent)
}
